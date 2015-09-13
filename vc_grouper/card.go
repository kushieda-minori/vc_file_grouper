package vc_grouper

//HD Images are located at the followinf URL Pattern:
//https://d2n1d3zrlbtx8o.cloudfront.net/download/CardHD.zip/CARDFILE.TIMESTAMP
//we have yet to determine how the timestamp is decided

// the card names match the ones listed in the MsgCardName_en.strb file
type Card struct {
	// card id
	_id,
	// card number, matches to the image file
	card_no,
	// card character id
	card_chara_id,
	// rarity of the card
	card_rare_id,
	// type of the card (Passion, Cool, Light, Dark)
	card_type_id,
	// unit cost
	deck_cost,
	// number of evolution statges available to the card
	last_evolution_rank,
	// this card current evolution stage
	evolution_rank,
	// id of the card that this card evolves into, -1 for no evolution
	evolution_card_id,
	// id of a possible turnover accident
	trans_card_id,
	// cost of the followers?
	follower_kind_id,
	// base soldiers
	default_follower,
	// max soldiers if evolved minimally
	max_follower,
	// base ATK
	default_offense,
	// max ATK if evolved minimally
	max_offense,
	// base DEF
	default_defense,
	// max DEF if evolved minimally
	max_defense,
	// First Skill
	skill_id_1,
	// second Skill
	skill_id_2,
	// Awakened Burst type (GSR or GUR)
	special_skill_id_1,
	// amount of medals can be traded for
	medal_rate,
	// amount of gold can be traded for
	price,
	// is closed
	is_closed int
	// name from the strings file
	name string
}

// List of possible Fusions (Amalgamations) from master file field "fusion_list"
type Amalgamation struct {
	// internal id
	_id,
	// card 1
	material_1,
	// card 2
	material_2,
	// card 3
	material_3,
	// card 4
	material_4,
	// resulting card
	fusion_card_id int
}

// list of possible card awakeneings and thier cost from master file field "card_awaken"
type CardAwaken struct {
	// awakening id
	_id,
	// case card
	base_card_id,
	// result card
	result_card_id,
	// chance of success
	percent,
	// material information
	material_1_item,
	material_1_count,
	material_2_item,
	material_2_count,
	material_3_item,
	material_3_count,
	material_4_item,
	material_4_count,
	// ? Order in the "Awakend Card List maybe?"
	order,
	// still available?
	is_closed int
}

// Card Character info from master_data field "card_character"
// These match up with all the MsgChara*_en.strb files
type CardCharacter struct {
	// card  charcter _id, matches to Card -> card_chara_id
	_id,
	// hidden param 1
	hidden_param_1,
	// `hidden param 2
	hidden_param_2,
	// hidden param 3
	hidden_param_3,
	// max friendship 0-30
	max_friendship int
	// text from Strings file
	description,
	friendship,
	login,
	meet,
	battle_start,
	battle_end,
	friendship_max,
	friendship_event string
}

// Follower kinds for soldier replenishment on cards
//these come from master file field "follower_kinds"
type FollowerKinds struct {
	_id,
	coin,
	iron,
	ether,
	// not really used
	speed int
}

var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}
var Rarity = [11]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR"}
