package vc

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
	Description     string `json:"description"`
	Friendship      string `json:"friendship"`
	Login           string `json:"login"`
	Meet            string `json:"meet"`
	BattleStart     string `json:"battleStart"`
	BattleEnd       string `json:"battleEnd"`
	FriendshipMax   string `json:"friendshipMax"`
	FriendshipEvent string `json:"friendshipEvent"`
	_cards          []Card
}

func (c *CardCharacter) Cards(v *VcFile) []Card {
	if c._cards == nil || len(c._cards) == 0 {
		c._cards = make([]Card, 0)
		for _, val := range v.Cards {
			//return the first one we find.
			if val.CardCharaId == c.Id {
				c._cards = append(c._cards, val)
			}
		}
	}
	return c._cards
}

func (c *CardCharacter) FirstEvoCard(v *VcFile) (card *Card) {
	card = nil
	for i, cd := range c.Cards(v) {
		if card == nil || cd.EvolutionRank <= card.EvolutionRank {
			card = &(c._cards[i])
		}
	}
	return
}

func CardCharacterScan(charId int, chars []CardCharacter) *CardCharacter {
	if charId > 0 {
		for k, val := range chars {
			//return the first one we find.
			if val.Id == charId {
				return &chars[k]
			}
		}
	}
	return nil
}
