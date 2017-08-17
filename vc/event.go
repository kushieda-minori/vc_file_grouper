package vc

import (
	"encoding/json"
)

type Event struct {
	Id                    int             `json:"_id"`
	EventTypeId           int             `json:"event_type_id"`
	BannerId              int             `json:"banner_id"`
	ProgressDisp          int             `json:"progress_disp"`
	CardId                int             `json:"card_id"`
	ItemrId               int             `json:"itemr_id"`
	JumpBanner            int             `json:"jump_banner"`
	JumpButton            int             `json:"jump_button"`
	JumpBannerOther       int             `json:"jump_banner_other"`
	SortOrder             int             `json:"sort_order"`
	TexIdCell             int             `json:"tex_id_cell"`
	TexIdImage            int             `json:"tex_id_image"`
	TexIdImage2           int             `json:"tex_id_image2"`
	MapId                 int             `json:"map_id"`
	KingSeriesId          int             `json:"king_series_id"`
	GuildBattleId         int             `json:"guild_battle_id"`
	KingId1               int             `json:"king_id1"`
	KingId2               int             `json:"king_id2"`
	KingId3               int             `json:"king_id3"`
	CardId1               int             `json:"card_id1"`
	CardId2               int             `json:"card_id2"`
	CardId3               int             `json:"card_id3"`
	CardId4               int             `json:"card_id4"`
	CardId5               int             `json:"card_id5"`
	CardId6               int             `json:"card_id6"`
	CardId7               int             `json:"card_id7"`
	CardId8               int             `json:"card_id8"`
	CardId9               int             `json:"card_id9"`
	CardId10              int             `json:"card_id10"`
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
	CollaboId             int             `json:"collabo_id"`
	MaintenanceTarget     int             `json:"maintenance_target"`
	TowerEventId          int             `json:"tower_event_id"`
	Name                  string          `json:"name"`        // MsgEventName_en.strb
	Description           string          `json:"description"` // MsgEventDesc_en.strb
	_map                  *Map            `json:"-"`
	_archwitches          []Archwitch     `json:"-"`
}

type EventBook struct {
	Id      int `json:"_id"`
	EventId int `json:"event_id"`
}

type EventCard struct {
	Id          int    `json:"_id"`
	EventBookId int    `json:"event_book_id"`
	CardId      int    `json:"card_id"`
	KindId      int    `json:"kind_id"`
	KindName    string `json:"kind_name"` // MsgEventCardKindName_en.strb
}

type RankReward struct {
	Id                       int       `json:"_id"`
	KingListId               int       `json:"king_list_id"` // same as King Series
	SheetId                  int       `json:"sheet_id"`     // maps to the reward sheet below
	GroupId                  int       `json:"group_id"`
	MidSheetId               int       `json:"mid_sheet_id"`
	MidBonusDistributionDate Timestamp `json:"mid_bonus_distribution_date"`
	IndividualPointReward    int       `json:"individual_point_reward"`
}

type RankRewardSheet struct {
	Id          int `json:"_id"`
	SheetId     int `json:"sheet_id"`
	RankFrom    int `json:"rank_from"`
	RankTo      int `json:"rank_to"`
	Cash        int `json:"cash"`
	FriendPoint int `json:"friend_point"`
	Coin        int `json:"coin"`
	Iron        int `json:"iron"`
	Ether       int `json:"ether"`
	Elixir      int `json:"elixir"`
	Exp         int `json:"exp"`
	ItemId      int `json:"item_id"`
	FragmentId  int `json:"fragment_id"`
	CardId      int `json:"card_id"`
	Num         int `json:"num"`
	Point       int `john:"point"`
}

func (e *Event) Map(v *VcFile) *Map {
	if e._map == nil && e.MapId > 0 {
		e._map = MapScan(e.MapId, v.Maps)
	}
	return e._map
}

func (e *Event) Tower(v *VcFile) *Tower {
	if e.TowerEventId <= 0 {
		return nil
	}

	return TowerScan(e.TowerEventId, v)
}

func (e *Event) Thor(v *VcFile) *ThorEvent {
	for k, te := range v.ThorEvents {
		if te.PublicStartDatetime == e.StartDatetime && te.PublicEndDatetime == e.EndDatetime {
			return &(v.ThorEvents[k])
		}
	}
	return nil
}

func (e *Event) Archwitches(v *VcFile) []Archwitch {
	if e.KingSeriesId > 0 {
		if e._archwitches == nil {

			// picks only unique Cards for the event
			set := make(map[int]Archwitch)
			for _, a := range v.Archwitches {
				if e.KingSeriesId == a.KingSeriesId {
					set[a.CardMasterId] = a
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

func (e *Event) RankRewards(v *VcFile) *RankReward {
	if e.KingSeriesId > 0 {
		for k, val := range v.RankRewards {
			if val.KingListId == e.KingSeriesId {
				return &v.RankRewards[k]
			}
		}
	}
	return nil
}

func (r *RankReward) MidRewards(v *VcFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if r.MidSheetId > 0 {
		for _, val := range v.RankRewardSheets {
			if val.SheetId == r.MidSheetId {
				set = append(set, val)
			}
		}
	}
	return set
}

func (r *RankReward) FinalRewards(v *VcFile) []RankRewardSheet {
	set := make([]RankRewardSheet, 0)
	if r.SheetId > 0 {
		for _, val := range v.RankRewardSheets {
			if val.SheetId == r.SheetId {
				set = append(set, val)
			}
		}
	}
	return set
}

func MaxEventId(events []Event) (max int) {
	max = 0
	for _, val := range events {
		if val.Id > max {
			max = val.Id
		}
	}
	return
}

func EventScan(eventId int, events []Event) *Event {
	if eventId > 0 {
		if eventId < len(events) && events[eventId-1].Id == eventId {
			return &events[eventId-1]
		}
		for k, val := range events {
			if val.Id == eventId {
				return &events[k]
			}
		}
	}
	return nil
}

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
