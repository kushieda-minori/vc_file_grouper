package handler

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/util"
	"zetsuboushita.net/vc_file_grouper/vc"
)

// CardHandler shows cards in order
func CardHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Cards</title></head><body>\n")
	for _, card := range vc.Data.Cards {
		fmt.Fprintf(w,
			"<div style=\"float: left; margin: 3px\"><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></div>",
			card.Image(),
			card.ID,
			card.Name)
	}
	io.WriteString(w, "</body></html>")
}

// CardLevelHandler shows card level information
func CardLevelHandler(w http.ResponseWriter, r *http.Request) {
	header := `{| class="article-table" style="float:left"
!Lvl!!To Next Lvl!!Total Needed`
	header2 := `{| class="article-table" style="float:left"
!Lvl!!{{Icon|gem}} Needed`
	header3 := `{| class="article-table" style="float:left"
!Lvl!!{{Icon|gold}} Needed!!{{Icon|iron}} Needed!!{{Icon|ether}} Needed!!{{Icon|gem}} Needed`

	genLevels := func(levels []vc.CardLevel) {
		io.WriteString(w, header)
		l := len(levels)
		for i, lvl := range levels {
			nxt := 0
			if i+1 < l {
				nxt = levels[i+1].Exp - lvl.Exp
			}
			fmt.Fprintf(w, `
|-
|%d||%d||%d`,
				lvl.ID,
				nxt,
				lvl.Exp,
			)
			if (i+1)%25 == 0 && i+1 < l {
				io.WriteString(w, "\n|}\n\n")
				io.WriteString(w, header)
			}
		}
		io.WriteString(w, "\n|}\n")
	}

	io.WriteString(w, "<html><head><title>Card Levels</title></head><body>\n")
	io.WriteString(w, "\nN-GUR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(vc.Data.CardLevels)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />XSR and XUR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(vc.Data.CardLevelsX)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />LR-GLR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(vc.Data.CardLevelsLR)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />XLR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(vc.Data.CardLevelsXLR)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />LR Resources<br/><textarea rows=\"25\" cols=\"80\">")

	io.WriteString(w, header2)
	l := len(vc.Data.LevelLRResources)
	for i, lvl := range vc.Data.LevelLRResources {
		fmt.Fprintf(w, `
|-
|%d||%d`,
			lvl.ID,
			lvl.Elixir,
		)
		if (i+1)%25 == 0 && i+1 < l {
			io.WriteString(w, "\n|}\n\n")
			io.WriteString(w, header2)
		}
	}
	io.WriteString(w, "\n|}\n</textarea>")

	io.WriteString(w, "\n<br />XSR & XUR Resources<br/><textarea rows=\"25\" cols=\"80\">")

	io.WriteString(w, header3)
	l = len(vc.Data.LevelXResources)
	for i, lvl := range vc.Data.LevelXResources {
		fmt.Fprintf(w, `
|-
|%d||%d||%d||%d||%d`,
			lvl.ID,
			lvl.Gold,
			lvl.Iron,
			lvl.Ether,
			lvl.Elixir,
		)
		if (i+1)%25 == 0 && i+1 < l {
			io.WriteString(w, "\n|}\n\n")
			io.WriteString(w, header3)
		}
	}

	io.WriteString(w, "\n|}\n</textarea>")
	io.WriteString(w, "\n<br />XLR Resources<br/><textarea rows=\"25\" cols=\"80\">")

	io.WriteString(w, header3)
	l = len(vc.Data.LevelXLRResources)
	for i, lvl := range vc.Data.LevelXLRResources {
		fmt.Fprintf(w, `
|-
|%d||%d||%d||%d||%d`,
			lvl.ID,
			lvl.Gold,
			lvl.Iron,
			lvl.Ether,
			lvl.Elixir,
		)
		if (i+1)%25 == 0 && i+1 < l {
			io.WriteString(w, "\n|}\n\n")
			io.WriteString(w, header3)
		}
	}
	io.WriteString(w, "\n|}\n</textarea>")
	io.WriteString(w, "</body></html>")
}

