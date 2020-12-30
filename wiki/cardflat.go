package wiki

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"vc_file_grouper/vc"
)

//CardFlat Main card info
type CardFlat struct {
	Element string `json:"element"`
	Rarity  string `json:"rarity"`
	Symbol  string `json:"symbol"`

	Skill         string `json:"skill"`
	SkillLv1      string `json:"skill lv1"`
	SkillLv10     string `json:"skill lv10"`
	SkillLv1Cost  string `json:"skill lv1 cost"`
	SkillLv10Cost string `json:"skill lv10 cost"`
	SkillProcs    string `json:"procs"`
	SkillRandom1  string `json:"random 1"`
	SkillRandom2  string `json:"random 2"`
	SkillRandom3  string `json:"random 3"`
	SkillRandom4  string `json:"random 4"`
	SkillRandom5  string `json:"random 5"`

	Skill2         string `json:"skill 2"`
	Skill2Lv1      string `json:"skill 2 lv1"`
	Skill2Lv10     string `json:"skill 2 lv10"`
	Skill2Lv1Cost  string `json:"skill 2 lv1 cost"`
	Skill2Lv10Cost string `json:"skill 2 lv10 cost"`
	Skill2Procs    string `json:"procs 2"`
	Skill2End      string `json:"skill 2 end"`
	Skill2Random1  string `json:"random 2 1"`
	Skill2Random2  string `json:"random 2 2"`
	Skill2Random3  string `json:"random 2 3"`
	Skill2Random4  string `json:"random 2 4"`
	Skill2Random5  string `json:"random 2 5"`

	Skill3         string `json:"skill 3"`
	Skill3Lv1      string `json:"skill 3 lv1"`
	Skill3Lv10     string `json:"skill 3 lv10"`
	Skill3Lv1Cost  string `json:"skill 3 lv1 cost"`
	Skill3Lv10Cost string `json:"skill 3 lv10 cost"`
	Skill3Procs    string `json:"procs 3"`
	Skill3End      string `json:"skill 3 end"`
	Skill3Random1  string `json:"random 3 1"`
	Skill3Random2  string `json:"random 3 2"`
	Skill3Random3  string `json:"random 3 3"`
	Skill3Random4  string `json:"random 3 4"`
	Skill3Random5  string `json:"random 3 5"`

	SkillA         string `json:"skill a"`
	SkillALv1      string `json:"skill a lv1"`
	SkillALv10     string `json:"skill a lv10"`
	SkillALv1Cost  string `json:"skill a lv1 cost"`
	SkillALv10Cost string `json:"skill a lv10 cost"`
	SkillAProcs    string `json:"procs a"`
	SkillAEnd      string `json:"skill a end"`
	SkillARandom1  string `json:"random a 1"`
	SkillARandom2  string `json:"random a 2"`
	SkillARandom3  string `json:"random a 3"`
	SkillARandom4  string `json:"random a 4"`
	SkillARandom5  string `json:"random a 5"`

	SkillG         string `json:"skill g"`
	SkillGLv1      string `json:"skill g lv1"`
	SkillGLv10     string `json:"skill g lv10"`
	SkillGLv1Cost  string `json:"skill g lv1 cost"`
	SkillGLv10Cost string `json:"skill g lv10 cost"`
	SkillGProcs    string `json:"procs g"`
	SkillGRandom1  string `json:"random g 1"`
	SkillGRandom2  string `json:"random g 2"`
	SkillGRandom3  string `json:"random g 3"`
	SkillGRandom4  string `json:"random g 4"`
	SkillGRandom5  string `json:"random g 5"`

	SkillG2         string `json:"skill g2"`
	SkillG2Lv1      string `json:"skill g2 lv1"`
	SkillG2Lv10     string `json:"skill g2 lv10"`
	SkillG2Lv1Cost  string `json:"skill g2 lv1 cost"`
	SkillG2Lv10Cost string `json:"skill g2 lv10 cost"`
	SkillG2Procs    string `json:"procs g2"`
	SkillG2End      string `json:"skill g2 end"`
	SkillG2Random1  string `json:"random g2 1"`
	SkillG2Random2  string `json:"random g2 2"`
	SkillG2Random3  string `json:"random g2 3"`
	SkillG2Random4  string `json:"random g2 4"`
	SkillG2Random5  string `json:"random g2 5"`

	SkillG3         string `json:"skill g3"`
	SkillG3Lv1      string `json:"skill g3 lv1"`
	SkillG3Lv10     string `json:"skill g3 lv10"`
	SkillG3Lv1Cost  string `json:"skill g3 lv1 cost"`
	SkillG3Lv10Cost string `json:"skill g3 lv10 cost"`
	SkillG3Procs    string `json:"procs g3"`
	SkillG3End      string `json:"skill g3 end"`
	SkillG3Random1  string `json:"random g3 1"`
	SkillG3Random2  string `json:"random g3 2"`
	SkillG3Random3  string `json:"random g3 3"`
	SkillG3Random4  string `json:"random g3 4"`
	SkillG3Random5  string `json:"random g3 5"`

	SkillGA         string `json:"skill ga"`
	SkillGALv1      string `json:"skill ga lv1"`
	SkillGALv10     string `json:"skill ga lv10"`
	SkillGALv1Cost  string `json:"skill ga lv1 cost"`
	SkillGALv10Cost string `json:"skill ga lv10 cost"`
	SkillGAProcs    string `json:"procs ga"`
	SkillGARandom1  string `json:"random ga 1"`
	SkillGARandom2  string `json:"random ga 2"`
	SkillGARandom3  string `json:"random ga 3"`
	SkillGARandom4  string `json:"random ga 4"`
	SkillGARandom5  string `json:"random ga 5"`

	SkillX         string `json:"skill x"`
	SkillXLv1      string `json:"skill x lv1"`
	SkillXLv10     string `json:"skill x lv10"`
	SkillXLv1Cost  string `json:"skill x lv1 cost"`
	SkillXLv10Cost string `json:"skill x lv10 cost"`
	SkillXProcs    string `json:"procs x"`
	SkillXRandom1  string `json:"random x 1"`
	SkillXRandom2  string `json:"random x 2"`
	SkillXRandom3  string `json:"random x 3"`
	SkillXRandom4  string `json:"random x 4"`
	SkillXRandom5  string `json:"random x 5"`

	SkillX2         string `json:"skill x2"`
	SkillX2Lv1      string `json:"skill x2 lv1"`
	SkillX2Lv10     string `json:"skill x2 lv10"`
	SkillX2Lv1Cost  string `json:"skill x2 lv1 cost"`
	SkillX2Lv10Cost string `json:"skill x2 lv10 cost"`
	SkillX2Procs    string `json:"procs x2"`
	SkillX2End      string `json:"skill x2 end"`
	SkillX2Random1  string `json:"random x2 1"`
	SkillX2Random2  string `json:"random x2 2"`
	SkillX2Random3  string `json:"random x2 3"`
	SkillX2Random4  string `json:"random x2 4"`
	SkillX2Random5  string `json:"random x2 5"`

	SkillX3         string `json:"skill x3"`
	SkillX3Lv1      string `json:"skill x3 lv1"`
	SkillX3Lv10     string `json:"skill x3 lv10"`
	SkillX3Lv1Cost  string `json:"skill x3 lv1 cost"`
	SkillX3Lv10Cost string `json:"skill x3 lv10 cost"`
	SkillX3Procs    string `json:"procs x3"`
	SkillX3End      string `json:"skill x3 end"`
	SkillX3Random1  string `json:"random x3 1"`
	SkillX3Random2  string `json:"random x3 2"`
	SkillX3Random3  string `json:"random x3 3"`
	SkillX3Random4  string `json:"random x3 4"`
	SkillX3Random5  string `json:"random x3 5"`

	SkillT      string `json:"skill t"`
	SkillTLv1   string `json:"skill t lv1"`
	SkillTLv10  string `json:"skill t lv10"`
	SkillTProcs string `json:"procs t"`
	SkillTEnd   string `json:"skill t end"`

	SkillGT      string `json:"skill gt"`
	SkillGTLv1   string `json:"skill gt lv1"`
	SkillGTLv10  string `json:"skill gt lv10"`
	SkillGTProcs string `json:"procs gt"`
	SkillGTEnd   string `json:"skill gt end"`

	SkillXT      string `json:"skill xt"`
	SkillXTLv1   string `json:"skill xt lv1"`
	SkillXTLv10  string `json:"skill xt lv10"`
	SkillXTProcs string `json:"procs xt"`
	SkillXTEnd   string `json:"skill xt end"`

	MaxLevel0 string `json:"max level 0"`
	Cost0     string `json:"cost 0"`
	Atk0      string `json:"atk 0"`
	Def0      string `json:"def 0"`
	Soldiers0 string `json:"soldiers 0"`
	Medals0   string `json:"medals 0"`
	Gold0     string `json:"gold 0"`

	MaxLevel1 string `json:"max level 1"`
	Cost1     string `json:"cost 1"`
	Atk1      string `json:"atk 1"`
	Def1      string `json:"def 1"`
	Soldiers1 string `json:"soldiers 1"`
	Medals1   string `json:"medals 1"`
	Gold1     string `json:"gold 1"`

	MaxLevel2 string `json:"max level 2"`
	Cost2     string `json:"cost 2"`
	Atk2      string `json:"atk 2"`
	Def2      string `json:"def 2"`
	Soldiers2 string `json:"soldiers 2"`
	Medals2   string `json:"medals 2"`
	Gold2     string `json:"gold 2"`

	MaxLevel3 string `json:"max level 3"`
	Cost3     string `json:"cost 3"`
	Atk3      string `json:"atk 3"`
	Def3      string `json:"def 3"`
	Soldiers3 string `json:"soldiers 3"`
	Medals3   string `json:"medals 3"`
	Gold3     string `json:"gold 3"`

	MaxLevel4 string `json:"max level 4"`
	Cost4     string `json:"cost 4"`
	Atk4      string `json:"atk 4"`
	Def4      string `json:"def 4"`
	Soldiers4 string `json:"soldiers 4"`
	Medals4   string `json:"medals 4"`
	Gold4     string `json:"gold 4"`

	MaxLevelA string `json:"max level a"`
	CostA     string `json:"cost a"`
	AtkA      string `json:"atk a"`
	DefA      string `json:"def a"`
	SoldiersA string `json:"soldiers a"`
	MedalsA   string `json:"medals a"`
	GoldA     string `json:"gold a"`

	MaxLevelG string `json:"max level g"`
	CostG     string `json:"cost g"`
	AtkG      string `json:"atk g"`
	DefG      string `json:"def g"`
	SoldiersG string `json:"soldiers g"`
	MedalsG   string `json:"medals g"`
	GoldG     string `json:"gold g"`

	MaxLevelGA string `json:"max level ga"`
	CostGA     string `json:"cost ga"`
	AtkGA      string `json:"atk ga"`
	DefGA      string `json:"def ga"`
	SoldiersGA string `json:"soldiers ga"`
	MedalsGA   string `json:"medals ga"`
	GoldGA     string `json:"gold ga"`

	MaxLevelX string `json:"max level x"`
	CostX     string `json:"cost x"`
	AtkX      string `json:"atk x"`
	DefX      string `json:"def x"`
	SoldiersX string `json:"soldiers x"`
	MedalsX   string `json:"medals x"`
	GoldX     string `json:"gold x"`

	FriendshipPoints string `json:"friendship points"`

	Login           string `json:"login"`
	Description     string `json:"description"`
	Friendship      string `json:"friendship"`
	Meet            string `json:"meet"`
	BattleStart     string `json:"battle start"`
	BattleEnd       string `json:"battle end"`
	FriendshipMax   string `json:"friendship max"`
	FriendshipEvent string `json:"friendship event"`
	Rebirth         string `json:"rebirth"`

	QuoteMisc1 string `json:"quote misc 1"`
	QuoteMisc2 string `json:"quote misc 2"`
	QuoteMisc3 string `json:"quote misc 3"`
	QuoteMisc4 string `json:"quote misc 4"`

	Likeability0 string `json:"likeability 0"`
	Likeability1 string `json:"likeability 1"`
	Likeability2 string `json:"likeability 2"`
	Likeability3 string `json:"likeability 3"`
	Likeability4 string `json:"likeability 4"`
	Likeability5 string `json:"likeability 5"`
	Likeability6 string `json:"likeability 6"`
	Likeability7 string `json:"likeability 7"`

	AwakenChance  string `json:"awaken chance"`
	AwakenCrystal string `json:"awaken crystal"`
	AwakenOrb     string `json:"awaken orb"`
	AwakenL       string `json:"awaken l"`
	AwakenM       string `json:"awaken m"`
	AwakenS       string `json:"awaken s"`

	RebirthChance     string `json:"rebirth chance"`
	RebirthItem1      string `json:"rebirth item 1"`
	RebirthItem1Count string `json:"rebirth item 1 count"`
	RebirthItem2      string `json:"rebirth item 2"`
	RebirthItem2Count string `json:"rebirth item 2 count"`
	RebirthItem3      string `json:"rebirth item 3"`
	RebirthItem3Count string `json:"rebirth item 3 count"`

	TurnOverTo   string `json:"turnoverto"`
	TurnOverFrom string `json:"turnoverfrom"`
	Availability string `json:"availability"`

	unknownFields map[string]string
}

