package vc

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

//HD Images are located at the following URL Pattern:
//https://d2n1d3zrlbtx8o.cloudfront.net/download/CardHD.zip/CARDFILE.TIMESTAMP
//we have yet to determine how the timestamp is decided

// the card names match the ones listed in the MsgCardName_en.strb file
type Card struct {
	Id                int    `json:"_id"`                 // card id
	CardNo            int    `json:"card_no"`             // card number, matches to the image file
	CardCharaId       int    `json:"card_chara_id"`       // card character id
	CardRareId        int    `json:"card_rare_id"`        // rarity of the card
	CardTypeId        int    `json:"card_type_id"`        // type of the card (Passion, Cool, Light, Dark)
	DeckCost          int    `json:"deck_cost"`           // unit cost
	LastEvolutionRank int    `json:"last_evolution_rank"` // number of evolution statges available to the card
	EvolutionRank     int    `json:"evolution_rank"`      // this card current evolution stage
	EvolutionCardId   int    `json:"evolution_card_id"`   // id of the card that this card evolves into, -1 for no evolution
	TransCardId       int    `json:"trans_card_id"`       // id of a possible turnover accident
	FollowerKindId    int    `json:"follower_kind_id"`    // cost of the followers?
	DefaultFollower   int    `json:"default_follower"`    // base soldiers
	MaxFollower       int    `json:"max_follower"`        // max soldiers if evolved minimally
	DefaultOffense    int    `json:"default_offense"`     // base ATK
	MaxOffense        int    `json:"max_offense"`         // max ATK if evolved minimally
	DefaultDefense    int    `json:"default_defense"`     // base DEF
	MaxDefense        int    `json:"max_defense"`         // max DEF if evolved minimally
	SkillId1          int    `json:"skill_id_1"`          // First Skill
	SkillId2          int    `json:"skill_id_2"`          // second Skill
	SkillId3          int    `json:"skill_id_3"`          // third Skill (LR)
	SpecialSkillId1   int    `json:"special_skill_id_1"`  // Awakened Burst type (GSR,GUR,GLR)
	ThorSkillId1      int    `json:"thor_skill_id_1"`     // no one knows
	MedalRate         int    `json:"medal_rate"`          // amount of medals can be traded for
	Price             int    `json:"price"`               // amount of gold can be traded for
	IsClosed          int    `json:"is_closed"`           // is closed
	Name              string `json:"name"`                // name from the strings file

	//Character Link
	character *CardCharacter
	archwitch *Archwitch
	//Skill Links
	skill1        *Skill
	skill2        *Skill
	skill3        *Skill
	specialSkill1 *Skill
	thorSkill1    *Skill
	// possible card evolutions
	prevEvo  *Card
	nextEvo  *Card
	_allEvos map[string]*Card
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
	Material5Item  int `json:"material_5_item"`
	Material5Count int `json:"material_5_count"`
	// ? Order in the "Awoken Card List maybe?"
	Order int `json:"order"`
	// still available?
	IsClosed int `json:"is_closed"`
}

func (ca *CardAwaken) Item(i int, data *VcFile) *Item {
	if i < 1 || i > 5 {
		return nil
	}
	switch i {
	case 1:
		if ca.Material1Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material1Item, data.Items)
	case 2:
		if ca.Material2Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material2Item, data.Items)
	case 3:
		if ca.Material3Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material3Item, data.Items)
	case 4:
		if ca.Material4Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material4Item, data.Items)
	case 5:
		if ca.Material5Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material5Item, data.Items)
	}
	return nil
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

type CardRarity struct {
	Id               int `json:"_id"`
	MaxCardLevel     int `json:"max_card_level"`
	SkillCoefficient int `json:"skill_coefficient"`
	// used to calculate the amount of exp this card gives when used as a material
	CardExpCoefficient     int    `json:"card_exp_coefficient"`
	EvolutionCoefficient   int    `json:"evolution_coefficient"`
	GuildbattleCoefficient int    `json:"guildbattle_coefficient"`
	Order                  int    `json:"order"`
	Signature              string `json:"signature"`
	CardLevelCoefficient   int    `json:"card_level_coefficient"`
	FragmentSlot           int    `json:"fragment_slot"`
	LimtOffense            int    `json:"limt_offense"`
	LimtDefense            int    `json:"limt_defense"`
	LimtMaxFollower        int    `json:"limt_max_follower"`
}

