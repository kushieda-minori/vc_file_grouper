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

//cleanVal repalces all double line breaks with single line breaks
func cleanVal(v string) string {
	regexp, _ := regexp.Compile(`([\r\n]\s*)+|(<br\s*[/]?>\s*)+`)
	return regexp.ReplaceAllString(v, "<br />")
	//return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v, "\r\n", "\n"), "\n\n", "\n"), "\n", "<br />")
}

//String outputs the struct as a Wiki Template:Card call
func (c *CardFlat) String() (ret string) {
	if c == nil {
		return ""
	}

	var inInterface map[string]string
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)

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
			if v != "" {
				keys = append(keys, k)
			}
		}
		if len(keys) > 0 {
			ret += "<!-- these fields were unknown to the bot, but have not been removed -->\n"
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

//Card converts the flat card representation to the structured representation
func (c CardFlat) Card() (ret Card) {
	ret = Card{
		IsUnReleased: false,
		Element:      c.Element,
		Rarity:       c.Rarity,
		Symbol:       c.Symbol,
		Description:  c.Description,
		Quotes: CardQuotes{
			Login:           c.Login,
			Friendship:      c.Friendship,
			Meet:            c.Meet,
			BattleStart:     c.BattleStart,
			BattleEnd:       c.BattleEnd,
			FriendshipMax:   c.FriendshipMax,
			FriendshipEvent: c.FriendshipEvent,
			Rebirth:         c.Rebirth,
		},
		TurnoverFrom: c.TurnOverFrom,
		TurnoverTo:   c.TurnOverTo,
		Availability: c.Availability,
		// dynamic items
		Evolutions:         make([]EvolutionDetails, 0),
		AwakeningMaterials: make([]AwakenMaterial, 0),
		RebirthMaterials:   make([]RebirthMaterial, 0),
		Amalgamations:      make([]Amalgamation, 0),
	}
	if c.Cost0 != "" {
		evo := EvolutionDetails{
			EvolutionKey: "0",
			Cost:         getInt(c.Cost0),
			MaxLevel:     getInt(c.MaxLevel0),
			Attack:       c.Atk0,
			Defense:      c.Def0,
			Soldiers:     c.Soldiers0,
			Medals:       getInt(c.Medals0),
			Gold:         getInt(c.Gold0),
			Skills:       make([]CardSkill, 0),
		}
		skill := CardSkill{
			EvoID:        "",
			IDMod:        "",
			Name:         c.Skill,
			Activations:  getInt(c.SkillProcs),
			MinEffect:    c.SkillLv1,
			MaxEffect:    c.SkillLv10,
			CostLv1:      getInt(c.SkillLv1Cost),
			CostLv10:     getInt(c.SkillLv10Cost),
			RandomSkills: []string{c.SkillRandom1, c.SkillRandom2, c.SkillRandom3, c.SkillRandom4, c.SkillRandom5},
		}
		evo.Skills = append(evo.Skills, skill)
		if c.Skill2 != "" {
			skill := CardSkill{
				EvoID:        "",
				IDMod:        "2",
				Name:         c.Skill2,
				Activations:  getInt(c.Skill2Procs),
				MinEffect:    c.Skill2Lv1,
				MaxEffect:    c.Skill2Lv10,
				CostLv1:      getInt(c.Skill2Lv1Cost),
				CostLv10:     getInt(c.Skill2Lv10Cost),
				RandomSkills: []string{c.Skill2Random1, c.Skill2Random2, c.Skill2Random3, c.Skill2Random4, c.Skill2Random5},
				Expiration:   getTime(c.Skill2End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.Skill3 != "" {
			skill := CardSkill{
				EvoID:        "",
				IDMod:        "3",
				Name:         c.Skill3,
				Activations:  getInt(c.Skill3Procs),
				MinEffect:    c.Skill3Lv1,
				MaxEffect:    c.Skill3Lv10,
				CostLv1:      getInt(c.Skill3Lv1Cost),
				CostLv10:     getInt(c.Skill3Lv10Cost),
				RandomSkills: []string{c.Skill3Random1, c.Skill3Random2, c.Skill3Random3, c.Skill3Random4, c.Skill3Random5},
				Expiration:   getTime(c.Skill3End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.SkillT != "" {
			skill := CardSkill{
				EvoID:        "",
				IDMod:        "t",
				Name:         c.SkillT,
				Activations:  getInt(c.SkillTProcs),
				MinEffect:    c.SkillTLv1,
				MaxEffect:    c.SkillTLv10,
				RandomSkills: []string{},
				Expiration:   getTime(c.SkillTEnd),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.Cost1 != "" {
		evo := EvolutionDetails{
			EvolutionKey: "1",
			Cost:         getInt(c.Cost1),
			MaxLevel:     getInt(c.MaxLevel1),
			Attack:       c.Atk1,
			Defense:      c.Def1,
			Soldiers:     c.Soldiers1,
			Medals:       getInt(c.Medals1),
			Gold:         getInt(c.Gold1),
			Skills:       make([]CardSkill, 0),
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.Cost2 != "" {
		evo := EvolutionDetails{
			EvolutionKey: "2",
			Cost:         getInt(c.Cost2),
			MaxLevel:     getInt(c.MaxLevel2),
			Attack:       c.Atk2,
			Defense:      c.Def2,
			Soldiers:     c.Soldiers2,
			Medals:       getInt(c.Medals2),
			Gold:         getInt(c.Gold2),
			Skills:       make([]CardSkill, 0),
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.Cost3 != "" {
		evo := EvolutionDetails{
			EvolutionKey: "3",
			Cost:         getInt(c.Cost3),
			MaxLevel:     getInt(c.MaxLevel3),
			Attack:       c.Atk3,
			Defense:      c.Def3,
			Soldiers:     c.Soldiers3,
			Medals:       getInt(c.Medals3),
			Gold:         getInt(c.Gold3),
			Skills:       make([]CardSkill, 0),
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.Cost4 != "" {
		evo := EvolutionDetails{
			EvolutionKey: "4",
			Cost:         getInt(c.Cost4),
			MaxLevel:     getInt(c.MaxLevel4),
			Attack:       c.Atk4,
			Defense:      c.Def4,
			Soldiers:     c.Soldiers4,
			Medals:       getInt(c.Medals4),
			Gold:         getInt(c.Gold4),
			Skills:       make([]CardSkill, 0),
		}
		if c.Cost0 == "" { // no initial evo
			skill := CardSkill{
				EvoID:        "",
				IDMod:        "",
				Name:         c.Skill,
				Activations:  getInt(c.SkillProcs),
				MinEffect:    c.SkillLv1,
				MaxEffect:    c.SkillLv10,
				CostLv1:      getInt(c.SkillLv1Cost),
				CostLv10:     getInt(c.SkillLv10Cost),
				RandomSkills: []string{c.SkillRandom1, c.SkillRandom2, c.SkillRandom3, c.SkillRandom4, c.SkillRandom5},
			}
			evo.Skills = append(evo.Skills, skill)
			if c.Skill2 != "" {
				skill := CardSkill{
					EvoID:        "",
					IDMod:        "2",
					Name:         c.Skill2,
					Activations:  getInt(c.Skill2Procs),
					MinEffect:    c.Skill2Lv1,
					MaxEffect:    c.Skill2Lv10,
					CostLv1:      getInt(c.Skill2Lv1Cost),
					CostLv10:     getInt(c.Skill2Lv10Cost),
					RandomSkills: []string{c.Skill2Random1, c.Skill2Random2, c.Skill2Random3, c.Skill2Random4, c.Skill2Random5},
					Expiration:   getTime(c.Skill2End),
				}
				evo.Skills = append(evo.Skills, skill)
			}
			if c.Skill3 != "" {
				skill := CardSkill{
					EvoID:        "",
					IDMod:        "3",
					Name:         c.Skill3,
					Activations:  getInt(c.Skill3Procs),
					MinEffect:    c.Skill3Lv1,
					MaxEffect:    c.Skill3Lv10,
					CostLv1:      getInt(c.Skill3Lv1Cost),
					CostLv10:     getInt(c.Skill3Lv10Cost),
					RandomSkills: []string{c.Skill3Random1, c.Skill3Random2, c.Skill3Random3, c.Skill3Random4, c.Skill3Random5},
					Expiration:   getTime(c.Skill3End),
				}
				evo.Skills = append(evo.Skills, skill)
			}
			if c.SkillT != "" {
				skill := CardSkill{
					EvoID:        "",
					IDMod:        "t",
					Name:         c.SkillT,
					Activations:  getInt(c.SkillTProcs),
					MinEffect:    c.SkillTLv1,
					MaxEffect:    c.SkillTLv10,
					RandomSkills: []string{},
					Expiration:   getTime(c.SkillTEnd),
				}
				evo.Skills = append(evo.Skills, skill)
			}
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.CostA != "" {
		evo := EvolutionDetails{
			EvolutionKey: "a",
			Cost:         getInt(c.CostA),
			MaxLevel:     getInt(c.MaxLevelA),
			Attack:       c.AtkA,
			Defense:      c.DefA,
			Soldiers:     c.SoldiersA,
			Medals:       getInt(c.MedalsA),
			Gold:         getInt(c.GoldA),
			Skills:       make([]CardSkill, 0),
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.CostG != "" {
		evo := EvolutionDetails{
			EvolutionKey: "g",
			Cost:         getInt(c.CostG),
			MaxLevel:     getInt(c.MaxLevelG),
			Attack:       c.AtkG,
			Defense:      c.DefG,
			Soldiers:     c.SoldiersG,
			Medals:       getInt(c.MedalsG),
			Gold:         getInt(c.GoldG),
			Skills:       make([]CardSkill, 0),
		}
		skill := CardSkill{
			EvoID:        "g",
			IDMod:        "",
			Name:         c.SkillG,
			Activations:  getInt(c.SkillGProcs),
			MinEffect:    c.SkillGLv1,
			MaxEffect:    c.SkillGLv10,
			CostLv1:      getInt(c.SkillGLv1Cost),
			CostLv10:     getInt(c.SkillGLv10Cost),
			RandomSkills: []string{c.SkillGRandom1, c.SkillGRandom2, c.SkillGRandom3, c.SkillGRandom4, c.SkillGRandom5},
		}
		evo.Skills = append(evo.Skills, skill)
		if c.Skill2 != "" {
			skill := CardSkill{
				EvoID:        "g",
				IDMod:        "2",
				Name:         c.SkillG2,
				Activations:  getInt(c.SkillG2Procs),
				MinEffect:    c.SkillG2Lv1,
				MaxEffect:    c.SkillG2Lv10,
				CostLv1:      getInt(c.SkillG2Lv1Cost),
				CostLv10:     getInt(c.SkillG2Lv10Cost),
				RandomSkills: []string{c.SkillG2Random1, c.SkillG2Random2, c.SkillG2Random3, c.SkillG2Random4, c.SkillG2Random5},
				Expiration:   getTime(c.SkillG2End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.Skill3 != "" {
			skill := CardSkill{
				EvoID:        "g",
				IDMod:        "3",
				Name:         c.SkillG3,
				Activations:  getInt(c.SkillG3Procs),
				MinEffect:    c.SkillG3Lv1,
				MaxEffect:    c.SkillG3Lv10,
				CostLv1:      getInt(c.SkillG3Lv1Cost),
				CostLv10:     getInt(c.SkillG3Lv10Cost),
				RandomSkills: []string{c.SkillG3Random1, c.SkillG3Random2, c.SkillG3Random3, c.SkillG3Random4, c.SkillG3Random5},
				Expiration:   getTime(c.SkillG3End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.SkillT != "" {
			skill := CardSkill{
				EvoID:        "g",
				IDMod:        "t",
				Name:         c.SkillGT,
				Activations:  getInt(c.SkillGTProcs),
				MinEffect:    c.SkillGTLv1,
				MaxEffect:    c.SkillGTLv10,
				RandomSkills: []string{},
				Expiration:   getTime(c.SkillGTEnd),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.CostGA != "" {
		evo := EvolutionDetails{
			EvolutionKey: "ga",
			Cost:         getInt(c.CostGA),
			MaxLevel:     getInt(c.MaxLevelGA),
			Attack:       c.AtkGA,
			Defense:      c.DefGA,
			Soldiers:     c.SoldiersGA,
			Medals:       getInt(c.MedalsGA),
			Gold:         getInt(c.GoldGA),
			Skills:       make([]CardSkill, 0),
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	if c.CostX != "" {
		evo := EvolutionDetails{
			EvolutionKey: "x",
			Cost:         getInt(c.CostX),
			MaxLevel:     getInt(c.MaxLevelX),
			Attack:       c.AtkX,
			Defense:      c.DefX,
			Soldiers:     c.SoldiersX,
			Medals:       getInt(c.MedalsX),
			Gold:         getInt(c.GoldX),
			Skills:       make([]CardSkill, 0),
		}
		skill := CardSkill{
			EvoID:        "x",
			IDMod:        "",
			Name:         c.SkillX,
			Activations:  getInt(c.SkillXProcs),
			MinEffect:    c.SkillXLv1,
			MaxEffect:    c.SkillXLv10,
			CostLv1:      getInt(c.SkillXLv1Cost),
			CostLv10:     getInt(c.SkillXLv10Cost),
			RandomSkills: []string{c.SkillXRandom1, c.SkillXRandom2, c.SkillXRandom3, c.SkillXRandom4, c.SkillXRandom5},
		}
		evo.Skills = append(evo.Skills, skill)
		if c.Skill2 != "" {
			skill := CardSkill{
				EvoID:        "x",
				IDMod:        "2",
				Name:         c.SkillX2,
				Activations:  getInt(c.SkillX2Procs),
				MinEffect:    c.SkillX2Lv1,
				MaxEffect:    c.SkillX2Lv10,
				CostLv1:      getInt(c.SkillX2Lv1Cost),
				CostLv10:     getInt(c.SkillX2Lv10Cost),
				RandomSkills: []string{c.SkillX2Random1, c.SkillX2Random2, c.SkillX2Random3, c.SkillX2Random4, c.SkillX2Random5},
				Expiration:   getTime(c.SkillX2End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.Skill3 != "" {
			skill := CardSkill{
				EvoID:        "x",
				IDMod:        "3",
				Name:         c.SkillX3,
				Activations:  getInt(c.SkillX3Procs),
				MinEffect:    c.SkillX3Lv1,
				MaxEffect:    c.SkillX3Lv10,
				CostLv1:      getInt(c.SkillX3Lv1Cost),
				CostLv10:     getInt(c.SkillX3Lv10Cost),
				RandomSkills: []string{c.SkillX3Random1, c.SkillX3Random2, c.SkillX3Random3, c.SkillX3Random4, c.SkillX3Random5},
				Expiration:   getTime(c.SkillX3End),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		if c.SkillT != "" {
			skill := CardSkill{
				EvoID:        "x",
				IDMod:        "t",
				Name:         c.SkillXT,
				Activations:  getInt(c.SkillXTProcs),
				MinEffect:    c.SkillXTLv1,
				MaxEffect:    c.SkillXTLv10,
				RandomSkills: []string{},
				Expiration:   getTime(c.SkillXTEnd),
			}
			evo.Skills = append(evo.Skills, skill)
		}
		ret.Evolutions = append(ret.Evolutions, evo)
	}
	return
}

//UpdateBaseData Updates the tempalte information from the VC data. Fields updated are `Element`, `Rarity`, `Symbol` and turn over information for evo accidents
func (c *CardFlat) UpdateBaseData(vcCard vc.Card) {
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
	for evoCode, evo := range evolutions {
		if evoCode == "H" {
			continue
		}
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
		if s == nil && ls == nil {
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
			min = s.FireMax() + " / 100% chance"
			max = s.FireMax() + " / 100% chance"
		}
		if ls != nil {
			// thor skills use the "Fire" text
			if mod == "t" {
				max = ls.FireMax() + " / 100% chance"
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
			if s.DefaultCost > 0 {
				rs.FieldByName(skillPrefix + "Lv1Cost").SetString(strconv.Itoa(s.DefaultCost))
			} else {
				rs.FieldByName(skillPrefix + "Lv1Cost").SetString("")
			}
			if s.MaxCost > 0 {
				rs.FieldByName(skillPrefix + "Lv10Cost").SetString(strconv.Itoa(s.MaxCost))
			} else {
				rs.FieldByName(skillPrefix + "Lv10Cost").SetString("")
			}
			rs.FieldByName(skillPrefix + "Procs").SetString(s.ActivationString())

			if s.EffectID == 36 {
				// Random Skill
				for i, v := range []int{s.EffectParam, s.EffectParam2, s.EffectParam3, s.EffectParam4, s.EffectParam5} {
					sr := vc.SkillScan(v)
					if sr != nil {
						rs.FieldByName(skillPrefix + "Random" + strconv.Itoa(i)).SetString(cleanVal(sr.FireMin()))
					}
				}
			}
			if s.Expires() {
				endField := rs.FieldByName(skillPrefix + "End")
				if endField.IsValid() {
					rs.SetString(fmt.Sprintf("%v", s.PublicEndDatetime))
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
			rs.FieldByName("Likeability" + strconv.Itoa(i)).SetString(cleanVal(q))
		}
	}
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
			log.Printf("Found `%s` : `%s`", currentKey, currentVal)
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
