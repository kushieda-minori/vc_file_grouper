package wiki

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//Card Main card info
type Card struct {
	Element string `json:"element"`
	Rarity  string `json:"rarity"`
	Symbol  string `json:"symbol"`

	Skill     string `json:"skill"`
	SkillLv1  string `json:"skill lv1"`
	SkillLv10 string `json:"skill lv10"`
	Procs     string `json:"procs"`
	Random1   string `json:"random 1"`
	Random2   string `json:"random 2"`
	Random3   string `json:"random 3"`
	Random4   string `json:"random 4"`
	Random5   string `json:"random 5"`

	Skill2     string `json:"skill 2"`
	Skill2Lv1  string `json:"skill 2 lv1"`
	Skill2Lv10 string `json:"skill 2 lv10"`
	Procs2     string `json:"procs 2"`
	Skill2End  string `json:"skill 2 end"`
	Random21   string `json:"random 2 1"`
	Random22   string `json:"random 2 2"`
	Random23   string `json:"random 2 3"`
	Random24   string `json:"random 2 4"`
	Random25   string `json:"random 2 5"`

	Skill3     string `json:"skill 3"`
	Skill3Lv1  string `json:"skill 3 lv1"`
	Skill3Lv10 string `json:"skill 3 lv10"`
	Procs3     string `json:"procs 3"`
	Skill3End  string `json:"skill 3 end"`
	Random31   string `json:"random 3 1"`
	Random32   string `json:"random 3 2"`
	Random33   string `json:"random 3 3"`
	Random34   string `json:"random 3 4"`
	Random35   string `json:"random 3 5"`

	SkillL     string `json:"skill l"`
	SkillLLv1  string `json:"skill l lv1"`
	SkillLLv10 string `json:"skill l lv10"`
	ProcsL     string `json:"procs l"`
	SkillLEnd  string `json:"skill l end"`

	SkillG     string `json:"skill g"`
	SkillGLv1  string `json:"skill g lv1"`
	SkillGLv10 string `json:"skill g lv10"`
	ProcsG     string `json:"procs g"`
	RandomG1   string `json:"random g 1"`
	RandomG2   string `json:"random g 2"`
	RandomG3   string `json:"random g 3"`
	RandomG4   string `json:"random g 4"`
	RandomG5   string `json:"random g 5"`

	SkillG2     string `json:"skill g2"`
	SkillG2Lv1  string `json:"skill g2 lv1"`
	SkillG2Lv10 string `json:"skill g2 lv10"`
	ProcsG2     string `json:"procs g2"`
	SkillG2End  string `json:"skill g2 end"`
	RandomG21   string `json:"random g2 1"`
	RandomG22   string `json:"random g2 2"`
	RandomG23   string `json:"random g2 3"`
	RandomG24   string `json:"random g2 4"`
	RandomG25   string `json:"random g2 5"`

	SkillG3     string `json:"skill g3"`
	SkillG3Lv1  string `json:"skill g3 lv1"`
	SkillG3Lv10 string `json:"skill g3 lv10"`
	ProcsG3     string `json:"procs g3"`
	SkillG3End  string `json:"skill g3 end"`
	RandomG31   string `json:"random g3 1"`
	RandomG32   string `json:"random g3 2"`
	RandomG33   string `json:"random g3 3"`
	RandomG34   string `json:"random g3 4"`
	RandomG35   string `json:"random g3 5"`

	SkillGL     string `json:"skill gl"`
	SkillGLLv1  string `json:"skill gl lv1"`
	SkillGLLv10 string `json:"skill gl lv10"`
	ProcsGL     string `json:"procs gl"`
	SkillGLEnd  string `json:"skill gl end"`

	SkillX     string `json:"skill x"`
	SkillXLv1  string `json:"skill x lv1"`
	SkillXLv10 string `json:"skill x lv10"`
	ProcsX     string `json:"procs x"`
	RandomX1   string `json:"random x 1"`
	RandomX2   string `json:"random x 2"`
	RandomX3   string `json:"random x 3"`
	RandomX4   string `json:"random x 4"`
	RandomX5   string `json:"random x 5"`

	SkillX2     string `json:"skill x2"`
	SkillX2Lv1  string `json:"skill x2 lv1"`
	SkillX2Lv10 string `json:"skill x2 lv10"`
	ProcsX2     string `json:"procs x2"`
	SkillX2End  string `json:"skill x2 end"`
	RandomX21   string `json:"random x2 1"`
	RandomX22   string `json:"random x2 2"`
	RandomX23   string `json:"random x2 3"`
	RandomX24   string `json:"random x2 4"`
	RandomX25   string `json:"random x2 5"`

	SkillX3     string `json:"skill x3"`
	SkillX3Lv1  string `json:"skill x3 lv1"`
	SkillX3Lv10 string `json:"skill x3 lv10"`
	ProcsX3     string `json:"procs x3"`
	SkillX3End  string `json:"skill x3 end"`
	RandomX31   string `json:"random x3 1"`
	RandomX32   string `json:"random x3 2"`
	RandomX33   string `json:"random x3 3"`
	RandomX34   string `json:"random x3 4"`
	RandomX35   string `json:"random x3 5"`

	SkillXL     string `json:"skill xl"`
	SkillXLLv1  string `json:"skill xl lv1"`
	SkillXLLv10 string `json:"skill xl lv10"`
	ProcsXL     string `json:"procs xl"`
	SkillXLEnd  string `json:"skill xl end"`

	SkillT     string `json:"skill t"`
	SkillTLv1  string `json:"skill t lv1"`
	SkillTLv10 string `json:"skill t lv10"`
	ProcsT     string `json:"procs t"`
	SkillTEnd  string `json:"skill t end"`

	SkillGT     string `json:"skill gt"`
	SkillGTLv1  string `json:"skill gt lv1"`
	SkillGTLv10 string `json:"skill gt lv10"`
	ProcsGT     string `json:"procs gt"`
	SkillGTEnd  string `json:"skill gt end"`

	SkillXT     string `json:"skill xt"`
	SkillXTLv1  string `json:"skill xt lv1"`
	SkillXTLv10 string `json:"skill xt lv10"`
	ProcsXT     string `json:"procs xt"`
	SkillXTEnd  string `json:"skill xt end"`

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

	AwakenChance string `json:"awaken chance"`
	AwakenOrb    string `json:"awaken orb"`
	AwakenL      string `json:"awaken l"`
	AwakenM      string `json:"awaken m"`
	AwakenS      string `json:"awaken s"`

	RebirthChance     string `json:"rebirth chance"`
	RebirthItem1      string `json:"rebirth item 1"`
	RebirthItem1Count string `json:"rebirth item 1 Count"`
	RebirthItem2      string `json:"rebirth item 2"`
	RebirthItem2Count string `json:"rebirth item 2 Count"`
	RebirthItem3      string `json:"rebirth item 3"`
	RebirthItem3Count string `json:"rebirth item 3 Count"`

	TurnOverTo   string `json:"turnoverto"`
	TurnOverFrom string `json:"turnoverfrom"`
	Availability string `json:"availability"`

	unknownFields map[string]string
	pageHeader    string
	pageFooter    string
}

