package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

	io.WriteString(w, "<a href=\"cards\" >Card List</a>\n")
	io.WriteString(w, "<a href=\"cards/csv\" >Card List as CSV</a>\n")

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
		fmt.Fprintf(w, "<tr><td>%d</td><td>%05d</td><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></td>\n",
			value.Id, value.CardNo, value.Name, value.EvolutionRank, value.GetRarity(), value.GetElement(),
			value.DeckCost, value.DefaultOffense, value.DefaultDefense, value.DefaultFollower, value.MaxOffense,
			value.MaxDefense, value.MaxFollower, value.GetSkill1Name(&VcData), value.GetSkillMin(&VcData), value.GetSkillMax(&VcData),
			value.GetSkillProcs(&VcData), value.GetSkillTarget(&VcData), value.GetSkillTargetLogic(&VcData), value.GetSkill2Name(&VcData),
			value.GetSpecialSkill1Name(&VcData), value.GetDescription(&VcData), value.GetFriendship(&VcData), value.GetLogin(&VcData),
			value.GetMeet(&VcData), value.GetBattleStart(&VcData), value.GetBattleEnd(&VcData), value.GetFriendshipMax(&VcData),
			value.GetFriendshipEvent(&VcData))
	}

	io.WriteString(w, "</tbody></table></body></html>")
}
