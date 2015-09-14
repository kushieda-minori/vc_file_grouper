package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"zetsuboushita.net/vc_file_grouper/vc_grouper"
)

var VcData vc_grouper.VcFile

func usage() {
	os.Stderr.WriteString("You must pass the location of the files.\n" +
		"Usage: " + os.Args[0] + " /path/to/com.nubee.valkyriecrusade/files\n")
}

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		usage()
		return
	}

	err := VcData.Read(os.Args[1])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	//main page
	http.HandleFunc("/", masterDataHandler)
	//image locations
	http.Handle("/cardimages/", http.StripPrefix("/cardimages/", http.FileServer(http.Dir(os.Args[1]+"/card/md"))))
	http.Handle("/cardthumbs/", http.StripPrefix("/cardthumbs/", http.FileServer(http.Dir(os.Args[1]+"/card/thumb"))))
	http.Handle("/cardimagesHD/", http.StripPrefix("/cardimagesHD/", http.FileServer(http.Dir(os.Args[1]+"/../hd/"))))
	http.Handle("/eventimages/", http.StripPrefix("/eventimages/", http.FileServer(http.Dir(os.Args[1]+"/event/largeimage"))))

	//dynamic pages
	http.HandleFunc("/cards/", cardHandler)

	http.ListenAndServe(":8080", nil)
}

func masterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	io.WriteString(w, "<html><body>\n")

	io.WriteString(w, "<a href=\"cards\" >Card List</a><br />\n")
	io.WriteString(w, "<a href=\"cards/table\" >Card List as a Table</a><br />\n")
	io.WriteString(w, "<a href=\"cards/csv\" >Card List as CSV</a><br />\n")

	io.WriteString(w, "</body></html>")
}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	if len(pathParts) > 1 {
		switch pathParts[1] {
		case "csv":
			cardCsvHandler(w, r)
		case "table":
			cardTableHandler(w, r)
		case "detail":
			cardDetailHandler(w, r)
		default:
			http.Error(w, "Unknown card option "+pathParts[1], http.StatusNotFound)
		}
		return
	}
	// render a card list here
	for key, value := range VcData.Cards {
		fmt.Fprintf(w,
			"<div style=\"float: left; margin: 3px\"><img src=\"/cardthumbs/%s\"/><br />%s</div>",
			value.GetImage(),
			value.Name)
		if key > 99 {
			break
		}
	}
}

func cardDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "cards/detail/id"
	if len(pathParts) < 3 {
		http.Error(w, "Invalid card id ", http.StatusNotFound)
		return
	}
	cardId, err := strconv.Atoi(pathParts[2])
	if err != nil || cardId < 1 || cardId > len(VcData.Cards) {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	card := &(VcData.Cards[cardId-1])
	evolutions := getEvolutions(*card)

	amalgamations := make([]vc_grouper.Amalgamation, 0)
	var turnOverTo, turnOverFrom *vc_grouper.Card
	for _, v := range evolutions {
		a := v.GetAmalgamations(&VcData)
		if len(a) > 0 {
			amalgamations = append(amalgamations, a...)
		}
		if v.TransCardId > 0 {
			turnOverTo = v.GetEvoAccident(VcData.Cards)
			break
		} else {
			turnOverFrom = v.IsEvoAccidentOf(VcData.Cards)
			if turnOverFrom != nil {
				break
			}
		}
	}

	var avail string
	if turnOverFrom != nil {
		avail += fmt.Sprintf("{{Card Icon|%s|Evolution Accident}}", turnOverFrom.Name)
	}

	io.WriteString(w, "<html><body>\n")
	fmt.Fprintf(w, "Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>", card.Name)
	io.WriteString(w, "<textarea style=\"width:100%;height:100%\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n|rarity = %s\n|skill = %s\n|skill lv1 = %s\n|skill lv10 = %s\n|procs = %s\n",
		card.GetElement(), card.GetRarity(), card.GetSkill1Name(&VcData), card.GetSkillMin(&VcData), card.GetSkillMax(&VcData), card.GetSkillProcs(&VcData))
	if evo, ok := evolutions["0"]; ok {
		fmt.Fprintf(w, "|cost 0 = %d\n|atk 0 = %d / %d\n|def 0 = %d / %d\n|soldiers 0 = %d / %d\n",
			evo.DeckCost, evo.DefaultOffense, evo.MaxOffense, evo.DefaultDefense, evo.MaxDefense, evo.DefaultFollower, evo.MaxFollower)
	}
	if evo, ok := evolutions["1"]; ok {
		fmt.Fprintf(w, "|cost 1 = %d\n|atk 1 = %d / ?\n|def 1 = %d / ?\n|soldiers 1 = %d / ?\n",
			evo.DeckCost, evo.DefaultOffense, evo.DefaultDefense, evo.DefaultFollower)
	}
	if evo, ok := evolutions["2"]; ok {
		fmt.Fprintf(w, "|cost 2 = %d\n|atk 2 = %d / ?\n|def 2 = %d / ?\n|soldiers 2 = %d / ?\n",
			evo.DeckCost, evo.DefaultOffense, evo.DefaultDefense, evo.DefaultFollower)
	}
	if evo, ok := evolutions["3"]; ok {
		fmt.Fprintf(w, "|cost 3 = %d\n|atk 3 = %d / ?\n|def 3 = %d / ?\n|soldiers 3 = %d / ?\n",
			evo.DeckCost, evo.DefaultOffense, evo.DefaultDefense, evo.DefaultFollower)
	}
	if evo, ok := evolutions["4"]; ok {
		fmt.Fprintf(w, "|cost 4 = %d\n|atk 4 = %d / ?\n|def 4 = %d / ?\n|soldiers 4 = %d / ?\n",
			evo.DeckCost, evo.DefaultOffense, evo.DefaultDefense, evo.DefaultFollower)
	}
	//TODO get amalgamation logic to work here
	// fmt.Fprintf(w,"|cost a = %s\n|atk a = %s\n|def a = %s\n|soldiers a = %s\n",)
	// fmt.Fprintf(w,"|skill a = %s\n|skill a lv1 = %s\n|skill a lv10 = %s\n|proc a = %s\n",)
	if evo, ok := evolutions["G"]; ok {
		fmt.Fprintf(w, "|cost g = %d\n|atk g = %d / ?\n|def g = %d / ?\n|soldiers g = %d / ?\n",
			evo.DeckCost, evo.DefaultOffense, evo.DefaultDefense, evo.DefaultFollower)
		fmt.Fprintf(w, "|skill g = %s\n|skill g lv1 = %s\n|skill g lv10 = %s\n|proc g = %s\n",
			evo.GetSkill1Name(&VcData), evo.GetSkillMin(&VcData), evo.GetSkillMax(&VcData), evo.GetSkillProcs(&VcData))
	}
	fmt.Fprintf(w, "|description = %s\n|friendship = %s\n|login = %s\n|meet = %s\n|battle start = %s\n|battle end = %s\n|friendship max = %s\n|friendship event = %s\n",
		html.EscapeString(card.GetDescription(&VcData)), html.EscapeString(card.GetFriendship(&VcData)),
		html.EscapeString(card.GetLogin(&VcData)), html.EscapeString(card.GetMeet(&VcData)),
		html.EscapeString(card.GetBattleStart(&VcData)), html.EscapeString(card.GetBattleEnd(&VcData)),
		html.EscapeString(card.GetFriendshipMax(&VcData)), html.EscapeString(card.GetFriendshipEvent(&VcData)))
	fmt.Fprintf(w, "|availability = %s\n", avail)
	// fmt.Fprintf(w,"|likeability 0 = %s\n|likeability 1 = %s\n|likeability 2 = %s\n|likeability 3 = %s\n|likeability 4 = %s\n|likeability 5 =%s\n",)
	io.WriteString(w, "}}")

	if turnOverTo != nil {
		fmt.Fprintf(w, "\n==Evolution Accident==\n{{Card Icon|%s}}\n", turnOverTo.Name)
	}

	//Write out amalgamations here

	io.WriteString(w, "</textarea></body></html>")
}

