package vc

import (
	"sort"
)

// Dungeon mst_dungeon
type Dungeon struct {
	ID int `json:"_id"`
	// PublicStartDatetime   Timestamp `json:"public_start_datetime"`
	// PublicEndDatetime     Timestamp `json:"public_end_datetime"`
	// RankingStartDatetime  Timestamp `json:"ranking_start_datetime"`
	// RankingEndDatetime    Timestamp `json:"ranking_end_datetime"`
	RankingRewardGroupID  int `json:"ranking_reward_group_id"`
	RankingArrivalGroupID int `json:"arrival_point_reward_group_id"`
	// ElementID             int       `json:"element_id"`
	// LevelSkillGroupID     int       `json:"level_skill_group_id"`
	// URLSchemeID           int       `json:"url_scheme_id"`
	Title string `json:"-"` // MsgDungeonTitle_en.strb
}

// DungeonAreaType mst_dungeon_area_type
type DungeonAreaType struct {
	ID   int    `json:"_id"`
	Name string `json:"-"` // MsgDungeonAreaTypeDesc_en.strb
}

// RankRewards rank rewards for the Dungeon
func (t *Dungeon) RankRewards(v *VFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingRewardGroupID > 0 {
		for _, val := range v.DungeonRewards {
			if val.SheetID == t.RankingRewardGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// ArrivalRewards rewards for arriving at a certain Dungeon floor
func (t *Dungeon) ArrivalRewards(v *VFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingArrivalGroupID > 0 {
		for _, val := range v.DungeonArrivalRewards {
			if val.SheetID == t.RankingArrivalGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// DungeonScan search for a Dungeon by ID
func DungeonScan(id int, v *VFile) *Dungeon {
	if id > 0 {
		l := len(v.Dungeons)
		i := sort.Search(l, func(i int) bool { return v.Dungeons[i].ID >= id })
		if i >= 0 && i < l && v.Dungeons[i].ID == id {
			return &(v.Dungeons[i])
		}
	}
	return nil

}
