package vc_grouper

//from the master_data file, this is the "levels" field
// this matches to the Kingdom level in the game
type Levels struct {
	// level
	_id,
	// energy
	energy,
	//exp
	exp,
	//number of friends
	friend_num,
	// deck cost
	deck_cost,
	// vitality
	npc_cost,
	// battle points
	king_cost int
}
