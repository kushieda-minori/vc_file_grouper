package vc

import ()

type Map struct {
	Id                  int       `json:"_id"`
	NameJp              string    `json:"name"`
	Order               int       `json:"order"`
	ImageId             int       `json:"image_id"`
	ReverseFlg          int       `json:"reverse_flg"`
	UnlockId            int       `json:"unlock_id"`
	LastAreaId          int       `json:"last_area_id"`
	AreaAddFlg          int       `json:"area_add_flg"`
	PublicStartDatetime Timestamp `json:"public_start_datetime"`
	PublicEndDatetime   Timestamp `json:"public_end_datetime"`
	EntryImageId        int       `json:"entry_image_id"`
	KingSeriesId        int       `json:"king_series_id"`
	KingId              int       `json:"king_id"`
	ElementalhallId     int       `json:"elementalhall_id"`
	Flags               int       `json:"flags"`
	ForBeginner         int       `json:"for_beginner"`
	NaviId              int       `json:"navi_id"`
	Name                string    `json:"name_tl"`   // MsgNPCMapName_en.strb
	StartMsg            string    `json:"start_msg"` // MsgNPCMapStart_en.strb
	areas               []Area    `json:"-"`
}

func (m *Map) Areas(v *VcFile) []Area {
	if m.areas == nil {
		m.areas = make([]Area, 0)
		for k, a := range v.Areas {
			if a.MapId == m.Id {
				m.areas = append(m.areas, v.Areas[k])
			}
		}
	}
	return m.areas
}

type Area struct {
	Id        int    `json:"_id"`
	MapId     int    `json:"map_id"`
	AreaNo    int    `json:"area_no"`
	Name      string `json:""` // MsgNPCAreaName_en.strb
	LongName  string `json:""` // MsgNPCAreaLongName_en.strb
	Start     string `json:""` // MsgNPCAreaStart_en.strb
	End       string `json:""` // MsgNPCAreaEnd_en.strb
	Story     string `json:""` // MsgNPCAreaStory_en.strb
	BossStart string `json:""` // MsgNPCBossStart_en.strb
	BossEnd   string `json:""` // MsgNPCBossEnd_en.strb
}

func MapScan(mapId int, maps []Map) *Map {
	if mapId > 0 {
		if mapId < len(maps) && maps[mapId-1].Id == mapId {
			return &maps[mapId-1]
		}
		for k, val := range maps {
			if val.Id == mapId {
				return &maps[k]
			}
		}
	}
	return nil
}