type CardSpecialCompose struct {
	Id           int `json:"_id"`
	CardMasterId int `json:"card_master_id"`
	Ratio        int `json:"ratio"` // same as CardRarity.CardExpCoefficient except for a specific card
}

func (c *Card) Image() string {
	return fmt.Sprintf("cd_%05d", c.CardNo)
}

func (c *Card) Rarity() string {
	if c.CardRareId >= 0 {
		return Rarity[c.CardRareId-1]
	} else {
		return ""
	}
}

func (c *Card) CardRarity(v *VcFile) *CardRarity {
	if c.CardRareId >= 0 {
		for _, cr := range v.CardRarities {
			if cr.Id == c.CardRareId {
				return &cr
			}
		}
	}
	return nil
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
		for k, val := range v.CardCharacters {
			if val.Id == c.CardCharaId {
				c.character = &v.CardCharacters[k]
				break
			}
		}
	}
	return c.character
}

func (c *Card) NextEvo(v *VcFile) *Card {
	if c.nextEvo == nil {
		if c.CardCharaId <= 0 || c.EvolutionCardId <= 0 || c.Rarity()[0] == 'H' {
			return nil
		}

		var tmp *Card
		for i, cd := range c.Character(v).Cards(v) {
			if cd.Id == c.EvolutionCardId {
				tmp = &(c.Character(v)._cards[i])
			}
		}

		// Terra -> Rhea evos to a different card
		if tmp == nil || tmp.CardCharaId != c.CardCharaId {
			return nil
		}
		c.nextEvo = tmp
		tmp.prevEvo = c
	}
	return c.nextEvo
}

func (c *Card) PrevEvo(v *VcFile) *Card {
	if c.prevEvo == nil {
		// no charcter ID or already lowest evo rank
		if c.CardCharaId <= 0 || c.EvolutionRank < 0 {
			return nil
		}

		var tmp *Card
		for i, cd := range c.Character(v).Cards(v) {
			if c.Id == cd.EvolutionCardId {
				tmp = &(c.Character(v)._cards[i])
			}
		}

		// Terra -> Rhea evos to a different card
		if tmp == nil || tmp.CardCharaId != c.CardCharaId {
			return nil
		}
		c.prevEvo = tmp
		tmp.nextEvo = c
	}
	return c.prevEvo
}

func (c *Card) calculateEvoStat(material1Stat, material2Stat, resultDefault, resultMax int) (ret int) {
	var evoRate float64
	if strings.HasSuffix(c.Rarity(), "UR") || strings.HasSuffix(c.Rarity(), "LR") {
		// LR and UR evo rate. basically no evo bonus applied to the result card.
		evoRate = 1.0
	} else if c.EvolutionRank == c.LastEvolutionRank {
		// 4* evo
		if c.CardCharaId == 250 || c.CardCharaId == 315 {
			// queen of ice, strategist
			evoRate = 1.209
		} else {
			// all other N-SR
			evoRate = 1.1
		}
	} else {
		//1*-3* evos
		if c.CardCharaId == 250 || c.CardCharaId == 315 {
			// queen of ice, strategist
			evoRate = 1.155
		} else {
			// all other N-SR
			evoRate = 1.05
		}
	}

	ret = (int(0.15 * float64(material1Stat))) +
		(int(0.15 * float64(material2Stat))) +
		(int(evoRate * float64(resultMax)))

	return
}

func calculateAmalStat(matStats []int, resultMax int) (ret int) {
	ret = resultMax
	for _, matStat := range matStats {
		ret += int(float64(matStat) * 0.08)
	}
	return
}

func (c *Card) calculateAwakeningStat(materialStat, resultMaxStat int) (ret int) {
	// ret = int(float64(materialStat)*0.20124178) + resultMaxStat
	return -1
}

/*
calculates the standard evolution stat but at max level. If this is a 4* card,
calculates for 5-card evo
*/
func (c *Card) EvoStandardMaxAttack(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoStandardMaxAttack(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxOffense)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoStandardMaxAttack(v), c.MaxOffense)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoStandardMaxAttack(v)))) * 2) + c.MaxOffense
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				ret = c.MaxOffense
			}
		}
	} else {
		materialStat := materialCard.EvoStandardMaxAttack(v)
		firstEvo := c.GetEvolutions(v)["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, firstEvo.MaxOffense, c.DefaultOffense, c.MaxOffense)
	}
	if ret > rarity.LimtOffense {
		return rarity.LimtOffense
	}
	return
}

