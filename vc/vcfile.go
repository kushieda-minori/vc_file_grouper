package vc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("-1"), nil
	}

	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

var location = time.FixedZone("JST", 32400) //time.LoadLocation("Asia/Tokyo")

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	if ts == -1 {
		t.Time = time.Time{}
	} else {
		if location != nil {
			t.Time = time.Unix(int64(ts), 0).In(location)
		} else {
			t.Time = time.Unix(int64(ts), 0)
		}
	}

	return nil
}

// Main Structure for the VC data file located in responce/maindata
type VcFile struct {
	Code   int `json:"code"`
	Common struct {
		UnixTime Timestamp `json:"unixtime"`
	} `json:"common"`
	Defs []struct {
		Id    int `json:"_id"`
		Value int `json:"value"`
	} `json:"defs"`
	DefsTune []struct {
		Id            int       `json:"_id"`
		MstDefsId     int       `json:"mst_defs_id"`
		Value         int       `json:"value"`
		PublicFlg     int       `json:"public_flg"`
		StartDateTime Timestamp `json:"start_datetime"`
		EndDateTime   Timestamp `json:"end_datetime"`
	} `json:"defs_tune"`
	ShortcutUrl          string                `json:"shortcut_url"`
	Version              int                   `json:"version"`
	Cards                []Card                `json:"cards"`
	Skills               []Skill               `json:"skills"`
	Amalgamations        []Amalgamation        `json:"fusion_list"`
	Awakenings           []CardAwaken          `json:"card_awaken"`
	CardCharacter        []CardCharacter       `json:"card_character"`
	FollowerKinds        []FollowerKind        `json:"follower_kinds"`
	Levels               []Level               `json:"levels"`
	CardLevels           []CardLevel           `json:"cardlevel"`
	DeckBonuses          []DeckBonus           `json:"deck_bonus"`
	DeckBonusConditions  []DeckBonusCond       `json:"deck_bonus_cond"`
	Archwitches          []Archwitch           `json:"kings"`
	ArchwitchSeries      []ArchwitchSeries     `json:"king_series"`
	ArchwitchFriendships []ArchwitchFriendship `json:"king_friendship"`
	Events               []Event               `json:"mst_event"`
	EventBooks           []EventBook           `json:"mst_event_book"`
	EventCards           []EventCard           `json:"mst_event_card"`
	RankRewards          []RankReward          `json:"ranking_bonus"`
	RankRewardSheets     []RankRewardSheet     `json:"ranking_bonussheet"`
	Maps                 []Map                 `json:"map"`
	Areas                []Area                `json:"area"`
	Items                []Item                `json:"items"`
	ThorEvents           []ThorEvent           `json:"mst_thorhammer"`
	ThorKings            []ThorKing            `json:"mst_thorhammer_king"`
	ThorKingCosts        []ThorKingCost        `json:"mst_thorhammer_king_cost"`
	ThorRankRewards      []ThorReward          `json:"mst_thorhammer_ranking_reward"`
	ThorPointRewards     []ThorReward          `json:"mst_thorhammer_point_reward"`
}

