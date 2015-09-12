package vc_grouper

// Unit Bonuses from master file field "deck_bonus"
// these match with the strings in MsgDeckBonusName_en.strb
// and MsgDeckBonusDesc_en.strb
type DeckBonus struct {
	// bonus id
	_id,
	// Affects ATK or DEF
	atk_def_flg,
	// ?
	value_type,
	// amount of the modifier
	value,
	// ?
	down_grade,
	// deck condition
	cond_type,
	// number of cards required
	req_num,
	// allows duplicates
	dup_flg int
}

// Deck Bonus Conditions from masfter file field "deck_bonus_cond"
type DeckBonusCond struct {
	// deck condition id
	_id,
	// deck bonus id
	deck_bonus_id,
	// group
	group,
	// type id
	cond_type_id,
	// reference
	ref_id int
}
