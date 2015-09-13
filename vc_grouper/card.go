package vc_grouper

//HD Images are located at the followinf URL Pattern:
//https://d2n1d3zrlbtx8o.cloudfront.net/download/CardHD.zip/CARDFILE.TIMESTAMP
//we have yet to determine how the timestamp is decided

// the card names match the ones listed in the MsgCardName_en.strb file
type Card struct {
	// card id
	Id int `json:"_id"`
	// card number, matches to the image file
	CardNo int `json:"card_no"`
	// card character id
	CardCharaId int `json:"card_chara_id"`
	// rarity of the card
	CardRareId int `json:"card_rare_id"`
	// type of the card (Passion, Cool, Light, Dark)
	CardTypeId int `json:"card_type_id"`
	// unit cost
	DeckCost int `json:"deck_cost"`
	// number of evolution statges available to the card
	LastEvolutionRank int `json:"last_evolution_rank"`
	// this card current evolution stage
	EvolutionRank int `json:"evolution_rank"`
	// id of the card that this card evolves into, -1 for no evolution
	EvolutionCardId int `json:"evolution_card_id"`
	// id of a possible turnover accident
	TransCardId int `json:"trans_card_id"`
	// cost of the followers?
	FollowerKindId int `json:"follower_kind_id"`
	// base soldiers
	DefaultFollower int `json:"default_follower"`
	// max soldiers if evolved minimally
	MaxFollower int `json:"max_follower"`
	// base ATK
	DefaultOffense int `json:"default_offense"`
	// max ATK if evolved minimally
	MaxOffense int `json:"max_offense"`
	// base DEF
	DefaultDefense int `json:"default_defense"`
	// max DEF if evolved minimally
	MaxDefense int `json:"max_defense"`
	// First Skill
	SkillId1 int `json:"skill_id_1"`
	// second Skill
	SkillId2 int `json:"skill_id_2"`
	// Awakened Burst type (GSR or GUR)
	SpecialSkillId1 int `json:"special_skill_id_1"`
	// amount of medals can be traded for
	MedalRate int `json:"medal_rate"`
	// amount of gold can be traded for
	Price int `json:"price"`
	// is closed
	IsClosed int `json:"is_closed"`
	// name from the strings file
	Name string `json:"-"`
}

// List of possible Fusions (Amalgamations) from master file field "fusion_list"
type Amalgamation struct {
	// internal id
	Id int `json:"_id"`
	// card 1
	Material1 int `json:"material_1"`
	// card 2
	Material2 int `json:"material_2"`
	// card 3
	Material3 int `json:"material_3"`
	// card 4
	Material4 int `json:"material_4"`
	// resulting card
	FusionCardId int `json:"fusion_card_id"`
}

// list of possible card awakeneings and thier cost from master file field "card_awaken"
type CardAwaken struct {
	// awakening id
	Id int `json:"_id"`
	// case card
	BaseCardId int `json:"base_card_id"`
	// result card
	ResultCardId int `json:"result_card_id"`
	// chance of success
	Percent int `json:"percent"`
	// material information
	Material1Item  int `json:"material_1_item"`
	Material1Count int `json:"material_1_count"`
	Material2Item  int `json:"material_2_item"`
	Material2Count int `json:"material_2_count"`
	Material3Item  int `json:"material_3_item"`
	Material3Count int `json:"material_3_count"`
	Material4Item  int `json:"material_4_item"`
	Material4Count int `json:"material_4_count"`
	// ? Order in the "Awakend Card List maybe?"
	Order int `json:"order"`
	// still available?
	IsClosed int `json:"is_closed"`
}

// Card Character info from master_data field "card_character"
// These match up with all the MsgChara*_en.strb files
type CardCharacter struct {
	// card  charcter Id, matches to Card -> card_chara_id
	Id int `json:"_id"`
	// hidden param 1
	HiddenParam1 int `json:"hidden_param_1"`
	// `hidden param 2
	HiddenParam2 int `json:"hidden_param_2"`
	// hidden param 3
	HiddenParam3 int `json:"hidden_param_3"`
	// max friendship 0-30
	MaxFriendship int `json:"max_friendship"`
	// text from Strings file
	Description      string `json:"-"`
	Friendship       string `json:"-"`
	Login            string `json:"-"`
	Meet             string `json:"-"`
	Battle_start     string `json:"-"`
	Battle_end       string `json:"-"`
	Friendship_max   string `json:"-"`
	Friendship_event string `json:"-"`
}

// Follower kinds for soldier replenishment on cards
//these come from master file field "follower_kinds"
type FollowerKind struct {
	Id    int `json:"_id"`
	Coin  int `json:"coin"`
	Iron  int `json:"iron"`
	Ether int `json:"ether"`
	// not really used
	Speed int `json:"speed"`
}

var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}
var Rarity = [11]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR"}
