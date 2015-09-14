package vc_grouper

//"kings" is the list of Archwitches
type Archwitch struct {
	Id            int `json:"_id"`
	KingSeriesId  int `json:"king_series_id"`
	CardMasterId  int `json:"card_master_id"`
	StatusGroupId int `json:"status_group_id"`
	PublicFlg     int `json:"public_flg"`
	RareFlg       int `json:"rare_flg"`
	RareIntensity int `json:"rare_intensity"`
	BattleTime    int `json:"battle_time"`
	Exp           int `json:"exp"`
	MaxFriendship int `json:"max_friendship"`
	SkillId1      int `json:"skill_id_1"`
	SkillId2      int `json:"skill_id_2"`
	WeatherId     int `json:"weather_id"`
	ModelName     int `json:"model_name"`
	ChainRatio2   int `json:"chain_ratio_2"`
	ServantId1    int `json:"servant_id_1"`
	ServantId2    int `json:"servant_id_2"`
}

//"king_series" is the list of AW events and the RR prizes
//Indicates what "event" the archwitch was a part of
type ArchwitchSeries struct {
	Id                   int `json:"_id"`
	RewardCardId         int `json:"reward_card_id"`
	PublicStartDatetime  int `json:"public_start_datetime"`
	PublicEndDatetime    int `json:"public_end_datetime"`
	ReceiveLimitDatetime int `json:"receive_limit_datetime"`
	IsBeginnerKing       int `json:"is_beginner_king"`
}

//"king_friendship" is the chance of friendship increasing on an AW
type ArchwitchFriendship struct {
	Id          int    `json:"_id"`
	KingId      int    `json:"king_id"`
	Friendship  int    `json:"friendship"`
	UpRate      int    `json:"up_rate"`
	Likability0 string `json:"-"`
	Likability1 string `json:"-"`
	Likability2 string `json:"-"`
	Likability3 string `json:"-"`
	Likability4 string `json:"-"`
	Likability5 string `json:"-"`
}