//OldNew old and new value
type OldNew struct {
	Old string `json:"old"`
	New string `json:"new"`
}

//cleanVal repalces all double line breaks with single line breaks
func cleanVal(v string) string {
	linebreakRegEx, _ := regexp.Compile(`(\s*[\r\n]\s*)+|(\s*<br\s*[/]?>\s*)+`)
	return linebreakRegEx.ReplaceAllString(strings.TrimSpace(v), "<br />")
}

func (c CardFlat) asMap() map[string]string {
	var inInterface map[string]string
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

//String outputs the struct as a Wiki Template:Card call
func (c *CardFlat) String() (ret string) {
	if c == nil {
		return ""
	}

	inInterface := c.asMap()

	// begin template
	ret += "{{Card\n"
	// iterate through record fields
	for _, field := range cardFieldOrder {
		if val, ok := inInterface[field]; ok {
			if strings.TrimSpace(val) != "" {
				ret += fmt.Sprintf("|%s = %s\n", field, cleanVal(val))
			}
		}
	}

	if len(c.unknownFields) > 0 {
		keys := make([]string, 0)
		for k, v := range c.unknownFields {
			// skip blank values
			if strings.TrimSpace(v) != "" {
				keys = append(keys, k)
			}
		}
		if len(keys) > 0 {
			//ret += "<!-- these fields were unknown to the bot, but have not been removed -->\n"
			ret += "\n"
			sort.Strings(keys)
			for _, field := range keys {
				ret += fmt.Sprintf("|%s = %s\n", field, cleanVal(c.unknownFields[field]))
			}
		}
	}
	// end template
	ret += "}}\n"

	return ret
}

//Equals validates that one card flat equals another
func (c CardFlat) Equals(that CardFlat) bool {
	thism := c.asMap()
	thatm := that.asMap()

	if len(thism) != len(thatm) {
		return false
	}
	for k, thisv := range thism {
		thatv, ok := thatm[k]
		if !ok || thisv != thatv {
			return false
		}
	}
	return true
}

//Differences validates that one card flat equals another
func (c CardFlat) Differences(that CardFlat) (ret map[string]OldNew) {
	thism := c.asMap()
	thatm := that.asMap()

	ret = make(map[string]OldNew, 0)
	seen := make(map[string]bool)

	for k, thisv := range thism {
		seen[k] = true
		thatv, ok := thatm[k]
		if !ok {
			ret[k] = OldNew{thisv, ""}
		} else if thisv != thatv {
			ret[k] = OldNew{thisv, thatv}
		}
	}

	if len(thism) != len(thatm) {
		for k, thatv := range thatm {
			// check if we've already recorded the change
			if _, ok := seen[k]; !ok {
				// check that it's in this
				thisv, ok := thism[k]
				if !ok {
					ret[k] = OldNew{"", thatv}
				} else if thisv != thatv {
					ret[k] = OldNew{thisv, thatv}
				}
			}
		}
	}
	return
}

func getInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func getTime(s string) *time.Time {
	// 2020-12-16 12:00:00 +0900 JST
	t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", s)
	if err == nil {
		return &t
	}
	log.Println(err.Error())
	return nil
}

//UpdateAll updates all data in the template as possible. The "availablity" is only updated if it's not already populated.
func (c *CardFlat) UpdateAll(vcCard *vc.Card, availability string) {
	c.UpdateBaseData(vcCard)
	c.UpdateEvoStats(vcCard.GetEvolutions())
	c.UpdateExchangeInfo(vcCard.GetEvolutions())
	c.UpdateSkills(vcCard.GetEvolutions())
	c.UpdateQuotes(vcCard)
	c.UpdateAwakenRebirthInfo(vcCard.GetEvolutions())
	if strings.TrimSpace(c.Availability) == "" {
		if availability != "" {
			c.Availability = availability
		} else {
			for _, evo := range vcCard.GetEvolutionCards() {
				if evo.IsAmalgamation() {
					c.Availability = "[[Amalgamation]]"
					break
				}
			}
		}
	}
}

//UpdateBaseData Updates the tempalte information from the VC data. Fields updated are `Element`, `Rarity`, `Symbol` and turn over information for evo accidents
func (c *CardFlat) UpdateBaseData(vcCard *vc.Card) {
	c.Element = vcCard.Element()
	c.Rarity = vcCard.MainRarity()
	c.Symbol = vcCard.Symbol()

	evolutions := vcCard.GetEvolutions()
	for _, k := range vc.EvoOrder {
		if evo, ok := evolutions[k]; ok {
			turnoverfrom := evo.EvoAccidentOf()
			if turnoverfrom != nil {
				c.TurnOverFrom = turnoverfrom.Name
			}
			turnoverto := evo.EvoAccident()
			if turnoverto != nil {
				c.TurnOverTo = turnoverto.Name
			}
		}
	}
}

//UpdateExchangeInfo Updates Gold and Medal exchange info for a specific evo
func (c *CardFlat) UpdateExchangeInfo(evolutions map[string]*vc.Card) {
	rs := reflect.ValueOf(c).Elem()
	lenEvos := len(evolutions)
	for evoCode, evo := range evolutions {
		evoCode = fixEvoKey(evoCode, evo, lenEvos)
		goldField := rs.FieldByName("Gold" + evoCode)
		if goldField.IsValid() {
			log.Printf("Evo %s: Gold %d, Medals %d", evoCode, evo.Price, evo.MedalRate)
			goldField.SetString(strconv.Itoa(evo.Price))
			rs.FieldByName("Medals" + evoCode).SetString(strconv.Itoa(evo.MedalRate))
		} else {
			log.Printf("Unknown evo: %s", evoCode)
		}
	}
}

//UpdateEvoStats updates stats for the evos provided
func (c *CardFlat) UpdateEvoStats(evolutions map[string]*vc.Card) {
	rs := reflect.ValueOf(c).Elem()

	lenEvos := len(evolutions)
	for evoCode, evo := range evolutions {
		evoCode = fixEvoKey(evoCode, evo, lenEvos)

		log.Printf("******Settings stats for evo %s", evoCode)
		a, d, s := maxStats(evo, len(evolutions))
		rs.FieldByName("MaxLevel" + evoCode).SetString(strconv.Itoa(evo.CardRarity().MaxCardLevel))
		rs.FieldByName("Cost" + evoCode).SetString(strconv.Itoa(evo.DeckCost))
		rs.FieldByName("Atk" + evoCode).SetString(fmt.Sprintf("%d%s", evo.DefaultOffense, a))
		rs.FieldByName("Def" + evoCode).SetString(fmt.Sprintf("%d%s", evo.DefaultDefense, d))
		rs.FieldByName("Soldiers" + evoCode).SetString(fmt.Sprintf("%d%s", evo.DefaultFollower, s))
	}
}

func fixEvoKey(evoCode string, evo *vc.Card, len int) string {
	if evoCode == "H" {
		if evo.EvolutionRank >= 0 {
			evoCode = strconv.Itoa(evo.EvolutionRank)
		} else if len == 1 {
			evoCode = "1"
		}
	}
	return evoCode
}

//UpdateAwakenRebirthInfo Updates awakening and rebirth information based on the cards provided. The cards provided are expected to only be valid evos of this card.
func (c *CardFlat) UpdateAwakenRebirthInfo(evolutions map[string]*vc.Card) {
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
			c.AwakenChance = strconv.Itoa(awakenInfo.Percent)
			for i := 1; i <= 5; i++ {
				item, count := awakenInfo.ItemAndCount(i)
				if item != nil {
					sCount := strconv.Itoa(count)
					if strings.Contains(item.NameEng, "Crystal") {
						c.AwakenCrystal = sCount
					} else if strings.Contains(item.NameEng, "Orb") {
						c.AwakenOrb = sCount
					} else if strings.Contains(item.NameEng, "(L)") {
						c.AwakenL = sCount
					} else if strings.Contains(item.NameEng, "(M)") {
						c.AwakenM = sCount
					} else if strings.Contains(item.NameEng, "(S)") {
						c.AwakenS = sCount
					} else {
						log.Printf("***Unknown Awakening item: %s\n", item.NameEng)
					}
				}
			}
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
			c.RebirthChance = strconv.Itoa(rebirthInfo.Percent)
			item, count := rebirthInfo.ItemAndCount(1)
			if item != nil {
				c.RebirthItem1 = item.NameEng
				c.RebirthItem1Count = strconv.Itoa(count)
			}
			item, count = rebirthInfo.ItemAndCount(2)
			if item != nil {
				c.RebirthItem2 = item.NameEng
				c.RebirthItem2Count = strconv.Itoa(count)
			}
			item, count = rebirthInfo.ItemAndCount(3)
			if item != nil {
				c.RebirthItem3 = item.NameEng
				c.RebirthItem3Count = strconv.Itoa(count)
			}
		}
	}
}

