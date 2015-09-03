package vc_decoder

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
}

// Skills info from master data field "skills" These match to the string in the files:
// MsgSkillName_en.strb
// MsgSkillDesc_en.strb - shown on the card
// MsgSkillFire_en.strb - used during battle
type Skills struct {
	// skill id
	_id,
	// level type for skill upgrade costs
	level_type,
	// skill type
	_type,
	// id for timing
	timing_id,
	// max procs
	max_count,
	// cond scene
	cond_scene_id,
	// cond side
	cond_side_id,
	// cond
	cond_id,
	// king series
	king_series_id,
	// king id
	king_id,
	// cond param
	cond_param,
	// default proc rate
	default_ratio,
	// max proc rate
	max_ratio,
	// date accessible
	public_start_datetime,
	public_end_datetime,
	// effect info
	effect_id,
	effect_param,
	effect_param_2,
	effect_param_3,
	effect_param_4,
	effect_param_5,
	effect_default_value,
	effect_max_value,
	// target info
	target_scope_id,
	target_logic_id,
	target_param,
	// animation info
	animation_id int
}

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

// List of possible Fusions (Amalgamations) from master file field "fusion_list"
type FusionList struct {
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

//"kings" is the list of Archwitches
// "kings": [{
// 		"_id": 1,
// 		"king_series_id": 1,
// 		"card_master_id": 317,
// 		"status_group_id": 4,
// 		"public_flg": 1,
// 		"rare_flg": 0,
// 		"rare_intensity": 1,
// 		"battle_time": 120,
// 		"exp": 40,
// 		"max_friendship": -1,
// 		"skill_id_1": 145,
// 		"skill_id_2": -1,
// 		"weather_id": 1,
// 		"model_name": "btl0001_v1",
// 		"chain_ratio_2": 0,
// 		"servant_id_1": -1,
// 		"servant_id_2": -1
// 	}

//"king_series" is the list of AW events and the RR prizes
// "king_series": [{
// 		"_id": 1,
// 		"reward_card_id": 317,
// 		"public_start_datetime": 1361113200,
// 		"public_end_datetime": 1362365940,
// 		"receive_limit_datetime": 1365044399,
// 		"is_beginner_king": 0
// 	}

//"king_friendship" is the chance of friendship increasing on an AW
// "king_friendship": [{
// 		"_id": 1,
// 		"king_id": 91,
// 		"friendship": 0,
// 		"up_rate": 0
// 	}

//"series" is the list of scared relics and the prizes for completeing them.
// "series": [{
// 		"_id": 1,
// 		"treasure_id_1": 1,
// 		"treasure_id_2": 2,
// 		"treasure_id_3": 3,
// 		"treasure_id_4": 4,
// 		"treasure_id_5": 5,
// 		"treasure_id_6": 6,
// 		"bonus_card_id_1": 77,
// 		"bonus_card_id_2": 262,
// 		"bonus_card_id_3": 77,
// 		"event_flg": 0,
// 		"area_attribute": -1,
// 		"public_flg": 1,
// 		"order": 2999
// 	}

// Follower kinds for soldier replenishment on cards
//these come from master file field "follower_kinds"
type follower_kinds struct {
	_id,
	coin,
	iron,
	ether,
	// not really used
	speed int
}

// the "garden" field lists some details about the kindoms available to the players
// "garden": [{
// 		"_id": 1,
// 		"block_x": 6,
// 		"block_y": 6,
// 		"unlock_block_x": 3,
// 		"unlock_block_y": 3,
// 		"bg_id": 1,
// 		"debris": 0,
// 		"castle_id": 7
// 	}, {
// 		"_id": 2,
// 		"block_x": 6,
// 		"block_y": 6,
// 		"unlock_block_x": 6,
// 		"unlock_block_y": 6,
// 		"bg_id": 2,
// 		"debris": 1,
// 		"castle_id": 66
// 	}]

//"garden_debris" lists information about clearing debris from your kingdom
// "garden_debris": [{
// 		"_id": 1,
// 		"garden_id": 2,
// 		"structure_id": 71,
// 		"x": 24,
// 		"y": 10,
// 		"level_cap": 1,
// 		"unlock_area_id": -1,
// 		"time": 1800,
// 		"coin": 8000,
// 		"iron": 8000,
// 		"ether": 8000,
// 		"cash": 0,
// 		"exp": 100
// 	}

// "structures" gives information about availability of for buildinds.
// The names of the structions in this list match those in the MsgBuildingName_en.strb file
// "structures": [{
// 		"_id": 1,
// 		"structure_type_id": 1,
// 		"max_lv": 10,
// 		"unlock_castle_id": 7,
// 		"unlock_castle_lv": -1,
// 		"unlock_area_id": -1,
// 		"base_num": 2,
// 		"size_x": 2,
// 		"size_y": 2,
// 		"order": 1000,
// 		"event_id": -1,
// 		"visitable": 0,
// 		"step": 0,
// 		"passable": 0,
// 		"connectable": 0,
// 		"enable": 1,
// 		"stockable": 1,
// 		"flag": 48,
// 		// 1 for kingdom 1, 2 for kingdom 2, 3 for both
// 		"garden_flag": 3
// 	}

// "event_structures" lists any structures available in the current event

//structure_level lists the level for the available structures
// "structure_level": [{
// 		"_id": 1,
// 		"structure_id": 28,
// 		"level": 1,
// 		"tex_id": 55,
// 		"level_cap": 1,
// 		"unlock_area_id": -1,
// 		"time": 0,
// 		"beginner_time": 0,
// 		"coin": 0,
// 		"iron": 0,
// 		"ether": 0,
// 		"cash": 500,
// 		"price": 0,
// 		"exp": 300
// 	}
