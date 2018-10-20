package vc

import (
	"sort"
)

// ThorEvent mst_thorhammer
type ThorEvent struct {
	ID                                     int       `json:"_id"`
	PublicStartDatetime                    Timestamp `json:"public_start_datetime"`
	PublicEndDatetime                      Timestamp `json:"public_end_datetime"`
	RankingStartDatetime                   Timestamp `json:"ranking_start_datetime"`
	RankingEndDatetime                     Timestamp `json:"ranking_end_datetime"`
	RankingRewardGroupID                   int       `json:"ranking_reward_group_id"`
	RankingRewardDestributionStartDatetime Timestamp `json:"ranking_reward_destribution_start_datetime"`
	PointRewardGroupID                     int       `json:"jump_button"`
	Title                                  string    `json:""` // MsgThorhammerTitle_en.strb
	_archwitches                           []ThorKing
}

// ThorKing mst_thorhammer_king
type ThorKing struct {
	ID            int `json:"_id"`
	ThorhammerID  int `json:"thorhammer_id"`
	CostGroupID   int `json:"cost_group_id"`
	CardMasterID1 int `json:"card_master_id_1"`
	CardMasterID2 int `json:"card_master_id_2"`
	StatusGroupID int `json:"status_group_id"`
	PublicFlg     int `json:"public_flg"`
	SkillID1      int `json:"skill_id_1"`
	SkillID2      int `json:"skill_id_2"`
	ServantID1    int `json:"servant_id_1"`
	ServantID2    int `json:"servant_id_2"`
}

// ThorKingCost mst_thorhammer_king_cost
type ThorKingCost struct {
	ID           int `json:"_id"`
	GroupID      int `json:"group_id"`
	Cost         int `json:"cost"`
	OffenseRatio int `json:"offense_ratio"`
	DefenseRatio int `json:"defense_ratio"`
}

// ThorReward mst_thorhammer_ranking_reward and mst_thorhammer_point_reward
type ThorReward struct {
	ID          int `json:"_id"`
	GroupID     int `json:"group_id"`
	RankFrom    int `json:"rank_from"`
	RankTo      int `json:"rank_to"`
	Cash        int `json:"cash"`
	FriendPoint int `json:"friend_point"`
	Coin        int `json:"coin"`
	Iron        int `json:"iron"`
	Ether       int `json:"ether"`
	Exp         int `json:"exp"`
	ItemID      int `json:"item_id"`
	CardID      int `json:"card_id"`
	Num         int `json:"num"`
}

// Archwitches for the event
func (e *ThorEvent) Archwitches(v *VFile) []ThorKing {
	if e._archwitches == nil {
		// picks only unique Cards for the event
		set := make(map[int]ThorKing)
		for _, a := range v.ThorKings {
			if e.ID == a.ThorhammerID {
				set[a.ID] = a
			}
		}

		e._archwitches = make([]ThorKing, 0)
		for _, a := range set {
			e._archwitches = append(e._archwitches, a)
		}
	}
	return e._archwitches
}

// MaxThorEventID Max thor event ID
func MaxThorEventID(events []ThorEvent) (max int) {
	max = 0
	for _, val := range events {
		if val.ID > max {
			max = val.ID
		}
	}
	return
}

// ThorEventScan searches for a thor event by ID
func ThorEventScan(id int, events []ThorEvent) *ThorEvent {
	if id > 0 {
		l := len(events)
		i := sort.Search(l, func(i int) bool { return events[i].ID >= id })
		if i >= 0 && i < l && events[i].ID == id {
			return &(events[i])
		}
	}
	return nil
}
