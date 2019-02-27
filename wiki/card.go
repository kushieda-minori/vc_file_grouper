package wiki

import (
	"encoding/json"
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
	Quotes             CardQuotes
	Evolutions         []EvolutionDetails
	AwakeningMaterials []AwakenMaterial
	RebirthMaterials   []RebirthMaterial
	TurnoverFrom       string
	TurnoverTo         string
	Availability       string
	Amalgamations      []Amalgamation
}

// CardQuotes Quotes that can appear on the cards
type CardQuotes struct {
	Friendship       string
	Login            string
	Meet             string
	BattleStart      string
	BattleEnd        string
	FriendshipMax    string
	FriendshipEvent  string
	Rebirth          string
	Miscellaneous    []string
	LikabilityQuotes []string
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
		IsUnReleased: c.IsClosed != 0,
		Element:      c.Element(),
		Rarity:       c.MainRarity(),
		Description:  c.Description(),
		Quotes: CardQuotes{
			Friendship:       c.Friendship(),
			Login:            c.Login(),
			Meet:             c.Meet(),
			BattleStart:      c.BattleStart(),
			BattleEnd:        c.BattleEnd(),
			FriendshipMax:    c.FriendshipMax(),
			FriendshipEvent:  c.FriendshipEvent(),
			Rebirth:          c.RebirthEvent(),
			LikabilityQuotes: getLikabilityQuotes(c),
		},
		TurnoverFrom:       getTurnoverFrom(c),
		TurnoverTo:         getTurnoverTo(c),
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
		aws := evo.ArchwitchesWithLikeabilityQuotes()
		for _, aw := range aws {
			for _, like := range aw.Likeability() {
				if _, ok := exists[like.Likability]; !ok {
					ret = append(ret, like.Likability)
					exists[like.Likability] = struct{}{}
				}
			}
		}
	}
	log.Printf("Found %d likeability Quotes for card %d:%s", len(ret), c.ID, c.Name)
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
			s := evo.AmalgamationPerfect()
			atk = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", s.Attack)
			def = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", s.Defense)
			sol = fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", s.Soldiers)
		}
		return
	}
	if evo.IsAmalgamation() {
		s := evo.EvoStandard()
		atk = " / " + strconv.Itoa(s.Attack)
		def = " / " + strconv.Itoa(s.Defense)
		sol = " / " + strconv.Itoa(s.Soldiers)
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			pStat := evo.AmalgamationPerfect()
			if s.NotEquals(pStat) {
				if evo.PossibleMixedEvo() {
					mStat := evo.EvoMixed()
					if s.NotEquals(mStat) {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Attack)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Defense)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Soldiers)
					}
				}
				lrStat1 := evo.AmalgamationLRStaticLvl1()
				if lrStat1.NotEquals(pStat) {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Soldiers)
				}
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Soldiers)
			}
		} else {
			pStat := evo.AmalgamationPerfect()
			if s.NotEquals(pStat) {
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Soldiers)
			}
		}
	} else if evo.EvolutionRank < 2 {
		// not an amalgamation.
		s := evo.EvoStandard()
		atk = " / " + strconv.Itoa(s.Attack)
		def = " / " + strconv.Itoa(s.Defense)
		sol = " / " + strconv.Itoa(s.Soldiers)
		pStat := evo.EvoPerfect()
		if s.NotEquals(pStat) {
			var evoType string
			if evo.PossibleMixedEvo() {
				evoType = "Amalgamation"
				mStat := evo.EvoMixed()
				if s.NotEquals(mStat) {
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Soldiers)
				}
			} else {
				evoType = "Evolution"
			}
			atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Attack, evoType)
			def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Defense, evoType)
			sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Soldiers, evoType)
		}
	} else {
		// not an amalgamation, Evo Rank >=2 (Awoken cards or 4* evos).
		s := evo.EvoStandard()
		atk = " / " + strconv.Itoa(s.Attack)
		def = " / " + strconv.Itoa(s.Defense)
		sol = " / " + strconv.Itoa(s.Soldiers)
		printedMixed := false
		printedPerfect := false
		if strings.HasSuffix(evo.Rarity(), "LR") {
			// print LR level1 static material amal
			if evo.PossibleMixedEvo() {
				mStat := evo.EvoMixed()
				if s.NotEquals(mStat) {
					printedMixed = true
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Soldiers)
				}
			}
			pStat := evo.AmalgamationPerfect()
			if s.NotEquals(pStat) {
				lrStat1 := evo.AmalgamationLRStaticLvl1()
				if lrStat1.NotEquals(pStat) {
					atk += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|LR 'Special' material Lvl-1, other materials Perfect Amalgamation}}", lrStat1.Soldiers)
				}
				printedPerfect = true
				atk += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Attack)
				def += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Defense)
				sol += fmt.Sprintf(" / {{tooltip|%d|Perfect Amalgamation}}", pStat.Soldiers)
			}
		}

		if !strings.HasSuffix(evo.Rarity(), "LR") &&
			evo.Rarity()[0] != 'G' && evo.Rarity()[0] != 'X' {
			// TODO need more logic here to check if it's an Amalg vs evo only.
			// may need different options depending on the type of card.
			if evo.EvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				s6 := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s6.Attack, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s6.Defense, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s6.Soldiers, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 card evo stats
					s9 := evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s9.Attack, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s9.Defense, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s9.Soldiers, 9)
				}
			}
			//If 4* card, calculate 16 card evo stats
			pStat := evo.EvoPerfect()
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
			atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", pStat.Attack, cards)
			def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", pStat.Defense, cards)
			sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", pStat.Soldiers, cards)
		}
		if evo.Rarity()[0] == 'G' || evo.Rarity()[0] == 'X' {
			evo.EvoStandardLvl1() // just to print out the level 1 G stats
			pStat := evo.EvoPerfect()
			if s.NotEquals(pStat) {
				var evoType string
				if !printedMixed && evo.PossibleMixedEvo() {
					evoType = "Amalgamation"
					mStat := evo.EvoMixed()
					if s.NotEquals(mStat) {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Attack)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Defense)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStat.Soldiers)
					}
				} else {
					evoType = "Evolution"
				}
				if !printedPerfect {
					atk += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Attack, evoType)
					def += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Defense, evoType)
					sol += fmt.Sprintf(" / {{tooltip|%d|Perfect %s}}", pStat.Soldiers, evoType)
				}
			}
			awakensFrom := evo.AwakensFrom()
			if awakensFrom == nil && evo.RebirthsFrom() != nil {
				awakensFrom = evo.RebirthsFrom().AwakensFrom()
			}
			if awakensFrom != nil && awakensFrom.LastEvolutionRank == 4 {
				//If 4* card, calculate 6 card evo stats
				s := evo.Evo6Card()
				atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Attack, 6)
				def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Defense, 6)
				sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Soldiers, 6)
				if evo.Rarity() != "GLR" {
					//If SR card, calculate 9 and 16 card evo stats
					s = evo.Evo9Card()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Attack, 9)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Defense, 9)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Soldiers, 9)
					s = evo.EvoPerfect()
					atk += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Attack, 16)
					def += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Defense, 16)
					sol += fmt.Sprintf(" / {{tooltip|%d|%d Card Evolution}}", s.Soldiers, 16)
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

//String converts the card data to a wiki string
func (c Card) String() (ret string) {
	if c.IsUnReleased {
		ret += "{{Unreleased}}"
	}
	data, err := json.MarshalIndent(c, "", " ")
	if err == nil {
		ret += fmt.Sprintf(`{{#invoke:Card|detail
|<nowiki>%s</nowiki>
|availability=%s
}}`,
			string(data),
			c.Availability,
		)
	} else {
		ret = err.Error()
	}
	return
}
