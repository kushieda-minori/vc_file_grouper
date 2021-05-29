package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"vc_file_grouper/structout"
	"vc_file_grouper/util"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
	"vc_file_grouper/wiki/api"
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
	if err != nil || cardID < 1 {
		http.Error(w, "Invalid card id "+pathParts[2], http.StatusNotFound)
		return
	}

	card := vc.CardScan(cardID)
	if card == nil {
		http.Error(w, "Invalid card id "+pathParts[2]+"\nCard not found.", http.StatusNotFound)
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
	var firstEvo *vc.Card
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

	cardName := firstEvo.Name
	if len(cardName) == 0 {
		cardName = firstEvo.Image()
	}

	prevCard, prevCardName := getPreviousCard(card)
	nextCard, nextCardName := getNextCard(card)

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
	fmt.Fprintf(w, "<div style=\"clear:left;float:left\"><a href=\"?action=uploadImages\">Upload Missing Images</a>\n<br /></div>")
	fmt.Fprintf(w, "<div style=\"float:left;padding-left:15px;\"><a href=\"?action=createOnWiki\">Create New Wiki Page</a></div>")

	if action := r.FormValue("action"); action != "" {
		if action == "uploadImages" {
			err = api.UploadNewCardUniqueImages(card)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
		if action == "createOnWiki" {
			err = api.CreateCardPage(card, "New Card")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
	}

	io.WriteString(w, "<div><textarea readonly=\"readonly\" style=\"width:100%;height:450px\">")

	cardPage := wiki.CardPage{}
	if card.IsClosed != 0 {
		cardPage.PageHeader = "{{Unreleased}}"
	}
	cardPage.CardInfo.UpdateAll(firstEvo, avail)

	io.WriteString(w, cardPage.String())

	//Write out amalgamations here
	if len(amalgamations) > 0 {
		io.WriteString(w, "\n==''[[Amalgamation]]''==\n")
		for _, v := range amalgamations {
			mats := v.Materials()
			l := len(mats)
			if l > 2 {
				fmt.Fprintf(w, "{{Amalgamation|matcount = %d\n|name 1 = %s|rarity 1 = %s\n|name 2 = %s|rarity 2 = %s\n|name 3 = %s|rarity 3 = %s\n",
					l-1, mats[0].Name, mats[0].Rarity(), mats[1].Name, mats[1].Rarity(), mats[2].Name, mats[2].Rarity())
				if l > 3 {
					fmt.Fprintf(w, "|name 4 = %s|rarity 4 = %s\n", mats[3].Name, mats[3].Rarity())
					if l > 4 {
						fmt.Fprintf(w, "|name 5 = %s|rarity 5 = %s\n", mats[4].Name, mats[4].Rarity())
					}
				}
				io.WriteString(w, "}}\n")
			}
		}
	}
	io.WriteString(w, "</textarea></div>")
	// show images here
	io.WriteString(w, "<div style=\"float:left\">")
	iconEvos := card.EvosWithDistinctImages(true)
	lenEvoKeys := len(evokeys)
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
		if _, err := os.Stat(filepath.Join(vc.FilePath, "card", "hd", evo.Image())); err == nil {
			pathPart = "cardHD"
		} else if _, err := os.Stat(filepath.Join(vc.FilePath, "card", "md", evo.Image())); err == nil {
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
		"Base DEF", "Base Sol", "Max ATK", "Max DEF", "Max Sold",
		"Skill 1 Name", "Skill 1 Fire Min", "Skill 1 Fire Max", "Skill 1 Min", "Skill 1 Max", "Skill 1 Procs", "Target Scope 1", "Target Logic 1",
		"Skill 2", "Skill 2 Fire Min", "Skill 2 Fire Max", "Skill 2 Min", "Skill 2 Max", "Skill 2 Procs", "Target Scope 2", "Target Logic 2",
		"Skill 3", "Skill 3 Fire Min", "Skill 3 Fire Max", "Skill 3 Min", "Skill 3 Max", "Skill 3 Procs", "Target Scope 3", "Target Logic 3",
		"Skill Special", "Skill SP Fire Min", "Skill SP Fire Max", "Skill SP Min", "Skill SP Max", "Skill SP Procs", "Target Scope SP", "Target Logic SP",
		"Thor Skill 1", "Skill TH Fire",
		"Description", "Friendship", "Login", "Meet",
		"Battle Start", "Battle End", "Friendship Max", "Friendship Event",
		"Is Closed"})
	for _, card := range vc.Data.Cards {
		err := cw.Write([]string{strconv.Itoa(card.ID), fmt.Sprintf("cd_%05d", card.CardNo), card.Name, strconv.Itoa(card.EvolutionRank),
			strconv.Itoa(card.TransCardID), card.Rarity(), card.Element(), strconv.Itoa(card.DeckCost), strconv.Itoa(card.DefaultOffense),
			strconv.Itoa(card.DefaultDefense), strconv.Itoa(card.DefaultFollower), strconv.Itoa(card.MaxOffense), strconv.Itoa(card.MaxDefense), strconv.Itoa(card.MaxFollower),
			card.Skill1Name(), card.Skill1().FireMin(), card.Skill1().FireMax(), card.SkillMin(), card.SkillMax(), card.SkillProcs(), card.SkillTarget(), card.SkillTargetLogic(),
			card.Skill2Name(), card.Skill2().FireMin(), card.Skill2().FireMax(), card.Skill2().SkillMin(), card.Skill2().SkillMax(), card.Skill2().ActivationString(), card.Skill2().TargetScope(), card.Skill2().TargetLogic(),
			card.Skill3Name(), card.Skill3().FireMin(), card.Skill3().FireMax(), card.Skill3().SkillMin(), card.Skill3().SkillMax(), card.Skill3().ActivationString(), card.Skill3().TargetScope(), card.Skill3().TargetLogic(),
			card.SpecialSkill1Name(), card.SpecialSkill1().FireMin(), card.SpecialSkill1().FireMax(), card.SpecialSkill1().SkillMin(), card.SpecialSkill1().SkillMax(), card.SpecialSkill1().ActivationString(), card.SpecialSkill1().TargetScope(), card.SpecialSkill1().TargetLogic(),
			card.ThorSkill1Name(),
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

//CardJSONStatHandler outputs card stats as JSON
func CardJSONStatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-glr-cards-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".json\"")
	w.Header().Set("Content-Type", "application/json")

	statCards := make([]structout.CardStatInfo, 0)

	for _, card := range vc.Data.Cards {
		if card.Rarity() == "GLR" || card.EvoIsReborn() {
			log.Printf("found GLR or Reborn card %d:%s - %s", card.ID, card.Name, card.Rarity())
			statCards = append(statCards, structout.ToCardStatInfo(card))
		}
	}
	jsonenc, err := json.MarshalIndent(statCards, "", "  ")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(jsonenc)
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
		"Max Rarity Sol",
		"Rebirth base Atk",
		"Rebirth base Def",
		"Rebirth base Sol",
		"Rebirth max Atk",
		"Rebirth max Def",
		"Rebirth max Sol",
		"Rebirth gain Atk",
		"Rebirth gain Def",
		"Rebirth gain Sol",
		"Rebirth gain/lvl Atk",
		"Rebirth gain/lvl Def",
		"Rebirth gain/lvl Sol",
		"Rebirth Max Rarity Atk",
		"Rebirth Max Rarity Def",
		"Rebirth Max Rarity Sol",
		"Rebirth Max Level",
	})

	cards := vc.Data.Cards.Copy()
	log.Printf("Total card count: %d", len(cards))
	// sort by name A-Z
	sort.Slice(cards, func(i, j int) bool {
		first := cards[i]
		second := cards[j]

		return first.Name < second.Name || (first.Name == second.Name && first.CardRareID < second.CardRareID)
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

			var RebirthbaseAtk, RebirthbaseDef, RebirthbaseSol,
				RebirthmaxAtk, RebirthmaxDef, RebirthmaxSol,
				RebirthgainAtk, RebirthgainDef, RebirthgainSol,
				RebirthgainlvlAtk, RebirthgainlvlDef, RebirthgainlvlSol,
				RebirthMaxRarityAtk, RebirthMaxRarityDef, RebirthMaxRaritySol,
				RebirthMaxLevel string
			if rb := card.RebirthsTo(); rb != nil {
				rbRare := rb.CardRarity()
				rmaxLevel := float64(rbRare.MaxCardLevel - 1)
				ratkGain := rb.MaxOffense - rb.DefaultOffense
				rdefGain := rb.MaxDefense - rb.DefaultDefense
				rsolGain := rb.MaxFollower - rb.DefaultFollower

				RebirthbaseAtk = strconv.Itoa(rb.DefaultOffense)
				RebirthbaseDef = strconv.Itoa(rb.DefaultDefense)
				RebirthbaseSol = strconv.Itoa(rb.DefaultFollower)
				RebirthmaxAtk = strconv.Itoa(rb.MaxOffense)
				RebirthmaxDef = strconv.Itoa(rb.MaxDefense)
				RebirthmaxSol = strconv.Itoa(rb.MaxFollower)
				RebirthgainAtk = strconv.Itoa(ratkGain)
				RebirthgainDef = strconv.Itoa(rdefGain)
				RebirthgainSol = strconv.Itoa(rsolGain)
				RebirthgainlvlAtk = fmt.Sprintf("%.4f", float64(ratkGain)/rmaxLevel)
				RebirthgainlvlDef = fmt.Sprintf("%.4f", float64(rdefGain)/rmaxLevel)
				RebirthgainlvlSol = fmt.Sprintf("%.4f", float64(rsolGain)/rmaxLevel)
				RebirthMaxRarityAtk = strconv.Itoa(rbRare.LimtOffense)
				RebirthMaxRarityDef = strconv.Itoa(rbRare.LimtDefense)
				RebirthMaxRaritySol = strconv.Itoa(rbRare.LimtMaxFollower)
				RebirthMaxLevel = strconv.Itoa(rbRare.MaxCardLevel)
			}

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
				RebirthbaseAtk,
				RebirthbaseDef,
				RebirthbaseSol,
				RebirthmaxAtk,
				RebirthmaxDef,
				RebirthmaxSol,
				RebirthgainAtk,
				RebirthgainDef,
				RebirthgainSol,
				RebirthgainlvlAtk,
				RebirthgainlvlDef,
				RebirthgainlvlSol,
				RebirthMaxRarityAtk,
				RebirthMaxRarityDef,
				RebirthMaxRaritySol,
				RebirthMaxLevel,
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
			match = match && card.MainRarity() == rarity
		}
		if element := qs.Get("element"); element != "" {
			match = match && card.Element() == element
		}
		if symbols := qs.Get("symbol"); symbols != "" {
			symbol, err := strconv.Atoi(symbols)
			if err == nil {
				match = match && card.CardSymbolID == symbol
			}
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
		if isClosed := qs.Get("isClosed"); isClosed != "" {
			match = match && card.IsClosed > 0
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
				s1 = skill1.Fire != "" &&
					(strings.Contains(strings.ToLower(skill1.Fire), strings.ToLower(skilldesc)) ||
						strings.Contains(strings.ToLower(skill1.SkillMin()), strings.ToLower(skilldesc)))
				//log.Printf(skill1.Fire + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(); skill2 != nil {
				s2 = skill2.Fire != "" &&
					(strings.Contains(strings.ToLower(skill2.Fire), strings.ToLower(skilldesc)) ||
						strings.Contains(strings.ToLower(skill2.SkillMin()), strings.ToLower(skilldesc)))
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
`+rarityToOptions(qs.Get("rarity"))+`
</select>
<label for="f_element">Element:</label><select id="f_element" name="element" value="%s">
<option value=""></option>
`+elementToOptions(qs.Get("element"))+`
</select>
<label for="f_symbol">Symbol:</label><select id="f_symbol" name="symbol" value="%s">
<option value=""></option>
`+symbolNamesToOptions(qs.Get("symbol"))+`
</select>
<label for="f_evos">Evo:</label><select id="f_evos" name="evos" value="%s">
<option value=""></option>
`+evosToOptions(qs.Get("evos"))+`
</select>
<label for="f_skillname">Skill Name:</label><input id="f_skillname" name="skillname" value="%s" />
<label for="f_skilldesc">Skill Description:</label><input id="f_skilldesc" name="skilldesc" value="%s" />
<label for="f_skillisthor">Has Thor Skill:</label><input id="f_skillisthor" name="isThor" type="checkbox" value="checked" %s />
<label for="f_hasrebirth">Has Rebirth:</label><input id="f_hasrebirth" name="hasRebirth" type="checkbox" value="checked" %s />
<label for="f_isclosed">Closed (not released):</label><input id="f_isclosed" name="isClosed" type="checkbox" value="checked" %s />
<button type="submit">Submit</button>
</form>
`,
		qs.Get("name"),
		qs.Get("rarity"),
		qs.Get("element"),
		qs.Get("symbol"),
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
		"Symbol",
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
		"Leader Skill",
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
			card.CardSymbolID,
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
			strconv.Itoa(card.LeaderSkillID),
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

func rarityToOptions(selected string) (ret string) {
	for _, v := range []string{"X", "N", "R", "SR", "UR", "LR", "VR"} {
		if v == selected {
			ret += fmt.Sprintf(`<option value="%s" selected="selected">%s</option>`, v, v)
		} else {
			ret += fmt.Sprintf(`<option value="%s">%s</option>`, v, v)
		}
	}
	return
}

func elementToOptions(selected string) (ret string) {
	for _, v := range []string{"Special", "Cool", "Passion", "Light", "Dark"} {
		if v == selected {
			ret += fmt.Sprintf(`<option value="%s" selected="selected">%s</option>`, v, v)
		} else {
			ret += fmt.Sprintf(`<option value="%s">%s</option>`, v, v)
		}
	}
	return
}

func symbolNamesToOptions(selected string) (ret string) {
	for i, v := range vc.Data.SymbolNames {
		if v == selected || strconv.Itoa(i) == selected {
			ret += fmt.Sprintf(`<option value="%d" selected="selected">%s</option>`, i, v)
		} else {
			ret += fmt.Sprintf(`<option value="%d">%s</option>`, i, v)
		}
	}
	return
}

func evosToOptions(selected string) (ret string) {
	evos := map[string]string{"1": "1★", "4": "4★"}
	for k, v := range evos {
		if k == selected {
			ret += fmt.Sprintf(`<option value="%s" selected="selected">%s</option>`, k, v)
		} else {
			ret += fmt.Sprintf(`<option value="%s">%s</option>`, k, v)
		}
	}
	return
}

// getPreviousCard finds the card before the earliest evolution of this card. Does
// not take into account cards that were released for awakening at a later date
// maybe should do this by character ID instead?
func getPreviousCard(card *vc.Card) (prev *vc.Card, prevName string) {
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

// getNextCard finds the card before the latest evolution of this card. Does
// not take into account cards that were released for awakening at a later date
// maybe should do this by character ID instead?
func getNextCard(card *vc.Card) (next *vc.Card, nextName string) {
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
