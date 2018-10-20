package vc

// Map information
type Map struct {
	ID                  int       `json:"_id"`
	NameJp              string    `json:"name"`
	Order               int       `json:"order"`
	ImageID             int       `json:"image_id"`
	ReverseFlg          int       `json:"reverse_flg"`
	UnlockID            int       `json:"unlock_id"`
	LastAreaID          int       `json:"last_area_id"`
	AreaAddFlg          int       `json:"area_add_flg"`
	PublicStartDatetime Timestamp `json:"public_start_datetime"`
	PublicEndDatetime   Timestamp `json:"public_end_datetime"`
	EntryImageID        int       `json:"entry_image_id"`
	KingSeriesID        int       `json:"king_series_id"`
	KingID              int       `json:"king_id"`
	ElementalhallID     int       `json:"elementalhall_id"`
	ElementalhallStart  Timestamp `json:"elementalhall_start_datetime"`
	Flags               int       `json:"flags"`
	ForBeginner         int       `json:"for_beginner"`
	NaviID              int       `json:"navi_id"`
	ExchangeItemID      int       `json:"exchange_item_id"`
	Name                string    `json:"name_tl"`   // MsgNPCMapName_en.strb
	StartMsg            string    `json:"start_msg"` // MsgNPCMapStart_en.strb
	areas               []Area
}

// Areas of the map
func (m *Map) Areas(v *VFile) []Area {
	if m.areas == nil {
		m.areas = make([]Area, 0)
		for k, a := range v.Areas {
			if a.MapID == m.ID {
				m.areas = append(m.areas, v.Areas[k])
			}
		}
	}
	return m.areas
}

// Area on a map
type Area struct {
	ID        int    `json:"_id"`
	MapID     int    `json:"map_id"`
	AreaNo    int    `json:"area_no"`
	Name      string `json:"-"` // MsgNPCAreaName_en.strb
	LongName  string `json:"-"` // MsgNPCAreaLongName_en.strb
	Start     string `json:"-"` // MsgNPCAreaStart_en.strb
	End       string `json:"-"` // MsgNPCAreaEnd_en.strb
	Story     string `json:"-"` // MsgNPCAreaStory_en.strb
	BossStart string `json:"-"` // MsgNPCBossStart_en.strb
	BossEnd   string `json:"-"` // MsgNPCBossEnd_en.strb
}

// MapScan Searches for a Map by ID
func MapScan(mapID int, maps []Map) *Map {
	if mapID > 0 {
		if mapID < len(maps) && maps[mapID-1].ID == mapID {
			return &maps[mapID-1]
		}
		for k, val := range maps {
			if val.ID == mapID {
				return &maps[k]
			}
		}
	}
	return nil
}