// CardDetailHandler shows details of a single card formatted for use in the Wiki
func CardDetailHandler(w http.ResponseWriter, r *http.Request) {
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
	cardID, err := strconv.Atoi(pathParts[2])
	if err != nil || cardID < 1 || cardID > len(vc.Data.Cards) {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	card := vc.CardScan(cardID)
	if card == nil {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	log.Printf("Found Card %s: (%d)", card.Name, card.ID)
	evolutions := card.GetEvolutions()
	amalgamations := getAmalgamations(evolutions)

	log.Printf("Card %s: (%d) has %d Evos and %d Amalgamations",
		card.Name,
		card.ID,
		len(evolutions),
		len(amalgamations),
	)
	var firstEvo, lastEvo *vc.Card
	var evokeys []string // cache of actual evos for this card
	for _, k := range vc.EvoOrder {
		if evo, ok := evolutions[k]; ok {
			evokeys = append(evokeys, k)
			if firstEvo == nil {
				log.Printf("Found Card %s: (%d) First Evo: %s: (%d)",
					card.Name,
					card.ID,
					evo.Name,
					evo.ID,
				)
				firstEvo = evo
			}
			if k == "H" {
				log.Printf("Found Card %s: (%d) Last Evo: %s: (%d)",
					card.Name,
					card.ID,
					evo.Name,
					evo.ID,
				)
				lastEvo = evo
			}
		}
	}

	var turnOverTo, turnOverFrom *vc.Card
	if firstEvo.ID > 0 {
		if firstEvo.TransCardID > 0 {
			turnOverTo = firstEvo.EvoAccident()
		} else {
			turnOverFrom = firstEvo.EvoAccidentOf()
		}
	}

	var avail string

	if _, ok := evolutions["A"]; ok {
		avail += " [[Amalgamation]]"
	} else if _, ok := evolutions["GA"]; ok {
		avail += " [[Amalgamation]]"
	} else if _, ok := evolutions["XA"]; ok {
		avail += " [[Amalgamation]]"
	} else {
		for _, evo := range evolutions {
			if evo.IsAmalgamation() {
				avail += " [[Amalgamation]]"
				break
			}
		}
	}

	cardName := card.Name
	if len(cardName) == 0 {
		cardName = firstEvo.Image()
	}

	prevCard, prevCardName := getPrevious(card)
	nextCard, nextCardName := getNext(card)

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", cardName)

	if prevCardName != "" {
		fmt.Fprintf(w, "<div style=\"float:left; width: 33%%;\"><a href=\"%d\">&lt;&lt; %s &lt;&lt;</a></div>\n", prevCard.ID, prevCardName)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}
	fmt.Fprint(w, "<div style=\"float:left; width: 33%;text-align:center;\"><a href=\"../table/\">All Cards</a></div>\n")
	if nextCardName != "" {
		fmt.Fprintf(w, "<div style=\"float:right; width: 33%%;;text-align:right;\"><a href=\"%d\">&gt;&gt; %s &gt;&gt;</a></div>\n", nextCard.ID, nextCardName)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}

	fmt.Fprintf(w, "<div style=\"clear:both;float:left\">Edit on the <a href=\"https://valkyriecrusade.fandom.com/wiki/%s?action=edit\">fandom</a>\n<br /></div>", cardName)
	io.WriteString(w, "<div><textarea readonly=\"readonly\" style=\"width:100%;height:450px\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n|rarity = %s\n", card.Element(), card.MainRarity())

	skillMap := make(map[string]string)

	skillEvoMod := ""
	if firstEvo.Rarity()[0] == 'G' {
		skillEvoMod = "g"
	} else if len(firstEvo.Rarity()) > 2 && firstEvo.Rarity()[0] == 'X' {
		skillEvoMod = "x"
	} else if evo, ok := evolutions["A"]; ok && firstEvo.ID == evo.ID {
		// if the first evo is the amalgamation evo...
		skillEvoMod = ""
	}

	skillMap[skillEvoMod] = printWikiSkill(firstEvo.Skill1(), lastEvo.Skill1(), skillEvoMod)
	skillMap[skillEvoMod+"2"] = printWikiSkill(firstEvo.Skill2(), lastEvo.Skill2(), skillEvoMod+"2")
	skillMap[skillEvoMod+"3"] = printWikiSkill(firstEvo.Skill3(), lastEvo.Skill3(), skillEvoMod+"3")
	skillMap[skillEvoMod+"3"] = printWikiSkill(firstEvo.Skill3(), lastEvo.Skill3(), skillEvoMod+"3")
	skillMap[skillEvoMod+"t"] = printWikiSkill(firstEvo.ThorSkill1(), lastEvo.ThorSkill1(), skillEvoMod+"t")

	// add amal skills as long as the first evo wasn't the amal
	if evo, ok := evolutions["A"]; ok && firstEvo.ID != evo.ID {
		aSkillName := evo.Skill1Name()
		if aSkillName != firstEvo.Skill1Name() {
			skillMap["a"] = printWikiSkill(evo.Skill1(), nil, "a")
		}
		if _, ok := evolutions["GA"]; ok {
			skillMap["ga"] = printWikiSkill(evo.Skill1(), nil, "ga")
		}
		if _, ok := evolutions["XA"]; ok {
			skillMap["xa"] = printWikiSkill(evo.Skill1(), nil, "xa")
		}
	}
	// add awoken skills as long as the first evo wasn't awoken
	if evo, ok := evolutions["G"]; ok && firstEvo.ID != evo.ID {
		skillMap["g"] = printWikiSkill(evo.Skill1(), nil, "g")
		skillMap["g2"] = printWikiSkill(evo.Skill2(), nil, "g2")
		skillMap["g3"] = printWikiSkill(evo.Skill3(), nil, "g3")
		skillMap["gt"] = printWikiSkill(evo.ThorSkill1(), nil, "gt")
	}
	// add rebirth skills as long as the first evo wasn't rebirth
	if evo, ok := evolutions["X"]; ok && firstEvo.ID != evo.ID {
		skillMap["x"] = printWikiSkill(evo.Skill1(), nil, "x")
		skillMap["x2"] = printWikiSkill(evo.Skill2(), nil, "x2")
		skillMap["x3"] = printWikiSkill(evo.Skill3(), nil, "x3")
		skillMap["xt"] = printWikiSkill(evo.ThorSkill1(), nil, "xt")
	}
	// order that we want to print the skills
	skillEvos := []string{"", "2", "3", "a", "g", "g2", "g3", "ga", "x", "x2", "x3", "xa", "t", "gt", "xt"}

	// actually print the skills now...
	for _, skillEvo := range skillEvos {
		if val, ok := skillMap[skillEvo]; ok && val != "" {
			io.WriteString(w, val)
		}
	}

	//traverse evolutions in order
	lenEvoKeys := len(evokeys)
	fmt.Fprint(w, "<!-- Not all evolution stats here will be applicable. Please delete any that are not needed. -->\n")
	for _, k := range evokeys {
		evo := evolutions[k]
		if k == "H" {
			if evo.EvolutionRank >= 0 {
				k = strconv.Itoa(evo.EvolutionRank)
			} else if lenEvoKeys == 1 {
				k = "1"
			}
		}
		maxAtk, maxDef, maxSol := maxStats(evo, len(evolutions))
		fmt.Fprintf(w, "|max level %[1]s = %[2]d\n|cost %[1]s = %[3]d\n|atk %[1]s = %[4]d%s\n|def %[1]s = %[6]d%s\n|soldiers %[1]s = %[8]d%s\n",
			strings.ToLower(k),
			evo.CardRarity().MaxCardLevel,
			evo.DeckCost,
			evo.DefaultOffense, maxAtk,
			evo.DefaultDefense, maxDef,
			evo.DefaultFollower, maxSol)
	}

	for _, k := range evokeys {
		evo := evolutions[k]
		if k == "H" {
			if evo.EvolutionRank >= 0 {
				k = strconv.Itoa(evo.EvolutionRank)
			} else if lenEvoKeys == 1 {
				k = "1"
			}
		}
		if evo.MedalRate > 0 {
			fmt.Fprintf(w, "|medals %[1]s = %[2]d\n", strings.ToLower(k), evo.MedalRate)
		}

		if evo.Price > 0 {
			fmt.Fprintf(w, "|gold %[1]s = %[2]d\n", strings.ToLower(k), evo.Price)
		}
	}

	fmt.Fprintf(w, "|description = %s\n|friendship = %s\n",
		html.EscapeString(card.Description()), html.EscapeString(strings.Replace(card.Friendship(), "\n", "<br />", -1)))
	login := card.Login()
	if len(strings.TrimSpace(login)) > 0 {
		fmt.Fprintf(w, "|login = %s\n", html.EscapeString(strings.Replace(login, "\n", "<br />", -1)))
	}
	fmt.Fprintf(w, "|meet = %s\n|battle start = %s\n|battle end = %s\n|friendship max = %s\n|friendship event = %s\n|rebirth = %s\n",
		html.EscapeString(strings.Replace(card.Meet(), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleStart(), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleEnd(), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipMax(), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipEvent(), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.RebirthEvent(), "\n", "<br />", -1)),
	)

	gevo, ok := evolutions["G"]
	if !ok {
		gevo, ok = evolutions["GA"]
	}
	if gevo != nil {
		var awakenInfo *vc.CardAwaken
		for idx, val := range vc.Data.Awakenings {
			if gevo.ID == val.ResultCardID {
				awakenInfo = &vc.Data.Awakenings[idx]
				break
			}
		}
		if awakenInfo != nil {
			printAwakenMaterials(w, awakenInfo)
		}
	}

	xevo, ok := evolutions["X"]
	if !ok {
		xevo, ok = evolutions["XA"]
	}
	if xevo != nil {
		var rebirthInfo *vc.CardAwaken
		for idx, val := range vc.Data.Rebirths {
			if xevo.ID == val.ResultCardID {
				rebirthInfo = &vc.Data.Rebirths[idx]
				break
			}
		}
		if rebirthInfo != nil {
			printRebirthMaterials(w, rebirthInfo)
		}
	}

	awl := make(vc.ArchwitchList, 0)
	for _, evo := range evolutions {
		aws := evo.ArchwitchesWithLikeabilityQuotes()
		log.Printf("Evo AW records: %v", aws)
		if len(aws) > 0 {
			awl = append(awl, aws...)
			log.Printf("Found %d AW records found for %s", len(awl), card.Name)
		}
	}
	log.Printf("AW records: %v", awl)
	if len(awl) > 0 {
		aw := awl.Earliest()
		log.Printf("AW %d found for card %s", aw.ID, card.Name)
		for _, like := range aw.Likeability() {
			fmt.Fprintf(w, "|likeability %d = %s\n",
				like.Friendship,
				html.EscapeString(strings.Replace(like.Likability, "\n", "<br />", -1)),
			)
		}
	}

	if turnOverFrom != nil {
		fmt.Fprintf(w, "|turnoverfrom = %s\n", turnOverFrom.Name)
	} else if turnOverTo != nil {
		fmt.Fprintf(w, "|turnoverto = %s\n", turnOverTo.Name)
		fmt.Fprintf(w, "|availability = %s\n", avail)
	} else {
		fmt.Fprintf(w, "|availability = %s\n", avail)
	}
	io.WriteString(w, "}}")

	//Write out amalgamations here
	if len(amalgamations) > 0 {
		io.WriteString(w, "\n==''[[Amalgamation]]''==\n")
		for _, v := range amalgamations {
			mats := v.Materials()
			l := len(mats)
			fmt.Fprintf(w, "{{Amalgamation|matcount = %d\n|name 1 = %s|rarity 1 = %s\n|name 2 = %s|rarity 2 = %s\n|name 3 = %s|rarity 3 = %s\n",
				l-1, mats[0].Name, mats[0].Rarity(), mats[1].Name, mats[1].Rarity(), mats[2].Name, mats[2].Rarity())
			if l > 3 {
				fmt.Fprintf(w, "|name 4 = %s|rarity 4 = %s\n", mats[3].Name, mats[3].Rarity())
			}
			if l > 4 {
				fmt.Fprintf(w, "|name 5 = %s|rarity 5 = %s\n", mats[4].Name, mats[4].Rarity())
			}
			io.WriteString(w, "}}\n")
		}
	}
	io.WriteString(w, "</textarea></div>")
	// show images here
	io.WriteString(w, "<div style=\"float:left\">")
	iconEvos := card.EvosWithDistinctImages(true)
	for _, k := range evokeys {
		evo := evolutions[k]
		// check if we should force non-H status
		r := evo.Rarity()[0]
		if lenEvoKeys == 1 && (evo.EvolutionRank == 1 || evo.EvolutionRank < 0) && r != 'H' && r != 'G' {
			k = "0"
		}
		dupWarn := ""
		if !util.Contains(iconEvos, k) {
			dupWarn = "Duplicate Image<br />"
		}
		fmt.Fprintf(w,
			`<div style="float: left; margin: 3px"><a href="/images/cardthumb/%s">%[4]s<img src="/images/cardthumb/%[1]s"/></a><br />%s : %s☆</div>`,
			evo.Image(),
			evo.Name,
			k,
			dupWarn,
		)
	}
	io.WriteString(w, "<div style=\"clear: both\">")
	imageEvos := card.EvosWithDistinctImages(false)
	for _, k := range evokeys {
		evo := evolutions[k]
		// check if we should force non-H status
		r := evo.Rarity()[0]
		if lenEvoKeys == 1 && (evo.EvolutionRank == 1 || evo.EvolutionRank < 0) && r != 'H' && r != 'G' {
			k = "0"
		}
		dupWarn := ""
		if !util.Contains(imageEvos, k) {
			dupWarn = "Duplicate Image<br />"
		}

		pathPart := ""
		if _, err := os.Stat(vc.FilePath + "/card/hd/" + evo.Image()); err == nil {
			pathPart = "cardHD"
		} else if _, err := os.Stat(vc.FilePath + "/card/md/" + evo.Image()); err == nil {
			pathPart = "card"
		} else {
			pathPart = "cardSD"
		}
		fmt.Fprintf(w,
			`<div style="float: left; margin: 3px"><a href="/images/%[4]s/%[1]s.png">%[5]s<img src="/images/%[4]s/%[1]s.png"/></a><br />%s : %s☆</div>`,
			evo.Image(),
			evo.Name,
			k,
			pathPart,
			dupWarn,
		)
	}
	io.WriteString(w, "</div>")
	io.WriteString(w, "</body></html>")
}

// CardCsvHandler outputs the cards as a CSV doc
func CardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-cards-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".csv\"")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	cw.Write([]string{"ID", "Card #", "Name", "Evo Rank", "TranscardID", "Rarity", "Element", "Deck Cost", "Base ATK",
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold", "Skill 1 Name", "Skill Min",
		"Skill Max", "Skill Procs", "Target Scope", "Target Logic", "Skill 2", "Skill 3", "Thor Skill 1", "Skill Special", "Description", "Friendship",
		"Login", "Meet", "Battle Start", "Battle End", "Friendship Max", "Friendship Event", "Is Closed"})
	for _, card := range vc.Data.Cards {
		err := cw.Write([]string{strconv.Itoa(card.ID), fmt.Sprintf("cd_%05d", card.CardNo), card.Name, strconv.Itoa(card.EvolutionRank),
			strconv.Itoa(card.TransCardID), card.Rarity(), card.Element(), strconv.Itoa(card.DeckCost), strconv.Itoa(card.DefaultOffense),
			strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower), strconv.Itoa(card.MaxOffense),
			strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower), card.Skill1Name(),
			card.SkillMin(), card.SkillMax(), card.SkillProcs(), card.SkillTarget(),
			card.SkillTargetLogic(), card.Skill2Name(), card.Skill3Name(), card.ThorSkill1Name(), card.SpecialSkill1Name(),
			card.Description(), card.Friendship(), card.Login(), card.Meet(),
			card.BattleStart(), card.BattleEnd(), card.FriendshipMax(), card.FriendshipEvent(),
			strconv.Itoa(card.IsClosed),
		})
		if err != nil {
			log.Printf(err.Error() + "\n")
		}
	}
	cw.Flush()
}

