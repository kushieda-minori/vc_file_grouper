package vc

import ()

type Event struct {
	Id                    int       `json:"_id"`
	EventTypeId           int       `json:"event_type_id"`
	BannerId              int       `json:"banner_id"`
	ProgressDisp          int       `json:"progress_disp"`
	CardId                int       `json:"card_id"`
	ItemrId               int       `json:"itemr_id"`
	JumpBanner            int       `json:"jump_banner"`
	JumpButton            int       `json:"jump_button"`
	JumpBannerOther       int       `json:"jump_banner_other"`
	SortOrder             int       `json:"sort_order"`
	TexIdCell             int       `json:"tex_id_cell"`
	TexIdImage            int       `json:"tex_id_image"`
	TexIdImage2           int       `json:"tex_id_image2"`
	MapId                 int       `json:"map_id"`
	KingSeriesId          int       `json:"king_series_id"`
	GuildBattleId         int       `json:"guild_battle_id"`
	KingId1               int       `json:"king_id1"`
	KingId2               int       `json:"king_id2"`
	KingId3               int       `json:"king_id3"`
	CardId1               int       `json:"card_id1"`
	CardId2               int       `json:"card_id2"`
	CardId3               int       `json:"card_id3"`
	CardId4               int       `json:"card_id4"`
	CardId5               int       `json:"card_id5"`
	CardId6               int       `json:"card_id6"`
	CardId7               int       `json:"card_id7"`
	CardId8               int       `json:"card_id8"`
	CardId9               int       `json:"card_id9"`
	CardId10              int       `json:"card_id10"`
	Keyword1              string    `json:"keyword1"`
	Keyword2              string    `json:"keyword2"`
	Keyword3              string    `json:"keyword3"`
	StartDatetime         Timestamp `json:"start_datetime"`
	EndDatetime           Timestamp `json:"end_datetime"`
	GardenBannerViewOrder int       `json:"garden_banner_view_order"`
	Param                 int       `json:"param"`
	Category              int       `json:"category"`
	TickerFlg             int       `json:"ticker_flg"`
	Name                  string    `json:"name"`        // MsgEventName_en.strb
	Description           string    `json:"description"` // MsgEventDesc_en.strb
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
	3:  "?",
	4:  "Sale",
	6:  "Summon",
	7:  "Special Event",
	8:  "Background Present",
	9:  "BINGO",
	10: "Alliance Battle",
	11: "Jewels Special Campain",
	12: "Alliance Duel",
	13: "Alliance Ultimate Battle",
	14: "Beginners Campaign",
}
