package vc

//from the master_data file, this is the "levels" field
// this matches to the Kingdom level in the game
type Level struct {
	Id        int `json:"_id"`        // level
	Energy    int `json:"energy"`     // energy
	Exp       int `json:"exp"`        //exp
	FriendNum int `json:"friend_num"` //number of friends
	DeckCost  int `json:"deck_cost"`  // deck cost
	NpcCost   int `json:"npc_cost"`   // vitality
	KingCost  int `json:"king_cost"`  // battle points
}
