package vc

import (
	"sort"
)

// Dungeon mst_dungeon
type Dungeon struct {
	SubEvent
	ElementID      int    `json:"element_id"`
	ExchangeItemID int    `json:"exchange_item_id"`
	Title          string `json:"-"` // MsgDungeonTitle_en.strb
}

// DungeonAreaType mst_dungeon_area_type
type DungeonAreaType struct {
	ID         int    `json:"_id"`
	AreaTypeID int    `json:"area_type_id"`
	Name       string `json:"-"` // MsgDungeonAreaTypeDesc_en.strb
}

// RankRewards rank rewards for the Dungeon
func (t *Dungeon) RankRewards() []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t.RankingRewardGroupID > 0 {
		for _, val := range Data.DungeonRewards {
			if val.GroupID == t.RankingRewardGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// ArrivalRewards rewards for arriving at a certain Dungeon floor
func (t *Dungeon) ArrivalRewards() []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if t == nil {
		return set
	}
	if t.ArrivalRewardGroupID > 0 {
		for _, val := range Data.DungeonArrivalRewards {
			if val.GroupID == t.ArrivalRewardGroupID {
				set = append(set, val)
			}
		}
	}
	return set
}

// DungeonScan search for a Dungeon by ID
func DungeonScan(id int) *Dungeon {
	if id > 0 {
		l := len(Data.Dungeons)
		i := sort.Search(l, func(i int) bool { return Data.Dungeons[i].ID >= id })
		if i >= 0 && i < l && Data.Dungeons[i].ID == id {
			return &(Data.Dungeons[i])
		}
	}
	return nil

}

//EventName Name of this event
func (d *Dungeon) EventName() string {
	if d == nil {
		return ""
	}
	for _, evt := range Data.Events {
		if evt.DungeonEventID == d.ID {
			return evt.Name
		}
	}
	return ""
}

//ScenarioHtml ScenarioHtml
func (d *Dungeon) ScenarioHtml() (string, error) {
	if d == nil {
		return "", nil
	}
	return d.SubEvent.GetScenarioHtml(d.EventName(), "dungeon")
}