/*
calculates the standard evolution stat but at max level. If this is a 4* card,
calculates for 5-card evo
*/
func (c *Card) EvoStandardMaxDefense(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoStandardMaxDefense(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxDefense)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoStandardMaxDefense(v), c.MaxDefense)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoStandardMaxDefense(v)))) * 2) + c.MaxDefense
				return
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				return c.MaxDefense
			}
		}
	} else {
		materialStat := materialCard.EvoStandardMaxDefense(v)
		firstEvo := c.GetEvolutions(v)["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, firstEvo.MaxDefense, c.DefaultDefense, c.MaxDefense)
	}
	if ret > rarity.LimtDefense {
		return rarity.LimtDefense
	}
	return
}

/*
calculates the standard evolution stat but at max level. If this is a 4* card,
calculates for 5-card evo
*/
func (c *Card) EvoStandardMaxSoldier(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoStandardMaxSoldier(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxFollower)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoStandardMaxSoldier(v), c.MaxFollower)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoStandardMaxSoldier(v)))) * 2) + c.MaxFollower
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				ret = c.MaxFollower
			}
		}
	} else {
		materialStat := materialCard.EvoStandardMaxSoldier(v)
		firstEvo := c.GetEvolutions(v)["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, firstEvo.MaxFollower, c.DefaultFollower, c.MaxFollower)
	}
	if ret > rarity.LimtMaxFollower {
		return rarity.LimtMaxFollower
	}
	return
}

/*
calculates the perfect evolution stat. If this is a 4* card, calculates for 5-card evo
*/
func (c *Card) EvoPerfectMaxAttack(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoPerfectMaxAttack(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxOffense)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoPerfectMaxAttack(v), c.MaxOffense)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoPerfectMaxAttack(v)))) * 2) + c.MaxOffense
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				ret = c.MaxOffense
			}
		}
	} else {
		rarity := c.CardRarity(v)
		materialStat := materialCard.EvoPerfectMaxAttack(v)
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, materialStat, c.DefaultOffense, c.MaxOffense)
		if ret > rarity.LimtOffense {
			return rarity.LimtOffense
		}
	}
	return
}

/*
calculates the perfect evolution stat. If this is a 4* card, calculates for 5-card evo
*/
func (c *Card) EvoPerfectMaxDefense(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoPerfectMaxDefense(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxDefense)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoPerfectMaxDefense(v), c.MaxDefense)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoPerfectMaxDefense(v)))) * 2) + c.MaxDefense
				return
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				return c.MaxDefense
			}
		}
	} else {
		rarity := c.CardRarity(v)
		materialStat := materialCard.EvoPerfectMaxDefense(v)
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, materialStat, c.DefaultOffense, c.MaxOffense)
		if ret > rarity.LimtOffense {
			return rarity.LimtOffense
		}
	}
	return
}

/*
calculates the perfect evolution stat. If this is a 4* card, calculates for 5-card evo
*/
func (c *Card) EvoPerfectMaxSoldier(v *VcFile) (ret int) {
	materialCard := c.PrevEvo(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.Id))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			amalgs := c.Amalgamations(v)
			var myAmal Amalgamation
			for _, a := range amalgs {
				if c.Id == a.FusionCardId {
					myAmal = a
					break
				}
			}
			mats := myAmal.Materials(v)
			matStats := make([]int, 0)
			for _, mat := range mats {
				if mat.Id != c.Id {
					matStats = append(matStats, mat.EvoPerfectMaxSoldier(v))
				}
			}
			ret = calculateAmalStat(matStats, c.MaxFollower)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			ret = c.calculateAwakeningStat(mat.EvoPerfectMaxSoldier(v), c.MaxFollower)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				ret = ((int(0.15 * float64(turnOver.EvoPerfectMaxSoldier(v)))) * 2) + c.MaxFollower
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				ret = c.MaxFollower
			}
		}
	} else {
		rarity := c.CardRarity(v)
		materialStat := materialCard.EvoPerfectMaxSoldier(v)
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
		ret = c.calculateEvoStat(materialStat, materialStat, c.DefaultOffense, c.MaxOffense)
		if ret > rarity.LimtOffense {
			return rarity.LimtOffense
		}
	}
	return
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

func (c *Card) Skill3(v *VcFile) *Skill {
	if c.skill3 == nil && c.SkillId3 > 0 {
		c.skill3 = SkillScan(c.SkillId3, v.Skills)
	}
	return c.skill3
}

