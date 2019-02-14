package wiki

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

//Card Main card info
type Card struct {
	IsUnReleased       bool // true to mark as NOT released
	Element            string
	Rarity             string
	Description        string
	Friendship         string
	Login              string
	Meet               string
	BattleStart        string
	BattleEnd          string
	FriendshipMax      string
	FriendshipEvent    string
	Rebirth            string
	TurnoverFrom       string
	TurnoverTo         string
	LikabilityQuotes   []string
	Evolutions         []EvolutionDetails
	Amalgamations      []Amalgamation
	AwakeningMaterials []AwakenMaterial
	RebirthMaterials   []RebirthMaterial
	Availability       string
}

//EvolutionDetails Indivicual evolution specific information
type EvolutionDetails struct {
	EvolutionKey string
	MaxLevel     int
	Cost         int
	BaseAttack   int
	BaseDefense  int
	BaseSoldiers int
	Attack       string
	Defense      string
	Soldiers     string
	Medals       int
	Gold         int
	Skills       []CardSkill
}

//AwakeningMaterialSize Enumeration of awakening material sizes
type AwakeningMaterialSize int

const (
	//Crystal Crystal awakening material
	Crystal AwakeningMaterialSize = 1 << iota // 0
	//Orb Orb awakening material
	Orb // 2
	//Large Large awakening material
	Large // 4
	//Medium Medium awakening material
	Medium // 8
	//Small Small awakening material
	Small // 16
)

//AwakenMaterial Amalgamation Material information
type AwakenMaterial struct {
	Size  AwakeningMaterialSize // Crystal, Orb, L, M, S,
	Count int
}

//RebirthMaterial Rebirth Material information
type RebirthMaterial struct {
	ItemName string
	Count    int
}

// Amalgamation Single amalgamation
type Amalgamation struct {
	Materials []AmalgamationMaterial
	Result    AmalgamationMaterial
}

//AmalgamationMaterial Sinalg material used in amalgamation
type AmalgamationMaterial struct {
	Name   string
	Rarity string
}

//From Creates a template card from a VC card
func (Card) From(c *vc.Card) Card {
	availability := ""
	for _, evo := range c.GetEvolutionCards() {
		if evo.IsAmalgamation() {
			availability += " [[Amalgamation]]"
			break
		}
	}
	skillsSeen := make(tmpSkillsSeen)
	tc := Card{
		IsUnReleased:       c.IsClosed != 0,
		Element:            c.Element(),
		Rarity:             c.MainRarity(),
		Description:        c.Description(),
		Friendship:         c.Friendship(),
		Login:              c.Login(),
		Meet:               c.Meet(),
		BattleStart:        c.BattleStart(),
		BattleEnd:          c.BattleEnd(),
		FriendshipMax:      c.FriendshipMax(),
		FriendshipEvent:    c.FriendshipEvent(),
		Rebirth:            c.RebirthEvent(),
		TurnoverFrom:       getTurnoverFrom(c),
		TurnoverTo:         getTurnoverTo(c),
		LikabilityQuotes:   getLikabilityQuotes(c),
		Evolutions:         getEvolutions(c, &skillsSeen),
		Amalgamations:      getAmalgamations(c),
		AwakeningMaterials: getAwakeningMaterials(c),
		RebirthMaterials:   getRebirthMaterials(c),
		Availability:       availability,
	}
	return tc
}

func getTurnoverFrom(c *vc.Card) string {
	source := c.EvoAccidentOf()
	if source == nil {
		return ""
	}
	return source.Name
}
func getTurnoverTo(c *vc.Card) string {
	result := c.EvoAccident()
	if result == nil {
		return ""
	}
	return result.Name
}
func getLikabilityQuotes(c *vc.Card) (ret []string) {
	ret = make([]string, 0, 5)
	exists := make(map[string]struct{}, 0)
	evolutions := c.GetEvolutionCards()
	for _, evo := range evolutions {
		aw := evo.Archwitch()
		if aw != nil {
			for _, like := range aw.Likeability() {
				if _, ok := exists[like.Likability]; !ok {
					ret = append(ret, like.Likability)
					exists[like.Likability] = struct{}{}
				}
			}
			if len(ret) > 0 {
				// if the AW record had quotes, don't look for more
				return
			}
		}
	}
	return
}

