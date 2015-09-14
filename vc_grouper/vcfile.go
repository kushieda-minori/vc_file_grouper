package vc_grouper

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
	// "strconv"
	"strings"
)

// Main Structure for the VC data file located in responce/maindata
type VcFile struct {
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
	Archwitches          []Archwitch           `json"kings"`
	ArchwitchSeriess     []ArchwitchSeries     `json"king_series"`
	ArchwitchFriendships []ArchwitchFriendship `json"king_friendship"`
}

// This reads the main data file and all associated files for strings
// the data is inserted directly into the struct.
func (v *VcFile) Read(root string) error {
	filename := root + "/response/master_all.dat"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		debug.PrintStack()
		return errors.New("no such file or directory: " + filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	dataLen := len(data)

	for data[dataLen-1] == 0 {
		dataLen--
	}

	// decode the main file
	err = json.Unmarshal(data[:dataLen], v)
	if err != nil {
		debug.PrintStack()
		return err
	}

	// card names
	names, err := readStringFile(root + "/string/MsgCardName_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.Cards) != len(names) {
		fmt.Fprintln(os.Stdout, "names: %v", names)
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character Names", len(v.Cards), len(names))
	}
	for key, _ := range v.Cards {
		v.Cards[key].Name = strings.Title(strings.ToLower(names[key]))
	}
	names = nil

	description, err := readStringFile(root + "/string/MsgCharaDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(description) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character descriptions", len(v.CardCharacter), len(description))
	}

	friendship, err := readStringFile(root + "/string/MsgCharaFriendship_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(friendship) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship", len(v.CardCharacter), len(friendship))
	}

	login, err := readStringFile(root + "/string/MsgCharaWelcome_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}

	meet, err := readStringFile(root + "/string/MsgCharaMeet_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(meet) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character meet", len(v.CardCharacter), len(meet))
	}

	battle_start, err := readStringFile(root + "/string/MsgCharaBtlStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(battle_start) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_start", len(v.CardCharacter), len(battle_start))
	}

	battle_end, err := readStringFile(root + "/string/MsgCharaBtlEnd_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(battle_end) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_end", len(v.CardCharacter), len(battle_end))
	}

	friendship_max, err := readStringFile(root + "/string/MsgCharaFriendshipMax_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(friendship_max) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_max", len(v.CardCharacter), len(friendship_max))
	}

	friendship_event, err := readStringFile(root + "/string/MsgCharaBonds_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.CardCharacter) != len(friendship_event) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_event", len(v.CardCharacter), len(friendship_event))
	}

	for key, _ := range v.CardCharacter {
		v.CardCharacter[key].Description = strings.Replace(description[key], "<br />", " ", -1)
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
		return err
	}

	description, err = readStringFile(root + "/string/MsgSkillDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}

	fire, err := readStringFile(root + "/string/MsgSkillFire_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}

	for key, _ := range v.Skills {
		if key < len(names) {
			v.Skills[key].Name = names[key]
		}
		if key < len(description) {
			v.Skills[key].Description = filterSkill(description[key])
		}
		if key < len(fire) {
			v.Skills[key].Fire = fire[key]
		}
	}

	return nil
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
	ret := strings.Replace(s, "\n", "<br />", -1)
	ret = strings.Replace(ret, "ÔºÖ", "%", -1)
	ret = strings.Replace(ret, "<i><break>", "<br />", -1)
	ret = strings.Replace(ret, "‚ô™", "♪", -1)
	ret = strings.Replace(ret, "‚Ä¶‚Ä¶", "..... ", -1)
	for strings.Contains(ret, "<br /><br />") {
		ret = strings.Replace(ret, "<br /><br />", "<br />", -1)
	}
	return ret
}

var regExSlash, _ = regexp.Compile("\\s*/\\s*")

func filterSkill(s string) string {
	//element icons
	ret := strings.Replace(s, "<img=24>", "{{Passion}}", -1)
	ret = strings.Replace(ret, "<img=25>", "{{Cool}}", -1)
	ret = strings.Replace(ret, "<img=26>", "{{Dark}}", -1)
	ret = strings.Replace(ret, "<img=27>", "{{Light}}", -1)
	// clean up '/' spacing
	//for regExSlash.MatchString(ret) {
	ret = regExSlash.ReplaceAllString(ret, " / ")
	// }
	return ret
}