//UpdateSkills Update skills for the card
func (c *CardFlat) UpdateSkills(evolutions map[string]*vc.Card) {
	skillsSeen := make(tmpSkillsSeen)
	rs := reflect.ValueOf(c).Elem()

	setSkill := func(s *vc.Skill, ls *vc.Skill, evoKey string, num int, mod string) {
		if s == nil {
			return
		}
		if _, seen := (skillsSeen)[s]; seen {
			return
		}
		(skillsSeen)[s] = tmpSkillHolder{
			Skill:         s,
			SkillNum:      num,
			SkillFirstEvo: evoKey,
		}
		// need to find if this is an evo-maxed skill
		min := s.SkillMin()
		max := s.SkillMax()
		// thor skills use the "Fire" text
		if mod == "t" {
			min = s.FireMax()
			max = s.FireMax()
		}
		if ls != nil {
			// thor skills use the "Fire" text
			if mod == "t" {
				max = ls.FireMax()
			} else {
				max = ls.SkillMax()
			}
			(skillsSeen)[ls] = tmpSkillHolder{
				Skill:         ls,
				SkillNum:      num,
				SkillFirstEvo: evoKey,
			}
		}

		if min == max {
			max = ""
		}

		if evoKey == "0" {
			evoKey = ""
		}

		// set the value of the skill
		skillPrefix := "Skill" + strings.ToUpper(evoKey+mod)
		skillNameField := rs.FieldByName(skillPrefix)
		if skillNameField.IsValid() {
			skillNameField.SetString(cleanVal(s.Name))
			rs.FieldByName(skillPrefix + "Lv1").SetString(cleanVal(min))
			rs.FieldByName(skillPrefix + "Lv10").SetString(cleanVal(max))
			f := rs.FieldByName(skillPrefix + "Lv1Cost")
			if f.IsValid() {
				if s.DefaultCost > 0 {
					f.SetString(strconv.Itoa(s.DefaultCost))
				} else {
					f.SetString("")
				}
			}

			f = rs.FieldByName(skillPrefix + "Lv10Cost")
			if f.IsValid() {
				if s.MaxCost > 0 {
					f.SetString(strconv.Itoa(s.MaxCost))
				} else {
					f.SetString("")
				}
			}
			rs.FieldByName(skillPrefix + "Procs").SetString(s.ActivationString())

			if s.EffectID == 36 {
				// Random Skill
				for i, v := range []int{s.EffectParam, s.EffectParam2, s.EffectParam3, s.EffectParam4, s.EffectParam5} {
					sr := vc.SkillScan(v)
					if sr != nil {
						rs.FieldByName(skillPrefix + "Random" + strconv.Itoa(i+1)).SetString(cleanVal(sr.FireMin()))
					}
				}
			}
			if s.Expires() {
				endField := rs.FieldByName(skillPrefix + "End")
				if endField.IsValid() {
					ed := fmt.Sprintf("%v", s.PublicEndDatetime)
					endField.SetString(ed)
				} else {
					log.Printf(skillPrefix + " had an end date set, but the template is not configured to accept end dates for that skill.")
				}
			}
		} else {
			log.Println(skillPrefix + " expected but not found in template")
		}
	}

	for _, evoID := range vc.EvoOrder {
		if evo, ok := evolutions[evoID]; ok {
			var lastEvo *vc.Card
			if evo.EvoIsFirst() {
				lastEvo = evo.LastEvo()
			}
			setSkill(evo.Skill1(), lastEvo.Skill1(), evoID, 1, "")
			setSkill(evo.Skill2(), lastEvo.Skill2(), evoID, 2, "2")
			setSkill(evo.Skill3(), lastEvo.Skill3(), evoID, 3, "3")
			setSkill(evo.ThorSkill1(), lastEvo.ThorSkill1(), evoID, 4, "t")
		}
	}
}