// CardCsvGLRHandler outputs GLR and Rebirth cards in a format usable for stat calcuations
func CardCsvGLRHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-glr-cards-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".csv\"")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	cw.Write([]string{"VC ID", "Name", "name-rare", "Element", "Rarity", "Base ATK",
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold", "ATK Gain",
		"DEF Gain", "SOL Gain", "ATK Gain/lvl", "DEF Gain/lvl", "SOL Gain/lvl",
		"Sol Perfect Min", "Sol Perfect Max",
		"Is Closed",
		"Max Level",
		"Max Rarity Atk",
		"Max Rarity Def",
		"Max Rarity Sold",
	})

	cards := vc.Data.Cards.Copy()
	log.Printf("Total card count: %d", len(cards))
	// sort by name A-Z
	sort.Slice(cards, func(i, j int) bool {
		first := cards[i]
		second := cards[j]

		return first.Name < second.Name
	})

	for _, card := range cards {
		if card.Rarity() == "GLR" || card.EvoIsReborn() {
			log.Printf("found GLR or Reborn card %d:%s - %s", card.ID, card.Name, card.Rarity())
			cardRare := card.CardRarity()
			maxLevel := float64(cardRare.MaxCardLevel - 1)
			atkGain := card.MaxOffense - card.DefaultOffense
			defGain := card.MaxDefense - card.DefaultDefense
			solGain := card.MaxFollower - card.DefaultFollower
			minStat := card.EvoPerfectLvl1()
			maxStat := card.EvoPerfect()
			err := cw.Write([]string{strconv.Itoa(card.ID), card.Name, card.Name + " - " + card.Rarity(), card.Element(), card.Rarity(),
				strconv.Itoa(card.DefaultOffense), strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower),
				strconv.Itoa(card.MaxOffense), strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower),
				strconv.Itoa(atkGain), strconv.Itoa(defGain), strconv.Itoa(solGain),
				fmt.Sprintf("%.4f", float64(atkGain)/maxLevel), fmt.Sprintf("%.4f", float64(defGain)/maxLevel), fmt.Sprintf("%.4f", float64(solGain)/maxLevel),
				strconv.Itoa(minStat.Soldiers),
				strconv.Itoa(maxStat.Soldiers),
				strconv.Itoa(card.IsClosed),
				strconv.Itoa(cardRare.MaxCardLevel),
				strconv.Itoa(cardRare.LimtOffense),
				strconv.Itoa(cardRare.LimtDefense),
				strconv.Itoa(cardRare.LimtMaxFollower),
			})
			if err != nil {
				log.Printf(err.Error() + "\n")
			}
		}
	}
	cw.Flush()
}

