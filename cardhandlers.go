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
			"<div style=\"float: left; margin: 3px\"><img src=\"/cardthumbs/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></div>",
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
	evolutions := getEvolutions(*card)

	amalgamations := make([]vc.Amalgamation, 0)
	var turnOverTo, turnOverFrom *vc.Card
	for _, evo := range evolutions {
		// os.Stdout.WriteString(fmt.Sprintf("Evo: %d Accident: %d\n", evo.Id, evo.TransCardId))
		a := evo.Amalgamations(&VcData)
		if len(a) > 0 && len(amalgamations) == 0 {
			amalgamations = append(amalgamations, a...)
		}
		if evo.TransCardId > 0 && turnOverTo == nil {
			turnOverTo = evo.EvoAccident(VcData.Cards)
		} else {
			if turnOverFrom == nil {
				turnOverFrom = evo.EvoAccidentOf(VcData.Cards)
			}
		}
	}
	sort.Sort(vc.ByMaterialCount(amalgamations))

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
	fmt.Fprintf(w, "<div style=\"float:left\">Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br />", card.Name)
	io.WriteString(w, "<textarea style=\"width:800px;height:760px\">")
	if card.IsClosed != 0 {
		io.WriteString(w, "{{Unreleased}}")
	}
	fmt.Fprintf(w, "{{Card\n|element = %s\n", card.Element())
	var firstEvo vc.Card
	var ok bool
	if firstEvo, ok = evolutions["0"]; ok {
		fmt.Fprintf(w, "|rarity = %s\n|skill = %s\n|skill lv1 = %s\n|skill lv10 = %s\n|procs = %s\n",
			firstEvo.Rarity(),
			html.EscapeString(firstEvo.Skill1Name(&VcData)),
			html.EscapeString(strings.Replace(firstEvo.SkillMin(&VcData), "\n", "<br />", -1)),
			html.EscapeString(strings.Replace(firstEvo.SkillMax(&VcData), "\n", "<br />", -1)),
			firstEvo.SkillProcs(&VcData))
	} else if firstEvo, ok = evolutions["1"]; ok {
		fmt.Fprintf(w, "|rarity = %s\n|skill = %s\n|skill lv1 = %s\n|skill lv10 = %s\n|procs = %s\n",
			firstEvo.Rarity(),
			html.EscapeString(firstEvo.Skill1Name(&VcData)),
			html.EscapeString(strings.Replace(firstEvo.SkillMin(&VcData), "\n", "<br />", -1)),
			html.EscapeString(strings.Replace(firstEvo.SkillMax(&VcData), "\n", "<br />", -1)),
			firstEvo.SkillProcs(&VcData))
	} else {
		io.WriteString(w, "|rarity = \n|skill = \n|skill lv1 = \n|skill lv10 = \n|procs = \n")
	}
	if evo, ok := evolutions["A"]; ok {
		aSkillName := evo.Skill1Name(&VcData)
		if aSkillName != firstEvo.Skill1Name(&VcData) {
			fmt.Fprintf(w, "|skill a = %s\n", html.EscapeString(aSkillName))
			fmt.Fprintf(w, "|skill a lv1 = %s\n|skill a lv10 = %s\n|proc a = %s\n",
				html.EscapeString(strings.Replace(evo.SkillMin(&VcData), "\n", "<br />", -1)),
				html.EscapeString(strings.Replace(evo.SkillMax(&VcData), "\n", "<br />", -1)),
				evo.SkillProcs(&VcData))
		}
		if gevo, ok := evolutions["GA"]; ok {
			gSkillName := strings.Replace(gevo.Skill1Name(&VcData), "☆", "", 1)
			if gSkillName != aSkillName {
				fmt.Fprintf(w, "|skill ga = %s\n", html.EscapeString(gSkillName))
			}
			fmt.Fprintf(w, "|skill ga lv1 = %s\n|skill ga lv10 = %s\n|proc ga = %s\n",
				html.EscapeString(strings.Replace(gevo.SkillMin(&VcData), "\n", "<br />", -1)),
				html.EscapeString(strings.Replace(gevo.SkillMax(&VcData), "\n", "<br />", -1)),
				gevo.SkillProcs(&VcData))
		}
	}
	if evo, ok := evolutions["G"]; ok {
		gSkillName := strings.Replace(evo.Skill1Name(&VcData), "☆", "", 1)
		if gSkillName != firstEvo.Skill1Name(&VcData) {
			fmt.Fprintf(w, "|skill g = %s\n", html.EscapeString(gSkillName))
		}
		fmt.Fprintf(w, "|skill g lv1 = %s\n|skill g lv10 = %s\n|proc g = %s\n",
			html.EscapeString(strings.Replace(evo.SkillMin(&VcData), "\n", "<br />", -1)),
			html.EscapeString(strings.Replace(evo.SkillMax(&VcData), "\n", "<br />", -1)),
			evo.SkillProcs(&VcData))
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
		html.EscapeString(card.Description(&VcData)), html.EscapeString(strings.Replace(card.Friendship(&VcData), "\n", "<br />", -1)))
	login := card.Login(&VcData)
	if len(strings.TrimSpace(login)) > 0 {
		fmt.Fprintf(w, "|login = %s\n", html.EscapeString(strings.Replace(login, "\n", "<br />", -1)))
	}
	fmt.Fprintf(w, "|meet = %s\n|battle start = %s\n|battle end = %s\n|friendship max = %s\n|friendship event = %s\n", html.EscapeString(strings.Replace(card.Meet(&VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.BattleStart(&VcData), "\n", "<br />", -1)), html.EscapeString(strings.Replace(card.BattleEnd(&VcData), "\n", "<br />", -1)),
		html.EscapeString(strings.Replace(card.FriendshipMax(&VcData), "\n", "<br />", -1)), html.EscapeString(strings.Replace(card.FriendshipEvent(&VcData), "\n", "<br />", -1)))
	// fmt.Fprintf(w,"|likeability 0 = %s\n|likeability 1 = %s\n|likeability 2 = %s\n|likeability 3 = %s\n|likeability 4 = %s\n|likeability 5 =%s\n",)
	if turnOverFrom != nil {
		fmt.Fprintf(w, "|turnoverfrom = %s\n", turnOverFrom.Name)
	} else if turnOverTo != nil {
		fmt.Fprintf(w, "|turnoverto = %s\n", turnOverTo.Name)
	} else {
		fmt.Fprintf(w, "|availability = %s\n", avail)
	}
	io.WriteString(w, "}}")

	//Write out amalgamations here
	if len(amalgamations) > 0 {
		io.WriteString(w, "\n==''[[Amalgamation]]''==\n")
		for _, v := range amalgamations {
			mats := v.Materials(&VcData)
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
	if card.Skill1(&VcData) != nil && card.Skill1(&VcData).EffectId == 36 {
		// Random Skill
		bs := card.Skill1(&VcData)
		io.WriteString(w, "\n==''Notes''==\nThe list of random skills she may fire are as follows:\n")
		for _, v := range []int{bs.EffectParam, bs.EffectParam2, bs.EffectParam3, bs.EffectParam4, bs.EffectParam5} {
			rs := vc.SkillScan(v, VcData.Skills)
			if rs != nil {
				fmt.Fprintf(w, "* %s \n", rs.SkillMin())
			}
		}
	}
	io.WriteString(w, "</textarea></div>")
	// show images here
	io.WriteString(w, "<div style=\"float:left\">")
	for _, k := range evokeys {
		evo := evolutions[k]
		fmt.Fprintf(w,
			"<div style=\"float: left; margin: 3px\"><img src=\"/cardthumbs/%s\"/><br />%s : %s☆</div>",
			evo.Image(),
			evo.Name, k)
	}
	io.WriteString(w, "<div style=\"clear: both\">")
	for _, k := range evokeys {
		evo := evolutions[k]
		fmt.Fprintf(w,
			"<div style=\"float: left; margin: 3px\"><img src=\"/cardimages/%s\"/><br />%s : %s☆</div>",
			evo.Image(),
			evo.Name, k)
	}
	io.WriteString(w, "</div>")
	io.WriteString(w, "</body></html>")
}

func cardCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-cards-"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
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
			strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower), card.Skill1Name(&VcData),
			card.SkillMin(&VcData), card.SkillMax(&VcData), card.SkillProcs(&VcData), card.SkillTarget(&VcData),
			card.SkillTargetLogic(&VcData), card.Skill2Name(&VcData), card.SpecialSkill1Name(&VcData),
			card.Description(&VcData), card.Friendship(&VcData), card.Login(&VcData), card.Meet(&VcData),
			card.BattleStart(&VcData), card.BattleEnd(&VcData), card.FriendshipMax(&VcData), card.FriendshipEvent(&VcData),
			strconv.Itoa(card.IsClosed),
		})
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}
	cw.Flush()
}

func cardTableHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	io.WriteString(w, "<html><head><title>All Cards</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead>\n")
	io.WriteString(w, "<tr><th>_id</th><th>card_no</th><th>name</th><th>evolution_rank</th><th>Rarity</th><th>Element</th><th>Character ID</th><th>deck_cost</th><th>default_offense</th><th>default_defense</th><th>default_follower</th><th>max_offense</th><th>max_defense</th><th>max_follower</th><th>Skill 1 Name</th><th>Skill Min</th><th>Skill Max</th><th>Skill Procs</th><th>Min Effect</th><th>Min Rate</th><th>Max Effect</th><th>Max Rate</th><th>Target Scope</th><th>Target Logic</th><th>Skill 2</th><th>Skill Special</th><th>Description</th><th>Friendship</th><th>Login</th><th>Meet</th><th>Battle Start</th><th>Battle End</th><th>Friendship Max</th><th>Friendship Event</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, card := range VcData.Cards {
		skill1 := card.Skill1(&VcData)
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2(&VcData)
		// skillS1 := card.SpecialSkill1(&VcData)
		fmt.Fprintf(w, "<tr><td>%d</td><td>%05d</td><td><a href=\"/cards/detail/%[1]d\">%[3]s</a></td><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></td>\n",
			card.Id, card.CardNo, card.Name, card.EvolutionRank, card.Rarity(), card.Element(), card.CardCharaId,
			card.DeckCost, card.DefaultOffense, card.DefaultDefense, card.DefaultFollower, card.MaxOffense,
			card.MaxDefense, card.MaxFollower, card.Skill1Name(&VcData), card.SkillMin(&VcData), card.SkillMax(&VcData),
			card.SkillProcs(&VcData), skill1.EffectDefaultValue, skill1.DefaultRatio, skill1.EffectMaxValue, skill1.MaxRatio,
			card.SkillTarget(&VcData), card.SkillTargetLogic(&VcData), card.Skill2Name(&VcData),
			card.SpecialSkill1Name(&VcData), card.Description(&VcData), card.Friendship(&VcData), card.Login(&VcData),
			card.Meet(&VcData), card.BattleStart(&VcData), card.BattleEnd(&VcData), card.FriendshipMax(&VcData),
			card.FriendshipEvent(&VcData))
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

func getEvolutions(card vc.Card) map[string]vc.Card {
	ret := make(map[string]vc.Card)

	if card.CardCharaId < 1 {
		ret["0"] = card
		return ret
	}

	amalgs := 0
	gs := 0
	gas := 0
	for _, val := range VcData.Cards {
		if card.CardCharaId == val.CardCharaId {
			if val.Rarity()[0] == 'G' {
				baseCard := val.AwakensFrom(&VcData)
				// if base card is nil, that means it's not available yet
				if baseCard != nil && baseCard.IsAmalgamation(VcData.Amalgamations) {
					if gas == 0 {
						ret["GA"] = val
					} else {
						ret["GA"+strconv.Itoa(gas)] = val
					}
					gas++
				} else {
					if gs == 0 {
						ret["G"] = val
					} else {
						ret["G"+strconv.Itoa(gs)] = val
					}
					gs++
				}
			} else if val.IsAmalgamation(VcData.Amalgamations) {
				if amalgs == 0 {
					ret["A"] = val
				} else {
					ret["A"+strconv.Itoa(amalgs)] = val
				}
				amalgs++
			} else {
				iEvo := val.EvolutionRank
				evo := strconv.Itoa(iEvo)
				if _, ok := ret[evo]; ok {
					for ; ok; _, ok = ret[evo] {
						iEvo++
						evo = strconv.Itoa(iEvo)
					}
					ret[evo] = val
				} else {
					ret[evo] = val
				}
			}
		}
	}
	return ret
}
