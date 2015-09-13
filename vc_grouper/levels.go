package vc_grouper

//from the master_data file, this is the "levels" field
// this matches to the Kingdom level in the game
type Levels struct {
	// level
	Id int `json:"_id"`
	// energy
	Energy int `json:"energy"`
	//exp
	Exp int `json:"exp"`
	//number of friends
	FriendNum int `json:"friend_num"`
	// deck cost
	DeckCost int `json:"deck_cost"`
	// vitality
	NpcCost int `json:"npc_cost"`
	// battle points
	KingCost int `json:"king_cost"`
}
