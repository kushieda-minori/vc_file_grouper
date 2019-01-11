package vc

import (
	"sort"
)

// Dungeon mst_dungeon
type Dungeon struct {
	ID                    int       `json:"_id"`
	ElementID             int       `json:"element_id"`
	RankingRewardGroupID  int       `json:"ranking_reward_group_id"`
	RankingArrivalGroupID int       `json:"arrival_point_reward_group_id"`
	ScenarioID            int       `json:"scenario_id"`
	URLSchemeID           int       `json:"url_scheme_id"`
	ExchangeItemID        int       `json:"exchange_item_id"`
	PublicStartDatetime   Timestamp `json:"public_start_datetime"`
	PublicEndDatetime     Timestamp `json:"public_end_datetime"`
	RankingStartDatetime  Timestamp `json:"ranking_start_datetime"`
	RankingEndDatetime    Timestamp `json:"ranking_end_datetime"`
	Title                 string    `json:"-"` // MsgDungeonTitle_en.strb
}

// DungeonAreaType mst_dungeon_area_type
type DungeonAreaType struct {
	ID         int    `json:"_id"`
	AreaTypeID int    `json:"area_type_id"`
	Name       string `json:"-"` // MsgDungeonAreaTypeDesc_en.strb
}

// RankRewards rank rewards for the Dungeon
func (t *Dungeon) RankRewards(v *VFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingRewardGroupID > 0 {
		for _, val := range v.DungeonRewards {
			if val.GroupID == t.RankingRewardGroupID {
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
			if val.GroupID == t.RankingArrivalGroupID {
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
