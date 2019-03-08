package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// DbFileLocation location of an existing nobu db file
var DbFileLocation = ""

// DB actual data after being loaded
var DB *Db

// Card card as known by Nobu Bot
type Card struct {
	VCID            int                  `json:"vcID"`
	Name            string               `json:"name"`
	Element         string               `json:"element"`
	Rarity          string               `json:"rarity"`
	EvoAccidentFrom string               `json:"evoAccidentFrom"`
	EvoAccidentTo   string               `json:"evoAccidentTo"`
	Skills          []Skill              `json:"skill"`
	Amalgamations   []AmalgamationRecipe `json:"amalgamations"`
	Link            string               `json:"link"`
	Images          []string             `json:"images"` // contains all images (not icons)
}

// Db list of cards Nobu-bot knows about
type Db []Card

// LoadDb loads an existing Db
func LoadDb() error {
	// ensure our path has been set
	if DbFileLocation == "" {
		return errors.New("Nobu DB Location not set")
	}
	// check if the file exists
	if _, err := os.Stat(DbFileLocation); os.IsNotExist(err) {
		return errors.New("no such file or directory: " + DbFileLocation)
	}

	// load the existing data from disk
	data, err := ioutil.ReadFile(DbFileLocation)
	if err != nil {
		return err
	}

	v := make(Db, 0)

	// decode the main file
	err = json.Unmarshal(data[:], &v)
	if err != nil {
		debug.PrintStack()
		return err
	}
	DB = &v
	return nil
}

// NewCard Converts a VC card to a Nobu DB card
func NewCard(vcCard *vc.Card) Card {
	imgLocs, err := getWikiImageLocations(vcCard)
	if err != nil {
		log.Printf(err.Error() + "\n")
	}
	evos := vcCard.GetEvolutionCards()
	turnoverFrom, turnoverTo := getTurnOver(evos)
	rarity := evos.MinimumEvolutionRank()
	return Card{
		VCID:            vcCard.ID,
		Name:            vcCard.Name,
		Element:         vcCard.Element(),
		Rarity:          rarity,
		EvoAccidentFrom: turnoverFrom,
		EvoAccidentTo:   turnoverTo,
		Skills:          newSkills(vcCard),
		Amalgamations:   getAmals(evos),
		Link: fmt.Sprintf("https://valkyriecrusade.fandom.com/wiki/%s",
			url.PathEscape(vcCard.Name),
		),
		Images: imgLocs,
	}
}

func getTurnOver(evos vc.CardList) (turnoverFrom, turnoverTo string) {
	for _, evo := range evos {
		accTo := evo.EvoAccident()
		accFrom := evo.EvoAccidentOf()
		if accFrom != nil && turnoverFrom == "" {
			turnoverFrom = accFrom.Name
		}
		if accTo != nil && turnoverTo == "" {
			turnoverTo = accTo.Name
		}
	}
	return
}

func getAmals(evos vc.CardList) []AmalgamationRecipe {
	ret := make([]AmalgamationRecipe, 0)
	for _, evo := range evos {
		amals := evo.Amalgamations()
		l := len(amals)
		if l > 0 {
			for _, amal := range amals {
				ret = append(ret, newRecipe(amal))
			}
		}
	}
	return ret
}

