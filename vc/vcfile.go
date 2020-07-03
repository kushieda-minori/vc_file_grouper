package vc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

// FilePath path to the main VC file
var FilePath string

// Data Main data file
var Data *VFile

// MasterDataStr master data as read from the file as a string.
var MasterDataStr string

// LangPack language pack to use
var LangPack string

// Timestamp in the JSON file
type Timestamp struct {
	time.Time
}

//BinImage image information from a .BIN file
type BinImage struct {
	ID   int
	Name string
	Data []byte
}

// ReadMasterData Reads the master data from the file location
func ReadMasterData(file string) error {
	if Data == nil {
		Data = &(VFile{})
	}
	b, err := Read(file)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
	MasterDataStr = string(b)
	return nil
}

// MarshalJSON converts a JSON timestamp to a GO time
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("-1"), nil
	}

	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

var location = time.FixedZone("JST", 32400) //time.LoadLocation("Asia/Tokyo")

// old cards that are not available anymore
var retiredCards = []int{
	61, 62, // bandit
	55, 56, // beastmaster
	102,      // cutthroat
	323, 324, // cyborg
	121, 122, // dancer
	38, 39, // dark knight
	123, 124, // detective
	163, 164, // doll master
	41, 42, // dragon knight
	43,       // dragon slayer
	225, 226, // dragonewt
	119, 120, // druid
	174, 175, // empress
	157, 158, // farmer
	201, 202, // fox spirit
	229, 230, // gnome
	242, 243, // harpy
	86, 87, // hunter
	63,            // idol
	1, 2, 3, 4, 5, // knight
	90, 91, // kung-fu master
	195, 196, // lycaon
	88, 89, // martial artist
	321, 322, // mechanic
	348, 349, // mythic knight
	265, 266, // oni
	40,     // paladin
	94, 95, // rune knight
	73, 74, // sage
	84, 85, // strategist
	100, 101, // swordsman
	183, 184, // sylph
	304, 305, // trickster
	133, 134, // vampire hunter
}

// new cards that are named the same as an old card that is still active
var newCards = []int{
	4721, 4722, 4734, // Spinner
	4752, 4753, 4785, // Sparky
	4711, 4712, // Diana
}

// ID is the character ID, not the card ID
var characterNameOverride = map[int]string{
	62:   "Kung-Fu Master",            // Kung Fu Master
	181:  "Ariel (Light)",             // Ariel
	254:  "Duckling Look-a-Like",      //Duckling Look-A-Like
	352:  "Ariel (Dark)",              // Ariel
	452:  "Joker",                     // Joker
	465:  "Joker (Cane)",              // Joker
	466:  "Joker (Sickle)",            // Joker
	495:  "Snowman MK II",             // Snowman MKⅡ
	1319: "Jack-O'-Sisters",           // Jack-o'-Sisters (older card with newer name format)
	1173: "Al-mi'raj",                 // Al-Mi'Raj
	1536: "Valiant Bellona (Bronze)",  // Valiant Bellona
	1537: "Valiant Bellona (Silver)",  // Valiant Bellona
	1538: "Valiant Bellona (Gold)",    // Valiant Bellona
	1816: "Gold Girl (SR)",            // Gold Girl
	1846: "Medal Girl (SR)",           // Medal Girl
	1869: "Super Chimry (Passion)",    // Super Chimry
	1870: "Super Chimry (Cool)",       // Super Chimry
	1871: "Super Chimry (Light)",      // Super Chimry
	1872: "Super Chimry (Dark)",       // Super Chimry
	1874: "Hyper Chimry (Passion)",    // Hyper Chimry
	1875: "Hyper Chimry (Cool)",       // Hyper Chimry
	1876: "Hyper Chimry (Light)",      // Hyper Chimry
	1877: "Hyper Chimry (Dark)",       // Hyper Chimry
	2024: "Playful Hades (Red)",       // Playful Hades
	2025: "Playful Hades (Green)",     // Playful Hades
	2026: "Playful Hades (Blue)",      // Playful Hades
	2080: "Eunice",                    // X Eunice = Yunice? Why Mynet...
	2397: "PM Demise",                 // Pm Demise
	2479: "DIY Ninja",                 // Diy Ninja
	2549: "Thunder Stone Shard (L)",   // Thunderstone Shard (L)
	2550: "Thunder Stone Shard (D)",   // Thunderstone Shard (D)
	2554: "Lightning Stone Shard (L)", // Lightning Shard (L)
	2555: "Lightning Stone Shard (D)", // Lightning Shard (D)
	2978: "Etna & Flonne",             // fix spacing
	3060: "Xmas Kyulu Kyula",          // remove the slash
	3167: "Kiyohime (collab)",         // new Kiyo from a collab
}

