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
	"runtime/debug"
	"strconv"
	"strings"
)

// Main Structure for the VC data file located in responce/maindata
type VcFile struct {
	cards          []Card
	skills         []Skill
	fusion_list    []Amalgamation
	card_awaken    []CardAwaken
	card_character []CardCharacter
	follower_kinds []FollowerKinds
}

// This reads the maindata file and all associated files for strings
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

	if data[len(data)-1] == 0 {
		os.Stdout.WriteString("File ends with 0\n")
		dataLen--
	} else {
		os.Stdout.WriteString("File did not end with 0\n")
	}

	// decode the main file
	err = json.Unmarshal(data[:dataLen], v)
	if err != nil {
		debug.PrintStack()
		return err
	}

	os.Stdout.WriteString(" Length: " + strconv.Itoa(len(v.cards)) + "\n")

	// card names
	names, err := readStringFile(root + "/string/MsgCardName_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.cards) != len(names) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character Names", len(v.cards), len(names))
	}
	for key, value := range v.cards {
		value.name = names[key]
	}
	names = nil

	description, err := readStringFile(root + "/string/MsgCharaDesc_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(description) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character descriptions", len(v.card_character), len(description))
	}

	friendship, err := readStringFile(root + "/string/MsgCharaFriendship_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(friendship) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship", len(v.card_character), len(friendship))
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
	if len(v.card_character) != len(meet) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character meet", len(v.card_character), len(meet))
	}

	battle_start, err := readStringFile(root + "/string/MsgCharaBtlStart_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(battle_start) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_start", len(v.card_character), len(battle_start))
	}

	battle_end, err := readStringFile(root + "/string/MsgCharaBtlEnd_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(battle_end) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character battle_end", len(v.card_character), len(battle_end))
	}

	friendship_max, err := readStringFile(root + "/string/MsgCharaFriendshipMax_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(friendship_max) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_max", len(v.card_character), len(friendship_max))
	}

	friendship_event, err := readStringFile(root + "/string/MsgCharaBonds_en.strb")
	if err != nil {
		debug.PrintStack()
		return err
	}
	if len(v.card_character) != len(friendship_event) {
		debug.PrintStack()
		return fmt.Errorf("%s did not match data file. master: %d, strings: %d",
			"Character friendship_event", len(v.card_character), len(friendship_event))
	}

	for key, value := range v.card_character {
		value.description = description[key]
		value.friendship = friendship[key]
		if key < len(login) {
			value.login = login[key]
		}
		value.meet = meet[key]
		value.battle_start = battle_start[key]
		value.battle_end = battle_end[key]
		value.friendship_max = friendship_max[key]
		value.friendship_event = friendship_event[key]
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

	for key, value := range v.skills {
		if key < len(names) {
			value.name = names[key]
		}
		if key < len(description) {
			value.description = description[key]
		}
		if key < len(fire) {
			value.fire = fire[key]
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
	ret := make([]string, 10)
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
	ret := strings.Replace(s, "\n", "<br />", -1)
	// skill images
	ret = strings.Replace(ret, "<img=24>", "{{Passion}}", -1)
	ret = strings.Replace(ret, "<img=25>", "{{Cool}}", -1)
	ret = strings.Replace(ret, "<img=26>", "{{Dark}}", -1)
	ret = strings.Replace(ret, "<img=27>", "{{Light}}", -1)

	return ret
}