//CardTableHandler outputs the cards in an HTML table
func CardTableHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	filter := func(card *vc.Card) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		if rarity := qs.Get("rarity"); rarity != "" {
			match = match && (card.Rarity() == rarity ||
				card.Rarity() == ("H"+rarity) ||
				card.Rarity() == ("G"+rarity))
		}
		if evos := qs.Get("evos"); evos != "" {
			evon, err := strconv.Atoi(evos)
			if err == nil && (evon == 1 || evon == 4) {
				match = match && (card.LastEvolutionRank == evon)
			}
		}
		if isThor := qs.Get("isThor"); isThor != "" {
			match = match && card.ThorSkillID1 > 0
		}
		if hasRebirth := qs.Get("hasRebirth"); hasRebirth != "" {
			match = match && card.HasRebirth()
		}
		if name := qs.Get("name"); name != "" {
			match = match && strings.Contains(strings.ToLower(card.Name), strings.ToLower(name))
		}
		if skillname := qs.Get("skillname"); skillname != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(); skill1 != nil {
				s1 = skill1.Name != "" && strings.Contains(strings.ToLower(skill1.Name), strings.ToLower(skillname))
				//log.Printf(skill1.Name + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(); skill2 != nil {
				s2 = skill2.Name != "" && strings.Contains(strings.ToLower(skill2.Name), strings.ToLower(skillname))
				//log.Printf(skill2.Name + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		if skilldesc := qs.Get("skilldesc"); skilldesc != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(); skill1 != nil {
				s1 = skill1.Fire != "" && strings.Contains(strings.ToLower(skill1.Fire), strings.ToLower(skilldesc))
				//log.Printf(skill1.Fire + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(); skill2 != nil {
				s2 = skill2.Fire != "" && strings.Contains(strings.ToLower(skill2.Fire), strings.ToLower(skilldesc))
				//log.Printf(skill2.Fire + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		return
	}
	// File header
	fmt.Fprintf(w, `<html><head><title>All Cards</title>
<style>table, th, td {border: 1px solid black;};</style>
</head>
<body>
<form method="GET">
<label for="f_name">Name:</label><input id="f_name" name="name" value="%s" />
<label for="f_rarity">Rarity:</label><select id="f_rarity" name="rarity" value="%s">
<option value=""></option>
<option value="X">X</option>
<option value="N">N</option>
<option value="R">R</option>
<option value="SR">SR</option>
<option value="UR">UR</option>
<option value="LR">LR</option>
</select>
<label for="f_evos">Evo:</label><select id="f_evos" name="evos" value="%s">
<option value=""></option>
<option value="1">1★</option>
<option value="4">4★</option>
</select>
<label for="f_skillname">Skill Name:</label><input id="f_skillname" name="skillname" value="%s" />
<label for="f_skilldesc">Skill Description:</label><input id="f_skilldesc" name="skilldesc" value="%s" />
<label for="f_skillisthor">Has Thor Skill:</label><input id="f_skillisthor" name="isThor" type="checkbox" value="checked" %s />
<label for="f_hasrebirth">Has Rebirth:</label><input id="f_hasrebirth" name="hasRebirth" type="checkbox" value="checked" %s />
<button type="submit">Submit</button>
</form>
`,
		qs.Get("name"),
		qs.Get("rarity"),
		qs.Get("evos"),
		qs.Get("skillname"),
		qs.Get("skilldesc"),
		qs.Get("isThor"),
		qs.Get("hasRebirth"),
	)
	fmt.Fprintf(w, "<div>\n<table>\n")
	printHTMLTableHeader(w,
		"_id",
		"card_no",
		"name",
		"evolution rank",
		"max evolution rank",
		"Next Evo",
		"Rarity",
		"Element",
		"Character ID",
		"deck_cost",
		"default offense",
		"default defense",
		"default follower",
		"max offense",
		"max defense",
		"max follower",
		"Skill 1 Name",
		"Skill Min",
		"Skill Max",
		"Skill Procs",
		"Min Effect",
		"Min Rate",
		"Max Effect",
		"Max Rate",
		"Target Scope",
		"Target Logic",
		"Skill 2",
		"Skill 3",
		"Thor Skill",
		"Skill Special",
		"Description",
		"Friendship",
		"Login",
		"Meet",
		"Battle Start",
		"Battle End",
		"Friendship Max",
		"Friendship Event",
	)
	fmt.Fprintf(w, "\n<tbody>\n")

	for i := len(vc.Data.Cards) - 1; i >= 0; i-- {
		card := vc.Data.Cards[i]
		if !filter(card) {
			continue
		}
		skill1 := card.Skill1()
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2()
		// skillS1 := card.SpecialSkill1()
		printHTMLTableRow(w,
			card.ID,
			fmt.Sprintf("<a href=\"/cards/detail/%d\">%05d</a>", card.ID, card.CardNo),
			fmt.Sprintf("<a href=\"/cards/detail/%d\">%s</a>", card.ID, card.Name),
			card.EvolutionRank,
			card.LastEvolutionRank,
			card.EvolutionCardID,
			card.Rarity(),
			card.Element(),
			card.CardCharaID,
			card.DeckCost,
			card.DefaultOffense,
			card.DefaultDefense,
			card.DefaultFollower,
			card.MaxOffense,
			card.MaxDefense,
			card.MaxFollower,
			card.Skill1Name(),
			card.SkillMin(),
			card.SkillMax(),
			card.SkillProcs(),
			skill1.EffectDefaultValue,
			skill1.DefaultRatio,
			skill1.EffectMaxValue,
			skill1.MaxRatio,
			card.SkillTarget(),
			card.SkillTargetLogic(),
			card.Skill2Name(),
			card.Skill3Name(),
			card.ThorSkill1Name(),
			card.SpecialSkill1Name(),
			card.Description(),
			card.Friendship(),
			card.Login(),
			card.Meet(),
			card.BattleStart(),
			card.BattleEnd(),
			card.FriendshipMax(),
			card.FriendshipEvent(),
		)
	}

	io.WriteString(w, "</tbody></table></div></body></html>")
}