func getEvolutions(c *vc.Card, tSkillsSeen *tmpSkillsSeen) []EvolutionDetails {
	evos := c.GetEvolutions()
	ret := make([]EvolutionDetails, 0, len(evos))
	for _, evoKey := range vc.EvoOrder {
		if evo, ok := evos[evoKey]; ok && evo != nil {
			atk, def, soldier := getStats(evo)
			ret = append(ret, EvolutionDetails{
				EvolutionKey: evoKey,
				Cost:         evo.DeckCost,
				MaxLevel:     evo.CardRarity().MaxCardLevel,
				Gold:         evo.Price,
				Medals:       evo.MedalRate,
				BaseAttack:   evo.DefaultOffense,
				BaseDefense:  evo.DefaultDefense,
				BaseSoldiers: evo.DefaultFollower,
				Attack:       atk,
				Defense:      def,
				Soldiers:     soldier,
				Skills:       getSkills(evo, evoKey, tSkillsSeen),
			})
		}

	}
	return ret
}

func getStats(evo *vc.Card) (atk, def, sol string) {
	numOfEvos := len(evo.GetEvolutions())
	if evo.CardRarity().MaxCardLevel == 1 && numOfEvos == 1 {
		// only X cards have a max level of 1 and they don't evo
		// only possible amalgamations like Philosopher's Stones
		if evo.IsAmalgamation() {
			atkStat, defStat, solStat := evo.AmalgamationPerfect()
			atk = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkStat)
			def = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defStat)
			sol = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solStat)
		}
		return
	}
	if evo.IsAmalgamation() {
		atkStat, defStat, solStat := evo.EvoStandard()
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect()
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				if evo.PossibleMixedEvo() {
					atkMStat, defMStat, solMStat := evo.EvoMixed()
					if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
					}
				}
				atkLRStat, defLRStat, solLRStat := evo.AmalgamationLRStaticLvl1()
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
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect()
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", atkPStat)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", defPStat)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", solPStat)
			}
		}
	} else if evo.EvolutionRank < 2 {
		// not an amalgamation.
		atkStat, defStat, solStat := evo.EvoStandard()
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		atkPStat, defPStat, solPStat := evo.EvoPerfect()
		if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
			var evoType string
			if evo.PossibleMixedEvo() {
				evoType = "Amalgamation"
				atkMStat, defMStat, solMStat := evo.EvoMixed()
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
		atkStat, defStat, solStat := evo.EvoStandard()
		atk = " / " + strconv.Itoa(atkStat)
		def = " / " + strconv.Itoa(defStat)
		sol = " / " + strconv.Itoa(solStat)
		printedMixed := false
		printedPerfect := false
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			if evo.PossibleMixedEvo() {
				atkMStat, defMStat, solMStat := evo.EvoMixed()
				if atkStat != atkMStat || defStat != defMStat || solStat != solMStat {
					printedMixed = true
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", atkMStat)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", defMStat)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", solMStat)
				}
			}
			atkPStat, defPStat, solPStat := evo.AmalgamationPerfect()
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				atkLRStat, defLRStat, solLRStat := evo.AmalgamationLRStaticLvl1()
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
				atkStat, defStat, solStat := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 card evo stats
					atkStat, defStat, solStat = evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 9)
				}
			}
			//If 4* card, calculate 16 card evo stats
			atkStat, defStat, solStat := evo.EvoPerfect()
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
			evo.EvoStandardLvl1() // just to print out the level 1 G stats
			atkPStat, defPStat, solPStat := evo.EvoPerfect()
			if atkStat != atkPStat || defStat != defPStat || solStat != solPStat {
				var evoType string
				if !printedMixed && evo.PossibleMixedEvo() {
					evoType = "Amalgamation"
					atkMStat, defMStat, solMStat := evo.EvoMixed()
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
			awakensFrom := evo.AwakensFrom()
			if awakensFrom == nil && evo.RebirthsFrom() != nil {
				awakensFrom = evo.RebirthsFrom().AwakensFrom()
			}
			if awakensFrom != nil && awakensFrom.LastEvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				atkStat, defStat, solStat := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 and 16 card evo stats
					atkStat, defStat, solStat = evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 9)
					atkStat, defStat, solStat = evo.EvoPerfect()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", atkStat, 16)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", defStat, 16)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", solStat, 16)
				}
			}
		}
	}
	return
}

