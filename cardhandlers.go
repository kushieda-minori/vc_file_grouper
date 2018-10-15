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
			card.Id,
			card.Name)
	}
	io.WriteString(w, "</body></html>")
}

func cardLevelHandler(w http.ResponseWriter, r *http.Request) {
	header := `{| class="article-table" style="float:left"
!Lvl!!To Next Lvl!!Total Needed`
	header2 := `{| class="article-table" style="float:left"
!Lvl!!Gems Needed`

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
				lvl.Id,
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
	io.WriteString(w, "\n<br />LR-GLR<br/><textarea rows=\"25\" cols=\"80\">")
	genLevels(VcData.CardLevelsLR)
	io.WriteString(w, "</textarea>")
	io.WriteString(w, "\n<br />LR Resources<br/><textarea rows=\"25\" cols=\"80\">")

	io.WriteString(w, header2)
	l := len(VcData.LRResources)
	for i, lvl := range VcData.LRResources {
		fmt.Fprintf(w, `
|-
|%d||%d{{Icon|gem}}`,
			lvl.Id,
			lvl.Elixir,
		)
		if (i+1)%25 == 0 && i+1 < l {
			io.WriteString(w, "\n|}\n\n")
			io.WriteString(w, header2)
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
	cardId, err := strconv.Atoi(pathParts[2])
	if err != nil || cardId < 1 || cardId > len(VcData.Cards) {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	card := vc.CardScan(cardId, VcData.Cards)
	if card == nil {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}
	evolutions := card.GetEvolutions(VcData)
	amalgamations := getAmalgamations(evolutions)

	var firstEvo, lastEvo *vc.Card
	evoOrder := []string{"0", "1", "2", "3", "H", "A", "G", "GA"}
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
	if firstEvo.Id > 0 {
		if firstEvo.TransCardId > 0 {
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

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", cardName)
	fmt.Fprintf(w, "<div>Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br />", cardName)
	io.WriteString(w, "<textarea readonly=\"readonly\" style=\"width:100%;height:450px\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n|rarity = %s\n", card.Element(), fixRarity(card.Rarity()))

	skillMap := make(map[string]string)

	skillEvoMod := ""
	if firstEvo.Rarity()[0] == 'G' {
		skillEvoMod = "g"
	} else if evo, ok := evolutions["A"]; ok && firstEvo.Id == evo.Id {
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
	} else if lastEvo != nil && lastEvo.Id > 0 {
		skillMap[skillEvoMod+"2"] = printWikiSkill(lastEvo.Skill2(VcData), nil, skillEvoMod+"2")
		// print skill 3 if it exists
		skillMap[skillEvoMod+"3"] = printWikiSkill(lastEvo.Skill3(VcData), nil, skillEvoMod+"3")
		skillMap[skillEvoMod+"t"] = printWikiSkill(lastEvo.ThorSkill1(VcData), nil, skillEvoMod+"t")
	}

	// add amal skills as long as the first evo wasn't the amal
	if evo, ok := evolutions["A"]; ok && firstEvo.Id != evo.Id {
		aSkillName := evo.Skill1Name(VcData)
		if aSkillName != firstEvo.Skill1Name(VcData) {
			skillMap["a"] = printWikiSkill(evo.Skill1(VcData), nil, "a")
		}
		if _, ok := evolutions["GA"]; ok {
			skillMap["ga"] = printWikiSkill(evo.Skill1(VcData), nil, "ga")
		}
	}
	// add awoken skills as long as the first evo wasn't awoken
	if evo, ok := evolutions["G"]; ok && firstEvo.Id != evo.Id {
		skillMap["g"] = printWikiSkill(evo.Skill1(VcData), nil, "g")
		skillMap["g2"] = printWikiSkill(evo.Skill2(VcData), nil, "g2")
		skillMap["g3"] = printWikiSkill(evo.Skill3(VcData), nil, "g3")
		skillMap["gt"] = printWikiSkill(evo.ThorSkill1(VcData), nil, "gt")
	}
	// order that we want to print the skills
	skillEvos := []string{"", "2", "3", "a", "t", "g", "g2", "g3", "ga", "gt"}

	// actually print the skills now...
	for _, skillEvo := range skillEvos {
		if val, ok := skillMap[skillEvo]; ok && val != "" {
			io.WriteString(w, val)
		}
	}

	//traverse evolutions in order
	lenEvoKeys := len(evokeys)
	for _, k := range evokeys {
		evo := evolutions[k]
		if k == "H" {
			if evo.EvolutionRank >= 0 {
				k = strconv.Itoa(evo.EvolutionRank)
			} else if lenEvoKeys == 1 {
				k = "1"
			}
		}
		fmt.Fprintf(w, "|cost %[1]s = %[2]d\n|atk %[1]s = %[3]d / %s\n|def %[1]s = %[5]d / %s\n|soldiers %[1]s = %[7]d / %s\n",
			strings.ToLower(k),
			evo.DeckCost,
			evo.DefaultOffense, maxStatAtk(evo, len(evolutions)),
			evo.DefaultDefense, maxStatDef(evo, len(evolutions)),
			evo.DefaultFollower, maxStatFollower(evo, len(evolutions)))
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
	fmt.Fprintf(w, "|meet = %s\n|battle start = %s\n|battle end = %s\n|friendship max = %s\n|friendship event = %s\n", html.EscapeString(strings.Replace(card.Meet(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleStart(VcData), "\n", "<br />", -1)), html.EscapeString(strings.Replace(card.BattleEnd(VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipMax(VcData), "\n", "<br />", -1)), html.EscapeString(strings.Replace(card.FriendshipEvent(VcData), "\n", "<br />", -1)))

	gevo, ok := evolutions["G"]
	if !ok {
		gevo, ok = evolutions["GA"]
	}
	if gevo != nil {
		var awakenInfo *vc.CardAwaken
		for _, val := range VcData.Awakenings {
			if gevo.Id == val.ResultCardId {
				awakenInfo = &val
				break
			}
		}
		if awakenInfo != nil {
			printAwakenMaterials(w, awakenInfo)
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
	return strings.TrimPrefix(strings.TrimPrefix(s, "G"), "H")
}

func cardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{"Id", "Card #", "Name", "Evo Rank", "TransCardId", "Rarity", "Element", "Deck Cost", "Base ATK",
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold", "Skill 1 Name", "Skill Min",
		"Skill Max", "Skill Procs", "Target Scope", "Target Logic", "Skill 2", "Skill 3", "Thor Skill 1", "Skill Special", "Description", "Friendship",
		"Login", "Meet", "Battle Start", "Battle End", "Friendship Max", "Friendship Event", "Is Closed"})
	for _, card := range VcData.Cards {
		err := cw.Write([]string{strconv.Itoa(card.Id), fmt.Sprintf("cd_%05d", card.CardNo), card.Name, strconv.Itoa(card.EvolutionRank),
			strconv.Itoa(card.TransCardId), card.Rarity(), card.Element(), strconv.Itoa(card.DeckCost), strconv.Itoa(card.DefaultOffense),
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

func cardTableHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	filter := func(card *vc.Card) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		if isThor := qs.Get("isThor"); isThor != "" {
			match = match && card.ThorSkillId1 > 0
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
<label for="f_skillname">Skill Name:</label><input id="f_skillname" name="skillname" value="%s" />
<label for="f_skilldesc">Skill Description:</label><input id="f_skilldesc" name="skilldesc" value="%s" />
<label for="f_skillisthor">Has Thor Skill:</label><input id="f_skillisthor" name="isThor" type="checkbox" value="checked" %s />
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
		qs.Get("skillname"),
		qs.Get("skilldesc"),
		qs.Get("isThor"),
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
			card.Id,
			card.CardNo,
			card.Name,
			card.EvolutionRank,
			card.EvolutionCardId,
			card.Rarity(),
			card.Element(),
			card.CardCharaId,
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

func maxStatAtk(evo *vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 || (numOfEvos == 1 && !evo.IsAmalgamation(VcData.Amalgamations)) {
		return strconv.Itoa(evo.MaxOffense)
	}
	return "?"
	// ret = strconv.Itoa(evo.BestEvoMaxAttack(VcData))
	// if evo.EvolutionRank > 2 {
	// 	ret = "{{ToolTip|" + ret + "|Perfect Evolution}}"
	// }
	// if evo.IsAmalgamation(evo.Amalgamations(VcData)) {
	// 	ret += " / {{ToolTip|?|Perfect Amalgamation}}"
	// }
	// return
}
func maxStatDef(evo *vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 || (numOfEvos == 1 && !evo.IsAmalgamation(VcData.Amalgamations)) {
		return strconv.Itoa(evo.MaxDefense)
	}
	return "?"
}
func maxStatFollower(evo *vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 || (numOfEvos == 1 && !evo.IsAmalgamation(VcData.Amalgamations)) {
		return strconv.Itoa(evo.MaxFollower)
	}
	return "?"
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

	if ls == nil || ls.Id == s.Id {
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

	if s.EffectId == 36 {
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
		os.Stdout.WriteString(fmt.Sprintf("Card: %d, Name: %s, Evo: %s\n", evo.Id, evo.Name, idx))
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
