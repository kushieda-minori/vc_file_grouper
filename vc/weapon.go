package vc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// WeaponEvent mst_weapon_event
type WeaponEvent struct {
	SubEvent
	WeaponID        int       `json:"weapon_id"`
	EventGachaID    int       `json:"eventgacha_id"`
	MaterialRemoval Timestamp `json:"material_removal_date"`
	Title           string    `json:"-"` // MsgWeaponEventTitle_en.strb
}

// Weapon mst_weapon_character
type Weapon struct {
	ID            int      `json:"_id"`
	RarityGroupID int      `json:"rarity_group_id"`
	RankGroupID   int      `json:"rank_group_id"`
	StatusID      int      `json:"status_id"`
	Names         []string `json:"-"` // MsgWeaponName_en.strb
	Descriptions  []string `json:"-"` // MsgWeaponDesc_en.strb
}

// WeaponSkill mst_weapon_skill
type WeaponSkill struct {
	ID          int    `json:"_id"`
	SkillType   int    `json:"skill_type"`
	Level       int    `json:"lv"`
	Value       int    `json:"value"`
	Description string `json:"-"` // MsgWeaponSkillDesc_en.strb
}

// WeaponSkillUnlockRank mst_weapon_skill_unlock_rank
type WeaponSkillUnlockRank struct {
	ID         int `json:"_id"`
	WeaponID   int `json:"weapon_id"`
	UnlockRank int `json:"unlock_rank"`
	SkillType  int `json:"skill_type"`
	SkillLevel int `json:"skill_Level"`
}

// WeaponSkillUnlockRankList list of unlocks
type WeaponSkillUnlockRankList []WeaponSkillUnlockRank

// WeaponRank mst_weapon_rank
type WeaponRank struct {
	ID      int `json:"_id"`
	GroupID int `json:"group_id"`
	Rank    int `json:"rank"`
	NeedExp int `json:"need_exp"`
	Gold    int `json:"coin"`
	Iron    int `json:"iron"`
	Ether   int `json:"ether"`
	Gem     int `json:"elixir"`
}

// WeaponStatus mst_weapon_status
type WeaponStatus struct {
	ID          int `json:"_id"`
	AtkMin      int `json:"offense_min"`
	AtkMax      int `json:"offense_max"`
	DefMin      int `json:"defense_min"`
	DefMax      int `json:"defense_max"`
	SoldiersMin int `json:"followers_min"`
	SoldiersMax int `json:"followers_max"`
}

// WeaponRarity mst_weapon_rarity
type WeaponRarity struct {
	ID         int `json:"_id"`
	GroupID    int `json:"group_id"`
	UnlockRank int `json:"unlock_rank"`
	Rarity     int `json:"rarity"`
}

// WeaponMaterial mst_weapon_material
// Material item that can be applied to a weapon of a certain rank, and how much Exp the item provides
type WeaponMaterial struct {
	ID       int `json:"_id"`
	WeaponID int `json:"weapon_id"`
	Rarity   int `json:"rarity"`
	Exp      int `json:"exp"`
	ItemID   int `json:"item_id"`
}

// WeaponKiller mst_weapon_killer
type WeaponKiller struct {
	ID                  int    `json:"_id"`
	Rarity              int    `json:"rarity"`
	Damage              int    `json:"damage"`
	SpecialSkillGroupID int    `json:"special_skill_group_id"`
	Description         string `json:""` // MsgWeaponKillerDesc_en.strb
}

// SkillUnlocks for the weapon
func (w *Weapon) SkillUnlocks() WeaponSkillUnlockRankList {
	set := make(WeaponSkillUnlockRankList, 0)
	if w != nil {
		for _, val := range Data.WeaponSkillUnlockRanks {
			if val.WeaponID == w.ID {
				set = append(set, val)
			}
		}
	}
	return set
}

// UpgradeMaterials of the weapon
func (w *Weapon) UpgradeMaterials() []WeaponMaterial {
	set := make([]WeaponMaterial, 0)
	if w != nil {
		for _, val := range Data.WeaponMaterials {
			if val.WeaponID == w.ID {
				set = append(set, val)
			}
		}
	}
	return set
}

// Status of the weapon
func (w *Weapon) Status() *WeaponStatus {
	if w != nil {
		for i, val := range Data.WeaponStatuses {
			if val.ID == w.StatusID {
				return &Data.WeaponStatuses[i]
			}
		}
	}
	return nil
}

