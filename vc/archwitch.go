package vc

//"kings" is the list of Archwitches
type Archwitch struct {
	Id            int                   `json:"_id"`
	KingSeriesId  int                   `json:"king_series_id"`
	CardMasterId  int                   `json:"card_master_id"`
	StatusGroupId int                   `json:"status_group_id"`
	PublicFlg     int                   `json:"public_flg"`
	RareFlg       int                   `json:"rare_flg"`
	RareIntensity int                   `json:"rare_intensity"`
	BattleTime    int                   `json:"battle_time"`
	Exp           int                   `json:"exp"`
	MaxFriendship int                   `json:"max_friendship"`
	SkillId1      int                   `json:"skill_id_1"`
	SkillId2      int                   `json:"skill_id_2"`
	WeatherId     int                   `json:"weather_id"`
	ModelName     string                `json:"model_name"`
	ChainRatio2   int                   `json:"chain_ratio_2"`
	ServantId1    int                   `json:"servant_id_1"`
	ServantId2    int                   `json:"servant_id_2"`
	likeability   []ArchwitchFriendship `json:-`
}

//"king_series" is the list of AW events and the RR prizes
//Indicates what "event" the archwitch was a part of
type ArchwitchSeries struct {
	Id                   int         `json:"_id"`
	RewardCardId         int         `json:"reward_card_id"`
	PublicStartDatetime  Timestamp   `json:"public_start_datetime"`
	PublicEndDatetime    Timestamp   `json:"public_end_datetime"`
	ReceiveLimitDatetime Timestamp   `json:"receive_limit_datetime"`
	IsBeginnerKing       int         `json:"is_beginner_king"`
	Description          string      `json:-`
	archwitches          []Archwitch `json:-`
}

//"king_friendship" is the chance of friendship increasing on an AW
type ArchwitchFriendship struct {
	Id         int    `json:"_id"`
	KingId     int    `json:"king_id"`
	Friendship int    `json:"friendship"`
	UpRate     int    `json:"up_rate"`
	Likability string `json:"-"` // MsgKingFriendshipDesc_en.strb
}

/*
Limited items and enemies `lmtd`
*/
type Limited struct {
	Id               int       `json:"_id"`
	LmtdTypeId       int       `json:"lmtd_type_id"`
	DayTypeId        int       `json:"day_type_id"`
	StartDatetime    Timestamp `json:"start_datetime"`
	EndDatetime      Timestamp `json:"end_datetime"`
	GroupId          int       `json:"group_id"`
	Enable           int       `json:"enable"`
	JumpBanner       int       `json:"jump_banner"`
	MapId            int       `json:"map_id"`
	QuestId          int       `json:"quest_id"`
	JumpButton       int       `json:"jump_button"`
	TexIdCell        int       `json:"tex_id_cell"`
	TexIdImage       int       `json:"tex_id_image"`
	TexIdImage2      int       `json:"tex_id_image2"`
	EnableDoButton   int       `json:"enable_do_button"`
	EnableRewardDisp int       `json:"enable_reward_disp"`
	RewardCardId     int       `json:"reward_card_id"`
	RewardItemId     int       `json:"reward_item_id"`
}

func (a *Archwitch) Likeability(v *VcFile) []ArchwitchFriendship {
	if a.likeability == nil {
		a.likeability = make([]ArchwitchFriendship, 0)
		for _, af := range v.ArchwitchFriendships {
			if a.Id == af.KingId {
				a.likeability = append(a.likeability, af)
			}
		}
	}
	return a.likeability
}

func (as *ArchwitchSeries) Archwitches(v *VcFile) []Archwitch {
	if as.archwitches == nil {
		as.archwitches = make([]Archwitch, 0)
		for _, a := range v.Archwitches {
			if as.Id == a.KingSeriesId {
				as.archwitches = append(as.archwitches, a)
			}
		}
	}
	return as.archwitches
}
