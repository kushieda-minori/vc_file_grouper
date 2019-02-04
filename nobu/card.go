package nobu

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// Card card as known by Nobu Bot
type Card struct {
	Name    string  `json:"name"`
	Element string  `json:"element"`
	Rarity  string  `json:"rarity"`
	Skills  []Skill `json:"skill"`
	Image   string  `json:"image"`
	Link    string  `json:"link"`
}

// NewCard Converts a VC card to a Nobu DB card
func NewCard(c *vc.Card, v *vc.VFile) Card {
	imgLoc, err := getWikiImageLocation(c.GetEvoImageName(v, false) + ".png")
	//imgLoc := ""
	//var err error
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
	imgLoc = strings.Replace(imgLoc, "_G.png", "_H.png", -1)
	imgLoc = strings.Replace(imgLoc, "_A.png", "_H.png", -1)
	return Card{
		Name:    c.Name,
		Element: c.Element(),
		Rarity:  c.MainRarity(),
		Skills:  newSkills(c, v),
		Image:   imgLoc,
		Link: fmt.Sprintf("http://valkyriecrusade.fandom.com/wiki/%s",
			url.PathEscape(c.Name),
		),
	}
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
