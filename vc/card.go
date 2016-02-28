package vc

import (
	"fmt"
	"strconv"
	"strings"
)

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
	Name string `json:"name"`
	//Character Link
	character *CardCharacter `json:"-"`
	archwitch *Archwitch     `json:"-"`
	//Skill Links
	skill1        *Skill `json:"-"`
	skill2        *Skill `json:"-"`
	specialSkill1 *Skill `json:"-"`
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
	Description     string `json:"description"`
	Friendship      string `json:"friendship"`
	Login           string `json:"login"`
	Meet            string `json:"meet"`
	BattleStart     string `json:"battleStart"`
	BattleEnd       string `json:"battleEnd"`
	FriendshipMax   string `json:"friendshipMax"`
	FriendshipEvent string `json:"friendshipEvent"`
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

func (c *Card) Image() string {
	return fmt.Sprintf("cd_%05d.png", c.CardNo)
}

func (c *Card) Rarity() string {
	if c.CardRareId >= 0 {
		return Rarity[c.CardRareId-1]
	} else {
		return ""
	}
}

func (c *Card) Element() string {
	if c.CardTypeId >= 0 {
		return Elements[c.CardTypeId-1]
	} else {
		return ""
	}
}

func (c *Card) Character(v *VcFile) *CardCharacter {
	if c.character == nil && c.CardCharaId > 0 {
		c.character = &v.CardCharacter[c.CardCharaId-1]
	}
	return c.character
}

func (c *Card) Archwitch(v *VcFile) *Archwitch {
	if c.archwitch == nil {
		for _, aw := range v.Archwitches {
			if c.Id == aw.CardMasterId {
				c.archwitch = &aw
				break
			}
		}
	}
	return c.archwitch
}

func (c *Card) EvoAccident(cards []Card) *Card {
	return CardScan(c.TransCardId, cards)
}

func (c *Card) EvoAccidentOf(cards []Card) *Card {
	for key, val := range cards {
		if val.TransCardId == c.Id {
			return &(cards[key])
		}
	}
	return nil
}

func (c *Card) Amalgamations(v *VcFile) []Amalgamation {
	ret := make([]Amalgamation, 0)
	for _, a := range v.Amalgamations {
		if c.Id == a.FusionCardId ||
			c.Id == a.Material1 ||
			c.Id == a.Material2 ||
			c.Id == a.Material3 ||
			c.Id == a.Material4 {

			ret = append(ret, a)
		}
	}
	return ret
}

func (c *Card) AwakensTo(v *VcFile) *Card {
	for _, val := range v.Awakenings {
		if c.Id == val.BaseCardId {
			return CardScan(val.ResultCardId, v.Cards)
		}
	}
	return nil
}

func (c *Card) AwakensFrom(v *VcFile) *Card {
	for _, val := range v.Awakenings {
		if c.Id == val.ResultCardId {
			return CardScan(val.BaseCardId, v.Cards)
		}
	}
	return nil
}

func (c *Card) HasAmalgamation(a []Amalgamation) bool {
	for _, v := range a {
		if c.Id == v.Material1 ||
			c.Id == v.Material2 ||
			c.Id == v.Material3 ||
			c.Id == v.Material4 {
			return true
		}
	}
	return false
}

func (c *Card) IsAmalgamation(a []Amalgamation) bool {
	for _, v := range a {
		if c.Id == v.FusionCardId {
			return true
		}
	}
	return false
}

func (c *Card) Skill1(v *VcFile) *Skill {
	if c.skill1 == nil && c.SkillId1 > 0 {
		c.skill1 = SkillScan(c.SkillId1, v.Skills)
	}
	return c.skill1
}

func (c *Card) Skill2(v *VcFile) *Skill {
	if c.skill2 == nil && c.SkillId2 > 0 {
		c.skill2 = SkillScan(c.SkillId2, v.Skills)
	}
	return c.skill2
}

func (c *Card) SpecialSkill1(v *VcFile) *Skill {
	if c.specialSkill1 == nil && c.SpecialSkillId1 > 0 {
		c.specialSkill1 = SkillScan(c.SpecialSkillId1, v.Skills)
	}
	return c.specialSkill1
}

func CardScan(cardId int, cards []Card) *Card {
	if cardId > 0 {
		if cardId < len(cards) && cards[cardId-1].Id == cardId {
			return &cards[cardId-1]
		}
		for k, val := range cards {
			if val.Id == cardId {
				return &cards[k]
			}
		}
	}
	return nil
}

func CardScanCharacter(charId int, cards []Card) *Card {
	if charId > 0 {
		for k, val := range cards {
			//return the first one we find.
			if val.CardCharaId == charId {
				return &cards[k]
			}
		}
	}
	return nil
}

func CardScanImage(cardId string, cards []Card) *Card {
	if cardId != "" {
		i, err := strconv.Atoi(cardId)
		if err != nil {
			return nil
		}
		for k, val := range cards {
			if val.CardNo == i {
				return &cards[k]
			}
		}
	}
	return nil
}

func (c *Card) Skill1Name(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) SkillMin(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.SkillMin()
}

func (c *Card) SkillMax(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.SkillMax()
}

func (c *Card) SkillProcs(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	// battle start skills seem to have random Max Count values. Force it to 1
	// since they can only proc once anyway
	if strings.Contains(strings.ToLower(c.SkillMin(v)), "battle start") {
		return "1"
	}
	// -1 MaxCount indicates no limit
	if s.MaxCount < 0 {
		return "Infinite"
	}
	return strconv.Itoa(s.MaxCount)
}

func (c *Card) SkillTarget(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.TargetScope()
}

func (c *Card) SkillTargetLogic(v *VcFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.TargetLogic()
}

func (c *Card) Skill2Name(v *VcFile) string {
	s := c.Skill2(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) SpecialSkill1Name(v *VcFile) string {
	s := c.SpecialSkill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) Description(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Description
}

func (c *Card) Friendship(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Friendship
}

func (c *Card) Login(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Login
}

func (c *Card) Meet(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Meet
}

func (c *Card) BattleStart(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.BattleStart
}

func (c *Card) BattleEnd(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.BattleEnd
}

func (c *Card) FriendshipMax(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipMax
}

func (c *Card) FriendshipEvent(v *VcFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipEvent
}

func (a *Amalgamation) MaterialCount() int {
	if a.Material4 > 0 {
		return 4
	}
	if a.Material3 > 0 {
		return 3
	}
	return 2
}
func (a *Amalgamation) Materials(v *VcFile) []*Card {
	ret := make([]*Card, 0)
	ret = append(ret, CardScan(a.Material1, v.Cards))
	ret = append(ret, CardScan(a.Material2, v.Cards))
	if a.Material3 > 0 {
		ret = append(ret, CardScan(a.Material3, v.Cards))
	}
	if a.Material4 > 0 {
		ret = append(ret, CardScan(a.Material4, v.Cards))
	}
	ret = append(ret, CardScan(a.FusionCardId, v.Cards))
	return ret
}

type ByMaterialCount []Amalgamation

func (s ByMaterialCount) Len() int {
	return len(s)
}
func (s ByMaterialCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByMaterialCount) Less(i, j int) bool {
	return s[i].MaterialCount() < s[j].MaterialCount()
}

var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}
var Rarity = [11]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR"}