func (c Card) String() (ret string) {
	if c.pageHeader != "" {
		ret += c.pageHeader + "\n\n"
	}

	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(c)
	json.Unmarshal(inrec, &inInterface)

	// begin template
	ret += "{{Card\n"
	// iterate through record fields
	for field, val := range inInterface {
		ret += fmt.Sprintf("|%s = %s\n", field, val)
	}

	if len(c.unknownFields) > 0 {
		ret += "<!-- these fields were known to the bot, but have not been removed -->\n"
		for field, val := range c.unknownFields {
			ret += fmt.Sprintf("|%s = %s\n", field, val)
		}
	}
	// end template
	ret += "}}\n"

	if c.pageFooter != "" {
		ret += c.pageFooter + "\n\n"
	}

	return ret
}

//ParseWikiPage Parses a wiki page into a card. returns `nil` is there is no Card template definition in the page.
func ParseWikiPage(pageText string) (ret *Card, err error) {
	pageText = strings.TrimSpace(pageText)
	// lowercase all the text for comparison reasons.
	pageLower := strings.ToLower(pageText)
	cardIdx := strings.Index(pageLower, "{{card")
	if pageText == "" || cardIdx < 0 {
		return nil, nil
	}
	var pageHeader string
	if cardIdx > 0 {
		pageHeader = pageText[:cardIdx]
	}

	// convert the page Card template to a map
	pageContentMap, cardEndIdx, err := parsePage(pageText[cardIdx:])
	if err != nil {
		return nil, err
	}

	// convert the parsed map to JSON
	datajson, err := json.Marshal(pageContentMap)
	if err != nil {
		return nil, err
	}
	ret = &Card{}
	// unmarshall the JSON into our Card object
	err = json.Unmarshal(datajson, ret)
	if err != nil {
		return nil, err
	}

	for k, v := range pageContentMap {
		if !fieldIsKnown(k) {
			if ret.unknownFields == nil {
				ret.unknownFields = make(map[string]string)
			}
			ret.unknownFields[k] = v
		}
	}

	ret.pageHeader = pageHeader

	if cardEndIdx+cardIdx < len(pageText) {
		ret.pageFooter = pageText[cardEndIdx+cardIdx:]
	}

	return ret, nil
}

