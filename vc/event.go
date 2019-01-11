package vc

import (
	"encoding/json"
	"sort"
)

// Event information
type Event struct {
	ID                    int             `json:"_id"`
	EventTypeID           int             `json:"event_type_id"`
	BannerID              int             `json:"banner_id"`
	ProgressDisp          int             `json:"progress_disp"`
	CardID                int             `json:"card_id"`
	ItemrID               int             `json:"itemr_id"`
	JumpBanner            int             `json:"jump_banner"`
	JumpButton            int             `json:"jump_button"`
	JumpBannerOther       int             `json:"jump_banner_other"`
	TexIDCell             int             `json:"tex_id_cell"`
	TexIDImage            int             `json:"tex_id_image"`
	TexIDImage2           int             `json:"tex_id_image2"`
	MapID                 int             `json:"map_id"`
	KingSeriesID          int             `json:"king_series_id"`
	GuildBattleID         int             `json:"guild_battle_id"`
	TowerEventID          int             `json:"tower_event_id"`
	DungeonEventID        int             `json:"dungeon_event_id"`
	SortOrder             int             `json:"sort_order"`
	KingID1               int             `json:"king_id1"`
	KingID2               int             `json:"king_id2"`
	KingID3               int             `json:"king_id3"`
	CardID1               int             `json:"card_id1"`
	CardID2               int             `json:"card_id2"`
	CardID3               int             `json:"card_id3"`
	CardID4               int             `json:"card_id4"`
	CardID5               int             `json:"card_id5"`
	CardID6               int             `json:"card_id6"`
	CardID7               int             `json:"card_id7"`
	CardID8               int             `json:"card_id8"`
	CardID9               int             `json:"card_id9"`
	CardID10              int             `json:"card_id10"`
	Keyword1              string          `json:"keyword1"`
	Keyword2              string          `json:"keyword2"`
	Keyword3              string          `json:"keyword3"`
	StartDatetime         Timestamp       `json:"start_datetime"`
	EndDatetime           Timestamp       `json:"end_datetime"`
	GardenBannerViewOrder int             `json:"garden_banner_view_order"`
	Param                 int             `json:"param"`
	Category              int             `json:"category"`
	TickerFlg             int             `json:"ticker_flg"`
	TargetLanguages       string          `json:"target_languages"`
	TargetMarkets         json.RawMessage `json:"target_markets"`
	CollaboID             int             `json:"collabo_id"`
	MaintenanceTarget     int             `json:"maintenance_target"`
	Pickup                int             `json:"pickup"`
	Name                  string          `json:"name"`        // MsgEventName_en.strb
	Description           string          `json:"description"` // MsgEventDesc_en.strb
	_map                  *Map
	_archwitches          []Archwitch
}

// EventBook event book
type EventBook struct {
	ID      int `json:"_id"`
	EventID int `json:"event_id"`
}

// EventCard event card
type EventCard struct {
	ID          int    `json:"_id"`
	EventBookID int    `json:"event_book_id"`
	CardID      int    `json:"card_id"`
	KindID      int    `json:"kind_id"`
	KindName    string `json:"kind_name"` // MsgEventCardKindName_en.strb
}

// RankReward rank rewards for an event
type RankReward struct {
	ID                       int       `json:"_id"`
	KingListID               int       `json:"king_list_id"` // same as King Series
	SheetID                  int       `json:"sheet_id"`     // maps to the reward sheet below
	GroupID                  int       `json:"group_id"`
	MidSheetID               int       `json:"mid_sheet_id"`
	MidBonusDistributionDate Timestamp `json:"mid_bonus_distribution_date"`
	IndividualPointReward    int       `json:"individual_point_reward"`
}

// RankRewardSheet details of the rank reward
type RankRewardSheet struct {
	ID          int `json:"_id"`
	SheetID     int `json:"sheet_id"`
	GroupID     int `json:"group_id"`
	RankFrom    int `json:"rank_from"`
	RankTo      int `json:"rank_to"`
	Cash        int `json:"cash"`
	FriendPoint int `json:"friend_point"`
	Coin        int `json:"coin"`
	Iron        int `json:"iron"`
	Ether       int `json:"ether"`
	Elixir      int `json:"elixir"`
	Exp         int `json:"exp"`
	ItemID      int `json:"item_id"`
	FragmentID  int `json:"fragment_id"`
	CardID      int `json:"card_id"`
	Num         int `json:"num"`
	Point       int `john:"point"`
}

