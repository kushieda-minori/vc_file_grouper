package nobu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// DbFileLocation location of an existing nobu db file
var DbFileLocation = ""

// DB actual data after being loaded
var DB *Db

// Card card as known by Nobu Bot
type Card struct {
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
		os.Stderr.WriteString(err.Error() + "\n")
	}
	//imgLoc = strings.Replace(imgLoc, "_G.png", "_H.png", -1)
	//imgLoc = strings.Replace(imgLoc, "_A.png", "_H.png", -1)

	imgLocs, err := getWikiImageLocations(c)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
	return Card{
		Name:    c.Name,
		Element: c.Element(),
		Rarity:  c.MainRarity(),
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
	rarity := c.MainRarity()
	for i, card := range *n {
		if card.Name == name && card.Element == element && (card.Rarity == rarity || card.Rarity == c.Rarity()) {
			ref := (*n)[i]
			ref.Skills = newSkills(c)
			if ref.Image == "" {
				imgLoc, err := getWikiImageLocation(c.GetEvoImageName(false) + ".png")
				if err != nil {
					ref.Image = imgLoc
				}
			}
			if ref.Images == nil || len(ref.Images) == 0 {
				imgLocs, err := getWikiImageLocations(c)
				if err != nil {
					ref.Images = imgLocs
				}
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

	*n = append(*n, NewCard(c))

	return false
}

func getWikiImageLocations(c *vc.Card) ([]string, error) {
	ret := make([]string, 0)
	errorMsg := ""
	evos := c.GetEvolutions()
	for _, evoID := range vc.EvoOrder {
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
	if errorMsg != "" {
		return ret, errors.New(errorMsg)
	}
	return ret, nil
}

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
			os.Stderr.WriteString(fmt.Sprintf("Warning: Unable to locate image '%s', status: %d; retry: %d\n", cardImageName, resp.StatusCode, i))
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
