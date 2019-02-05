package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/nobu"
	"zetsuboushita.net/vc_file_grouper/vc"
)

var nobuDB *nobu.Db

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

func botConfigHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Bot Config</title></head><body>\n")

	// check form value and update if valid
	newpath := r.FormValue("path")
	if newpath != "" {
		if _, err := os.Stat(newpath); os.IsNotExist(err) {
			io.WriteString(w, "<div>Invalid new path specified</div>")
		} else {
			nobu.DbFileLocation = newpath
			if nobuDB, err = nobu.LoadDb(); err != nil {
				fmt.Fprintf(w, "<div>%s</div>", err.Error())
			} else {
				io.WriteString(w, "<div>Success</div>")
			}
		}
	}
	// write out the form
	fmt.Fprintf(w, `<form method="post">
<label for="f_path">Data Path</label>
<input id="f_path" name="path" value="%s" style="width:300px"/>
<button type="submit">Submit</button>
<p><a href="/">back</a></p>
</form>`,
		html.EscapeString(nobu.DbFileLocation),
	)
	io.WriteString(w, "</body></html>")
}

func botUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Bot DB</title></head><body>\n")
	cards := make([]vc.Card, 0)
	for _, card := range VcData.Cards {
		cardRare := card.CardRarity(VcData)
		evos := card.GetEvolutions(VcData)
		evosLen := len(evos)
		skill1 := card.Skill1(VcData)
		skill1Min := ""
		if skill1 != nil {
			skill1Min = skill1.SkillMin()
		}
		if card.IsClosed != 0 ||
			card.IsRetired() ||
			//card.EvolutionRank < 0 || // skip any cards with evo rank < 0
			(evosLen > 1 && card.EvolutionRank > 0) || // skip any card that is not the first evo
			card.Element() == "Special" ||
			cardRare.Signature == "n" ||
			cardRare.Signature == "hn" ||
			cardRare.Signature[0] == 'x' || // ignore normal X and all "Reborn"
			cardRare.Signature == "r" ||
			cardRare.Signature == "hr" ||
			card.AwakensFrom(VcData) != nil || // ignore G* cards that are actually awoken
			card.PrevEvo(VcData) != nil || // ignore cards that have a previous evolution
			strings.Contains(skill1Min, "Battle EXP +5%") {
			// don't output low rarities or non-final evos
			if card.IsClosed == 0 &&
				!card.IsRetired() &&
				(evosLen == 1 || card.EvolutionRank == 0) &&
				(card.MainRarity() == "SR" || card.MainRarity() == "UR" || card.MainRarity() == "LR") {
				os.Stderr.WriteString(fmt.Sprintf("Skipped and %s: %s - evo: %d/%d\n", cardRare.Signature, card.Name, card.EvolutionRank, evosLen))
			}
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

	for _, card := range cards {
		// to get the image location, we are going to ask Fandom for it:
		// https://valkyriecrusade.fandom.com/index.php?title=Special:FilePath&file=Image Name.jpg
		// this URL returns the actual image location in the HTTP Redirect Location header.
		nobuDB.AddOrUpdate(&card, VcData)
	}

	b, err := json.MarshalIndent(nobuDB, "", " ")
	if err != nil {
		fmt.Fprintf(w, "<div>%s</div>", err.Error())
	} else {
		// overwrite the old file
		err = ioutil.WriteFile(nobu.DbFileLocation, b, os.FileMode(0655))
		if err != nil {
			fmt.Fprintf(w, "<div>%s</div>", err.Error())
		} else {
			io.WriteString(w, "<div>Success</div>")
		}
	}
	io.WriteString(w, "</body></html>")
}
