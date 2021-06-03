package vc

import (
	"sort"
)

// Tower mst_tower
type Tower struct {
	SubEvent
	ElementID         int    `json:"element_id"`
	LevelSkillGroupID int    `json:"level_skill_group_id"`
	Title             string `json:"-"` // MsgTowerTitle_en.strb
}

// RankRewards rank rewards for the tower
func (t *Tower) RankRewards() []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingRewardGroupID > 0 {
		for _, val := range Data.TowerRewards {
			if val.SheetID == t.RankingRewardGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// ArrivalRewards rewards for arriving at a certain Tower floor
func (t *Tower) ArrivalRewards() []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.ArrivalRewardGroupID > 0 {
		for _, val := range Data.TowerArrivalRewards {
			if val.SheetID == t.ArrivalRewardGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// TowerScan search for a tower by ID
func TowerScan(id int) *Tower {
	if id > 0 {
		l := len(Data.Towers)
		i := sort.Search(l, func(i int) bool { return Data.Towers[i].ID >= id })
		if i >= 0 && i < l && Data.Towers[i].ID == id {
			return &(Data.Towers[i])
		}
	}
	return nil
}

//CleanedEventName event name cleaned for file systems
func (t *Tower) CleanedEventName() string {
	return cleanForFileName(t.EventName())
}

//EventName Name of this event
func (t *Tower) EventName() string {
	if t == nil {
		return ""
	}
	for _, evt := range Data.Events {
		if evt.TowerEventID == t.ID {
			return evt.Name
		}
	}
	return ""
}

//ScenarioHtml ScenarioHtml
func (t *Tower) ScenarioHtml() (string, error) {
	if t == nil {
		return "", nil
	}
	if t.ID == 26 {
		t.ScenarioID = 1
	}
	return t.SubEvent.GetScenarioHtml(t.EventName(), "tower")
}
