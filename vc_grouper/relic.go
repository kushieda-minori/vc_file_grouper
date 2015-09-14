package vc_grouper

//"series" is the list of scared relics and the prizes for completeing them.
type Relic struct {
	Id            int `json:"_id"`
	TreasureId1   int `json:"treasure_id_1"`
	TreasureId2   int `json:"treasure_id_2"`
	TreasureId3   int `json:"treasure_id_3"`
	TreasureId4   int `json:"treasure_id_4"`
	TreasureId5   int `json:"treasure_id_5"`
	TreasureId6   int `json:"treasure_id_6"`
	BonusCardId1  int `json:"bonus_card_id_1"`
	BonusCardId2  int `json:"bonus_card_id_2"`
	BonusCardId3  int `json:"bonus_card_id_3"`
	EventFlg      int `json:"event_flg"`
	AreaAttribute int `json:"area_attribute"`
	PublicFlg     int `json:"public_flg"`
	Order         int `json:"order"`
}