// getPrevious finds the card before the earliest evolution of this card. Does
// not take into account cards that were released for awakening at a later date
// maybe should do this by character ID instead?
func getPrevious(card *vc.Card) (prev *vc.Card, prevName string) {
	if card == nil {
		return nil, ""
	}

	minID := card.GetEvolutionCards().Earliest().ID - 1
	for prev = vc.CardScan(minID); prev == nil && minID > 0; prev = vc.CardScan(minID) {
		minID--
	}
	if prev != nil {
		prev = prev.GetEvolutionCards().Earliest()
		if prev.Name == "" {
			prevName = prev.Image()
		} else {
			prevName = prev.Name
		}
	}
	return
}

// getNext finds the card before the latest evolution of this card. Does
// not take into account cards that were released for awakening at a later date
// maybe should do this by character ID instead?
func getNext(card *vc.Card) (next *vc.Card, nextName string) {
	if card == nil {
		return nil, ""
	}
	maxID := card.GetEvolutionCards().Latest().ID + 1
	lastID := vc.Data.Cards.Latest().ID
	for next = vc.CardScan(maxID); next == nil && maxID < lastID; next = vc.CardScan(maxID) {
		maxID++
	}
	if next != nil {
		next = next.GetEvolutionCards().Earliest()
		if next.Name == "" {
			nextName = next.Image()
		} else {
			nextName = next.Name
		}
	}
	return
}