// AddOrUpdate Checks if the card exists in the DB or is not yet there.
// If the card exists, then the skill information is updated.
// If the card does not exists, it is added to the end of the array.
// Since the "DB" is not indexed, this call is O(N) scanning the "DB"
// for every add/updated.
// If the card is found, true is returned, if a new card is added, false is returned.
func (botDB *Db) AddOrUpdate(vcCard *vc.Card) bool {
	name := vcCard.Name
	element := vcCard.Element()
	mainRarity := vcCard.MainRarity()
	evos := vcCard.GetEvolutionCards()
	turnoverFrom, turnoverTo := getTurnOver(evos)
	rarity := evos.MinimumEvolutionRank()
	for i, botCard := range *botDB {
		if botCard.VCID == vcCard.ID || (strings.ToUpper(strings.TrimSpace(botCard.Name)) == strings.ToUpper(name) &&
			botCard.Element == element &&
			(botCard.Rarity == mainRarity || botCard.Rarity == vcCard.Rarity())) {

			log.Printf("Card %s already exists, updating.", vcCard.Name)
			ref := &((*botDB)[i]) // get a reference we can update

			ref.VCID = vcCard.ID
			ref.Name = vcCard.Name // ensure any oddities are taken care of
			ref.Rarity = rarity
			ref.EvoAccidentFrom = turnoverFrom
			ref.EvoAccidentTo = turnoverTo
			ref.Skills = newSkills(vcCard)
			ref.Amalgamations = getAmals(evos)
			newPath := fmt.Sprintf("https://valkyriecrusade.fandom.com/wiki/%s",
				url.PathEscape(name),
			)
			if ref.Link != newPath {
				ref.Link = newPath
			}
			if ref.Images == nil || len(ref.Images) == 0 || len(ref.Images) != len(vcCard.EvosWithDistinctImages(false)) {
				imgLocs, err := getWikiImageLocations(vcCard)
				if err != nil {
				}
				ref.Images = imgLocs
			}
			return true
		}
	}

	log.Printf("***** Card %s is new, adding. *****", vcCard.Name)
	*botDB = append(*botDB, NewCard(vcCard))

	return false
}

func getWikiImageLocations(vcCard *vc.Card) ([]string, error) {
	log.Printf("Locating images for Card %s", vcCard.Name)
	ret := make([]string, 0)
	errorMsg := ""
	evos := vcCard.GetEvolutions()
	evoImages := vcCard.EvosWithDistinctImages(false)
	imageNamesSeen := make(map[string]struct{})
	log.Printf("Looking up images for %d evolutions on card %d:%s", len(evoImages), vcCard.ID, vcCard.Name)
	for _, evoID := range evoImages {
		if evo, ok := evos[evoID]; ok {
			imgName := evo.GetEvoImageName(false) + ".png"
			if _, seen := imageNamesSeen[imgName]; !seen {
				imageNamesSeen[imgName] = struct{}{}
				imgLoc, err := getWikiImageLocation(imgName)
				if err != nil {
					errorMsg += "|" + err.Error()
				}
				if imgLoc != "" {
					ret = append(ret, imgLoc)
				}
			}
		}
	}
	log.Printf("Found %d images", len(ret))
	if errorMsg != "" {
		return ret, errors.New(errorMsg)
	}
	return ret, nil
}

// func getWikiImageLocation(cardImageName string) (string, error) {
// 	return "", nil
// }

func getWikiImageLocation(cardImageName string) (string, error) {
	// to get the image location, we are going to ask Fandom for it:
	// https://valkyriecrusade.fandom.com/index.php?title=Special:FilePath&file=Image Name.jpg
	// this URL returns the actual image location in the HTTP Redirect Location header.
	log.Printf("Looking up image %s", cardImageName)
	myURL := "https://valkyriecrusade.fandom.com/index.php?title=Special:FilePath&file=" + url.QueryEscape(cardImageName)
	nextURL := myURL
	var i int
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	for i < 10 {
		resp, err := client.Get(nextURL)

		if err != nil {
			return "", err
		}

		if resp.StatusCode == 200 {
			// found our last redirect
			break
		} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// there was a problem with the request itself (i.e. not found)
			return "", fmt.Errorf("Unable to locate image '%s', status: %d", cardImageName, resp.StatusCode)
		} else if resp.StatusCode >= 500 {
			// there was a problem with the server, we can retry these.
			log.Printf("Warning: Unable to locate image '%s', status: %d; retry: %d\n", cardImageName, resp.StatusCode, i)
			// in case the problem is rate limiting, slow down for a bit.
			time.Sleep(2 * time.Second)
		}
		if resp.Header.Get("Location") != "" {
			nextURL = resp.Header.Get("Location")
		}
		i++
	}
	if nextURL != myURL {
		u, err := url.Parse(nextURL)
		if err != nil {
			return "", err
		}
		u.RawQuery = ""
		return u.String(), nil
	}
	return "", errors.New("Unable to locate image")
}