// Map for an event if one exists (usually just AW events)
func (e *Event) Map(v *VFile) *Map {
	if e._map == nil && e.MapID > 0 {
		e._map = MapScan(e.MapID, v.Maps)
	}
	return e._map
}

// Tower information for the event if it's a tower event
func (e *Event) Tower(v *VFile) *Tower {
	if e.TowerEventID <= 0 {
		return nil
	}

	return TowerScan(e.TowerEventID, v)
}

// DemonRealm information for the event if it's a Demon Realm Voyage event
func (e *Event) DemonRealm(v *VFile) *Dungeon {
	if e.DungeonEventID <= 0 {
		return nil
	}

	return DungeonScan(e.DungeonEventID, v)
}

// GuildBattle information if it's an Alliance Battle
func (e *Event) GuildBattle(v *VFile) *GuildBattle {
	if e.GuildBattleID <= 0 {
		return nil
	}

	return GuildBattleScan(e.GuildBattleID, v.GuildBattles)
}

// Thor information for Thor events
func (e *Event) Thor(v *VFile) *ThorEvent {
	for k, te := range v.ThorEvents {
		if te.PublicStartDatetime == e.StartDatetime && te.PublicEndDatetime == e.EndDatetime {
			return &(v.ThorEvents[k])
		}
	}
	return nil
}

// Archwitches for this event.
func (e *Event) Archwitches(v *VFile) []Archwitch {
	if e.KingSeriesID > 0 {
		if e._archwitches == nil {

			// picks only unique Cards for the event
			set := make(map[int]Archwitch)
			for _, a := range v.Archwitches {
				if e.KingSeriesID == a.KingSeriesID {
					set[a.CardMasterID] = a
				}
			}

			e._archwitches = make([]Archwitch, 0)
			for _, a := range set {
				e._archwitches = append(e._archwitches, a)
			}
		}
	} else if e._archwitches == nil {
		e._archwitches = make([]Archwitch, 0)
	}
	return e._archwitches
}

// RankRewards for this event
func (e *Event) RankRewards(v *VFile) *RankReward {
	if e.KingSeriesID > 0 {
		for k, val := range v.RankRewards {
			if val.KingListID == e.KingSeriesID {
				return &v.RankRewards[k]
			}
		}
	}
	return nil
}

// MidRewards for this event
func (r *RankReward) MidRewards(v *VFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if r.MidSheetID > 0 {
		for _, val := range v.RankRewardSheets {
			if val.SheetID == r.MidSheetID {
				set = append(set, val)
			}
		}
	}
	return set
}

// FinalRewards for this event
func (r *RankReward) FinalRewards(v *VFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if r.SheetID > 0 {
		for _, val := range v.RankRewardSheets {
			if val.SheetID == r.SheetID {
				set = append(set, val)
			}
		}
	}
	return set
}

// MaxEventID for the events in the list
func MaxEventID(events []Event) (max int) {
	max = 0
	for _, val := range events {
		if val.ID > max {
			max = val.ID
		}
	}
	return
}

// EventScan searches for an event by ID
func EventScan(id int, events []Event) *Event {
	if id > 0 {
		l := len(events)
		i := sort.Search(l, func(i int) bool { return events[i].ID >= id })
		if i >= 0 && i < l && events[i].ID == id {
			return &(events[i])
		}
	}
	return nil
}

// EventType information
var EventType = map[int]string{
	1:  "Archwitch",
	2:  "?",
	3:  "?",
	4:  "Sale",
	5:  "?",
	6:  "Summon",
	7:  "Special Event",
	8:  "Background Present",
	9:  "BINGO",
	10: "Alliance Battle",
	11: "Special Campaign/Abyssal Archwitch",
	12: "Alliance Duel",
	13: "Alliance Ultimate Battle",
	14: "Beginners Campaign",
	15: "?",
	16: "Alliance Bingo Battle",
	17: "?tests?",
	18: "Tower Event",
}