func getAmalgamations(c *vc.Card) []Amalgamation {
	ret := make([]Amalgamation, 0)
	for _, evo := range c.GetEvolutionCards() {
		for _, amal := range evo.Amalgamations() {
			mats := amal.MaterialsOnly()
			myMats := make([]AmalgamationMaterial, 0, len(mats))
			for _, mat := range mats {
				if mat != nil {
					myMats = append(myMats, AmalgamationMaterial{
						Name:   mat.Name,
						Rarity: mat.Rarity(),
					})
				}
			}
			amalResult := amal.Result()
			if amalResult != nil {
				ret = append(ret, Amalgamation{
					Materials: myMats,
					Result: AmalgamationMaterial{
						Name:   amalResult.Name,
						Rarity: amalResult.Rarity(),
					},
				})
			}
		}
	}
	return ret
}

func getAwakeningMaterials(c *vc.Card) []AwakenMaterial {
	ret := make([]AwakenMaterial, 0, 5)
	for _, evo := range c.GetEvolutionCards() {
		if evo.EvoIsAwoken() {
			var awakenInfo *vc.CardAwaken
			for idx, val := range vc.Data.Awakenings {
				if evo.ID == val.ResultCardID && val.IsClosed == 0 {
					awakenInfo = &vc.Data.Awakenings[idx]
					break
				}
			}
			if awakenInfo != nil {
				items := awakenInfo.ItemCounts()
				for _, itemCount := range items {
					awMat := getAwakenMaterial(itemCount.Item, itemCount.Count)
					if awMat != nil {
						ret = append(ret, *awMat)
					}
				}
				if len(ret) > 0 {
					sort.Slice(ret, func(a, b int) bool {
						return ret[a].Size < ret[b].Size
					})
					return ret
				}
			}
		}
	}
	return ret
}

func getAwakenMaterial(item *vc.Item, count int) *AwakenMaterial {
	if count <= 0 {
		return nil
	}
	if strings.Contains(item.NameEng, "Crystal") {
		return &AwakenMaterial{
			Size:  Crystal,
			Count: count,
		}
	} else if strings.Contains(item.NameEng, "Orb") {
		return &AwakenMaterial{
			Size:  Orb,
			Count: count,
		}
	} else if strings.Contains(item.NameEng, "(L)") {
		return &AwakenMaterial{
			Size:  Large,
			Count: count,
		}
	} else if strings.Contains(item.NameEng, "(M)") {
		return &AwakenMaterial{
			Size:  Medium,
			Count: count,
		}
	} else if strings.Contains(item.NameEng, "(S)") {
		return &AwakenMaterial{
			Size:  Small,
			Count: count,
		}
	}
	log.Printf("Unknown Awakening Material: %d:%s", item.ID, item.NameEng)
	return nil
}

func getRebirthMaterials(c *vc.Card) []RebirthMaterial {
	ret := make([]RebirthMaterial, 0, 3)
	for _, evo := range c.GetEvolutionCards() {
		if evo.EvoIsReborn() {
			var awakenInfo *vc.CardAwaken
			for idx, val := range vc.Data.Rebirths {
				if evo.ID == val.ResultCardID && val.IsClosed == 0 {
					awakenInfo = &vc.Data.Rebirths[idx]
					break
				}
			}
			if awakenInfo != nil {
				items := awakenInfo.ItemCounts()
				for _, itemCount := range items {
					ret = append(ret, RebirthMaterial{
						ItemName: itemCount.Item.NameEng,
						Count:    itemCount.Count,
					})
				}
				if len(ret) > 0 {
					sort.Slice(ret, func(a, b int) bool {
						return ret[a].Count > ret[b].Count
					})
					return ret
				}
			}
		}
	}
	return ret
}
