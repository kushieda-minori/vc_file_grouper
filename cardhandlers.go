package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

func cardHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Cards</title></head><body>\n")
	for _, card := range VcData.Cards {
		fmt.Fprintf(w,
			"<div style=\"float: left; margin: 3px\"><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></div>",
			card.Image(),
			card.ID,
			card.Name)
	}
	io.WriteString(w, "</body></html>")
}

func cardLevelHandler(w http.ResponseWriter, r *http.Request) {
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
	genLevels(VcData.CardLevels)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />XSR and XUR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(VcData.CardLevelsX)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />LR-GLR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(VcData.CardLevelsLR)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />XLR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(VcData.CardLevelsXLR)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />LR Resources<br/><textarea rows=\"25\" cols=\"80\">")

	io.WriteString(w, header2)
	l := len(VcData.LevelLRResources)
	for i, lvl := range VcData.LevelLRResources {
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
	l = len(VcData.LevelXResources)
	for i, lvl := range VcData.LevelXResources {
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
	l = len(VcData.LevelXLRResources)
	for i, lvl := range VcData.LevelXLRResources {
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
	cardID, err := strconv.Atoi(pathParts[2])
	if err != nil || cardID < 1 || cardID > len(VcData.Cards) {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	card := vc.CardScan(cardID, VcData.Cards)
	if card == nil {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}
	evolutions := card.GetEvolutions(VcData)
	amalgamations := getAmalgamations(evolutions)

	var firstEvo, lastEvo *vc.Card
	evoOrder := []string{"0", "1", "2", "3", "H", "A", "G", "GA", "X", "XA"}
	var evokeys []string // cache of actual evos for this card
	for _, k := range evoOrder {
		evo, ok := evolutions[k]
		if ok {
			evokeys = append(evokeys, k)
			if firstEvo == nil {
				firstEvo = evo
			}
			if k == "H" {
				lastEvo = evo
			}
		}
	}

	var turnOverTo, turnOverFrom *vc.Card
	if firstEvo.ID > 0 {
		if firstEvo.TransCardID > 0 {
			turnOverTo = firstEvo.EvoAccident(VcData.Cards)
		} else {
			turnOverFrom = firstEvo.EvoAccidentOf(VcData.Cards)
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
			if evo.IsAmalgamation(VcData.Amalgamations) {
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

	fmt.Fprintf(w, "<div style=\"clear:both;float:left\">Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br /></div>", cardName)
	io.WriteString(w, "<div><textarea readonly=\"readonly\" style=\"width:100%;height:450px\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n|rarity = %s\n", card.Element(), fixRarity(card.Rarity()))

	skillMap := make(map[string]string)

	skillEvoMod := ""
	if firstEvo.Rarity()[0] == 'G' {
		skillEvoMod = "g"
	} else if firstEvo.Rarity()[0] == 'X' {
		skillEvoMod = "x"
	} else if evo, ok := evolutions["A"]; ok && firstEvo.ID == evo.ID {
		// if the first evo is the amalgamation evo...
		skillEvoMod = ""
	}

	skillMap[skillEvoMod] = printWikiSkill(firstEvo.Skill1(VcData), nil, skillEvoMod)

	skill2 := firstEvo.Skill2(VcData)
	if skill2 != nil {
		// look for skills that improve due to evolution
		if lastEvo == nil {
			skillMap[skillEvoMod+"2"] = printWikiSkill(skill2, nil, skillEvoMod+"2")
		} else {
			skillMap[skillEvoMod+"2"] = printWikiSkill(skill2, lastEvo.Skill2(VcData), skillEvoMod+"2")
		}
		// print skill 3 if it exists
		skillMap[skillEvoMod+"3"] = printWikiSkill(firstEvo.Skill3(VcData), nil, skillEvoMod+"3")
	} else if lastEvo != nil && lastEvo.ID > 0 {
		skillMap[skillEvoMod+"2"] = printWikiSkill(lastEvo.Skill2(VcData), nil, skillEvoMod+"2")
		// print skill 3 if it exists
		skillMap[skillEvoMod+"3"] = printWikiSkill(lastEvo.Skill3(VcData), nil, skillEvoMod+"3")
		skillMap[skillEvoMod+"t"] = printWikiSkill(lastEvo.ThorSkill1(VcData), nil, skillEvoMod+"t")
	}

	// add amal skills as long as the first evo wasn't the amal
	if evo, ok := evolutions["A"]; ok && firstEvo.ID != evo.ID {
		aSkillName := evo.Skill1Name(VcData)
		if aSkillName != firstEvo.Skill1Name(VcData) {
			skillMap["a"] = printWikiSkill(evo.Skill1(VcData), nil, "a")
		}
		if _, ok := evolutions["GA"]; ok {
			skillMap["ga"] = printWikiSkill(evo.Skill1(VcData), nil, "ga")
		}
		if _, ok := evolutions["XA"]; ok {
			skillMap["xa"] = printWikiSkill(evo.Skill1(VcData), nil, "xa")
		}
	}
	// add awoken skills as long as the first evo wasn't awoken
	if evo, ok := evolutions["G"]; ok && firstEvo.ID != evo.ID {
		skillMap["g"] = printWikiSkill(evo.Skill1(VcData), nil, "g")
		skillMap["g2"] = printWikiSkill(evo.Skill2(VcData), nil, "g2")
		skillMap["g3"] = printWikiSkill(evo.Skill3(VcData), nil, "g3")
		skillMap["gt"] = printWikiSkill(evo.ThorSkill1(VcData), nil, "gt")
	}
	// add rebirth skills as long as the first evo wasn't rebirth
	if evo, ok := evolutions["X"]; ok && firstEvo.ID != evo.ID {
		skillMap["x"] = printWikiSkill(evo.Skill1(VcData), nil, "x")
		skillMap["x2"] = printWikiSkill(evo.Skill2(VcData), nil, "x2")
		skillMap["x3"] = printWikiSkill(evo.Skill3(VcData), nil, "x3")
		skillMap["xt"] = printWikiSkill(evo.ThorSkill1(VcData), nil, "xt")
	}
	// order that we want to print the skills
	skillEvos := []string{"", "2", "3", "a", "t", "g", "g2", "g3", "ga", "gt", "x", "x2", "x3", "xa", "xt"}

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
			evo.CardRarity(VcData).MaxCardLevel,
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
		html.EscapeString(card.Description(VcData)), html.EscapeString(strings.Replace(card.Friendship(VcData), "\n", "<br />", -1)))
	login := card.Login(VcData)
	if len(strings.TrimSpace(login)) > 0 {
		fmt.Fprintf(w, "|login = %s\n", html.EscapeString(strings.Replace(login, "\n", "<br />", -1)))
	}
	fmt.Fprintf(w, "|meet = %s\n|battle start = %s\n|battle end = %s\n|friendship max = %s\n|friendship event = %s\n|rebirth = %s\n",
		html.EscapeString(strings.Replace(card.Meet(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleStart(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleEnd(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipMax(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipEvent(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.RebirthEvent(VcData), "\n", "<br />", -1)),
	)

	gevo, ok := evolutions["G"]
	if !ok {
		gevo, ok = evolutions["GA"]
	}
	if gevo != nil {
		var awakenInfo *vc.CardAwaken
		for idx, val := range VcData.Awakenings {
			if gevo.ID == val.ResultCardID {
				awakenInfo = &VcData.Awakenings[idx]
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
		for idx, val := range VcData.Rebirths {
			if xevo.ID == val.ResultCardID {
				rebirthInfo = &VcData.Rebirths[idx]
				break
			}
		}
		if rebirthInfo != nil {
			printRebirthMaterials(w, rebirthInfo)
		}
	}

	var aw *vc.Archwitch
	for _, evo := range evolutions {
		if nil != evo.Archwitch(VcData) {
			aw = evo.Archwitch(VcData)
			break
		}
	}
	if aw != nil {
		for _, like := range aw.Likeability(VcData) {
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
			mats := v.Materials(VcData)
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
	for _, k := range evokeys {
		evo := evolutions[k]
		// check if we should force non-H status
		r := evo.Rarity()[0]
		if lenEvoKeys == 1 && (evo.EvolutionRank == 1 || evo.EvolutionRank < 0) && r != 'H' && r != 'G' {
			k = "0"
		}
		fmt.Fprintf(w,
			`<div style="float: left; margin: 3px"><a href="/images/cardthumb/%s"><img src="/images/cardthumb/%[1]s"/></a><br />%s : %s☆</div>`,
			evo.Image(),
			evo.Name,
			k,
		)
	}
	io.WriteString(w, "<div style=\"clear: both\">")
	for _, k := range evokeys {
		evo := evolutions[k]
		// check if we should force non-H status
		r := evo.Rarity()[0]
		if lenEvoKeys == 1 && (evo.EvolutionRank == 1 || evo.EvolutionRank < 0) && r != 'H' && r != 'G' {
			k = "0"
		}

		if _, err := os.Stat(vcfilepath + "/card/hd/" + evo.Image()); err == nil {
			fmt.Fprintf(w,
				`<div style="float: left; margin: 3px"><a href="/images/cardHD/%s.png"><img src="/images/cardHD/%[1]s.png"/></a><br />%s : %s☆</div>`,
				evo.Image(),
				evo.Name, k)
		} else if _, err := os.Stat(vcfilepath + "/card/md/" + evo.Image()); err == nil {
			fmt.Fprintf(w,
				`<div style="float: left; margin: 3px"><a href="/images/card/%s.png"><img src="/images/card/%[1]s.png"/></a><br />%s : %s☆</div>`,
				evo.Image(),
				evo.Name, k)
		} else {
			fmt.Fprintf(w,
				`<div style="float: left; margin: 3px"><a href="/images/cardSD/%s.png"><img src="/images/cardSD/%[1]s.png"/></a><br />%s : %s☆</div>`,
				evo.Image(),
				evo.Name, k)
		}
	}
	io.WriteString(w, "</div>")
	io.WriteString(w, "</body></html>")
}

func fixRarity(s string) string {
	l := len(s)
	switch l {
	case 1:
		// N, X, R
		return s
	case 2:
		// HN, HX, HR, SR, UR, LR
		return strings.TrimPrefix(s, "H")
	case 3:
		// HSR, GSR, HUR, GUR, HLR, GLR, XSR, XUR, XLR
		return strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(s, "H"), "G"), "X")
	default:
		// not a known rarity!
		return s
	}
}

func cardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv\"")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	cw.Write([]string{"ID", "Card #", "Name", "Evo Rank", "TranscardID", "Rarity", "Element", "Deck Cost", "Base ATK",
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold", "Skill 1 Name", "Skill Min",
		"Skill Max", "Skill Procs", "Target Scope", "Target Logic", "Skill 2", "Skill 3", "Thor Skill 1", "Skill Special", "Description", "Friendship",
		"Login", "Meet", "Battle Start", "Battle End", "Friendship Max", "Friendship Event", "Is Closed"})
	for _, card := range VcData.Cards {
		err := cw.Write([]string{strconv.Itoa(card.ID), fmt.Sprintf("cd_%05d", card.CardNo), card.Name, strconv.Itoa(card.EvolutionRank),
			strconv.Itoa(card.TransCardID), card.Rarity(), card.Element(), strconv.Itoa(card.DeckCost), strconv.Itoa(card.DefaultOffense),
			strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower), strconv.Itoa(card.MaxOffense),
			strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower), card.Skill1Name(VcData),
			card.SkillMin(VcData), card.SkillMax(VcData), card.SkillProcs(VcData), card.SkillTarget(VcData),
			card.SkillTargetLogic(VcData), card.Skill2Name(VcData), card.Skill3Name(VcData), card.ThorSkill1Name(VcData), card.SpecialSkill1Name(VcData),
			card.Description(VcData), card.Friendship(VcData), card.Login(VcData), card.Meet(VcData),
			card.BattleStart(VcData), card.BattleEnd(VcData), card.FriendshipMax(VcData), card.FriendshipEvent(VcData),
			strconv.Itoa(card.IsClosed),
		})
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}
	cw.Flush()
}

func cardCsvGLRHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-glr-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv\"")
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

	cards := make([]vc.Card, len(VcData.Cards))
	copy(cards, VcData.Cards)

	// sort by name A-Z
	sort.Slice(cards, func(i, j int) bool {
		first := cards[i]
		second := cards[j]

		return first.Name < second.Name
	})

	for _, card := range cards {
		if card.Rarity() == "GLR" || (len(card.Rarity()) == 3 && card.Rarity()[0] == 'X') {
			cardRare := card.CardRarity(VcData)
			maxLevel := float64(cardRare.MaxCardLevel - 1)
			atkGain := card.MaxOffense - card.DefaultOffense
			defGain := card.MaxDefense - card.DefaultDefense
			solGain := card.MaxFollower - card.DefaultFollower
			_, _, minSol := (card.EvoPerfectLvl1(VcData))
			_, _, maxSol := (card.EvoPerfect(VcData))
			err := cw.Write([]string{strconv.Itoa(card.ID), card.Name, card.Name + " - " + card.Rarity(), card.Element(), card.Rarity(),
				strconv.Itoa(card.DefaultOffense), strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower),
				strconv.Itoa(card.MaxOffense), strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower),
				strconv.Itoa(atkGain), strconv.Itoa(defGain), strconv.Itoa(solGain),
				fmt.Sprintf("%.4f", float64(atkGain)/maxLevel), fmt.Sprintf("%.4f", float64(defGain)/maxLevel), fmt.Sprintf("%.4f", float64(solGain)/maxLevel),
				strconv.Itoa(minSol), strconv.Itoa(maxSol),
				strconv.Itoa(card.IsClosed),
				strconv.Itoa(cardRare.MaxCardLevel),
				strconv.Itoa(cardRare.LimtOffense),
				strconv.Itoa(cardRare.LimtDefense),
				strconv.Itoa(cardRare.LimtMaxFollower),
			})
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
			}
		}
	}
	cw.Flush()
}

func cardTableHandler(w http.ResponseWriter, r *http.Request) {
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
			match = match && card.HasRebirth(VcData)
		}
		if name := qs.Get("name"); name != "" {
			match = match && strings.Contains(strings.ToLower(card.Name), strings.ToLower(name))
		}
		if skillname := qs.Get("skillname"); skillname != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(VcData); skill1 != nil {
				s1 = skill1.Name != "" && strings.Contains(strings.ToLower(skill1.Name), strings.ToLower(skillname))
				//os.Stdout.WriteString(skill1.Name + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(VcData); skill2 != nil {
				s2 = skill2.Name != "" && strings.Contains(strings.ToLower(skill2.Name), strings.ToLower(skillname))
				//os.Stdout.WriteString(skill2.Name + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		if skilldesc := qs.Get("skilldesc"); skilldesc != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(VcData); skill1 != nil {
				s1 = skill1.Fire != "" && strings.Contains(strings.ToLower(skill1.Fire), strings.ToLower(skilldesc))
				//os.Stdout.WriteString(skill1.Fire + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(VcData); skill2 != nil {
				s2 = skill2.Fire != "" && strings.Contains(strings.ToLower(skill2.Fire), strings.ToLower(skilldesc))
				//os.Stdout.WriteString(skill2.Fire + " " + strconv.FormatBool(s2) + "\n")
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
<div>
<table><thead><tr>
<th>_id</th>
<th>card_no</th>
<th>name</th>
<th>evolution_rank</th>
<th>Next Evo</th>
<th>Rarity</th>
<th>Element</th>
<th>Character ID</th>
<th>deck_cost</th>
<th>default_offense</th>
<th>default_defense</th>
<th>default_follower</th>
<th>max_offense</th>
<th>max_defense</th>
<th>max_follower</th>
<th>Skill 1 Name</th>
<th>Skill Min</th>
<th>Skill Max</th>
<th>Skill Procs</th>
<th>Min Effect</th>
<th>Min Rate</th>
<th>Max Effect</th>
<th>Max Rate</th>
<th>Target Scope</th>
<th>Target Logic</th>
<th>Skill 2</th>
<th>Skill 3</th>
<th>Thor Skill</th>
<th>Skill Special</th>
<th>Description</th>
<th>Friendship</th>
<th>Login</th>
<th>Meet</th>
<th>Battle Start</th>
<th>Battle End</th>
<th>Friendship Max</th>
<th>Friendship Event</th>
</tr></thead>
<tbody>
`,
		qs.Get("name"),
		qs.Get("rarity"),
		qs.Get("evos"),
		qs.Get("skillname"),
		qs.Get("skilldesc"),
		qs.Get("isThor"),
		qs.Get("hasRebirth"),
	)
	for i := len(VcData.Cards) - 1; i >= 0; i-- {
		card := VcData.Cards[i]
		if !filter(&card) {
			continue
		}
		skill1 := card.Skill1(VcData)
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2(VcData)
		// skillS1 := card.SpecialSkill1(VcData)
		fmt.Fprintf(w, "<tr><td>%d</td>"+
			"<td><a href=\"/cards/detail/%[1]d\">%05[2]d</a></td>"+
			"<td><a href=\"/cards/detail/%[1]d\">%[3]s</a></td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td></tr>\n",
			card.ID,
			card.CardNo,
			card.Name,
			card.EvolutionRank,
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
			card.Skill1Name(VcData),
			card.SkillMin(VcData),
			card.SkillMax(VcData),
			card.SkillProcs(VcData),
			skill1.EffectDefaultValue,
			skill1.DefaultRatio,
			skill1.EffectMaxValue,
			skill1.MaxRatio,
			card.SkillTarget(VcData),
			card.SkillTargetLogic(VcData),
			card.Skill2Name(VcData),
			card.Skill3Name(VcData),
			card.ThorSkill1Name(VcData),
			card.SpecialSkill1Name(VcData),
			card.Description(VcData),
			card.Friendship(VcData),
			card.Login(VcData),
			card.Meet(VcData),
			card.BattleStart(VcData),
			card.BattleEnd(VcData),
			card.FriendshipMax(VcData),
			card.FriendshipEvent(VcData),
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

	minID := card.GetEvolutionCards(VcData).Earliest().ID - 1
	for prev = vc.CardScan(minID, VcData.Cards); prev == nil && minID > 0; prev = vc.CardScan(minID, VcData.Cards) {
		minID--
	}
	if prev != nil {
		prev = prev.GetEvolutionCards(VcData).Earliest()
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
	maxID := card.GetEvolutionCards(VcData).Latest().ID + 1
	lastID := vc.CardList(VcData.Cards).Latest().ID
	for next = vc.CardScan(maxID, VcData.Cards); next == nil && maxID < lastID; next = vc.CardScan(maxID, VcData.Cards) {
		maxID++
	}
	if next != nil {
		next = next.GetEvolutionCards(VcData).Earliest()
		if next.Name == "" {
			nextName = next.Image()
		} else {
			nextName = next.Name
		}
	}
	return
}

func maxStats(evo *vc.Card, numOfEvos int) (atk, def, sol string) {
	if evo.CardRarity(VcData).MaxCardLevel == 1 && numOfEvos == 1 {
		// only X cards have a max level of 1 and they don't evo
		// only possible amalgamations like Philosopher's Stones
		if evo.IsAmalgamation(VcData.Amalgamations) {
			atkStat, defStat, solStat := evo.AmalgamationPerfect(VcData)
			atk = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkStat)
			def = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defStat)
			sol = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solStat)
		}
		return
	}
	if evo.IsAmalgamation(VcData.Amalgamations) {
		atkStat, defStat, solStat := evo.EvoStandard(VcData)
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect(VcData)
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				if evo.PossibleMixedEvo(VcData) {
					atkMStat, defMStat, solMStat := evo.EvoMixed(VcData)
					if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
					}
				}
				atkLRStat, defLRStat, solLRStat := evo.AmalgamationLRStaticLvl1(VcData)
				if atkLRStat != atkPStat || defLRStat != defPStat || solLRStat != solPStat {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", atkLRStat)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", defLRStat)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", solLRStat)
				}
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkPStat)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defPStat)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solPStat)
			}
		} else {
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect(VcData)
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkPStat)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defPStat)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solPStat)
			}
		}
	} else if evo.EvolutionRank < 2 {
		// not an amalgamation.
		atkStat, defStat, solStat := evo.EvoStandard(VcData)
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		atkPStat, defPStat, solPStat := evo.EvoPerfect(VcData)
		if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
			var evoType string
			if evo.PossibleMixedEvo(VcData) {
				evoType = "Amalgamation"
				atkMStat, defMStat, solMStat := evo.EvoMixed(VcData)
				if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
				}
			} else {
				evoType = "Evolution"
			}
			atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", atkPStat, evoType)
			def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", defPStat, evoType)
			sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", solPStat, evoType)
		}
	} else {
		// not an amalgamation, Evo Rank >=2 (Awoken cards or 4* evos).
		atkStat, defStat, solStat := evo.EvoStandard(VcData)
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		printedMixed := false
		printedPerfect := false
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			if evo.PossibleMixedEvo(VcData) {
				atkMStat, defMStat, solMStat := evo.EvoMixed(VcData)
				if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
					printedMixed = true
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
				}
			}
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect(VcData)
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				atkLRStat, defLRStat, solLRStat := evo.AmalgamationLRStaticLvl1(VcData)
				if atkLRStat != atkPStat || defLRStat != defPStat || solLRStat != solPStat {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", atkLRStat)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", defLRStat)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", solLRStat)
				}
				printedPerfect = true
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkPStat)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defPStat)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solPStat)
			}
		}

		if !strings.HasSuffix(evo.Rarity(), "LR") &&
			evo.Rarity()[0] != 'G' && evo.Rarity()[0] != 'X' {
			// TODO need more logic here to check if it's an Amalg vs evo only.
			// may need different options depending on the type of card.
			if evo.EvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				atkStat, defStat, solStat := evo.Evo6Card(VcData)
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 card evo stats
					atkStat, defStat, solStat = evo.Evo9Card(VcData)
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 9)
				}
			}
			//If 4* card, calculate 16 card evo stats
			atkStat, defStat, solStat := evo.EvoPerfect(VcData)
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
			atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, cards)
			def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, cards)
			sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, cards)
		}
		if evo.Rarity()[0] == 'G' || evo.Rarity()[0] == 'X' {
			evo.EvoStandardLvl1(VcData) // just to print out the level 1 G stats
			atkPStat, defPStat, solPStat := evo.EvoPerfect(VcData)
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				var evoType string
				if !printedMixed && evo.PossibleMixedEvo(VcData) {
					evoType = "Amalgamation"
					atkMStat, defMStat, solMStat := evo.EvoMixed(VcData)
					if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
					}
				} else {
					evoType = "Evolution"
				}
				if !printedPerfect {
					atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", atkPStat, evoType)
					def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", defPStat, evoType)
					sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", solPStat, evoType)
				}
			}
			awakensFrom := evo.AwakensFrom(VcData)
			if awakensFrom == nil && evo.RebirthsFrom(VcData) != nil {
				awakensFrom = evo.RebirthsFrom(VcData).AwakensFrom(VcData)
			}
			if awakensFrom != nil && awakensFrom.LastEvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				atkStat, defStat, solStat := evo.Evo6Card(VcData)
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 and 16 card evo stats
					atkStat, defStat, solStat = evo.Evo9Card(VcData)
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 9)
					atkStat, defStat, solStat = evo.EvoPerfect(VcData)
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 16)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 16)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 16)
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

	printAwakenMaterial(w, awakenInfo.Item(1, VcData), awakenInfo.Material1Count)
	printAwakenMaterial(w, awakenInfo.Item(2, VcData), awakenInfo.Material2Count)
	printAwakenMaterial(w, awakenInfo.Item(3, VcData), awakenInfo.Material3Count)
	printAwakenMaterial(w, awakenInfo.Item(4, VcData), awakenInfo.Material4Count)
	printAwakenMaterial(w, awakenInfo.Item(5, VcData), awakenInfo.Material5Count)

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

	printRebirthMaterial(w, 1, awakenInfo.Item(1, VcData), awakenInfo.Material1Count)
	printRebirthMaterial(w, 2, awakenInfo.Item(2, VcData), awakenInfo.Material2Count)
	printRebirthMaterial(w, 3, awakenInfo.Item(3, VcData), awakenInfo.Material3Count)
}

func printRebirthMaterial(w http.ResponseWriter, matNum int, item *vc.Item, count int) {
	if item == nil || count <= 0 {
		return
	}
	if strings.Contains(item.NameEng, "Bud") {
		fmt.Fprintf(w, "|rebirth bud = %d\n", count)
	} else if strings.Contains(item.NameEng, "Bloom") {
		fmt.Fprintf(w, "|rebirth bloom = %d\n", count)
	} else if strings.Contains(item.NameEng, "Flora") {
		fmt.Fprintf(w, "|rebirth flora = %d\n", count)
	} else if strings.Contains(item.NameEng, "Secret Elixir") {
		fmt.Fprintf(w, "|rebirth elixir = %d\n", count)
	} else if strings.Contains(item.NameEng, "Medicinal Herb") {
		fmt.Fprintf(w, "|rebirth herb = %d\n", count)
	} else if strings.Contains(item.NameEng, "Zera") {
		fmt.Fprintf(w, "|rebirth zera = %d\n", count)
	} else {
		fmt.Fprintf(w, "*******Unknown Rebirth item: %s\n", item.NameEng)
	}
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

	if ls == nil || ls.ID == s.ID {
		// 1st skills only have lvl10, skill 2, 3, and thor do not.
		if len(s.Levels(VcData)) == 10 {
			skillLvl10 = html.EscapeString(strings.Replace(s.SkillMax(), "\n", "<br />", -1))
			lv10 = fmt.Sprintf("\n|skill %slv10 = %s",
				evoMod,
				skillLvl10,
			)
		}
	} else {
		// handles Great DMG skills where the last evo is the max skill
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
		html.EscapeString(s.Name),
		skillLvl1,
		lv10,
		s.MaxCount,
	)

	if s.EffectID == 36 {
		// Random Skill
		for k, v := range []int{s.EffectParam, s.EffectParam2, s.EffectParam3, s.EffectParam4, s.EffectParam5} {
			rs := vc.SkillScan(v, VcData.Skills)
			if rs != nil {
				ret += fmt.Sprintf("|random %s%d = %s \n", evoMod, k+1, rs.FireMin())
			}
		}
	}

	// Check if the second skill expires
	if (s.PublicEndDatetime.After(time.Time{})) {
		ret += fmt.Sprintf("|skill %send = %v\n", evoMod, s.PublicEndDatetime)
	}
	return
}

func getAmalgamations(evolutions map[string]*vc.Card) []vc.Amalgamation {
	amalgamations := make([]vc.Amalgamation, 0)
	seen := map[vc.Amalgamation]bool{}
	for idx, evo := range evolutions {
		os.Stdout.WriteString(fmt.Sprintf("Card: %d, Name: %s, Evo: %s\n", evo.ID, evo.Name, idx))
		as := evo.Amalgamations(VcData)
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
