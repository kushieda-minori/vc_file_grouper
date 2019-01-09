package vc

// CardCharacter info from master_data field "card_character"
// These match up with all the MsgChara*_en.strb files
type CardCharacter struct {
	// card  charcter ID, matches to Card -> card_chara_id
	ID int `json:"_id"`
	// hidden param 1
	HiddenParam1 int `json:"hidden_param_1"`
	// `hidden param 2
	HiddenParam2 int `json:"hidden_param_2"`
	// hidden param 3
	HiddenParam3 int `json:"hidden_param_3"`
	// max friendship 0-30
	MaxFriendship int `json:"max_friendship"`
	// text from Strings file
	Description     string `json:"-"` // MsgCharaDesc_en.strb
	Friendship      string `json:"-"` // MsgCharaBonds_en.strb
	Login           string `json:"-"` // MsgCharaWelcome_en.strb
	Meet            string `json:"-"` // MsgCharaMeet_en.strb
	BattleStart     string `json:"-"` // MsgCharaBtlStart_en.strb
	BattleEnd       string `json:"-"` // MsgCharaBtlEnd_en.strb
	FriendshipMax   string `json:"-"` // MsgCharaFriendshipMax_en.strb
	FriendshipEvent string `json:"-"` // MsgCharaBonds_en.strb
	Rebirth         string `json:"-"` // MsgCharaSuperAwaken_en.strb
	_cards          []Card
}

// Cards that are under this character
func (c *CardCharacter) Cards(v *VFile) []Card {
	if c._cards == nil || len(c._cards) == 0 {
		c._cards = make([]Card, 0)
		for _, val := range v.Cards {
			//return the first one we find.
			if val.CardCharaID == c.ID {
				c._cards = append(c._cards, val)
			}
		}
	}
	return c._cards
}

// FirstEvoCard first evolution of the cards under this character
func (c *CardCharacter) FirstEvoCard(v *VFile) (card *Card) {
	card = nil
	for i, cd := range c.Cards(v) {
		if card == nil || cd.EvolutionRank <= card.EvolutionRank {
			card = &(c._cards[i])
		}
	}
	return
}

// CardCharacterScan searches for a character by id
func CardCharacterScan(charID int, chars []CardCharacter) *CardCharacter {
	if charID > 0 {
		for k, val := range chars {
			//return the first one we find.
			if val.ID == charID {
				return &chars[k]
			}
		}
	}
	return nil
}