// StatusDescription Status of the weapon
func (w *Weapon) StatusDescription() string {
	if w == nil {
		return ""
	}

	if w.StatusID <= len(WeaponStatusTypes) {
		return WeaponStatusTypes[w.StatusID]
	}

	return strconv.Itoa(w.StatusID)
}

// Events the weapon has appeared in
func (w *Weapon) Events() []WeaponEvent {
	set := make([]WeaponEvent, 0)
	if w != nil {
		for _, val := range Data.WeaponEvents {
			if val.WeaponID == w.ID {
				set = append(set, val)
			}
		}
	}
	return set
}

// EventNames names/titles of the events the weapon appeared in
func (w *Weapon) EventNames() []string {
	set := make([]string, 0)
	if w != nil {
		for _, event := range w.Events() {
			set = append(set, event.Title)
		}
	}
	return set
}

// Ranks of the weapon
func (w *Weapon) Ranks() []WeaponRank {
	set := make([]WeaponRank, 0)
	if w != nil {
		for _, val := range Data.WeaponRanks {
			if val.GroupID == w.RankGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// MaxRank max rank of the weapon
func (w *Weapon) MaxRank() int {
	return len(w.Ranks())
}

// Rarities rarities of the weapon
func (w *Weapon) Rarities() []WeaponRarity {
	set := make([]WeaponRarity, 0)
	if w != nil {
		for _, val := range Data.WeaponRarities {
			if val.GroupID == w.RarityGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// RarityCount max rank of the weapon
func (w *Weapon) RarityCount() int {
	return len(w.Rarities())
}

// MaxRarity max rank of the weapon
func (w *Weapon) MaxRarity() int {
	return w.Rarities()[len(w.Rarities())-1].Rarity
}

// MaxRarityName Gets the last rarity name
func (w *Weapon) MaxRarityName() string {
	if w == nil {
		return ""
	}
	lNames := len(w.Names)
	if lNames == 0 {
		return "N/A"
	}
	name := w.Names[lNames-1]
	if firstCardWithName(name) != nil {
		return name + " (Weapon)"
	}
	return name
}

// Skill that is unlocked for a weapon's rank
func (w *WeaponSkillUnlockRank) Skill() *WeaponSkill {
	if w != nil {
		for i, val := range Data.WeaponSkills {
			if w.SkillType == val.SkillType && w.SkillLevel == val.Level {
				return &Data.WeaponSkills[i]
			}
		}
	}
	return nil
}

func (ws *WeaponSkill) String() string {
	if ws == nil {
		return ""
	}
	return fmt.Sprintf("%d: %d - %s",
		ws.SkillType,
		ws.Level,
		ws.DescriptionFormatted(),
	)
}

// DescriptionFormatted formats the description
func (ws *WeaponSkill) DescriptionFormatted() string {
	if ws == nil {
		return ""
	}
	return strings.ReplaceAll(ws.Description, "{1:value}", strconv.Itoa(ws.Value))
}

// TypeName formats the description
func (ws *WeaponSkill) TypeName() string {
	if ws == nil {
		return ""
	}
	lwst := len(WeaponSkillTypes)
	if ws.SkillType <= lwst {
		return WeaponSkillTypes[ws.SkillType]
	}
	return strconv.Itoa(ws.SkillType)
}

func (w *WeaponSkillUnlockRank) String() string {
	if w == nil {
		return ""
	}
	return fmt.Sprintf("Rank %d: Opens %s",
		w.UnlockRank,
		w.Skill().String(),
	)
}

// Item gets the item required to upgrade the weapon.
func (wm *WeaponMaterial) Item() *Item {
	if wm == nil {
		return nil
	}
	return ItemScan(wm.ItemID)
}

//ArrivalRewards Arrival rewards for soul weapon events
func (we *WeaponEvent) ArrivalRewards() []RankRewardSheet {
	rewards := make([]RankRewardSheet, 0)
	if we == nil {
		return rewards
	}
	for i, v := range Data.WeaponArrivalRewards {
		if v.GroupID == we.ArrivalRewardGroupID {
			rewards = append(rewards, Data.WeaponArrivalRewards[i])
		}
	}
	return rewards
}

//RankRewards Rank Rewards for soul weapon events
func (we *WeaponEvent) RankRewards() []RankRewardSheet {
	rewards := make([]RankRewardSheet, 0)
	if we == nil {
		return rewards
	}
	for i, v := range Data.WeaponRewards {
		if v.GroupID == we.RankingRewardGroupID {
			rewards = append(rewards, Data.WeaponRewards[i])
		}
	}
	return rewards
}

//CleanedEventName event name cleaned for file systems
func (we *WeaponEvent) CleanedEventName() string {
	return cleanForFileName(we.EventName())
}

//EventName Name of this event
func (we *WeaponEvent) EventName() string {
	if we == nil {
		return ""
	}
	for _, evt := range Data.Events {
		if evt.WeaponEventID == we.ID {
			return evt.Name
		}
	}
	return ""
}

// ScenarioHtml ScenarioHtml
func (we *WeaponEvent) ScenarioHtml() (string, error) {
	if we == nil {
		return "", nil
	}
	return we.SubEvent.GetScenarioHtml(we.EventName(), "weapon_event")
}

// WeaponScan searches for a weapon by ID
func WeaponScan(id int) *Weapon {
	if id > 0 {
		l := len(Data.Weapons)
		i := sort.Search(l, func(i int) bool { return Data.Weapons[i].ID >= id })
		if i >= 0 && i < l && Data.Weapons[i].ID == id {
			return &(Data.Weapons[i])
		}
	}
	return nil
}

// WeaponEventScan searches for a weapon event by ID
func WeaponEventScan(id int) *WeaponEvent {
	if id > 0 {
		l := len(Data.WeaponEvents)
		i := sort.Search(l, func(i int) bool { return Data.WeaponEvents[i].ID >= id })
		if i >= 0 && i < l && Data.WeaponEvents[i].ID == id {
			return &(Data.WeaponEvents[i])
		}
	}
	return nil
}

//Image gets the raw image name
func (w Weapon) Image(rarity int) string {
	return fmt.Sprintf("wp_%05[1]d_%02[2]d", w.ID, rarity)
}

//ImageName gets the nice image name
func (w Weapon) ImageName(rarity int, isThumb bool) string {
	if isThumb {
		return fmt.Sprintf("%s_%02[2]d_icon.png", strings.Replace(w.MaxRarityName(), " (Weapon)", "", -1), rarity)
	}
	return fmt.Sprintf("%s_%02[2]d.png", strings.Replace(w.MaxRarityName(), " (Weapon)", "", -1), rarity)
}

//GetImageData Gets the image data for weapons by rank
func (w Weapon) GetImageData(isThumb bool) map[string][]byte {
	ret := make(map[string][]byte)
	var sdPath string = filepath.Join(FilePath, "weapon", "sd")
	var mdPath string = filepath.Join(FilePath, "weapon", "md")
	var hdPath string = filepath.Join(FilePath, "weapon", "hd")
	var thumbPath string = filepath.Join(FilePath, "weapon", "thumb")

	for i := 1; i <= w.MaxRarity(); i++ {
		fileName := w.Image(i)
		var fullpath string
		var err error
		if isThumb {
			fullpath = filepath.Join(thumbPath, fileName)
			if _, err = os.Stat(fullpath); os.IsNotExist(err) {
				log.Printf("Unable to find Weapon icon %s : %s", w.MaxRarityName(), fileName)
				continue
			}
		} else {
			fullpath = filepath.Join(hdPath, fileName)
			if _, err = os.Stat(fullpath); os.IsNotExist(err) {
				fullpath = filepath.Join(mdPath, fileName)
				if _, err = os.Stat(fullpath); os.IsNotExist(err) {
					fullpath = filepath.Join(sdPath, fileName)
					if _, err = os.Stat(fullpath); os.IsNotExist(err) {
						log.Printf("Unable to find Weapon image %s : %s", w.MaxRarityName(), fileName)
						continue
					}
				}
			}
		}
		// decode the file
		var b []byte
		b, err = Decode(fullpath)
		if err != nil {
			return nil
		}
		imageName := w.ImageName(i, isThumb)
		ret[imageName] = b
	}
	return ret
}

//Copy returns a copy of this list. Useful for local sorting
func (d WeaponSkillUnlockRankList) Copy() WeaponSkillUnlockRankList {
	ret := make(WeaponSkillUnlockRankList, len(d))
	copy(ret, d)
	return ret
}

// WeaponSkillTypes types of weapon skills.
var WeaponSkillTypes = []string{"0", "KO Gauge", "Poison Aid", "Elemental Aid", "Elemental ATK", "Skill Boost", "Burst Chance"}

// WeaponStatusTypes descriptions for weapon statuses
var WeaponStatusTypes = []string{"", "Balanced", "Attack", "Defense", "Soldier"}
