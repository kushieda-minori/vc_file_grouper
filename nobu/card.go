package nobu

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
	VCID    int      `json:"vcID"`
	Name    string   `json:"name"`
	Element string   `json:"element"`
	Rarity  string   `json:"rarity"`
	Skills  []Skill  `json:"skill"`
	Image   string   `json:"image"`  // will be phased out
	Images  []string `json:"images"` // contains all images (not icons)
	Link    string   `json:"link"`
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
func NewCard(c *vc.Card) Card {
	imgLoc, err := getWikiImageLocation(c.GetEvoImageName(false) + ".png")
	//imgLoc := ""
	//var err error
	if err != nil {
		log.Printf(err.Error() + "\n")
	}
	//imgLoc = strings.Replace(imgLoc, "_G.png", "_H.png", -1)
	//imgLoc = strings.Replace(imgLoc, "_A.png", "_H.png", -1)

	imgLocs, err := getWikiImageLocations(c)
	if err != nil {
		log.Printf(err.Error() + "\n")
	}
	rarity := c.MainRarity()
	if len(c.GetEvolutions()) == 1 {
		rarity = c.Rarity()
	}
	return Card{
		VCID:    c.ID,
		Name:    c.Name,
		Element: c.Element(),
		Rarity:  rarity,
		Skills:  newSkills(c),
		Image:   imgLoc,
		Images:  imgLocs,
		Link: fmt.Sprintf("https://valkyriecrusade.fandom.com/wiki/%s",
			url.PathEscape(c.Name),
		),
	}
}

// AddOrUpdate Checks if the card exists in the DB or is not yet there.
// If the card exists, then the skill information is updated.
// If the card does not exists, it is added to the end of the array.
// Since the "DB" is not indexed, this call is O(N) scanning the "DB"
// for every add/updated.
// If the card is found, true is returned, if a new card is added, false is returned.
func (n *Db) AddOrUpdate(c *vc.Card) bool {
	name := c.Name
	element := c.Element()
	mainRarity := c.MainRarity()
	rarity := c.MainRarity()
	if len(c.GetEvolutions()) == 1 {
		rarity = c.Rarity()
	}
	for i, card := range *n {
		if card.VCID == c.ID || (card.VCID == 0 &&
			strings.ToUpper(strings.TrimSpace(card.Name)) == strings.ToUpper(name) &&
			card.Element == element &&
			(card.Rarity == mainRarity || card.Rarity == c.Rarity())) {

			log.Printf("Card %s already exists, updating.", c.Name)
			ref := &((*n)[i]) // get a reference we can update

			ref.VCID = c.ID
			ref.Name = c.Name // ensure any oddities are taken care of
			ref.Rarity = rarity
			ref.Skills = newSkills(c)
			if ref.Image == "" {
				imgLoc, err := getWikiImageLocation(c.GetEvoImageName(false) + ".png")
				if err != nil {
					ref.Image = imgLoc
				}
			}
			if ref.Images == nil || len(ref.Images) == 0 || len(ref.Images) != len(c.EvosWithDistinctImages(false)) {
				imgLocs, err := getWikiImageLocations(c)
				if err != nil {
				}
				ref.Images = imgLocs
			}
			newPath := fmt.Sprintf("https://valkyriecrusade.fandom.com/wiki/%s",
				url.PathEscape(name),
			)
			if ref.Link != newPath {
				ref.Link = newPath
			}
			return true
		}
	}

	log.Printf("Card %s is new, adding.", c.Name)
	*n = append(*n, NewCard(c))

	return false
}

func getWikiImageLocations(c *vc.Card) ([]string, error) {
	log.Printf("Locating images for Card %s", c.Name)
	ret := make([]string, 0)
	errorMsg := ""
	evos := c.GetEvolutions()
	evoImages := c.EvosWithDistinctImages(false)
	for _, evoID := range evoImages {
		if evo, ok := evos[evoID]; ok {
			imgName := evo.GetEvoImageName(false) + ".png"
			imgLoc, err := getWikiImageLocation(imgName)
			if err != nil {
				errorMsg += "|" + err.Error()
			}
			if imgLoc != "" {
				ret = append(ret, imgLoc)
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
