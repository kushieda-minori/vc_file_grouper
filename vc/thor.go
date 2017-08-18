package vc

import (
	"sort"
)

// mst_thorhammer
type ThorEvent struct {
	Id                                     int        `json:"_id"`
	PublicStartDatetime                    Timestamp  `json:"public_start_datetime"`
	PublicEndDatetime                      Timestamp  `json:"public_end_datetime"`
	RankingStartDatetime                   Timestamp  `json:"ranking_start_datetime"`
	RankingEndDatetime                     Timestamp  `json:"ranking_end_datetime"`
	RankingRewardGroupId                   int        `json:"ranking_reward_group_id"`
	RankingRewardDestributionStartDatetime Timestamp  `json:"ranking_reward_destribution_start_datetime"`
	PointRewardGroupId                     int        `json:"jump_button"`
	_archwitches                           []ThorKing `json:"-"`
	Title                                  string     `json:""` // MsgThorhammerTitle_en.strb
}

// mst_thorhammer_king
type ThorKing struct {
	Id            int `json:"_id"`
	ThorhammerId  int `json:"thorhammer_id"`
	CostGroupId   int `json:"cost_group_id"`
	CardMasterId1 int `json:"card_master_id_1"`
	CardMasterId2 int `json:"card_master_id_2"`
	StatusGroupId int `json:"status_group_id"`
	PublicFlg     int `json:"public_flg"`
	SkillId1      int `json:"skill_id_1"`
	SkillId2      int `json:"skill_id_2"`
	ServantId1    int `json:"servant_id_1"`
	ServantId2    int `json:"servant_id_2"`
}

// mst_thorhammer_king_cost
type ThorKingCost struct {
	Id           int `json:"_id"`
	GroupId      int `json:"group_id"`
	Cost         int `json:"cost"`
	OffenseRatio int `json:"offense_ratio"`
	DefenseRatio int `json:"defense_ratio"`
}

//mst_thorhammer_ranking_reward and mst_thorhammer_point_reward
type ThorReward struct {
	Id          int `json:"_id"`
	GroupId     int `json:"group_id"`
	RankFrom    int `json:"rank_from"`
	RankTo      int `json:"rank_to"`
	Cash        int `json:"cash"`
	FriendPoint int `json:"friend_point"`
	Coin        int `json:"coin"`
	Iron        int `json:"iron"`
	Ether       int `json:"ether"`
	Exp         int `json:"exp"`
	ItemId      int `json:"item_id"`
	CardId      int `json:"card_id"`
	Num         int `json:"num"`
}

func (e *ThorEvent) Archwitches(v *VcFile) []ThorKing {
	if e._archwitches == nil {
		// picks only unique Cards for the event
		set := make(map[int]ThorKing)
		for _, a := range v.ThorKings {
			if e.Id == a.ThorhammerId {
				set[a.Id] = a
			}
		}

		e._archwitches = make([]ThorKing, 0)
		for _, a := range set {
			e._archwitches = append(e._archwitches, a)
		}
	}
	return e._archwitches
}

func MaxThorEventId(events []ThorEvent) (max int) {
	max = 0
	for _, val := range events {
		if val.Id > max {
			max = val.Id
		}
	}
	return
}

func ThorEventScan(id int, events []ThorEvent) *ThorEvent {
	if id > 0 {
		l := len(events)
		i := sort.Search(l, func(i int) bool { return events[i].Id >= id })
		if i >= 0 && i < l && events[i].Id == id {
			return &(events[i])
		}
	}
	return nil
}
