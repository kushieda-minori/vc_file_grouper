package vc_grouper

// Unit Bonuses from master file field "deck_bonus"
// these match with the strings in MsgDeckBonusName_en.strb
// and MsgDeckBonusDesc_en.strb
type DeckBonus struct {
	// bonus id
	Id int `json:"_id"`
	// Affects ATK or DEF
	AtkDefFlg int `json:"atk_def_flg"`
	// ?
	ValueType int `json:"value_type"`
	// amount of the modifier
	Value int `json:"value"`
	// ?
	DownGrade int `json:"down_grade"`
	// deck condition
	CondType int `json:"cond_type"`
	// number of cards required
	ReqNum int `json:"req_num"`
	// allows duplicates
	DupFlg      int    `json:"dup_flg"`
	Name        string `json:"-"`
	Description string `json:"-"`
}

// Deck Bonus Conditions from masfter file field "deck_bonus_cond"
type DeckBonusCond struct {
	// deck condition id
	Id int `json:"_id"`
	// deck bonus id
	DeckBonusId int `json:"deck_bonus_id"`
	// group
	Group int `json:"group"`
	// type id
	CondTypeId int `json:"cond_type_id"`
	// reference
	RefId int `json:"ref_id"`
}