// This reads the main data file and all associated files for strings
// the data is inserted directly into the struct.
func (v *VcFile) Read(root string) ([]byte, error) {
	filename := root + "/response/master_all"

	var data []byte
	var err error
	if _, err = os.Stat(filename + ".json"); os.IsNotExist(err) {
		_, data, err = DecodeAndSave(filename)
		if err != nil {
			return nil, errors.New("no such file or directory: " + filename)
		}
	} else {
		data, err = ioutil.ReadFile(filename + ".json")
		if err != nil {
			return nil, err
		}
	}

	// decode the main file
	err = json.Unmarshal(data[:], v)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	// card names
	names, err := readStringFile(root + "/string/MsgCardName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.Cards) > len(names) {
		fmt.Fprintln(os.Stdout, "names: %v", names)
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character Names", len(v.Cards), len(names))
	}
	for key, _ := range v.Cards {
		v.Cards[key].Name = strings.Replace(strings.Title(strings.ToLower(names[key])), "'S", "'s", -1)
	}
	names = nil

	description, err := readStringFile(root + "/string/MsgCharaDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(description) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character descriptions", len(v.CardCharacter), len(description))
	}

	friendship, err := readStringFile(root + "/string/MsgCharaFriendship_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(friendship) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship", len(v.CardCharacter), len(friendship))
	}

	login, err := readStringFile(root + "/string/MsgCharaWelcome_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	meet, err := readStringFile(root + "/string/MsgCharaMeet_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(meet) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character meet", len(v.CardCharacter), len(meet))
	}

	battle_start, err := readStringFile(root + "/string/MsgCharaBtlStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(battle_start) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_start", len(v.CardCharacter), len(battle_start))
	}

	battle_end, err := readStringFile(root + "/string/MsgCharaBtlEnd_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(battle_end) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_end", len(v.CardCharacter), len(battle_end))
	}

	friendship_max, err := readStringFile(root + "/string/MsgCharaFriendshipMax_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(friendship_max) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_max", len(v.CardCharacter), len(friendship_max))
	}

	friendship_event, err := readStringFile(root + "/string/MsgCharaBonds_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	if len(v.CardCharacter) > len(friendship_event) {
		debug.PrintStack()
		return nil, fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_event", len(v.CardCharacter), len(friendship_event))
	}

	for key, _ := range v.CardCharacter {
		v.CardCharacter[key].Description = strings.Replace(description[key], "\n", " ", -1)
		v.CardCharacter[key].Friendship = friendship[key]
		if key < len(login) {
			v.CardCharacter[key].Login = login[key]
		}
		v.CardCharacter[key].Meet = meet[key]
		v.CardCharacter[key].BattleStart = battle_start[key]
		v.CardCharacter[key].BattleEnd = battle_end[key]
		v.CardCharacter[key].FriendshipMax = friendship_max[key]
		v.CardCharacter[key].FriendshipEvent = friendship_event[key]
	}
	description = nil
	friendship = nil
	login = nil
	meet = nil
	battle_start = nil
	battle_end = nil
	friendship_max = nil
	friendship_event = nil

	//Read Skill strings
	names, err = readStringFile(root + "/string/MsgSkillName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	description, err = readStringFile(root + "/string/MsgSkillDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	fire, err := readStringFile(root + "/string/MsgSkillFire_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key, _ := range v.Skills {
		if key < len(names) {
			v.Skills[key].Name = filterSkill(names[key])
		}
		if key < len(description) {
			v.Skills[key].Description = filterSkill(description[key])
		}
		if key < len(fire) {
			v.Skills[key].Fire = filterSkill(fire[key])
		}
	}

	// event strings
	evntNames, err := readStringFile(root + "/string/MsgEventName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	evntDescrs, err := readStringFile(root + "/string/MsgEventDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key, _ := range v.Events {
		if key < len(evntNames) {
			v.Events[key].Name = filter(evntNames[v.Events[key].Id-1])
		}
		if key < len(evntDescrs) {
			v.Events[key].Description = filter(filterColors(evntDescrs[v.Events[key].Id-1]))
		}
	}

	// map strings
	mapNames, err := readStringFile(root + "/string/MsgNPCMapName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	mapStart, err := readStringFile(root + "/string/MsgNPCMapStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key, _ := range v.Maps {
		if key < len(mapNames) {
			v.Maps[key].Name = mapNames[key]
		}
		if key < len(mapStart) {
			v.Maps[key].StartMsg = filter(filterColors(mapStart[key]))
		}
	}

	areaName, err := readStringFile(root + "/string/MsgNPCAreaName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaLongName, err := readStringFile(root + "/string/MsgNPCAreaLongName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaStart, err := readStringFile(root + "/string/MsgNPCAreaStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaEnd, err := readStringFile(root + "/string/MsgNPCAreaEnd_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaStory, err := readStringFile(root + "/string/MsgNPCAreaStory_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	bossStart, err := readStringFile(root + "/string/MsgNPCBossEnd_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	bossEnd, err := readStringFile(root + "/string/MsgNPCBossStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key, _ := range v.Areas {
		if key < len(bossStart) {
			v.Areas[key].BossStart = filterColors(bossStart[key])
		}
		if key < len(bossEnd) {
			v.Areas[key].BossEnd = filterColors(bossEnd[key])
		}
		if key < len(areaStart) {
			v.Areas[key].Start = filterColors(areaStart[key])
		}
		if key < len(areaEnd) {
			v.Areas[key].End = filterColors(areaEnd[key])
		}
		if key < len(areaName) {
			v.Areas[key].Name = filterColors(areaName[key])
		}
		if key < len(areaLongName) {
			v.Areas[key].LongName = filterColors(areaLongName[key])
		}
		if key < len(areaStory) {
			v.Areas[key].Story = filterColors(areaStory[key])
		}
	}

	awlikeability, err := readStringFile(root + "/string/MsgKingFriendshipDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	// Archwitch Likeability
	for key, _ := range v.ArchwitchFriendships {
		if key < len(awlikeability) {
			v.ArchwitchFriendships[key].Likability = filter(awlikeability[key])
		}
	}

	kingDescription, err := readStringFile(root + "/string/MsgKingTitle_en.strb")
	// king series descriptions
	for key, _ := range v.ArchwitchSeries {
		if key < len(kingDescription) {
			v.ArchwitchSeries[key].Description = filter(kingDescription[key])
		}
	}

	dbonusName, err := readStringFile(root + "/string/MsgDeckBonusName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	dbonusDesc, err := readStringFile(root + "/string/MsgDeckBonusDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	// Deck Bonuses
	for key, _ := range v.DeckBonuses {
		if key < len(dbonusName) {
			v.DeckBonuses[key].Name = filter(dbonusName[key])
		}
		if key < len(dbonusDesc) {
			v.DeckBonuses[key].Description = filter(dbonusDesc[key])
		}
	}

	//Items
	itemdsc, err := readStringFile(root + "/string/MsgShopItemDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemdscshp, err := readStringFile(root + "/string/MsgShopItemDescInShop_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemdscsub, err := readStringFile(root + "/string/MsgShopItemDescSub_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemname, err := readStringFile(root + "/string/MsgShopItemName_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemuse, err := readStringFile(root + "/string/MsgShopItemUseResult_en.strb")
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key, _ := range v.Items {
		if key < len(itemdsc) {
			v.Items[key].Description = filter(itemdsc[key])
		}
		if key < len(itemdscshp) {
			v.Items[key].DescriptionInShop = filter(itemdscshp[key])
		}
		if key < len(itemdscsub) {
			v.Items[key].DescriptionSub = filter(itemdscsub[key])
		}
		if key < len(itemname) {
			v.Items[key].NameEng = filter(itemname[key])
		}
		if key < len(itemuse) {
			v.Items[key].MsgUse = filter(itemuse[key])
		}
	}

	return data, nil
}

func readStringFile(filename string) ([]string, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		debug.PrintStack()
		return nil, errors.New("no such file or directory: " + filename)
	}
	f, err := os.Open(filename)
	if err != nil {
		debug.PrintStack()
		return nil, errors.New("Error opening: " + filename)
	}
	r := bufio.NewReader(f)

	//skip the 8 byte header
	_, err = r.Discard(8)
	if err != nil {
		debug.PrintStack()
		return nil, errors.New("Error skipping the file header for file " + filename)
	}

	// find the "null" seperator between the binary info and the strings
	null := []byte("null\000")
	var line []byte
	for {
		if line, err = r.ReadBytes('\000'); err != nil {
			debug.PrintStack()
			return nil, errors.New("Error reading the file " + filename)
		}
		if bytes.Equal(line, null) {
			break
		}
	}

	//read the strings
	ret := make([]string, 0)
	for {
		if line, err = r.ReadBytes('\000'); err == io.EOF {
			break
		}
		if err != nil {
			debug.PrintStack()
			return nil, errors.New("Error reading the file " + filename)
		}
		// remove the null terminator
		ret = append(ret, filter(string(line[:len(line)-1])))
	}
	return ret, nil
}

//User this to do common string replacements in the VC data files
func filter(s string) string {
	if s == "null" {
		return ""
	}
	ret := strings.TrimSpace(s)
	// standardize utf enocoded symbols
	ret = strings.Replace(ret, "％", "%", -1)
	ret = strings.Replace(ret, "　", " ", -1)
	ret = strings.Replace(ret, "／", "/", -1)
	ret = strings.Replace(ret, "＞", ">", -1)
	ret = strings.Replace(ret, "・", " • ", -1)
	// game controls that aren't needed for wikia
	ret = strings.Replace(ret, "<i><break>", "\n", -1)
	// remove duplicate newlines
	for strings.Contains(ret, "\n\n") {
		ret = strings.Replace(ret, "\n\n", "\n", -1)
	}
	//remove duplicate spaces
	for strings.Contains(ret, "  ") {
		ret = strings.Replace(ret, "  ", " ", -1)
	}
	//ret = strings.Replace(ret, "\n", "<br />", -1)

	ret = strings.Replace(ret, "<img=5>", "[[File:Gemstone.png]]", -1)

	return ret
}

var regexpSlash = regexp.MustCompile("\\s*[/]\\s*")

func filterSkill(s string) string {
	ret := strings.TrimSpace(s)
	//element icons
	ret = strings.Replace(ret, "<img=24>", "{{Passion}}", -1)
	ret = strings.Replace(ret, "<img=25>", "{{Cool}}", -1)
	ret = strings.Replace(ret, "<img=26>", "{{Dark}}", -1)
	ret = strings.Replace(ret, "<img=27>", "{{Light}}", -1)
	//atk def icons
	ret = strings.Replace(ret, "<img=48>", "{{Atk}}", -1)
	ret = strings.Replace(ret, "<img=51>", "{{Atkdef}}", -1)

	// clean up '/' spacing
	ret = regexpSlash.ReplaceAllString(ret, " / ")
	// make counter attack consistent
	ret = strings.Replace(ret, "% Counter", "%\nCounter", -1)
	ret = strings.Replace(ret, "%, Counter", "%\nCounter", -1)
	return ret
}

func filterColors(s string) string {
	ret := strings.TrimSpace(s)
	rc, _ := regexp.Compile("<col=(.+?)>\\n*")
	ret = rc.ReplaceAllString(ret, "<span class=\"vc_color$1\">")

	rc, _ = regexp.Compile("<colrgb=(.+?)>\\n*")
	ret = rc.ReplaceAllString(ret, "<span style=\"color:rgb($1);\">")

	ret = strings.Replace(ret, "</col>", "</span>", -1)

	// strip all size commands out
	rs, _ := regexp.Compile("<(/?)size(=.+?)?>")
	ret = rs.ReplaceAllLiteralString(ret, "")
	return ret
}