func init() {
	sort.Ints(retiredCards)
	sort.Ints(newCards)
}

// UnmarshalJSON converts a GO time to a JSON timestamp
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

// VFile Main Structure for the VC data file located in responce/maindata
type VFile struct {
	Code   int `json:"code"`
	Common struct {
		UnixTime Timestamp `json:"unixtime"`
	} `json:"common"`
	Defs []struct {
		ID    int `json:"_id"`
		Value int `json:"value"`
	} `json:"defs"`
	DefsTune []struct {
		ID            int       `json:"_id"`
		MstDefsID     int       `json:"mst_defs_id"`
		Value         int       `json:"value"`
		PublicFlg     int       `json:"public_flg"`
		StartDateTime Timestamp `json:"start_datetime"`
		EndDateTime   Timestamp `json:"end_datetime"`
	} `json:"defs_tune"`
	ShortcutURL                 string                      `json:"shortcut_url"`
	Version                     int                         `json:"version"`
	URLSchemes                  []URLScheme                 `json:"url_scheme"`
	Cards                       CardList                    `json:"cards"`
	Skills                      []Skill                     `json:"skills"`
	SkillLevels                 []SkillLevel                `json:"skill_level"`
	CustomSkillLevels           []CustomSkillLevel          `json:"custom_skill_level"`
	SkillCostIncrementPatterns  []SkillCostIncrementPattern `json:"skill_cost_increment_pattern"`
	Amalgamations               []Amalgamation              `json:"fusion_list"`
	Awakenings                  []CardAwaken                `json:"card_awaken"`
	Rebirths                    []CardAwaken                `json:"card_super_awaken"`
	CardCharacters              []CardCharacter             `json:"card_character"`
	FollowerKinds               []FollowerKind              `json:"follower_kinds"`
	CardRarities                []CardRarity                `json:"card_rares"`
	CardSpecialComposes         []CardSpecialCompose        `json:"card_special_compose"`
	Levels                      []Level                     `json:"levels"`
	LevelupBonuses              []LevelupBonus              `json:"levelup_bonus"`
	CardLevels                  []CardLevel                 `json:"cardlevel"`
	CardLevelsLR                []CardLevel                 `json:"cardlevel_lr"`
	CardLevelsX                 []CardLevel                 `json:"cardlevel_x"`
	CardLevelsXLR               []CardLevel                 `json:"cardlevel_xlr"`
	LevelLRResources            []LevelResource             `json:"card_compose_resource"`
	LevelXResources             []LevelResource             `json:"card_compose_resource_x"`
	LevelXLRResources           []LevelResource             `json:"card_compose_resource_xlr"`
	DeckBonuses                 []DeckBonus                 `json:"deck_bonus"`
	DeckBonusConditions         []DeckBonusCond             `json:"deck_bonus_cond"`
	Archwitches                 ArchwitchList               `json:"kings"`
	ArchwitchSeries             []ArchwitchSeries           `json:"king_series"`
	ArchwitchFriendships        []ArchwitchFriendship       `json:"king_friendship"`
	Events                      []Event                     `json:"mst_event"`
	EventBooks                  []EventBook                 `json:"mst_event_book"`
	EventCards                  []EventCard                 `json:"mst_event_card"`
	RankRewards                 []RankReward                `json:"ranking_bonus"`
	RankRewardSheets            []RankRewardSheet           `json:"ranking_bonussheet"`
	Maps                        []Map                       `json:"map"`
	Areas                       []Area                      `json:"area"`
	Items                       []Item                      `json:"items"`
	Structures                  []Structure                 `json:"structures"`
	StructureLevels             []StructureLevel            `json:"structure_level"`
	StructureNumCosts           []StructureCost             `json:"structure_num_cost"`
	ResourceLevels              []ResourceLevel             `json:"resource"`
	BankLevels                  []BankLevel                 `json:"bank_level"`
	CastleLevels                []CastleLevel               `json:"castle_level"`
	SpecialEffects              []SpecialEffect             `json:"special_effect"`
	ThorEvents                  []ThorEvent                 `json:"mst_thorhammer"`
	ThorKings                   []ThorKing                  `json:"mst_thorhammer_king"`
	ThorKingCosts               []ThorKingCost              `json:"mst_thorhammer_king_cost"`
	ThorRankRewards             []ThorReward                `json:"mst_thorhammer_ranking_reward"`
	ThorPointRewards            []ThorReward                `json:"mst_thorhammer_point_reward"`
	GuildBattles                []GuildBattle               `json:"mst_guildbattle_schedule"`
	GuildBingoBattles           []GuildBingoBattle          `json:"mst_guildbingo"`
	GuildBingoExchangeRewards   []GuildBingoExchangeReward  `json:"mst_guildbingo_exchange_reward"`
	GuildBingoPointCampaigns    []GuildBingoPointCampaign   `json:"mst_guildbingo_point_campaign"`
	GuildBattleRewardRefs       []GuildBattleRewardRef      `json:"mst_guildbattle_point_reward"`
	GuildBattleIndividualPoints []RankRewardSheet           `json:"mst_guildbattle_point_rewardsheet"`
	GuildBattleRankingRewards   []RankRewardSheet           `json:"mst_guildbattle_individual_ranking_reward"`
	GuildAUBWinRewards          []GuildAUBWinReward         `json:"mst_guildbattle_win_reward"`
	Towers                      []Tower                     `json:"mst_tower"`
	TowerRewards                []RankRewardSheet           `json:"mst_tower_ranking_reward"`
	TowerArrivalRewards         []RankRewardSheet           `json:"mst_tower_arrival_point_reward"`
	Dungeons                    []Dungeon                   `json:"mst_dungeon"`
	DungeonAreaTypes            []DungeonAreaType           `json:"mst_dungeon_area_type"`
	DungeonRewards              []RankRewardSheet           `json:"mst_dungeon_ranking_reward"`
	DungeonArrivalRewards       []RankRewardSheet           `json:"mst_dungeon_arrival_point_reward"`
	Weapons                     []Weapon                    `json:"mst_weapon_character"`
	WeaponEvents                []WeaponEvent               `json:"mst_weapon_event"`
	WeaponKillers               []WeaponKiller              `json:"mst_weapon_killer"`
	WeaponMaterials             []WeaponMaterial            `json:"mst_weapon_material"`
	WeaponRanks                 []WeaponRank                `json:"mst_weapon_rank"`
	WeaponRarities              []WeaponRarity              `json:"mst_weapon_rarity"`
	WeaponSkills                []WeaponSkill               `json:"mst_weapon_skill"`
	WeaponSkillUnlockRanks      []WeaponSkillUnlockRank     `json:"mst_weapon_skill_unlock_rank"`
	WeaponStatuses              []WeaponStatus              `json:"mst_weapon_status"`
	WeaponRewards               []RankRewardSheet           `json:"mst_weapon_ranking_reward"`
	WeaponArrivalRewards        []RankRewardSheet           `json:"mst_weapon_arrival_point_reward"`
}

