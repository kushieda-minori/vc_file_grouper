package vc

import ()

// mst_tower
type Tower struct {
	Id                    int       `json:"_id"`
	PublicStartDatetime   Timestamp `json:"public_start_datetime"`
	PublicEndDatetime     Timestamp `json:"public_end_datetime"`
	RankingStartDatetime  Timestamp `json:"ranking_start_datetime"`
	RankingEndDatetime    Timestamp `json:"ranking_end_datetime"`
	RankingRewardGroupId  int       `json:"ranking_reward_group_id"`
	RankingArrivalGroupId int       `json:"arrival_point_reward_group_id"`
	ElementId             int       `json:"element_id"`
	LevelSkillGroupId     int       `json:"level_skill_group_id"`
	UrlSchemeId           int       `json:"url_scheme_id"`
	Title                 string    `json:"-"` // MsgTowerTitle_en.strb
}

func (t *Tower) RankRewards(v *VcFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingRewardGroupId > 0 {
		for _, val := range v.TowerReward {
			if val.SheetId == t.RankingRewardGroupId {
				set = append(set, val)
			}
		}
	}
	return set
}

func (t *Tower) ArrivalRewards(v *VcFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingArrivalGroupId > 0 {
		for _, val := range v.TowerArrivalReward {
			if val.SheetId == t.RankingArrivalGroupId {
				set = append(set, val)
			}
		}
	}
	return set
}

func TowerScan(towerId int, v *VcFile) *Tower {
	if towerId > 0 {
		for k, val := range v.Tower {
			if val.Id == towerId {
				return &(v.Tower[k])
			}
		}
	}
	return nil

}
