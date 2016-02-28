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
	evolutions := getEvolutions(card)
	amalgamations := getAmalgamations(evolutions)

	lastEvo, ok := evolutions["H"]
	if ok {
		delete(evolutions, "H")
	}
	firstEvo, ok := evolutions["F"]
	if ok {
		delete(evolutions, "F")
	} else {
		firstEvo, ok = evolutions["0"]
		if !ok {
			firstEvo = evolutions["1"]
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

	for _, av := range amalgamations {
		// check to see if this card is the result of an amalgamation
		resultAmalgFound := false
		if !resultAmalgFound {
			for _, ev := range evolutions {
				if ev.Id == av.FusionCardId {
					avail += " [[Amalgamation]]"
					resultAmalgFound = true
					break
				}
			}
		}
		// look for a self amalgamation
		if _, ok := evolutions["A"]; !ok {
			var materialFuseId int
			for _, ev := range evolutions {
				switch ev.Id {
				case av.Material1, av.Material2, av.Material3, av.Material4:
					materialFuseId = av.FusionCardId
				}
			}
			if materialFuseId > 0 {
				fuseCard := vc.CardScan(materialFuseId, VcData.Cards)
				if fuseCard.Name == card.Name {
					evolutions["A"] = *fuseCard
					avail += " [[Amalgamation]]"
				}
			}
		}
		if _, ok := evolutions["A"]; ok && resultAmalgFound {
			break
		}
	}

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", card.Name)
	fmt.Fprintf(w, "<div>Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br />", card.Name)
	io.WriteString(w, "<textarea readonly=\"readonly\" style=\"width:100%;height:450px\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n", card.Element())
	if firstEvo.Id > 0 {
		fmt.Fprintf(w, "|rarity = %s\n|skill = %s\n|skill lv1 = %s\n|skill lv10 = %s\n|procs = %s\n%s",
			firstEvo.Rarity(),
			html.EscapeString(firstEvo.Skill1Name(VcData)),
			html.EscapeString(strings.Replace(firstEvo.SkillMin(VcData), "\n", "<br />", -1)),
			html.EscapeString(strings.Replace(firstEvo.SkillMax(VcData), "\n", "<br />", -1)),
			firstEvo.SkillProcs(VcData),
			randomSkillEffects(firstEvo.Skill1(VcData), ""),
		)

		skill2 := firstEvo.Skill2(VcData)
		if skill2 != nil {
			lastEvoSkill2Max := ""
			if lastEvo.Id > 0 {
				lastEvoSkill2 := lastEvo.Skill2(VcData)
				if lastEvoSkill2.Id != skill2.Id {
					lastEvoSkill2Max = "\n|skill 2 lv10 = " +
						html.EscapeString(strings.Replace(lastEvoSkill2.SkillMax(), "\n", "<br />", -1))
				}
			}
			fmt.Fprintf(w, "|skill 2 = %s\n|skill 2 lv1 = %s%s\n|procs 2 = %d\n%s",
				html.EscapeString(skill2.Name),
				html.EscapeString(strings.Replace(skill2.SkillMin(), "\n", "<br />", -1)),
				lastEvoSkill2Max,
				skill2.MaxCount,
				randomSkillEffects(skill2, "2"),
			)

			// Check if the second skill expires
			if (skill2.PublicEndDatetime.After(time.Time{})) {
				fmt.Fprintf(w, "|skill 2 end = %v\n", skill2.PublicEndDatetime)
			}
		} else if lastEvo.Id > 0 {
			skill2 = lastEvo.Skill2(VcData)
			if skill2 != nil {
				fmt.Fprintf(w, "|skill 2 = %s\n|skill 2 lv10 = %s\n|procs 2 = %d\n%s",
					html.EscapeString(skill2.Name),
					html.EscapeString(strings.Replace(skill2.SkillMin(), "\n", "<br />", -1)),
					skill2.MaxCount,
					randomSkillEffects(skill2, "2"),
				)

				// Check if the second skill expires
				if (skill2.PublicEndDatetime.After(time.Time{})) {
					fmt.Fprintf(w, "|skill 2 end = %v\n", skill2.PublicEndDatetime)
				}
			}
		}
	} else {
		io.WriteString(w, "|rarity = \n|skill = \n|skill lv1 = \n|skill lv10 = \n|procs = \n")
	}
	if evo, ok := evolutions["A"]; ok {
		aSkillName := evo.Skill1Name(VcData)
		if aSkillName != firstEvo.Skill1Name(VcData) {
			fmt.Fprintf(w, "|skill a = %s\n", html.EscapeString(aSkillName))
			fmt.Fprintf(w, "|skill a lv1 = %s\n|skill a lv10 = %s\n|procs a = %s\n%s",
				html.EscapeString(strings.Replace(evo.SkillMin(VcData), "\n", "<br />", -1)),
				html.EscapeString(strings.Replace(evo.SkillMax(VcData), "\n", "<br />", -1)),
				evo.SkillProcs(VcData),
				randomSkillEffects(evo.Skill1(VcData), "a"),
			)
		}
		if gevo, ok := evolutions["GA"]; ok {
			gSkillName := strings.Replace(gevo.Skill1Name(VcData), "☆", "", 1)
			if gSkillName != aSkillName {
				fmt.Fprintf(w, "|skill ga = %s\n", html.EscapeString(gSkillName))
			}
			fmt.Fprintf(w, "|skill ga lv1 = %s\n|skill ga lv10 = %s\n|procs ga = %s\n%s",
				html.EscapeString(strings.Replace(gevo.SkillMin(VcData), "\n", "<br />", -1)),
				html.EscapeString(strings.Replace(gevo.SkillMax(VcData), "\n", "<br />", -1)),
				gevo.SkillProcs(VcData),
				randomSkillEffects(gevo.Skill1(VcData), "ga"),
			)
		}
	}
	if evo, ok := evolutions["G"]; ok {
		gSkillName := strings.Replace(evo.Skill1Name(VcData), "☆", "", 1)
		if gSkillName != firstEvo.Skill1Name(VcData) {
			fmt.Fprintf(w, "|skill g = %s\n", html.EscapeString(gSkillName))
		}
		fmt.Fprintf(w, "|skill g lv1 = %s\n|skill g lv10 = %s\n|procs g = %s\n%s",
			html.EscapeString(strings.Replace(evo.SkillMin(VcData), "\n", "<br />", -1)),
			html.EscapeString(strings.Replace(evo.SkillMax(VcData), "\n", "<br />", -1)),
			evo.SkillProcs(VcData),
			randomSkillEffects(evo.Skill1(VcData), "g"),
		)
		skill2 := evo.Skill2(VcData)
		if skill2 != nil {
			fs2 := firstEvo.Skill2(VcData)
			if fs2 == nil || fs2.Id != skill2.Id {
				fmt.Fprintf(w, "|skill g2 = %s\n", html.EscapeString(skill2.Name))
				fmt.Fprintf(w, "|skill g2 lv1 = %s\n|procs g2 = %d\n%s",
					html.EscapeString(strings.Replace(skill2.SkillMin(), "\n", "<br />", -1)),
					skill2.MaxCount,
					randomSkillEffects(skill2, "g2"),
				)

				// Check if the second skill expires
				if (skill2.PublicEndDatetime.After(time.Time{})) {
					fmt.Fprintf(w, "|skill g2 end = %v\n", skill2.PublicEndDatetime)
				}
			}
		}

	}

	//traverse evolutions in order
	var evokeys []string
	for k := range evolutions {
		evokeys = append(evokeys, k)
	}
	sort.Strings(evokeys)
	for _, k := range evokeys {
		evo := evolutions[k]
		fmt.Fprintf(w, "|cost %[1]s = %[2]d\n|atk %[1]s = %[3]d / %s\n|def %[1]s = %[5]d / %s\n|soldiers %[1]s = %[7]d / %s\n",
			strings.ToLower(k),
			evo.DeckCost,
			evo.DefaultOffense, maxStatAtk(evo, len(evolutions)),
			evo.DefaultDefense, maxStatDef(evo, len(evolutions)),
			evo.DefaultFollower, maxStatFollower(evo, len(evolutions)))
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

	var awakenInfo *vc.CardAwaken
	for _, val := range VcData.Awakenings {
		if lastEvo.Id == val.BaseCardId {
			awakenInfo = &val
			break
		}
	}
	if awakenInfo != nil {
		fmt.Fprintf(w, "|awaken chance = %d\n|awaken orb = %d\n|awaken l = %d\n|awaken m = %d\n|awaken s = %d\n",
			awakenInfo.Percent,
			awakenInfo.Material1Count,
			awakenInfo.Material2Count,
			awakenInfo.Material3Count,
			awakenInfo.Material4Count,
		)
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
		fmt.Fprintf(w,
			`<div style="float: left; margin: 3px"><a href="/images/card/%s"><img src="/images/card/%[1]s"/></a><br />%s : %s☆</div>`,
			evo.Image(),
			evo.Name, k)
	}
	io.WriteString(w, "</div>")
	io.WriteString(w, "</body></html>")
}

func cardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-cards-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{"Id", "Card #", "Name", "Evo Rank", "TransCardId", "Rarity", "Element", "Deck Cost", "Base ATK",
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold", "Skill 1 Name", "Skill Min",
		"Skill Max", "Skill Procs", "Target Scope", "Target Logic", "Skill 2", "Skill Special", "Description", "Friendship",
		"Login", "Meet", "Battle Start", "Battle End", "Friendship Max", "Friendship Event", "Is Closed"})
	for _, card := range VcData.Cards {
		err := cw.Write([]string{strconv.Itoa(card.Id), fmt.Sprintf("cd_%05d", card.CardNo), card.Name, strconv.Itoa(card.EvolutionRank),
			strconv.Itoa(card.TransCardId), card.Rarity(), card.Element(), strconv.Itoa(card.DeckCost), strconv.Itoa(card.DefaultOffense),
			strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower), strconv.Itoa(card.MaxOffense),
			strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower), card.Skill1Name(VcData),
			card.SkillMin(VcData), card.SkillMax(VcData), card.SkillProcs(VcData), card.SkillTarget(VcData),
			card.SkillTargetLogic(VcData), card.Skill2Name(VcData), card.SpecialSkill1Name(VcData),
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
<button type="submit">Submit</button>
</form>
<div>
<table><thead><tr>
<th>_id</th><th>card_no</th><th>name</th><th>evolution_rank</th><th>Next Evo</th><th>Rarity</th><th>Element</th><th>Character ID</th><th>deck_cost</th><th>default_offense</th><th>default_defense</th><th>default_follower</th><th>max_offense</th><th>max_defense</th><th>max_follower</th><th>Skill 1 Name</th><th>Skill Min</th><th>Skill Max</th><th>Skill Procs</th><th>Min Effect</th><th>Min Rate</th><th>Max Effect</th><th>Max Rate</th><th>Target Scope</th><th>Target Logic</th><th>Skill 2</th><th>Skill Special</th><th>Description</th><th>Friendship</th><th>Login</th><th>Meet</th><th>Battle Start</th><th>Battle End</th><th>Friendship Max</th><th>Friendship Event</th>
</tr></thead>
<tbody>
`,
		qs.Get("name"),
		qs.Get("skillname"),
		qs.Get("skilldesc"),
	)
	for _, card := range VcData.Cards {
		if !filter(&card) {
			continue
		}
		skill1 := card.Skill1(VcData)
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2(VcData)
		// skillS1 := card.SpecialSkill1(VcData)
		fmt.Fprintf(w, "<tr><td>%d</td><td><a href=\"/cards/detail/%[1]d\">%05[2]d</a></td><td><a href=\"/cards/detail/%[1]d\">%[3]s</a></td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></td></tr>\n",
			card.Id, card.CardNo, card.Name, card.EvolutionRank, card.EvolutionCardId, card.Rarity(), card.Element(), card.CardCharaId,
			card.DeckCost, card.DefaultOffense, card.DefaultDefense, card.DefaultFollower, card.MaxOffense,
			card.MaxDefense, card.MaxFollower, card.Skill1Name(VcData), card.SkillMin(VcData), card.SkillMax(VcData),
			card.SkillProcs(VcData), skill1.EffectDefaultValue, skill1.DefaultRatio, skill1.EffectMaxValue, skill1.MaxRatio,
			card.SkillTarget(VcData), card.SkillTargetLogic(VcData), card.Skill2Name(VcData),
			card.SpecialSkill1Name(VcData), card.Description(VcData), card.Friendship(VcData), card.Login(VcData),
			card.Meet(VcData), card.BattleStart(VcData), card.BattleEnd(VcData), card.FriendshipMax(VcData),
			card.FriendshipEvent(VcData))
	}

	io.WriteString(w, "</tbody></table></div></body></html>")
}

func maxStatAtk(evo vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 {
		return strconv.Itoa(evo.MaxOffense)
	}
	return "?"
}
func maxStatDef(evo vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 {
		return strconv.Itoa(evo.MaxDefense)
	}
	return "?"
}
func maxStatFollower(evo vc.Card, numOfEvos int) string {
	if evo.EvolutionRank == 0 {
		return strconv.Itoa(evo.MaxFollower)
	}
	return "?"
}

func randomSkillEffects(bs *vc.Skill, evoMod string) (ret string) {
	ret = ""
	if evoMod != "" {
		evoMod += " "
	}
	if bs != nil && bs.EffectId == 36 {
		// Random Skill
		for k, v := range []int{bs.EffectParam, bs.EffectParam2, bs.EffectParam3, bs.EffectParam4, bs.EffectParam5} {
			rs := vc.SkillScan(v, VcData.Skills)
			if rs != nil {
				ret += fmt.Sprintf("|random %s%d = %s \n", evoMod, k+1, rs.FireMin())
			}
		}
	}
	return ret
}

func getEvolutions(card *vc.Card) map[string]vc.Card {
	ret := make(map[string]vc.Card)

	// handle cards like Chimrey and Time Traveler (enemy)
	if card.CardCharaId < 1 {
		ret["0"] = *card
		ret["F"] = *card
		ret["H"] = *card
		return ret
	}

	getAmalBaseCard := func(card *vc.Card) map[string]vc.Card {
		if card.IsAmalgamation(VcData.Amalgamations) {
			// check for a base amalgamation with the same name
			// if there is one, use that for the base card
			for _, amal := range card.Amalgamations(VcData) {
				if card.Id == amal.FusionCardId {
					// material 1
					ac := vc.CardScan(amal.Material1, VcData.Cards)
					if ac.Id != card.Id && ac.Name == card.Name {
						return getEvolutions(ac)
					}
					// material 2
					ac = vc.CardScan(amal.Material2, VcData.Cards)
					if ac.Id != card.Id && ac.Name == card.Name {
						return getEvolutions(ac)
					}
					// material 3
					ac = vc.CardScan(amal.Material3, VcData.Cards)
					if ac != nil && ac.Id != card.Id && ac.Name == card.Name {
						return getEvolutions(ac)
					}
					// material 4
					ac = vc.CardScan(amal.Material4, VcData.Cards)
					if ac != nil && ac.Id != card.Id && ac.Name == card.Name {
						return getEvolutions(ac)
					}
				}
			}
		}
		return nil
	}

	// find the lowest evolution and work from there.
	if card.Rarity()[0] == 'G' {
		bc := card.AwakensFrom(VcData)
		if bc != nil {
			// look for self amalgamation (like sulis)
			amalBaseCard := getAmalBaseCard(bc)
			if amalBaseCard != nil {
				return amalBaseCard
			}
			return getEvolutions(bc)
		}
	} else if card.EvolutionRank != 0 {
		// check for a previous evolution
		for _, c := range VcData.Cards {
			if c.EvolutionCardId == card.Id {
				return getEvolutions(&c)
			}
		}
	} else {
		// check for self amalgamation (like sulis)
		amalBaseCard := getAmalBaseCard(card)
		if amalBaseCard != nil {
			return amalBaseCard
		}
	}

	// get the actual evolution list

	ret[strconv.Itoa(card.EvolutionRank)] = *card
	ret["F"] = *card
	nextId := card.EvolutionCardId
	lastEvo := card
	for nextId > 1 {
		nextCard := vc.CardScan(nextId, VcData.Cards)
		// verify that we haven't switched characters
		if card.CardCharaId == nextCard.CardCharaId {
			nextEvo := strconv.Itoa(nextCard.EvolutionRank)
			ret[nextEvo] = *nextCard
			lastEvo = nextCard
		}
		nextId = nextCard.EvolutionCardId
	}
	ret["H"] = *lastEvo

	// check if the card has a known awakening:
	cardg := lastEvo.AwakensTo(VcData)
	// this doesn't mean that the card has an awakening, it just makes it easier to find
	if cardg != nil {
		ret["G"] = *cardg
	} else {
		// look for the awakening the hard way now, based on the character id
		gs := 0
		for _, val := range VcData.Cards {
			if card.CardCharaId == val.CardCharaId && val.Rarity()[0] == 'G' {
				if gs == 0 {
					ret["G"] = val
				} else {
					ret["G"+strconv.Itoa(gs)] = val
				}
				gs++
			}
		}
	}

	return ret
}

func getAmalgamations(evolutions map[string]vc.Card) []vc.Amalgamation {
	amalgamations := make([]vc.Amalgamation, 0)
	for idx, evo := range evolutions {
		if idx == "H" || idx == "F" {
			continue
		}
		os.Stdout.WriteString(fmt.Sprintf("Card: %d, Name: %s, Evo: %s", evo.Id, evo.Name, idx))
		a := evo.Amalgamations(VcData)
		if len(a) > 0 {
			os.Stdout.WriteString("	Found\n")
			amalgamations = append(amalgamations, a...)
		} else {
			os.Stdout.WriteString("	Not Found\n")
		}
	}
	sort.Sort(vc.ByMaterialCount(amalgamations))
	return amalgamations
}