// Read This reads the main data file and all associated files for strings
// the data is inserted directly into the struct.
func Read(root string) ([]byte, error) {
	filename := filepath.Join(root, "response", "master_all")

	var data []byte
	var err error
	var jsonFileInfo os.FileInfo
	if jsonFileInfo, err = os.Stat(filename + ".json"); os.IsNotExist(err) {
		_, data, err = DecodeAndSave(filename)
		if err != nil {
			return nil, errors.New("no such file or directory: " + filename)
		}
	} else {
		md, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}
		// check the timestamp on the saved file and verify the master data has not been updated
		if jsonFileInfo.ModTime().Unix() >= md.ModTime().Unix() {
			data, err = ioutil.ReadFile(filename + ".json")
			if err != nil {
				return nil, err
			}
		} else {
			_, data, err = DecodeAndSave(filename)
			if err != nil {
				return nil, errors.New("no such file or directory: " + filename)
			}
		}
	}

	// decode the main file
	err = json.Unmarshal(data[:], Data)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	strRoot := filepath.Join(root, "string")
	// card names
	names, err := ReadStringFile(filepath.Join(strRoot, "MsgCardName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key := range Data.Cards {
		card := Data.Cards[key]
		if card.ID <= len(names) {
			card.Name = cleanCardName(names[card.ID-1], card)
		}
	}

	renameSpecialAmalCardsWithDupNames()

	// initialize the evolutions
	for key := range Data.Cards {
		card := Data.Cards[key]
		card.GetEvolutions()
	}

	description, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	friendship, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaFriendship_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	login, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaWelcome_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	meet, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaMeet_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	battleStart, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaBtlStart_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	battleEnd, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaBtlEnd_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	friendshipMax, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaFriendshipMax_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	friendshipEvent, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaBonds_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	rebirthEvent, err := ReadStringFile(filepath.Join(strRoot, "MsgCharaSuperAwaken_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	for key := range Data.CardCharacters {
		chara := &Data.CardCharacters[key]
		if chara.ID <= len(description) {
			chara.Description = strings.ReplaceAll(description[chara.ID-1], "\n", " ")
		}
		if chara.ID <= len(friendship) {
			chara.Friendship = friendship[chara.ID-1]
		}
		if chara.ID <= len(login) {
			chara.Login = login[chara.ID-1]
		}
		if chara.ID <= len(meet) {
			chara.Meet = meet[chara.ID-1]
		}
		if chara.ID <= len(battleStart) {
			chara.BattleStart = battleStart[chara.ID-1]
		}
		if chara.ID <= len(battleEnd) {
			chara.BattleEnd = battleEnd[chara.ID-1]
		}
		if chara.ID <= len(friendshipMax) {
			chara.FriendshipMax = friendshipMax[chara.ID-1]
		}
		if chara.ID <= len(friendshipEvent) {
			chara.FriendshipEvent = friendshipEvent[chara.ID-1]
		}
		if chara.ID <= len(rebirthEvent) {
			chara.Rebirth = rebirthEvent[chara.ID-1]
		}
	}

	//Read Skill strings
	names, err = ReadStringFile(filepath.Join(strRoot, "MsgSkillName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	description, err = ReadStringFile(filepath.Join(strRoot, "MsgSkillDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	fire, err := ReadStringFile(filepath.Join(strRoot, "MsgSkillFire_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	for key := range Data.Skills {
		skill := &Data.Skills[key]
		if skill.ID <= len(names) {
			skill.Name = filterSkill(names[skill.ID-1])
		}
		if skill.ID <= len(description) {
			skill.Description = filterSkill(description[skill.ID-1])
		}
		if skill.ID <= len(fire) {
			skill.Fire = filterSkill(fire[skill.ID-1])
		}
	}

	// event strings
	names, err = ReadStringFile(filepath.Join(strRoot, "MsgEventName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	description, err = ReadStringFile(filepath.Join(strRoot, "MsgEventDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	for key := range Data.Events {
		evnt := &Data.Events[key]
		if evnt.ID <= len(names) {
			evnt.Name = filter(names[evnt.ID-1])
		}
		if evnt.ID <= len(description) {
			evnt.Description = filterElementImages(filter(filterColors(description[evnt.ID-1])))
		}
	}

	// map strings
	mapNames, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCMapName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	mapStart, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCMapStart_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key := range Data.Maps {
		m := &Data.Maps[key]
		if m.ID <= len(mapNames) {
			m.Name = mapNames[m.ID-1]
		}
		if m.ID <= len(mapStart) {
			m.StartMsg = filter(filterColors(mapStart[m.ID-1]))
		}
	}

	areaName, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCAreaName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaLongName, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCAreaLongName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaStart, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCAreaStart_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaEnd, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCAreaEnd_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	areaStory, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCAreaStory_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	bossStart, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCBossEnd_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	bossEnd, err := ReadStringFile(filepath.Join(strRoot, "MsgNPCBossStart_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key := range Data.Areas {
		area := &Data.Areas[key]
		if area.ID <= len(bossStart) {
			area.BossStart = filterColors(bossStart[area.ID-1])
		}
		if area.ID <= len(bossEnd) {
			area.BossEnd = filterColors(bossEnd[area.ID-1])
		}
		if area.ID <= len(areaStart) {
			area.Start = filterColors(areaStart[area.ID-1])
		}
		if area.ID <= len(areaEnd) {
			area.End = filterColors(areaEnd[area.ID-1])
		}
		if area.ID <= len(areaName) {
			area.Name = filterColors(areaName[area.ID-1])
		}
		if area.ID <= len(areaLongName) {
			area.LongName = filterColors(areaLongName[area.ID-1])
		}
		if area.ID <= len(areaStory) {
			area.Story = filterColors(areaStory[area.ID-1])
		}
	}

	awlikeability, err := ReadStringFile(filepath.Join(strRoot, "MsgKingFriendshipDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	// Archwitch Likeability
	for key := range Data.ArchwitchFriendships {
		awf := &Data.ArchwitchFriendships[key]
		if awf.ID <= len(awlikeability) {
			awf.Likability = filter(awlikeability[awf.ID-1])
		}
	}

	kingDescription, err := ReadStringFile(filepath.Join(strRoot, "MsgKingTitle_en.strb"))
	// king series descriptions
	for key := range Data.ArchwitchSeries {
		aws := &Data.ArchwitchSeries[key]
		if aws.ID <= len(kingDescription) {
			aws.Description = filter(kingDescription[aws.ID-1])
		}
	}

	dbonusName, err := ReadStringFile(filepath.Join(strRoot, "MsgDeckBonusName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	dbonusDesc, err := ReadStringFile(filepath.Join(strRoot, "MsgDeckBonusDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	// Deck Bonuses
	for key := range Data.DeckBonuses {
		db := &Data.DeckBonuses[key]
		if db.ID <= len(dbonusName) {
			db.Name = filter(dbonusName[db.ID-1])
		}
		if db.ID <= len(dbonusDesc) {
			db.Description = filter(dbonusDesc[db.ID-1])
		}
	}

	//Items
	itemdsc, err := ReadStringFile(filepath.Join(strRoot, "MsgShopItemDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemdscshp, err := ReadStringFile(filepath.Join(strRoot, "MsgShopItemDescInShop_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemdscsub, err := ReadStringFile(filepath.Join(strRoot, "MsgShopItemDescSub_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemname, err := ReadStringFile(filepath.Join(strRoot, "MsgShopItemName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	itemuse, err := ReadStringFile(filepath.Join(strRoot, "MsgShopItemUseResult_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key := range Data.Items {
		item := &Data.Items[key]
		if item.ID <= len(itemdsc) {
			item.Description = filter(itemdsc[item.ID-1])
		}
		if item.ID <= len(itemdscshp) {
			item.DescriptionInShop = filter(itemdscshp[item.ID-1])
		}
		if item.ID <= len(itemdscsub) {
			item.DescriptionSub = filter(itemdscsub[item.ID-1])
		}
		if item.ID <= len(itemname) {
			item.NameEng = filter(itemname[item.ID-1])
		}
		if item.ID <= len(itemuse) {
			item.MsgUse = filter(itemuse[item.ID-1])
		}
	}

	buildname, err := ReadStringFile(filepath.Join(strRoot, "MsgBuildingName_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	builddesc, err := ReadStringFile(filepath.Join(strRoot, "MsgBuildingDesc_en.strb"))
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	for key := range Data.Structures {
		s := &Data.Structures[key]
		if s.ID <= len(buildname) {
			s.Name = filter(buildname[s.ID-1])
		}
		if key <= len(builddesc) {
			s.Description = filter(builddesc[s.ID-1])
		}
	}

	if Data.ThorEvents != nil {
		thorTitle, err := ReadStringFile(filepath.Join(strRoot, "MsgThorhammerTitle_en.strb"))
		if err != nil {
			debug.PrintStack()
			return data, err
		}
		for key := range Data.ThorEvents {
			te := &Data.ThorEvents[key]
			if te.ID <= len(thorTitle) {
				te.Title = filter(thorTitle[te.ID-1])
			}
		}
	}

	if Data.Weapons != nil {
		weaponName, err := ReadStringFile(filepath.Join(strRoot, "MsgWeaponName_en.strb"))
		if err != nil {
			debug.PrintStack()
			return data, err
		}
		weaponDesc, err := ReadStringFile(filepath.Join(strRoot, "MsgWeaponDesc_en.strb"))
		if err != nil {
			debug.PrintStack()
			return data, err
		}
		lwn := len(weaponName)
		lwd := len(weaponDesc)
		ridx := 0
		for key := range Data.Weapons {
			weap := &Data.Weapons[key]
			lastridx := ridx + weap.MaxRarity()
			for i := ridx; i < lastridx; i++ {
				if i < lwn {
					weap.Names = append(weap.Names, cleanWeaponName(weaponName[i]))
				}
				if i < lwd {
					weap.Descriptions = append(weap.Descriptions, filter(weaponDesc[i]))
				}
			}
			ridx = lastridx
		}
	}

	if Data.WeaponSkills != nil {
		weaponSkill, err := ReadStringFile(filepath.Join(strRoot, "MsgWeaponSkillDesc_en.strb"))
		if err != nil {
			debug.PrintStack()
			return data, err
		}
		for key := range Data.WeaponSkills {
			wskill := &Data.WeaponSkills[key]
			if wskill.ID <= len(weaponSkill) {
				wskill.Description = filter(weaponSkill[wskill.ID-1])
			}
		}
	}

	if Data.WeaponEvents != nil {
		weaponEvent, err := ReadStringFile(filepath.Join(strRoot, "MsgWeaponEventTitle_en.strb"))
		if err != nil {
			debug.PrintStack()
			return data, err
		}
		for key := range Data.WeaponEvents {
			wevent := &Data.WeaponEvents[key]
			if wevent.ID <= len(weaponEvent) {
				wevent.Title = filter(weaponEvent[wevent.ID-1])
			}
		}
	}

	return data, nil
}

//ReadStringFile Reads a binary string file filtering common issues out
func ReadStringFile(fname string) ([]string, error) {
	return ReadStringFileFilter(fname, true)
}

//ReadStringFileFilter Reads a binary string file with the strings optionally filtered
func ReadStringFileFilter(fname string, filtered bool) ([]string, error) {
	filename := strings.Replace(fname, "_en.strb", "_"+LangPack+".strb", 1)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		debug.PrintStack()
		return nil, errors.New("no such file or directory: " + filename)
	}
	f, err := os.Open(filename)
	if err != nil {
		debug.PrintStack()
		return nil, errors.New("Error opening: " + filename)
	}
	defer f.Close()

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
		if filtered {
			ret = append(ret, filter(string(line[:len(line)-1])))
		} else {
			ret = append(ret, string(line[:len(line)-1]))
		}
	}
	return ret, nil
}

//ReadBinFileImages reads a binary file and returns the image data (PNG only)
func ReadBinFileImages(filename string) ([]BinImage, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	l := len(data)

	nameStart := []byte("\x00\x00\x00")
	lnameStart := len(nameStart)
	nameEnd := byte('\000')

	pngStart := []byte("\x89PNG")
	lpngStart := len(pngStart)
	pngEnd := []byte("IEND\xAEB`\x82")
	lpngEnd := len(pngEnd)

	findNameStart := func(data []byte, startIdx int) int {
		for i := startIdx; i < (l - lnameStart); i++ {
			if bytes.Equal(data[i:i+lnameStart], nameStart) {
				if i+lnameStart+1 < l {
					c := data[i+lnameStart]
					if c >= 'a' && c <= 'z' {
						return i + lnameStart // exclude the 3 null bytes
					}
				}
			}
		}
		return -1
	}
	findNameEnd := func(data []byte, startIdx int) int {
		for i := startIdx; i < (l - 1); i++ {
			if data[i] == nameEnd {
				return i
			}
		}
		return -1
	}
	findPngStart := func(data []byte, startIdx int) int {
		for i := startIdx; i < (l - lpngStart); i++ {
			if bytes.Equal(data[i:i+lpngStart], pngStart) {
				return i
			}
		}
		return -1
	}
	findPngEnd := func(data []byte, startIdx int) int {
		for i := startIdx; i < (l - lpngEnd); i++ {
			if bytes.Equal(data[i:i+lpngEnd], pngEnd) {
				return i + lpngEnd
			}
		}
		return -1
	}
	//start of the PNG image
	firstPng := findPngStart(data, 0)
	if firstPng < 0 {
		return nil, errors.New("unable to locate any images")
	}
	//parse names
	start := 0
	names := make([]string, 0)
	for start < firstPng {
		start = findNameStart(data, start)
		if start < 0 || start > firstPng {
			break
		}
		end := findNameEnd(data, start)
		if end < 0 || end > firstPng {
			break
		}
		if end-start < 5 { // none of the names are shorter than 5 characters
			start = end
			continue
		}
		name := string(data[start:end])
		names = append(names, name)
		//log.Printf("found image name '%s', idx: %d-%d\n", name, start, end)
		start = end
	}

	lnames := len(names)
	getImageName := func(idx int) string {
		if idx < lnames && names[idx] != "" {
			return names[idx] + ".png"
		}
		return fmt.Sprintf("structure_%05d.png", idx+1)
	}

	start = firstPng
	// look for PNG images
	ret := make([]BinImage, 0)
	i := 0 // skip the "dummy" name
	for start < (l - (lpngStart + lpngEnd)) {
		i++
		start = findPngStart(data, start)
		if start < 0 {
			break
		}
		end := findPngEnd(data, start)
		if end < 0 {
			return nil, errors.New("unable to locate the end of an image")
		}

		ret = append(ret, BinImage{ID: i, Name: getImageName(i), Data: data[start:end]})
		//log.Printf("found image, idx: %d\n", start)
		start = end
	}
	log.Printf("found %d image names and %d images\n",
		lnames,
		len(ret),
	)

	return ret, nil
}

func cleanCardName(name string, card *Card) string {
	ret := ""
	if newName, ok := characterNameOverride[card.CardCharaID]; ok {
		// use an overridden hard-coded name
		ret = newName
	} else {
		ret = strings.ReplaceAll(strings.Title(strings.ToLower(name)), "'S", "'s")
		ret = strings.ReplaceAll(ret, "(Sr)", "(SR)")
		ret = strings.ReplaceAll(ret, "(Ur)", "(UR)")
		ret = strings.ReplaceAll(ret, "(Lr)", "(LR)")
		if card.CardCharaID < 1450 {
			// use lowecase prepositions and articles as these are cards in the wiki before this program.
			ret = strings.ReplaceAll(ret, " Of ", " of ")
			ret = strings.ReplaceAll(ret, "-Of-", "-of-")
			ret = strings.ReplaceAll(ret, " The ", " the ")
			ret = strings.ReplaceAll(ret, "-The-", "-the-")
			ret = strings.ReplaceAll(ret, " In ", " in ")
			ret = strings.ReplaceAll(ret, "-In-", "-in-")
			ret = strings.ReplaceAll(ret, " O'", " o'")
			ret = strings.ReplaceAll(ret, "-O'", "-o'")
			ret = strings.ReplaceAll(ret, " Du ", " du ") // french "of"
		}
	}
	// old cards
	if card.IsRetired() {
		ret += " (Old)"
	} else {
		// new cards that are named the same as an old card that is still active
		newIDx := sort.SearchInts(newCards, card.ID)
		if newIDx >= 0 && newIDx < len(newCards) && newCards[newIDx] == card.ID {
			ret += " (New)"
		}
	}
	return ret
}

func cleanWeaponName(name string) string {
	return strings.ReplaceAll(strings.Title(strings.ToLower(name)), "'S", "'s")
}

// GetBinFileImages gets a subset of images from the bin index. 1-based index.
func GetBinFileImages(filename string, idxs ...int) ([]BinImage, error) {
	if len(idxs) == 0 {
		return nil, errors.New("Index out of bounds")
	}
	images, err := ReadBinFileImages(filename)
	if err != nil {
		return nil, err
	}
	ret := make([]BinImage, 0, len(idxs))
	for _, idx := range idxs {
		if idx < 1 || idx > len(images) {
			return nil, errors.New("Index out of bounds")
		}
		ret = append(ret, images[idx-1])
	}
	return ret, nil
}

//Use this to do common string replacements in the VC data files
func filter(s string) string {
	if s == "null" {
		return ""
	}
	ret := strings.TrimSpace(s)
	// standardize utf enocoded symbols
	ret = strings.ReplaceAll(ret, "％", "%")
	ret = strings.ReplaceAll(ret, "　", " ")
	ret = strings.ReplaceAll(ret, "／", "/")
	ret = strings.ReplaceAll(ret, "＞", ">")
	ret = strings.ReplaceAll(ret, "・", " • ")
	ret = strings.ReplaceAll(ret, "（", " (")
	ret = strings.ReplaceAll(ret, "）", ") ")
	// game controls that aren't needed for fandom
	ret = strings.ReplaceAll(ret, "<i><break>", " ")
	// remove duplicate newlines
	for strings.Contains(ret, "\n\n") {
		ret = strings.ReplaceAll(ret, "\n\n", "\n")
	}
	//remove duplicate spaces
	for strings.Contains(ret, "  ") {
		ret = strings.ReplaceAll(ret, "  ", " ")
	}
	//ret = strings.ReplaceAll(ret, "\n", "<br />")

	ret = strings.ReplaceAll(ret, "<img=1>Gold", "{{Icon|gold}}")
	ret = strings.ReplaceAll(ret, "<img=4>Iron", "{{Icon|iron}}")
	ret = strings.ReplaceAll(ret, "<img=3>Ether", "{{Icon|ether}}")
	ret = strings.ReplaceAll(ret, "<img=56>Gem", "{{Icon|gem}}")
	ret = strings.ReplaceAll(ret, "<img=1>", "{{Icon|gold}}")
	ret = strings.ReplaceAll(ret, "<img=4>", "{{Icon|iron}}")
	ret = strings.ReplaceAll(ret, "<img=3>", "{{Icon|ether}}")
	ret = strings.ReplaceAll(ret, "<img=56>", "{{Icon|gem}}")
	ret = strings.ReplaceAll(ret, "<img=5>", "{{Icon|jewel}}")

	return strings.TrimSpace(ret)
}

func filterElementImages(s string) string {
	ret := strings.TrimSpace(s)
	//element icons
	ret = strings.ReplaceAll(ret, "<img=24>", "{{Passion}}")
	ret = strings.ReplaceAll(ret, "<img=25>", "{{Cool}}")
	ret = strings.ReplaceAll(ret, "<img=26>", "{{Dark}}")
	ret = strings.ReplaceAll(ret, "<img=27>", "{{Light}}")
	return ret
}

var regexpSlash = regexp.MustCompile("\\s*[/]\\s*")

func filterSkill(s string) string {
	ret := filterElementImages(s)

	//atk def icons
	ret = strings.ReplaceAll(ret, "<img=48>", "{{Atk}}")
	ret = strings.ReplaceAll(ret, "<img=51>", "{{Atkdef}}")

	// clean up '/' spacing
	ret = regexpSlash.ReplaceAllString(ret, " / ")
	// make counter attack consistent
	ret = strings.ReplaceAll(ret, "% Counter", "%\nCounter")
	ret = strings.ReplaceAll(ret, "%, Counter", "%\nCounter")
	return ret
}

func filterColors(s string) string {
	ret := strings.TrimSpace(s)
	rc, _ := regexp.Compile("<col=(.+?)>\\n*")
	ret = rc.ReplaceAllString(ret, "<span class=\"vc_color$1\">")

	rc, _ = regexp.Compile("<colrgb=(.+?)>\\n*")
	ret = rc.ReplaceAllString(ret, "<span style=\"color:rgb($1);\">")

	ret = strings.ReplaceAll(ret, "</col>", "</span>")

	// strip all size commands out
	rs, _ := regexp.Compile("<(/?)size(=.+?)?>")
	ret = rs.ReplaceAllLiteralString(ret, "")
	return ret
}