func (c *Card) SpecialSkill1(v *VcFile) *Skill {
	if c.specialSkill1 == nil && c.SpecialSkillId1 > 0 {
		c.specialSkill1 = SkillScan(c.SpecialSkillId1, v.Skills)
	}
	return c.specialSkill1
}

func (c *Card) ThorSkill1(v *VcFile) *Skill {
	if c.thorSkill1 == nil && c.ThorSkillId1 > 0 {
		c.thorSkill1 = SkillScan(c.ThorSkillId1, v.Skills)
	}
	return c.thorSkill1
}

func CardScan(id int, cards []Card) *Card {
	if id <= 0 {
		return nil
	}
	l := len(cards)
	i := sort.Search(l, func(i int) bool { return cards[i].Id >= id })
	if i >= 0 && i < l && cards[i].Id == id {
		return &(cards[i])
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

func (c *Card) Skill3Name(v *VcFile) string {
	s := c.Skill3(v)
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
	if s.Name != "" {
		return s.Name
	}
	i := strconv.Itoa(s.Id)
	return i
}

func (c *Card) ThorSkill1Name(v *VcFile) string {
	s := c.ThorSkill1(v)
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

type CardList []Card

func (d CardList) Earliest() *Card {
	var min *Card = nil
	for _, card := range d {
		if min == nil || min.Id > card.Id {
			min = &card
		}
	}
	return min
}

func (d CardList) Latest() *Card {
	var max *Card = nil
	for _, card := range d {
		if max == nil || max.Id < card.Id {
			max = &card
		}
	}
	return max
}

func (card *Card) GetEvolutions(v *VcFile) map[string]*Card {
	if card._allEvos == nil {
		ret := make(map[string]*Card)

		// handle cards like Chimrey and Time Traveler (enemy)
		if card.CardCharaId < 1 {
			os.Stdout.WriteString(fmt.Sprintf("No character info Card: %d, Name: %s, Evo: %d\n", card.Id, card.Name, card.EvolutionRank))
			ret["0"] = card
			card._allEvos = ret
			return ret
		}

		c := card
		// check if this is an awoken card
		if c.Rarity()[0] == 'G' {
			tmp := c.AwakensFrom(v)
			if tmp == nil {
				ch := c.Character(v)
				if ch != nil && ch.Cards(v)[0].Name == c.Name {
					c = &(ch.Cards(v)[0])
				}
				// the name changed, so we'll keep this card
			} else {
				c = tmp
			}
		}

		// get earliest evo
		for tmp := c.PrevEvo(v); tmp != nil; tmp = tmp.PrevEvo(v) {
			c = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c.Id, c.Name, c.EvolutionRank))
		}

		// check for a base amalgamation with the same name
		// if there is one, use that for the base card
		var getAmalBaseCard func(card *Card) *Card
		getAmalBaseCard = func(card *Card) *Card {
			if card.IsAmalgamation(v.Amalgamations) {
				os.Stdout.WriteString(fmt.Sprintf("Checking Amalgamation base for Card: %d, Name: %s, Evo: %d\n", card.Id, card.Name, card.EvolutionRank))
				for _, amal := range card.Amalgamations(v) {
					if card.Id == amal.FusionCardId {
						// material 1
						ac := CardScan(amal.Material1, v.Cards)
						if ac.Id != card.Id && ac.Name == card.Name {
							if ac.IsAmalgamation(v.Amalgamations) {
								return getAmalBaseCard(ac)
							} else {
								return ac
							}
						}
						// material 2
						ac = CardScan(amal.Material2, v.Cards)
						if ac.Id != card.Id && ac.Name == card.Name {
							if ac.IsAmalgamation(v.Amalgamations) {
								return getAmalBaseCard(ac)
							} else {
								return ac
							}
						}
						// material 3
						ac = CardScan(amal.Material3, v.Cards)
						if ac != nil && ac.Id != card.Id && ac.Name == card.Name {
							if ac.IsAmalgamation(v.Amalgamations) {
								return getAmalBaseCard(ac)
							} else {
								return ac
							}
						}
						// material 4
						ac = CardScan(amal.Material4, v.Cards)
						if ac != nil && ac.Id != card.Id && ac.Name == card.Name {
							if ac.IsAmalgamation(v.Amalgamations) {
								return getAmalBaseCard(ac)
							} else {
								return ac
							}
						}
					}
				}
			}
			return card
		}

		// at this point we should have the first card in the evolution path
		c = getAmalBaseCard(c)

		// get earliest evo (again...)
		for tmp := c.PrevEvo(v); tmp != nil; tmp = tmp.PrevEvo(v) {
			c = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c.Id, c.Name, c.EvolutionRank))
		}

		os.Stdout.WriteString(fmt.Sprintf("Base Card: %d, Name: '%s', Evo: %d\n", c.Id, c.Name, c.EvolutionRank))

		// populate the actual evos.

		checkEndCards := func(c *Card) (awakening, amalCard, amalAwakening *Card) {
			awakening = c.AwakensTo(v)
			// check for Amalgamation
			if c.HasAmalgamation(v.Amalgamations) {
				amals := c.Amalgamations(v)
				for _, amal := range amals {
					// get the result card
					tamalCard := CardScan(amal.FusionCardId, v.Cards)
					if tamalCard != nil && tamalCard.Id != c.Id {
						os.Stdout.WriteString(fmt.Sprintf("Found amalgamation: %d, Name: '%s', Evo: %d\n", tamalCard.Id, tamalCard.Name, tamalCard.EvolutionRank))
						if tamalCard.Name == c.Name {
							amalCard = tamalCard
							// check for amal awakening
							amalAwakening = amalCard.AwakensTo(v)
							return // awakening, amalCard, amalAwakening
						}
					}
				}
			}
			return // awakening, amalCard, amalAwakening
		}

		for nextEvo := c; nextEvo != nil; nextEvo = nextEvo.NextEvo(v) {
			os.Stdout.WriteString(fmt.Sprintf("Next Evo is Card: %d, Name: '%s', Evo: %d\n", nextEvo.Id, nextEvo.Name, nextEvo.EvolutionRank))
			if nextEvo.EvolutionRank <= 0 {
				evoRank := "0"
				if nextEvo.Rarity()[0] == 'H' {
					evoRank = "H"
				}
				ret[evoRank] = nextEvo
				if nextEvo.LastEvolutionRank < 0 {
					// check for awakening
					awakening, amalCard, amalAwakening := checkEndCards(nextEvo)
					if awakening != nil {
						ret["G"] = awakening
					}
					if amalCard != nil {
						ret["A"] = amalCard
					}
					if amalAwakening != nil {
						ret["GA"] = amalAwakening
					}
				}
			} else if nextEvo.Rarity()[0] == 'G' {
				// for some reason we hit a G during Evo traversal. Probably a G originating
				// from amalgamation
				ret["G"] = nextEvo
			} else if nextEvo.EvolutionRank == c.LastEvolutionRank || nextEvo.Rarity()[0] == 'H' || nextEvo.LastEvolutionRank < 0 {
				ret["H"] = nextEvo
				// check for awakening
				awakening, amalCard, amalAwakening := checkEndCards(nextEvo)
				if awakening != nil {
					ret["G"] = awakening
				}
				if amalCard != nil {
					ret["A"] = amalCard
				}
				if amalAwakening != nil {
					ret["GA"] = amalAwakening
				}
			} else {
				// not the last evo. These never awaken or have amalgamations
				ret[strconv.Itoa(nextEvo.EvolutionRank)] = nextEvo
			}
		}

		// if we have a GA with no H and no G, just change GA -> G for simplicity
		if _, ok := ret["GA"]; ok {
			_, hasH := ret["H"]
			_, hasG := ret["G"]
			if !hasH && !hasG {
				ret["G"] = ret["GA"]
				delete(ret, "GA")
			}
		}

		// normalize X cards
		lenEvoKeys := len(ret)
		if lenEvoKeys == 1 {
			for k, evo := range ret {
				r := evo.Rarity()[0]
				if k != "0" && (evo.EvolutionRank == 1 || evo.EvolutionRank < 0) && r != 'H' && r != 'G' {
					ret["0"] = evo
					delete(ret, k)
				}
			}
		}

		os.Stdout.WriteString("Found Evos: ")
		for key, card := range ret {
			card._allEvos = ret
			os.Stdout.WriteString(fmt.Sprintf("(%s: %d) ", key, card.Id))
		}
		os.Stdout.WriteString("\n")
		card._allEvos = ret
		return ret
	}
	return card._allEvos
}

var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}
var Rarity = [14]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR", "LR", "HLR", "GLR"}