func maxStats(evo *vc.Card, numOfEvos int) (atk, def, sol string) {
	// stats := evo.EvoStandard()
	// atk += fmt.Sprintf(" / %d", stats.Attack)
	// def += fmt.Sprintf(" / %d", stats.Defense)
	// sol += fmt.Sprintf(" / %d", stats.Soldiers)

	// stats = evo.EvoMixed()
	// atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", stats.Attack)
	// def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", stats.Defense)
	// sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", stats.Soldiers)

	// stats = evo.AmalgamationPerfect()
	// atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Attack)
	// def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Defense)
	// sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Soldiers)
	// return

	if evo.CardRarity().MaxCardLevel == 1 && numOfEvos == 1 {
		// only X cards have a max level of 1 and they don't evo
		// only possible amalgamations like Philosopher's Stones
		if evo.IsAmalgamation() {
			stats := evo.AmalgamationPerfect()
			atk = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Attack)
			def = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Defense)
			sol = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", stats.Soldiers)
		}
		return
	}
	if evo.IsAmalgamation() {
		stats := evo.EvoStandard()
		atk = " / " + strconv.Itoa(stats.Attack)
		def = " / " + strconv.Itoa(stats.Defense)
		sol = " / " + strconv.Itoa(stats.Soldiers)
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			pStats := evo.AmalgamationPerfect()
			if stats.Attack != pStats.Attack || stats.Defense != pStats.Defense || stats.Soldiers != pStats.Soldiers {
				if evo.PossibleMixedEvo() {
					mStats := evo.EvoMixed()
					if stats.Attack != mStats.Attack || stats.Defense != mStats.Defense || stats.Soldiers != mStats.Soldiers {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
					}
				}
				lrStats := evo.AmalgamationLRStaticLvl1()
				if lrStats.Attack != pStats.Attack || lrStats.Defense != pStats.Defense || lrStats.Soldiers != pStats.Soldiers {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Soldiers)
				}
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Soldiers)
			}
		} else {
			pStats := evo.AmalgamationPerfect()
			if stats.Attack != pStats.Attack || stats.Defense != pStats.Defense || stats.Soldiers != pStats.Soldiers {
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Soldiers)
			}
		}
	} else if evo.EvolutionRank < 2 {
		// not an amalgamation.
		stats := evo.EvoStandard()
		atk = " / " + strconv.Itoa(stats.Attack)
		def = " / " + strconv.Itoa(stats.Defense)
		sol = " / " + strconv.Itoa(stats.Soldiers)
		pStats := evo.EvoPerfect()
		if stats.Attack != pStats.Attack || stats.Defense != pStats.Defense || stats.Soldiers != pStats.Soldiers {
			var evoType string
			if evo.PossibleMixedEvo() {
				evoType = "Amalgamation"
				mStats := evo.EvoMixed()
				if stats.Attack != mStats.Attack || stats.Defense != mStats.Defense || stats.Soldiers != mStats.Soldiers {
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
				}
			} else {
				evoType = "Evolution"
			}
			atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Attack, evoType)
			def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Defense, evoType)
			sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Soldiers, evoType)
		}
	} else {
		// not an amalgamation, Evo Rank >=2 (Awoken cards or 4* evos).
		stats := evo.EvoStandard()
		atk = " / " + strconv.Itoa(stats.Attack)
		def = " / " + strconv.Itoa(stats.Defense)
		sol = " / " + strconv.Itoa(stats.Soldiers)
		printedMixed := false
		printedPerfect := false
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			if evo.PossibleMixedEvo() {
				mStats := evo.EvoMixed()
				if stats.Attack != mStats.Attack || stats.Defense != mStats.Defense || stats.Soldiers != mStats.Soldiers {
					printedMixed = true
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
				}
			}
			pStats := evo.AmalgamationPerfect()
			if stats.Attack != pStats.Attack || stats.Defense != pStats.Defense || stats.Soldiers != pStats.Soldiers {
				lrStats := evo.AmalgamationLRStaticLvl1()
				if lrStats.Attack != pStats.Attack || lrStats.Defense != pStats.Defense || lrStats.Soldiers != pStats.Soldiers {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStats.Soldiers)
				}
				printedPerfect = true
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStats.Soldiers)
			}
		}

		if !strings.HasSuffix(evo.Rarity(), "LR") &&
			evo.Rarity()[0] != 'G' && evo.Rarity()[0] != 'X' {
			// TODO need more logic here to check if it's an Amalg vs evo only.
			// may need different options depending on the type of card.
			if evo.EvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				stats := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 card evo stats
					stats = evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, 9)
				}
			}
			//If 4* card, calculate 16 card evo stats
			stats := evo.EvoPerfect()
			var cards int
			switch evo.EvolutionRank {
			case 1:
				cards = 2
			case 2:
				cards = 4
			case 3:
				cards = 8
			case 4:
				cards = 16
			}
			atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, cards)
			def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, cards)
			sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, cards)
		}
		if evo.Rarity()[0] == 'G' || evo.Rarity()[0] == 'X' {
			evo.EvoStandardLvl1() // just to print out the level 1 G stats
			pStats := evo.EvoPerfect()
			if stats.Attack != pStats.Attack || stats.Defense != pStats.Defense || stats.Soldiers != pStats.Soldiers {
				var evoType string
				if !printedMixed && evo.PossibleMixedEvo() {
					evoType = "Amalgamation"
					mStats := evo.EvoMixed()
					if stats.Attack != mStats.Attack || stats.Defense != mStats.Defense || stats.Soldiers != mStats.Soldiers {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
					}
				} else {
					evoType = "Evolution"
				}
				if !printedPerfect {
					atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Attack, evoType)
					def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Defense, evoType)
					sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStats.Soldiers, evoType)
				}
			}
			awakensFrom := evo.AwakensFrom()
			if awakensFrom == nil && evo.RebirthsFrom() != nil {
				awakensFrom = evo.RebirthsFrom().AwakensFrom()
			}
			if awakensFrom != nil && awakensFrom.LastEvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				stats := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 and 16 card evo stats
					stats = evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, 9)
					stats = evo.EvoPerfect()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Attack, 16)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Defense, 16)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", stats.Soldiers, 16)
				}
			}
		}
	}
	return
}