func parsePage(pageText string) (map[string]string, int, error) {
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
			currentKey, err := getNextKey(&runes, &rPos)
			if err != nil {
				return nil, 0, err
			}
			if currentKey == "" {
				// positional param
				positionalParamNum++
				currentKey = strconv.Itoa(positionalParamNum)
			}
			currentVal, err := getNextValue(&runes, &rPos)
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

func getNextKey(runes *[]rune, rPos *int) (ret string, err error) {
	start := *rPos
	for ; (*rPos) < len(*runes); (*rPos)++ {
		r := (*runes)[*rPos]
		if r == '=' {
			ret = strings.TrimSpace(string((*runes)[start:(*rPos)]))
			(*rPos)++ // skip past the '='
			return
		}
		if r == '|' || (r == '}' && (*runes)[(*rPos)+1] == '}') {
			// position parameter instead of keyed.
			*rPos = start
			return "", nil
		}
	}
	return "", errors.New("Invalid page format. Unable to locate key separator")
}

func getNextValue(runes *[]rune, rPos *int) (ret string, err error) {
	bracketStack := list.List{}
	start := *rPos
	for ; (*rPos) < len(*runes); (*rPos)++ {
		r := (*runes)[*rPos]
		if bracketStack.Len() == 0 && (r == '|' || (r == '}' && (*runes)[(*rPos)+1] == '}')) {
			ret = strings.TrimSpace(string((*runes)[start:(*rPos)]))
			(*rPos)-- // step back so the main parser sees the separator
			return
		}
		if r == '{' {
			bracketStack.PushFront(r)
		} else if r == '}' {
			if bracketStack.Len() > 0 {
				toRemove := bracketStack.Front()
				bracketStack.Remove(toRemove)
			} else {
				// return "", errors.New("Invalid page format. Unable to locate key separator")
			}
		}
	}
	return "", errors.New("Invalid page format. Unable to locate key separator")
}

func fieldIsKnown(field string) bool {
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
	"procs",
	"random 1",
	"random 2",
	"random 3",
	"random 4",
	"random 5",
	"skill 2",
	"skill 2 lv1",
	"skill 2 lv10",
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
	"procs 3",
	"skill 3 end",
	"random 3 1",
	"random 3 2",
	"random 3 3",
	"random 3 4",
	"random 3 5",
	"skill l",
	"skill l lv1",
	"skill l lv10",
	"procs l",
	"skill l end",
	"skill g",
	"skill g lv1",
	"skill g lv10",
	"procs g",
	"random g 1",
	"random g 2",
	"random g 3",
	"random g 4",
	"random g 5",
	"skill g2",
	"skill g2 lv1",
	"skill g2 lv10",
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
	"procs g3",
	"skill g3 end",
	"random g3 1",
	"random g3 2",
	"random g3 3",
	"random g3 4",
	"random g3 5",
	"skill gl",
	"skill gl lv1",
	"skill gl lv10",
	"procs gl",
	"skill gl end",
	"skill x",
	"skill x lv1",
	"skill x lv10",
	"procs x",
	"random x 1",
	"random x 2",
	"random x 3",
	"random x 4",
	"random x 5",
	"skill x2",
	"skill x2 lv1",
	"skill x2 lv10",
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
	"procs x3",
	"skill x3 end",
	"random x3 1",
	"random x3 2",
	"random x3 3",
	"random x3 4",
	"random x3 5",
	"skill xl",
	"skill xl lv1",
	"skill xl lv10",
	"procs xl",
	"skill xl end",
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
	"login",
	"description",
	"friendship",
	"meet",
	"battle start",
	"battle end",
	"friendship max",
	"friendship event",
	"rebirth",
	"quote misc 1",
	"quote misc 2",
	"quote misc 3",
	"quote misc 4",
	"likeability 0",
	"likeability 1",
	"likeability 2",
	"likeability 3",
	"likeability 4",
	"likeability 5",
	"awaken chance",
	"awaken orb",
	"awaken l",
	"awaken m",
	"awaken s",
	"rebirth chance",
	"rebirth item 1",
	"rebirth item 1 Count",
	"rebirth item 2",
	"rebirth item 2 Count",
	"rebirth item 3",
	"rebirth item 3 Count",
	"turnoverto",
	"turnoverfrom",
	"availability",
}