func getEvolutions(card vc_grouper.Card) map[string]vc_grouper.Card {
	ret := make(map[string]vc_grouper.Card)
	for _, val := range VcData.Cards {
		if card.CardCharaId == val.CardCharaId {
			if val.GetRarity()[0] == 'G' {
				ret["G"] = val
			} else {
				ret[strconv.Itoa(val.EvolutionRank)] = val
			}
		}
	}
	return ret
}

func cardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	io.WriteString(w, "<html><body>\n")
	io.WriteString(w, "_id,card_no,name,evolution_rank,Rarity,Element,deck_cost,default_offense,default_defense,default_follower,max_offense,max_defense,max_follower,Skill 1 Name,Skill Min,Skill Max,Skill Procs,Target Scope,Target Logic,Skill 2,Skill Special,Description,Friendship,Login,Meet,Battle Start,Battle End,Friendship Max,Friendship Event<br>\n")
	for _, value := range VcData.Cards {
		fmt.Fprintf(w, "%d, %05d, %s, %d, %s, %s, %d, %d, %d, %d, %d, %d, %d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s <br />",
			value.Id, value.CardNo, value.Name, value.EvolutionRank, value.GetRarity(), value.GetElement(),
			value.DeckCost, value.DefaultOffense, value.DefaultDefense, value.DefaultFollower, value.MaxOffense,
			value.MaxDefense, value.MaxFollower, value.GetSkill1Name(&VcData), value.GetSkillMin(&VcData), value.GetSkillMax(&VcData),
			value.GetSkillProcs(&VcData), value.GetSkillTarget(&VcData), value.GetSkillTargetLogic(&VcData), value.GetSkill2Name(&VcData),
			value.GetSpecialSkill1Name(&VcData), value.GetDescription(&VcData), value.GetFriendship(&VcData), value.GetLogin(&VcData),
			value.GetMeet(&VcData), value.GetBattleStart(&VcData), value.GetBattleEnd(&VcData), value.GetFriendshipMax(&VcData),
			value.GetFriendshipEvent(&VcData))
	}

	io.WriteString(w, "</body></html>")
}

func cardTableHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	io.WriteString(w, "<html><body>\n")
	io.WriteString(w, "<table><thead><tr><th>_id</th><th>card_no</th><th>name</th><th>evolution_rank</th><th>Rarity</th><th>Element</th><th>deck_cost</th><th>default_offense</th><th>default_defense</th><th>default_follower</th><th>max_offense</th><th>max_defense</th><th>max_follower</th><th>Skill 1 Name</th><th>Skill Min</th><th>Skill Max</th><th>Skill Procs</th><th>Target Scope</th><th>Target Logic</th><th>Skill 2</th><th>Skill Special</th><th>Description</th><th>Friendship</th><th>Login</th><th>Meet</th><th>Battle Start</th><th>Battle End</th><th>Friendship Max</th><th>Friendship Event</th></tr></thead><tbody>\n")
	for _, value := range VcData.Cards {
		fmt.Fprintf(w, "<tr><td>%d</td><td>%05d</td><td><a href=\"/cards/detail/%d\">%s</a></td><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></td>\n",
			value.Id, value.CardNo, value.Id, value.Name, value.EvolutionRank, value.GetRarity(), value.GetElement(),
			value.DeckCost, value.DefaultOffense, value.DefaultDefense, value.DefaultFollower, value.MaxOffense,
			value.MaxDefense, value.MaxFollower, value.GetSkill1Name(&VcData), value.GetSkillMin(&VcData), value.GetSkillMax(&VcData),
			value.GetSkillProcs(&VcData), value.GetSkillTarget(&VcData), value.GetSkillTargetLogic(&VcData), value.GetSkill2Name(&VcData),
			value.GetSpecialSkill1Name(&VcData), value.GetDescription(&VcData), value.GetFriendship(&VcData), value.GetLogin(&VcData),
			value.GetMeet(&VcData), value.GetBattleStart(&VcData), value.GetBattleEnd(&VcData), value.GetFriendshipMax(&VcData),
			value.GetFriendshipEvent(&VcData))
	}

	io.WriteString(w, "</tbody></table></body></html>")
}
