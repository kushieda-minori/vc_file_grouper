package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/nobu"
	"zetsuboushita.net/vc_file_grouper/vc"
)

// BotHandler Generates a new Bot database
func BotHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "filename=\"vcData-nobu-bot-cards-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".json\"")
	w.Header().Set("Content-Type", "application/json")

	cards := make([]vc.Card, 0)
	for _, card := range vc.Data.Cards {
		if shouldExcludeCard(&card) {
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
		nobuCards = append(nobuCards, nobu.NewCard(&card))
	}
	b, err := json.MarshalIndent(nobuCards, "", " ")

	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(b[:]))
	}
}

// BotConfigHandler configures the data location for an existing bot DB
func BotConfigHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Bot Config</title></head><body>\n")

	// check form value and update if valid
	newpath := r.FormValue("path")
	if newpath != "" {
		if _, err := os.Stat(newpath); os.IsNotExist(err) {
			io.WriteString(w, "<div>Invalid new path specified</div>")
		} else {
			nobu.DbFileLocation = newpath
			if err = nobu.LoadDb(); err != nil {
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

// BotUpdateHandler Updates the existing bot DB with new/missing card info
func BotUpdateHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Bot DB</title></head><body>\n")
	namesOfCards := vc.CardsByName()

	bl := len(*nobu.DB)
	lNames := len(namesOfCards)
	i := 0
	for name, cards := range namesOfCards {
		i++
		card := firstCard(&cards)
		if card == nil {
			log.Printf("%d/%d ***********Skipping %d Cards with name %s\n", i, lNames, len(cards), name)
		} else {
			log.Printf("%d/%d Adding/Updating Bot Card %s\n", i, lNames, name)
			nobu.DB.AddOrUpdate(card)
		}
	}
	log.Printf("Bot DB changed from %d records to %d", bl, len(*nobu.DB))

	// sort the cards by VC release IDs
	sort.Slice(*nobu.DB, func(i, j int) bool {
		first := (*nobu.DB)[i]
		second := (*nobu.DB)[j]

		return first.VCID < second.VCID
	})

	b, err := json.MarshalIndent(nobu.DB, "", " ")
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
	log.Println("finished updating Bot Cards")
}

func firstCard(cards *[]vc.Card) *vc.Card {
	sort.Slice(*cards, func(a, b int) bool {
		evoCmp := (*cards)[a].EvolutionRank == (*cards)[b].EvolutionRank
		if evoCmp {
			return (*cards)[a].ID < (*cards)[b].ID
		}
		return (*cards)[a].EvolutionRank < (*cards)[b].EvolutionRank
	})
	for i := range *cards {
		card := &((*cards)[i])
		if !shouldExcludeCard(card) {
			return card
		}
	}
	return nil
}

func shouldExcludeCard(card *vc.Card) bool {
	cardRare := card.CardRarity()
	evos := card.GetEvolutions()
	evosLen := len(evos)
	// skill1 := card.Skill1()
	// skill1Min := ""
	// if skill1 != nil {
	// 	skill1Min = skill1.SkillMin()
	// }
	return card.IsClosed != 0 ||
		card.IsRetired() ||
		//card.EvolutionRank < 0 || // skip any cards with evo rank < 0
		(evosLen > 1 && card.PrevEvo() != nil) || // skip any card that is not the first evo
		(card.Element() == "Special" && nameNotAllowed(card.Name)) ||
		cardRare.Signature == "n" ||
		cardRare.Signature == "hn" ||
		cardRare.Signature == "x" || // ignore normal X
		//cardRare.Signature == "r" ||
		//cardRare.Signature == "hr" ||
		(card.EvoIsAwoken() && card.AwakensFrom() != nil) || // ignore G* cards that are actually awoken
		(card.EvoIsReborn() && card.RebirthsFrom() != nil) || // ignore Rebirth cards that are actually reborn
		//strings.Contains(skill1Min, "Battle EXP +5%") ||
		cardIsOnlyAmalMaterial(card) ||
		false
}

func nameNotAllowed(name string) bool {
	name = strings.TrimSpace(strings.ToUpper(name))
	return !(strings.Contains(name, "SLIME") ||
		strings.Contains(name, "GOLD GIRL") ||
		strings.Contains(name, "MEDAL GIRL") ||
		strings.Contains(name, "MIRROR MAIDEN") ||
		false)
}

func cardIsOnlyAmalMaterial(card *vc.Card) bool {
	return card.HasAmalgamation() &&
		((card.DefaultDefense == card.MaxDefense &&
			card.DefaultOffense == card.MaxOffense &&
			card.DefaultFollower == card.MaxFollower) ||
			(card.SkillMin() == card.SkillMax() && strings.Contains(card.SkillMin(), "Battle EXP +5%")))
}
