package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

func botHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "filename=\"vcData-nobu-bot-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".json\"")
	w.Header().Set("Content-Type", "application/json")

	cards := make([]vc.Card, 0)
	for _, card := range VcData.Cards {
		cardRare := card.CardRarity(VcData)
		if card.IsClosed != 0 ||
			cardRare.Signature == "n" ||
			cardRare.Signature == "hn" ||
			cardRare.Signature == "x" ||
			//cardRare.Signature == "r" ||
			//cardRare.Signature == "hr" ||
			card.AwakensFrom(VcData) != nil ||
			card.PrevEvo(VcData) != nil {
			// don't output low rarities or non-final evos
			continue
		}
		cards = append(cards, card)
	}

	// sort by ID
	sort.Slice(cards, func(i, j int) bool {
		first := cards[i]
		second := cards[j]

		return first.ID < second.ID
	})

	io.WriteString(w, "    [\n")
	cLen := len(cards)
	for idx, card := range cards {
		// to get the image location, we are going to ask Fandom for it:
		// https://valkyriecrusade.fandom.com/index.php?title=Special:FilePath&file=Image Name.jpg
		// this URL returns the actual image location in the HTTP Redirect Location header.
		imgLoc, err := getWikiImageLocation(card.GetEvoImageName(VcData, false) + ".png")
		//imgLoc := ""
		//var err error
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		}
		imgLoc = strings.Replace(imgLoc, "_G.png", "_H.png", -1)
		imgLoc = strings.Replace(imgLoc, "_A.png", "_H.png", -1)

		//lastEvo := card.LastEvo(VcData)

		skills := "" //make([]vc.Skill, 0)

		// for the skills, we specifically want:
		// 1. unique skills (non-awoken and awoken)
		// 2. skills that don't expire (event/thor)

		comma := ""
		if idx+1 < cLen {
			comma = ",\n"
		}

		fmt.Fprintf(w,
			`      {
        "name": "%s",
        "element": "%s",
        "rarity": "%s",
        "skill": [%s],
        "image": "%s",
        "link": "http://valkyriecrusade.fandom.com/wiki/%s"
      }%s`,
			card.Name,
			card.Element(),
			card.Rarity(),
			skills,
			imgLoc,
			url.PathEscape(card.Name),
			comma,
		)

	}
	io.WriteString(w, "\n    ]\n")
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
