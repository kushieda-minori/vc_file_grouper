package vc

// Level from the master_data file, this is the "levels" field
// this matches to the Kingdom level in the game
type Level struct {
	ID        int `json:"_id"`        // level
	Energy    int `json:"energy"`     // energy
	Exp       int `json:"exp"`        //exp
	FriendNum int `json:"friend_num"` //number of friends
	DeckCost  int `json:"deck_cost"`  // deck cost
	NpcCost   int `json:"npc_cost"`   // vitality
	KingCost  int `json:"king_cost"`  // battle points
}

// LevelResource the amount of "resources" needed to upgrade a card to a
// card level. The ID matches to the card level
type LevelResource struct {
	ID     int `json:"_id"`
	Gold   int `json:"coin"`
	Iron   int `json:"iron"`
	Ether  int `json:"ether"`
	Elixir int `json:"elixir"`
}

// LevelupBonus for kingdom level increases.
type LevelupBonus struct {
	ID          int `json:"_id"`          // ID
	Level       int `json:"level"`        // Kingdom level
	Jewels      int `json:"cash"`         // jewels
	FriendPoint int `json:"friend_point"` // friendship points
	Gold        int `json:"coin"`         // gold
	Iron        int `json:"iron"`         // iron
	Ether       int `json:"ether"`        // ether
	Gem         int `json:"elixir"`       // gems
	Exp         int `json:"exp"`          // kingdom exp
	ItemID      int `json:"item_id"`      // Item rewarded
	FragmentID  int `json:"fragment_id"`  //Fragment? rewarded
	CardID      int `json:"card_id"`      // card rewarded
	Num         int `json:"num"`          // number of either Items, Cards, or Fragments?
	PublicFlag  int `json:"public_flg"`   // ?
}