//UpdateQuotes Updates the card quotes
func (c *CardFlat) UpdateQuotes(card *vc.Card) {
	c.FriendshipPoints = strconv.Itoa(card.Character().MaxFriendship)
	c.Description = cleanVal(card.Description())
	c.Friendship = cleanVal(card.Friendship())
	c.Login = cleanVal(card.Login())
	c.Meet = cleanVal(card.Meet())
	c.BattleStart = cleanVal(card.BattleStart())
	c.BattleEnd = cleanVal(card.BattleEnd())
	c.FriendshipMax = cleanVal(card.FriendshipMax())
	c.FriendshipEvent = cleanVal(card.FriendshipEvent())
	c.Rebirth = cleanVal(card.RebirthEvent())

	lQuotes := getLikabilityQuotes(card)
	if len(lQuotes) > 0 {
		rs := reflect.ValueOf(c).Elem()
		for i, q := range lQuotes {
			f := rs.FieldByName("Likeability" + strconv.Itoa(i))
			if f.IsValid() {
				f.SetString(cleanVal(q))
			} else {
				log.Fatalf("Unable to find liability quote for %d", i)
			}
		}
	}
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

	if evo.CardRarity() != nil && evo.CardRarity().MaxCardLevel == 1 && numOfEvos == 1 {
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
			if stats.NotEquals(pStats) {
				if evo.PossibleMixedEvo() {
					mStats := evo.EvoMixed()
					if stats.NotEquals(mStats) {
						atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
						def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
						sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
					}
				}
				lrStats := evo.AmalgamationLRStaticLvl1()
				if lrStats.NotEquals(pStats) {
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
			if stats.NotEquals(pStats) {
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
		if stats.NotEquals(pStats) {
			var evoType string
			if evo.PossibleMixedEvo() {
				evoType = "Amalgamation"
				mStats := evo.EvoMixed()
				if stats.NotEquals(mStats) {
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
				if stats.NotEquals(mStats) {
					printedMixed = true
					atk += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Attack)
					def += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Defense)
					sol += fmt.Sprintf(" / {{tooltip|%d|Mixed Evolution}}", mStats.Soldiers)
				}
			}
			pStats := evo.AmalgamationPerfect()
			if stats.NotEquals(pStats) {
				lrStats := evo.AmalgamationLRStaticLvl1()
				if lrStats.NotEquals(pStats) {
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
		if evo.EvoIsAwoken() || evo.EvoIsReborn() {
			evo.EvoStandardLvl1() // just to print out the level 1 G stats
			pStats := evo.EvoPerfect()
			if stats.NotEquals(pStats) {
				var evoType string
				if !printedMixed && evo.PossibleMixedEvo() {
					evoType = "Amalgamation"
					mStats := evo.EvoMixed()
					if stats.NotEquals(mStats) {
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
				if evo.MainRarity() != "LR" {
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

func parseCard(pageText string) (map[string]string, int, error) {
	positionalParamNum := 0

	ret := make(map[string]string)
	rPos := -1
	runes := []rune(pageText)
	rLen := len(runes)
	for rPos = len("{{card"); rPos < rLen; rPos++ {
		r := runes[rPos]
		if r == '}' && runes[rPos+1] == '}' {
			// at the end of the card info. break;
			rPos++
			break
		}
		if r == '|' {
			rPos++
			currentKey, err := getNextTempalteKey(&runes, &rPos)
			if err != nil {
				return nil, 0, err
			}
			if currentKey == "" {
				// positional param
				positionalParamNum++
				currentKey = strconv.Itoa(positionalParamNum)
			}
			currentVal, err := getNextTemplateValue(&runes, &rPos)
			if err != nil {
				return nil, 0, err
			}
			//log.Printf("Found `%s` : `%s`", currentKey, currentVal)
			ret[currentKey] = currentVal
		}
	}

	var cPos int
	if rPos >= rLen {
		log.Printf("no footer found")
		cPos = len(pageText)
	} else {
		cPos = len(string(runes[:rPos+1]))
	}

	return ret, cPos, nil
}

func cardFieldIsKnown(field string) bool {
	for _, f := range cardFieldOrder {
		if field == f {
			return true
		}
	}
	return false
}

var cardFieldOrder []string = []string{
	"element",
	"rarity",
	"symbol",
	"skill",
	"skill lv1",
	"skill lv10",
	"skill lv1 cost",
	"skill lv10 cost",
	"procs",
	"random 1",
	"random 2",
	"random 3",
	"random 4",
	"random 5",
	"skill 2",
	"skill 2 lv1",
	"skill 2 lv10",
	"skill 2 lv1 cost",
	"skill 2 lv10 cost",
	"procs 2",
	"skill 2 end",
	"random 2 1",
	"random 2 2",
	"random 2 3",
	"random 2 4",
	"random 2 5",
	"skill 3",
	"skill 3 lv1",
	"skill 3 lv10",
	"skill 3 lv1 cost",
	"skill 3 lv10 cost",
	"procs 3",
	"skill 3 end",
	"random 3 1",
	"random 3 2",
	"random 3 3",
	"random 3 4",
	"random 3 5",
	"skill g",
	"skill g lv1",
	"skill g lv10",
	"skill g lv1 cost",
	"skill g lv10 cost",
	"procs g",
	"random g 1",
	"random g 2",
	"random g 3",
	"random g 4",
	"random g 5",
	"skill g2",
	"skill g2 lv1",
	"skill g2 lv10",
	"skill g2 lv1 cost",
	"skill g2 lv10 cost",
	"procs g2",
	"skill g2 end",
	"random g2 1",
	"random g2 2",
	"random g2 3",
	"random g2 4",
	"random g2 5",
	"skill g3",
	"skill g3 lv1",
	"skill g3 lv10",
	"skill g3 lv1 cost",
	"skill g3 lv10 cost",
	"procs g3",
	"skill g3 end",
	"random g3 1",
	"random g3 2",
	"random g3 3",
	"random g3 4",
	"random g3 5",
	"skill a",
	"skill a lv1",
	"skill a lv10",
	"skill a lv1 cost",
	"skill a lv10 cost",
	"procs a",
	"skill a end",
	"random a 1",
	"random a 2",
	"random a 3",
	"random a 4",
	"random a 5",
	"skill ga",
	"skill ga lv1",
	"skill ga lv10",
	"skill ga lv1 cost",
	"skill ga lv10 cost",
	"procs ga",
	"random ga 1",
	"random ga 2",
	"random ga 3",
	"random ga 4",
	"random ga 5",
	"skill x",
	"skill x lv1",
	"skill x lv10",
	"skill x lv1 cost",
	"skill x lv10 cost",
	"procs x",
	"random x 1",
	"random x 2",
	"random x 3",
	"random x 4",
	"random x 5",
	"skill x2",
	"skill x2 lv1",
	"skill x2 lv10",
	"skill x2 lv1 cost",
	"skill x2 lv10 cost",
	"procs x2",
	"skill x2 end",
	"random x2 1",
	"random x2 2",
	"random x2 3",
	"random x2 4",
	"random x2 5",
	"skill x3",
	"skill x3 lv1",
	"skill x3 lv10",
	"skill x3 lv1 cost",
	"skill x3 lv10 cost",
	"procs x3",
	"skill x3 end",
	"random x3 1",
	"random x3 2",
	"random x3 3",
	"random x3 4",
	"random x3 5",
	"skill t",
	"skill t lv1",
	"skill t lv10",
	"procs t",
	"skill t end",
	"skill gt",
	"skill gt lv1",
	"skill gt lv10",
	"procs gt",
	"skill gt end",
	"skill xt",
	"skill xt lv1",
	"skill xt lv10",
	"procs xt",
	"skill xt end",
	"max level 0",
	"cost 0",
	"atk 0",
	"def 0",
	"soldiers 0",
	"medals 0",
	"gold 0",
	"max level 1",
	"cost 1",
	"atk 1",
	"def 1",
	"soldiers 1",
	"medals 1",
	"gold 1",
	"max level 2",
	"cost 2",
	"atk 2",
	"def 2",
	"soldiers 2",
	"medals 2",
	"gold 2",
	"max level 3",
	"cost 3",
	"atk 3",
	"def 3",
	"soldiers 3",
	"medals 3",
	"gold 3",
	"max level 4",
	"cost 4",
	"atk 4",
	"def 4",
	"soldiers 4",
	"medals 4",
	"gold 4",
	"max level a",
	"cost a",
	"atk a",
	"def a",
	"soldiers a",
	"medals a",
	"gold a",
	"max level g",
	"cost g",
	"atk g",
	"def g",
	"soldiers g",
	"medals g",
	"gold g",
	"max level ga",
	"cost ga",
	"atk ga",
	"def ga",
	"soldiers ga",
	"medals ga",
	"gold ga",
	"max level x",
	"cost x",
	"atk x",
	"def x",
	"soldiers x",
	"medals x",
	"gold x",
	"friendship points",
	"description",
	"friendship",
	"login",
	"meet",
	"battle start",
	"battle end",
	"friendship max",
	"friendship event",
	"rebirth",
	"likeability 0",
	"likeability 1",
	"likeability 2",
	"likeability 3",
	"likeability 4",
	"likeability 5",
	"likeability 6",
	"likeability 7",
	"quote misc 1",
	"quote misc 2",
	"quote misc 3",
	"quote misc 4",
	"awaken chance",
	"awaken crystal",
	"awaken orb",
	"awaken l",
	"awaken m",
	"awaken s",
	"rebirth chance",
	"rebirth item 1",
	"rebirth item 1 count",
	"rebirth item 2",
	"rebirth item 2 count",
	"rebirth item 3",
	"rebirth item 3 count",
	"turnoverto",
	"turnoverfrom",
	"availability",
}
