package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"zetsuboushita.net/vc_file_grouper/nobu"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func botHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "filename=\"vcData-nobu-bot-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".json\"")
	w.Header().Set("Content-Type", "application/json")

	cards := make([]vc.Card, 0)
	for _, card := range VcData.Cards {
		cardRare := card.CardRarity(VcData)
		evos := card.GetEvolutions(VcData)
		if card.IsClosed != 0 ||
			card.IsRetired() ||
			(len(evos) > 1 && card.EvolutionRank != 0) ||
			cardRare.Signature == "n" ||
			cardRare.Signature == "hn" ||
			cardRare.Signature[0] == 'x' || // ignore normal X and all "Reborn"
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

	nobuCards := make([]nobu.Card, 0)
	for _, card := range cards {
		// to get the image location, we are going to ask Fandom for it:
		// https://valkyriecrusade.fandom.com/index.php?title=Special:FilePath&file=Image Name.jpg
		// this URL returns the actual image location in the HTTP Redirect Location header.
		nobuCards = append(nobuCards, nobu.NewCard(&card, VcData))
	}
	b, err := json.MarshalIndent(nobuCards, "", " ")

	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(b[:]))
	}
}