func printAwakenMaterials(w http.ResponseWriter, awakenInfo *vc.CardAwaken) {
	if awakenInfo == nil {
		return
	}

	fmt.Fprintf(w, "|awaken chance = %d\n", awakenInfo.Percent)

	printAwakenMaterial(w, awakenInfo.Item(1), awakenInfo.Material1Count)
	printAwakenMaterial(w, awakenInfo.Item(2), awakenInfo.Material2Count)
	printAwakenMaterial(w, awakenInfo.Item(3), awakenInfo.Material3Count)
	printAwakenMaterial(w, awakenInfo.Item(4), awakenInfo.Material4Count)
	printAwakenMaterial(w, awakenInfo.Item(5), awakenInfo.Material5Count)

}

func printAwakenMaterial(w http.ResponseWriter, item *vc.Item, count int) {
	if item == nil {
		return
	}
	if strings.Contains(item.NameEng, "Crystal") {
		fmt.Fprintf(w, "|awaken crystal = %d\n", count)
	} else if strings.Contains(item.NameEng, "Orb") {
		fmt.Fprintf(w, "|awaken orb = %d\n", count)
	} else if strings.Contains(item.NameEng, "(L)") {
		fmt.Fprintf(w, "|awaken l = %d\n", count)
	} else if strings.Contains(item.NameEng, "(M)") {
		fmt.Fprintf(w, "|awaken m = %d\n", count)
	} else if strings.Contains(item.NameEng, "(S)") {
		fmt.Fprintf(w, "|awaken s = %d\n", count)
	} else {
		fmt.Fprintf(w, "*******Unknown item: %s\n", item.NameEng)
	}
}

