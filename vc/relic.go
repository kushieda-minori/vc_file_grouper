package vc

// Relic "series" is the list of scared relics and the prizes for completeing them.
type Relic struct {
	ID            int `json:"_id"`
	TreasureID1   int `json:"treasure_id_1"`
	TreasureID2   int `json:"treasure_id_2"`
	TreasureID3   int `json:"treasure_id_3"`
	TreasureID4   int `json:"treasure_id_4"`
	TreasureID5   int `json:"treasure_id_5"`
	TreasureID6   int `json:"treasure_id_6"`
	BonusCardID1  int `json:"bonus_card_id_1"`
	BonusCardID2  int `json:"bonus_card_id_2"`
	BonusCardID3  int `json:"bonus_card_id_3"`
	EventFlg      int `json:"event_flg"`
	AreaAttribute int `json:"area_attribute"`
	PublicFlg     int `json:"public_flg"`
	Order         int `json:"order"`
}
