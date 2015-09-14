package vc_grouper

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
	Name string `json:"-"`
	//Character Link
	character *CardCharacter
	//Skill Links
	skill1        *Skill
	skill2        *Skill
	specialSkill1 *Skill
}

func (c *Card) GetImage() string {
	return fmt.Sprintf("cd_%05d.png", c.CardNo)
}

func (c *Card) GetRarity() string {
	if c.CardRareId >= 0 {
		return Rarity[c.CardRareId-1]
	} else {
		return ""
	}
}

func (c *Card) GetElement() string {
	if c.CardTypeId >= 0 {
		return Elements[c.CardTypeId-1]
	} else {
		return ""
	}
}

func (c *Card) GetCharacter(v *VcFile) *CardCharacter {
	if c.character == nil && c.CardCharaId > 0 {
		c.character = &v.CardCharacter[c.CardCharaId-1]
	}
	return c.character
}

func (c *Card) GetEvoAccident(cards []Card) *Card {
	if c.TransCardId > 0 {
		return &(cards[c.TransCardId-1])
	}
	return nil
}

func (c *Card) IsEvoAccidentOf(cards []Card) *Card {
	for key, val := range cards {
		if val.TransCardId == c.Id {
			return &(cards[key])
		}
	}
	return nil
}

func (c *Card) GetAmalgamations(v *VcFile) []Amalgamation {
	ret := make([]Amalgamation, 0)
	for _, v := range v.Amalgamations {
		if c.Id == v.FusionCardId ||
			c.Id == v.Material1 ||
			c.Id == v.Material2 ||
			c.Id == v.Material3 ||
			c.Id == v.Material4 {
			ret = append(ret, v)
		}
	}
	return ret
}

func (c *Card) GetSkill1(v *VcFile) *Skill {
	if c.skill1 == nil && c.SkillId1 > 0 {
		c.skill1 = skillScan(c.SkillId1, v.Skills)
	}
	return c.skill1
}

func (c *Card) GetSkill2(v *VcFile) *Skill {
	if c.skill2 == nil && c.SkillId2 > 0 {
		c.skill2 = skillScan(c.SkillId2, v.Skills)
	}
	return c.skill2
}

func (c *Card) GetSpecialSkill1(v *VcFile) *Skill {
	if c.specialSkill1 == nil && c.SpecialSkillId1 > 0 {
		c.specialSkill1 = skillScan(c.SpecialSkillId1, v.Skills)
	}
	return c.specialSkill1
}

func skillScan(id int, skills []Skill) *Skill {
	if id <= 0 {
		return nil
	}
	if id < len(skills) && id == skills[id-1].Id {
		return &(skills[id-1])
	}
	for k, v := range skills {
		if id == v.Id {
			return &(skills[k])
		}
	}
	return nil
}

func (c *Card) GetSkill1Name(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) GetSkillMin(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	return s.GetSkillMin()
}

func (c *Card) GetSkillMax(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	return s.GetSkillMax()
}

func (c *Card) GetSkillProcs(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	// battle start skills seem to have random Max Count values. Force it to 1
	// since they can only proc once anyway
	if strings.Contains(strings.ToLower(c.GetSkillMin(v)), "battle start") {
		return "1"
	}
	// -1 MaxCount indicates no limit
	if s.MaxCount < 0 {
		return "Infinite"
	}
	return strconv.Itoa(s.MaxCount)
}

func (c *Card) GetSkillTarget(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	return s.GetTargetScope()
}

func (c *Card) GetSkillTargetLogic(v *VcFile) string {
	s := c.GetSkill1(v)
	if s == nil {
		return ""
	}
	return s.GetTargetLogic()
}

func (c *Card) GetSkill2Name(v *VcFile) string {
	s := c.GetSkill2(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) GetSpecialSkill1Name(v *VcFile) string {
	s := c.GetSpecialSkill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

func (c *Card) GetDescription(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.Description
}

func (c *Card) GetFriendship(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.Friendship
}

func (c *Card) GetLogin(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.Login
}

func (c *Card) GetMeet(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.Meet
}

func (c *Card) GetBattleStart(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.BattleStart
}

func (c *Card) GetBattleEnd(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.BattleEnd
}

func (c *Card) GetFriendshipMax(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipMax
}

func (c *Card) GetFriendshipEvent(v *VcFile) string {
	ch := c.GetCharacter(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipEvent
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
	Description     string `json:"-"`
	Friendship      string `json:"-"`
	Login           string `json:"-"`
	Meet            string `json:"-"`
	BattleStart     string `json:"-"`
	BattleEnd       string `json:"-"`
	FriendshipMax   string `json:"-"`
	FriendshipEvent string `json:"-"`
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