func printRebirthMaterials(w http.ResponseWriter, awakenInfo *vc.CardAwaken) {
	if awakenInfo == nil {
		return
	}

	fmt.Fprintf(w, "|rebirth chance = %d\n", awakenInfo.Percent)

	printRebirthMaterial(w, 1, awakenInfo.Item(1), awakenInfo.Material1Count)
	printRebirthMaterial(w, 2, awakenInfo.Item(2), awakenInfo.Material2Count)
	printRebirthMaterial(w, 3, awakenInfo.Item(3), awakenInfo.Material3Count)
}

func printRebirthMaterial(w http.ResponseWriter, matNum int, item *vc.Item, count int) {
	if item == nil || count <= 0 {
		return
	}
	fmt.Fprintf(w, "|rebirth item %d = %s\n", matNum, item.NameEng)
	fmt.Fprintf(w, "|rebirth item %d count = %d\n", matNum, count)
}

func printWikiSkill(s *vc.Skill, ls *vc.Skill, evoMod string) (ret string) {
	ret = ""
	if s == nil {
		return
	}

	if evoMod != "" {
		evoMod += " "
	}

	lv10 := ""
	skillLvl10 := ""
	sName := s.Name
	if ls == nil || ls.ID == s.ID {
		// 1st skills only have lvl10, skill 2, 3, and thor do not.
		if len(s.Levels()) == 10 {
			skillLvl10 = html.EscapeString(strings.Replace(s.SkillMax(), "\n", "<br />", -1))
			lv10 = fmt.Sprintf("\n|skill %slv10 = %s",
				evoMod,
				skillLvl10,
			)
		}
	} else {
		// handles Great DMG skills where the last evo is the max skill
		sName = ls.Name
		skillLvl10 = html.EscapeString(strings.Replace(ls.SkillMax(), "\n", "<br />", -1))
		lv10 = fmt.Sprintf("\n|skill %slv10 = %s",
			evoMod,
			skillLvl10,
		)
	}

	skillLvl1 := ""
	if strings.HasSuffix(evoMod, "t ") {
		skillLvl1 = s.FireMax() + " / 100% chance"
	} else {
		skillLvl1 = s.SkillMin()
	}
	skillLvl1 = html.EscapeString(strings.Replace(skillLvl1, "\n", "<br />", -1))

	if skillLvl1 == skillLvl10 {
		lv10 = ""
	}

	ret = fmt.Sprintf(`|skill %[1]s= %[2]s
|skill %[1]slv1 = %[3]s%[4]s
|procs %[1]s= %[5]d
`,
		evoMod,
		html.EscapeString(sName),
		skillLvl1,
		lv10,
		s.MaxCount,
	)

	if s.EffectID == 36 {
		// Random Skill
		for k, v := range []int{s.EffectParam, s.EffectParam2, s.EffectParam3, s.EffectParam4, s.EffectParam5} {
			rs := vc.SkillScan(v)
			if rs != nil {
				ret += fmt.Sprintf("|random %s%d = %s \n", evoMod, k+1, rs.FireMin())
			}
		}
	}

	// Check if the second skill expires
	if s.Expires() {
		ret += fmt.Sprintf("|skill %send = %v\n", evoMod, s.PublicEndDatetime)
	}
	return
}

func getAmalgamations(evolutions map[string]*vc.Card) []vc.Amalgamation {
	amalgamations := make([]vc.Amalgamation, 0)
	seen := map[vc.Amalgamation]bool{}
	for idx, evo := range evolutions {
		log.Printf("Card: %d, Name: %s, Evo: %s\n", evo.ID, evo.Name, idx)
		as := evo.Amalgamations()
		if len(as) > 0 {
			for _, a := range as {
				if _, ok := seen[a]; !ok {
					amalgamations = append(amalgamations, a)
					seen[a] = true
				}
			}
		}
	}
	sort.Sort(vc.ByMaterialCount(amalgamations))
	return amalgamations
}
