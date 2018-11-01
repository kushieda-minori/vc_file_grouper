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

// Card is a distinct card in the game. The card names match the ones listed in the MsgCardName_en.strb file
type Card struct {
	ID                int    `json:"_id"`                 // card id
	CardNo            int    `json:"card_no"`             // card number, matches to the image file
	CardCharaID       int    `json:"card_chara_id"`       // card character id
	CardRareID        int    `json:"card_rare_id"`        // rarity of the card
	CardTypeID        int    `json:"card_type_id"`        // type of the card (Passion, Cool, Light, Dark)
	DeckCost          int    `json:"deck_cost"`           // unit cost
	LastEvolutionRank int    `json:"last_evolution_rank"` // number of evolution statges available to the card
	EvolutionRank     int    `json:"evolution_rank"`      // this card current evolution stage
	EvolutionCardID   int    `json:"evolution_card_id"`   // id of the card that this card evolves into, -1 for no evolution
	TransCardID       int    `json:"trans_card_id"`       // id of a possible turnover accident
	FollowerKindID    int    `json:"follower_kind_id"`    // cost of the followers?
	DefaultFollower   int    `json:"default_follower"`    // base soldiers
	MaxFollower       int    `json:"max_follower"`        // max soldiers if evolved minimally
	DefaultOffense    int    `json:"default_offense"`     // base ATK
	MaxOffense        int    `json:"max_offense"`         // max ATK if evolved minimally
	DefaultDefense    int    `json:"default_defense"`     // base DEF
	MaxDefense        int    `json:"max_defense"`         // max DEF if evolved minimally
	SkillID1          int    `json:"skill_id_1"`          // First Skill
	SkillID2          int    `json:"skill_id_2"`          // second Skill
	SkillID3          int    `json:"skill_id_3"`          // third Skill (LR)
	SpecialSkillID1   int    `json:"special_skill_id_1"`  // Awakened Burst type (GSR,GUR,GLR)
	ThorSkillID1      int    `json:"thor_skill_id_1"`     // no one knows
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

// FollowerKind for soldier replenishment on cards
//these come from master file field "follower_kinds"
type FollowerKind struct {
	ID    int `json:"_id"`
	Coin  int `json:"coin"`
	Iron  int `json:"iron"`
	Ether int `json:"ether"`
	// not really used
	Speed int `json:"speed"`
}

// CardRarity information about a single Card Rarity
type CardRarity struct {
	ID               int `json:"_id"`
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

// CardSpecialCompose special information regaurding a cards use as material during fusion (leveling up)
type CardSpecialCompose struct {
	ID           int `json:"_id"`
	CardMasterID int `json:"card_master_id"`
	Ratio        int `json:"ratio"` // same as CardRarity.CardExpCoefficient except for a specific card
}

// Image name of the card
func (c *Card) Image() string {
	return fmt.Sprintf("cd_%05d", c.CardNo)
}

// Rarity of the card as a plain string
func (c *Card) Rarity() (ret string) {
	if c.CardRareID >= 0 {
		ret = Rarity[c.CardRareID-1]
		// need to handle X cards that have actual Evolutions (Philospher Stone)
		if ret == "X" && c.EvolutionRank > 0 && c.EvolutionRank == c.LastEvolutionRank && len(c._allEvos) > 1 {
			ret = "HX"
		}
	}
	return

}

// CardRarity with full rarity information
func (c *Card) CardRarity(v *VFile) *CardRarity {
	if c.CardRareID >= 0 {
		for _, cr := range v.CardRarities {
			if cr.ID == c.CardRareID {
				return &cr
			}
		}
	}
	return nil
}

// Element of the card
func (c *Card) Element() string {
	if c.CardTypeID >= 0 {
		return Elements[c.CardTypeID-1]
	}
	return ""

}

// Character information of the card
func (c *Card) Character(v *VFile) *CardCharacter {
	if c.character == nil && c.CardCharaID > 0 {
		for k, val := range v.CardCharacters {
			if val.ID == c.CardCharaID {
				c.character = &v.CardCharacters[k]
				break
			}
		}
	}
	return c.character
}

// NextEvo is the next evolution of this card, or nil if no further evolutions are possible.
// Amalgamations or Awakenings may still be possible.
func (c *Card) NextEvo(v *VFile) *Card {
	if c.ID == c.EvolutionCardID {
		// bad data
		return nil
	}
	if c.nextEvo == nil {
		if c.CardCharaID <= 0 || c.EvolutionCardID <= 0 || c.Rarity()[0] == 'H' {
			return nil
		}

		var tmp *Card
		for i, cd := range c.Character(v).Cards(v) {
			if cd.ID == c.EvolutionCardID {
				tmp = &(c.Character(v)._cards[i])
			}
		}

		// Terra -> Rhea evos to a different card
		if tmp == nil || tmp.CardCharaID != c.CardCharaID {
			return nil
		}
		if tmp.ID == c.ID {
			// bad data...
			return nil
		}
		c.nextEvo = tmp
		tmp.prevEvo = c
	}
	return c.nextEvo
}

// PrevEvo is the previous evolution of this card, or nil if no further evolutions are possible.
// Evo Accidents or amalgamtions may still be possible.
func (c *Card) PrevEvo(v *VFile) *Card {
	if c.prevEvo == nil {
		// no charcter ID or already lowest evo rank
		if c.CardCharaID <= 0 || c.EvolutionRank < 0 {
			return nil
		}

		var tmp *Card
		for i, cd := range c.Character(v).Cards(v) {
			if c.ID == cd.EvolutionCardID {
				tmp = &(c.Character(v)._cards[i])
			}
		}

		// Terra -> Rhea evos to a different card
		if tmp == nil || tmp.CardCharaID != c.CardCharaID {
			return nil
		}
		if tmp.ID == c.ID {
			// bad data...
			return nil
		}
		c.prevEvo = tmp
		tmp.nextEvo = c
	}
	return c.prevEvo
}

// PossibleMixedEvo Checks if this card has a possible mixed evo:
// Amalgamation at evo[0] with a standard evo after it. This may issue
// false positives for cards that can only be obtained through
// amalgamation at evo[0] and there is no drop/RR
func (c *Card) PossibleMixedEvo(v *VFile) bool {
	firstEvo := c.GetEvolutions(v)["0"]
	secondEvo := c.GetEvolutions(v)["1"]
	if secondEvo == nil {
		secondEvo = c.GetEvolutions(v)["H"]
	}
	return firstEvo != nil && secondEvo != nil &&
		firstEvo.IsAmalgamation(v.Amalgamations) &&
		firstEvo.EvolutionCardID == secondEvo.ID
}

// calculateEvoStat calculates the evo stats.
func (c *Card) calculateEvoStat(material1Stat, material2Stat, resultMax int) (ret int) {
	if c.EvolutionRank == c.LastEvolutionRank && c.EvolutionRank == 1 {
		// single evo gets no final stage bonus
		ret = resultMax
	} else if c.EvolutionRank == c.LastEvolutionRank {
		// 4* evo bonus for final stage
		if c.CardCharaID == 250 || c.CardCharaID == 315 {
			// queen of ice, strategist
			ret = (int(1.209 * float64(resultMax)))
		} else {
			// all other N-SR (UR and LR? ATM, there are no UR 4*)
			ret = (int(1.1 * float64(resultMax)))
		}
	} else {
		//4* evo bonus for intermediate stage
		if c.CardCharaID == 250 || c.CardCharaID == 315 {
			// queen of ice, strategist
			ret = (int(1.155 * float64(resultMax)))
		} else {
			// all other N-SR (UR and LR? ATM, there are no UR 4*)
			ret = (int(1.05 * float64(resultMax)))
		}
	}
	transferRate := 0.15
	ret += (int(transferRate * float64(material1Stat))) +
		(int(transferRate * float64(material2Stat)))

	return
}

// calculateEvoAccidentStat calculated the stats if an evo accident happens
func calculateEvoAccidentStat(materialStat, resultMax int) int {
	return ((int(0.15 * float64(materialStat))) * 2) + resultMax
}

// calculateLrAmalStat calculates amalgamation where one of the material is not max stat.
// this is generaly used for GUR+LR=LR or GLR+LR=GLR
func calculateLrAmalStat(sourceMatStat, lrMatStat, resultMax int) (ret int) {
	ret = resultMax
	ret += int(float64(sourceMatStat)*0.08) + int(float64(lrMatStat)*0.03)
	return
}

// calculateAwakeningStat calculated the stats after awakening
func (c *Card) calculateAwakeningStat(mat *Card, materialAtkAtLvl1, materialDefAtLvl1, materialSolAtLvl1 int, atLevel1 bool) (atk, def, soldier int) {
	// Awakening calculation thanks to Elle (https://docs.google.com/spreadsheets/d/1CT41xSuHyibfDSHQOyON4DkCPRlwW2WCy_ag6Pb76z4/edit?usp=sharing)
	if atLevel1 {
		atk = c.DefaultOffense + (materialAtkAtLvl1 - mat.DefaultOffense)
		def = c.DefaultDefense + (materialDefAtLvl1 - mat.DefaultDefense)
		soldier = c.DefaultFollower + (materialSolAtLvl1 - mat.DefaultFollower)
	} else {
		atk = c.MaxOffense + (materialAtkAtLvl1 - mat.DefaultOffense)
		def = c.MaxDefense + (materialDefAtLvl1 - mat.DefaultDefense)
		soldier = c.MaxFollower + (materialSolAtLvl1 - mat.DefaultFollower)
	}
	return
}

// AmalgamationStandard calculated the stats if the materials have been evo'd using
// standard processes (4* 5 card evos).
func (c *Card) AmalgamationStandard(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoStandard(v)
	}
	mats := myAmal.Materials(v)
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.EvoStandard(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// AmalgamationStandardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 5 card evos).
func (c *Card) AmalgamationStandardLvl1(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoStandardLvl1(v)
	}
	mats := myAmal.Materials(v)
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.EvoStandard(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Amalgamation6Card calculated the stats if the materials have been evo'd using
// standard processes (4* 6 card evos).
func (c *Card) Amalgamation6Card(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo6Card(v)
	}
	mats := myAmal.Materials(v)
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating 6-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.Evo6Card(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Amalgamation6CardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 6 card evos).
func (c *Card) Amalgamation6CardLvl1(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo6CardLvl1(v)
	}
	mats := myAmal.Materials(v)
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating 6-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.Evo6Card(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Amalgamation9Card calculated the stats if the materials have been evo'd using
// standard processes (4* 9 card evos).
func (c *Card) Amalgamation9Card(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo9Card(v)
	}
	mats := myAmal.Materials(v)
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating 9-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.Evo9Card(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Amalgamation9CardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 9 card evos).
func (c *Card) Amalgamation9CardLvl1(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo9CardLvl1(v)
	}
	mats := myAmal.Materials(v)
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating 9-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.Evo9Card(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// AmalgamationPerfect calculates the amalgamation stats if the material cards
// have all been evo'd perfectly (4* 16-card evos)
func (c *Card) AmalgamationPerfect(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoPerfect(v)
	}
	mats := myAmal.Materials(v)
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.EvoPerfect(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// AmalgamationPerfectLvl1 calculates the amalgamation stats at level 1 if the material cards
// have all been evo'd perfectly (4* 16-card evos)
func (c *Card) AmalgamationPerfectLvl1(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoPerfectLvl1(v)
	}
	mats := myAmal.Materials(v)
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
			matAtk, matDef, matSoldier := mat.EvoPerfect(v)
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// AmalgamationLRStaticLvl1 calculates the amalgamation stats if the material cards
// have all been evo'd perfectly (4* 16-card evos) with the exception of LR materials
// that have unchanging stats (i.e. 9999 / 9999)
func (c *Card) AmalgamationLRStaticLvl1(v *VFile) (atk, def, soldier int) {
	amalgs := c.Amalgamations(v)
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		if strings.HasSuffix(c.Rarity(), "LR") && c.Element() == "Special" {
			return c.EvoPerfectLvl1(v)
		}
		return c.EvoPerfect(v)
	}
	mats := myAmal.Materials(v)
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		if mat.ID != c.ID {
			os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
			if strings.HasSuffix(mat.Rarity(), "LR") && mat.Element() == "Special" {
				matAtk, matDef, matSoldier := mat.AmalgamationLRStaticLvl1(v)
				// no 5% bonus for max level
				atk += int(float64(matAtk) * 0.03)
				def += int(float64(matDef) * 0.03)
				soldier += int(float64(matSoldier) * 0.03)
			} else {
				matAtk, matDef, matSoldier := mat.AmalgamationLRStaticLvl1(v)
				atk += int(float64(matAtk) * 0.08)
				def += int(float64(matDef) * 0.08)
				soldier += int(float64(matSoldier) * 0.08)
			}
		}
	}
	rarity := c.CardRarity(v)
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// EvoStandard calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 5-card evo
func (c *Card) EvoStandard(v *VFile) (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard(v)
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.EvoStandardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Awakening standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.EvoStandard(v)
					atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
					def = calculateEvoAccidentStat(matDef, c.MaxDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
					os.Stdout.WriteString(fmt.Sprintf("Using Evo Accident stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.MaxOffense
					def = c.MaxDefense
					soldier = c.MaxFollower
					os.Stdout.WriteString(fmt.Sprintf("Using base Max stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoStandard(v)
		firstEvo := c.GetEvolutions(v)["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, firstEvo.MaxOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, firstEvo.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, firstEvo.MaxFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, c.MaxOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, c.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, c.MaxFollower)
		}
		os.Stdout.WriteString(fmt.Sprintf("Using Evo stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	os.Stdout.WriteString(fmt.Sprintf("Final stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	return
}

// EvoStandardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 5-card evo
func (c *Card) EvoStandardLvl1(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandardLvl1(v)
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.EvoStandardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Awakening StandardLvl1 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.EvoStandard(v)
					atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
					def = calculateEvoAccidentStat(matDef, c.MaxDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.DefaultOffense
					def = c.DefaultDefense
					soldier = c.DefaultFollower
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoStandard(v)
		firstEvo := c.GetEvolutions(v)["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, firstEvo.DefaultOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, firstEvo.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, firstEvo.DefaultFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, c.DefaultOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, c.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, c.DefaultFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// EvoMixed calculates the standard evolution for cards with a evo[0] amalgamation it calculates
// the evo[1] using 1 amalgamated, and one as a "drop"
func (c *Card) EvoMixed(v *VFile) (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Mixed Evo for %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard(v)
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.EvoMixedLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Awakening Mixed card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.EvoStandard(v)
					atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
					def = calculateEvoAccidentStat(matDef, c.MaxDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
					os.Stdout.WriteString(fmt.Sprintf("Using Evo Accident stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.MaxOffense
					def = c.MaxDefense
					soldier = c.MaxFollower
					os.Stdout.WriteString(fmt.Sprintf("Using base Max stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		var matAtk, matDef, matSoldier int
		if c.EvolutionRank == 1 && c.PossibleMixedEvo(v) && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matAtk, matDef, matSoldier = materialCard.AmalgamationStandard(v)
		} else {
			matAtk, matDef, matSoldier = materialCard.EvoStandard(v)
		}
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, firstEvo.MaxOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, firstEvo.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, firstEvo.MaxFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, c.MaxOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, c.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, c.MaxFollower)
		}
		os.Stdout.WriteString(fmt.Sprintf("Using Evo stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	os.Stdout.WriteString(fmt.Sprintf("Final stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	return
}

// EvoMixedLvl1 calculates the standard evolution for cards with a evo[0] amalgamation it calculates
// the evo[1] using 1 amalgamated, and one as a "drop"
func (c *Card) EvoMixedLvl1(v *VFile) (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Mixed Evo for Lvl1 %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard(v)
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.EvoMixedLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Awakening Mixed card Lvl1 %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.EvoStandard(v)
					atk = calculateEvoAccidentStat(matAtk, c.DefaultOffense)
					def = calculateEvoAccidentStat(matDef, c.DefaultDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.DefaultFollower)
					os.Stdout.WriteString(fmt.Sprintf("Using Evo Accident stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.DefaultOffense
					def = c.DefaultDefense
					soldier = c.DefaultFollower
					os.Stdout.WriteString(fmt.Sprintf("Using base Max stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		var matAtk, matDef, matSoldier int
		if c.EvolutionRank == 1 && c.PossibleMixedEvo(v) && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matAtk, matDef, matSoldier = materialCard.AmalgamationStandard(v)
		} else {
			matAtk, matDef, matSoldier = materialCard.EvoStandard(v)
		}
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, firstEvo.DefaultOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, firstEvo.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, firstEvo.DefaultFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, firstEvo.MaxOffense, c.DefaultOffense)
			def = c.calculateEvoStat(matDef, firstEvo.MaxDefense, c.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, firstEvo.MaxFollower, c.DefaultFollower)
		}
		os.Stdout.WriteString(fmt.Sprintf("Using Evo stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	os.Stdout.WriteString(fmt.Sprintf("Final stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
	return
}

// Evo6Card calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 6-card evo
func (c *Card) Evo6Card(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation6Card(v)
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.Evo6CardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Awakening 6 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.Evo6Card(v)
					atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
					def = calculateEvoAccidentStat(matDef, c.MaxDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.MaxOffense
					def = c.MaxDefense
					soldier = c.MaxFollower
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		mat1Atk, mat1Def, mat1Soldier := materialCard.Evo6Card(v)
		var mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank == 4 {
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["1"].Evo6Card(v)
		} else {
			mat2Atk, mat2Def, mat2Soldier = firstEvo.Evo6Card(v)
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, firstEvo.MaxOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, firstEvo.MaxDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, firstEvo.MaxFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, c.MaxOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, c.MaxDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, c.MaxFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Evo6CardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 6-card evo
func (c *Card) Evo6CardLvl1(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation6CardLvl1(v)
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.Evo6CardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Awakening 6 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.Evo6Card(v)
					atk = calculateEvoAccidentStat(matAtk, c.DefaultOffense)
					def = calculateEvoAccidentStat(matDef, c.DefaultDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.DefaultFollower)
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.DefaultOffense
					def = c.DefaultDefense
					soldier = c.DefaultFollower
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		mat1Atk, mat1Def, mat1Soldier := materialCard.Evo6Card(v)
		var mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank == 4 {
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["1"].Evo6Card(v)
		} else {
			mat2Atk, mat2Def, mat2Soldier = firstEvo.Evo6Card(v)
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, firstEvo.DefaultOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, firstEvo.DefaultDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, firstEvo.DefaultFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, c.DefaultOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, c.DefaultDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, c.DefaultFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Evo9Card calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 9-card evo
// this one should only be used for <= SR
func (c *Card) Evo9Card(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation9Card(v)
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.Evo9CardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Awakening 9 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.Evo9Card(v)
					atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
					def = calculateEvoAccidentStat(matDef, c.MaxDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.MaxOffense
					def = c.MaxDefense
					soldier = c.MaxFollower
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		var mat1Atk, mat1Def, mat1Soldier, mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Atk, mat1Def, mat1Soldier = materialCard.EvoPerfect(v)
			mat2Atk, mat2Def, mat2Soldier = mat1Atk, mat1Def, mat1Soldier
		} else if c.EvolutionRank == 3 {
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions(v)["2"].EvoStandard(v)
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["1"].EvoStandard(v)
		} else {
			// this would be the materials to get 4*
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions(v)["2"].Evo9Card(v)
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["3"].Evo9Card(v)
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, firstEvo.MaxOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, firstEvo.MaxDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, firstEvo.MaxFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, c.MaxOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, c.MaxDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, c.MaxFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Evo9CardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 9-card evo
// this one should only be used for <= SR
func (c *Card) Evo9CardLvl1(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation9CardLvl1(v)
		} else {
			mat := c.AwakensFrom(v)
			if mat != nil {
				// if this is an awakwening, calculate the max...
				matAtk, matDef, matSoldier := mat.Evo9CardLvl1(v)
				atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Awakening 9 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				// check for Evo Accident
				mat = c.EvoAccidentOf(v.Cards)
				if mat != nil {
					// calculate the transfered stats of the 2 material cards
					// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
					matAtk, matDef, matSoldier := mat.Evo9Card(v)
					atk = calculateEvoAccidentStat(matAtk, c.DefaultOffense)
					def = calculateEvoAccidentStat(matDef, c.DefaultDefense)
					soldier = calculateEvoAccidentStat(matSoldier, c.DefaultFollower)
				} else {
					// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
					atk = c.DefaultOffense
					def = c.DefaultDefense
					soldier = c.DefaultFollower
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions(v)["0"]
		var mat1Atk, mat1Def, mat1Soldier, mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Atk, mat1Def, mat1Soldier = materialCard.EvoPerfect(v)
			mat2Atk, mat2Def, mat2Soldier = mat1Atk, mat1Def, mat1Soldier
		} else if c.EvolutionRank == 3 {
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions(v)["2"].Evo9Card(v)
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["1"].Evo9Card(v)
		} else {
			// this would be the materials to get 4*
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions(v)["2"].Evo9Card(v)
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions(v)["3"].Evo9Card(v)
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, firstEvo.DefaultOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, firstEvo.DefaultDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, firstEvo.DefaultFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(mat1Atk, mat2Atk, c.DefaultOffense)
			def = c.calculateEvoStat(mat1Def, mat2Def, c.DefaultDefense)
			soldier = c.calculateEvoStat(mat1Soldier, mat2Soldier, c.DefaultFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// EvoPerfect calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 16-card evo
func (c *Card) EvoPerfect(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.AmalgamationPerfect(v)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1(v)
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
			os.Stdout.WriteString(fmt.Sprintf("Awakening Perfect card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matAtk, matDef, matSoldier := turnOver.EvoPerfect(v)
				atk = calculateEvoAccidentStat(matAtk, c.MaxOffense)
				def = calculateEvoAccidentStat(matDef, c.MaxDefense)
				soldier = calculateEvoAccidentStat(matSoldier, c.MaxFollower)
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				atk = c.MaxOffense
				def = c.MaxDefense
				soldier = c.MaxFollower
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoPerfect(v)
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions(v)["0"]
			atk = c.calculateEvoStat(matAtk, matAtk, firstEvo.MaxOffense)
			def = c.calculateEvoStat(matDef, matDef, firstEvo.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, matSoldier, firstEvo.MaxFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, matAtk, c.MaxOffense)
			def = c.calculateEvoStat(matDef, matDef, c.MaxDefense)
			soldier = c.calculateEvoStat(matSoldier, matSoldier, c.MaxFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// EvoPerfectLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 16-card evo
func (c *Card) EvoPerfectLvl1(v *VFile) (atk, def, soldier int) {
	materialCard := c.PrevEvo(v)
	rarity := c.CardRarity(v)
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation(v.Amalgamations) {
			// calculate the amalgamation stats here
			atk, def, soldier = c.AmalgamationPerfectLvl1(v)
		} else if c.AwakensFrom(v) != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom(v)
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1(v)
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
			os.Stdout.WriteString(fmt.Sprintf("Awakening Perfect card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf(v.Cards)
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matAtk, matDef, matSoldier := turnOver.EvoPerfect(v)
				atk = calculateEvoAccidentStat(matAtk, c.DefaultOffense)
				def = calculateEvoAccidentStat(matDef, c.DefaultDefense)
				soldier = calculateEvoAccidentStat(matSoldier, c.DefaultFollower)
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				atk = c.DefaultOffense
				def = c.DefaultDefense
				soldier = c.DefaultFollower
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoPerfect(v)
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions(v)["0"]
			atk = c.calculateEvoStat(matAtk, matAtk, firstEvo.DefaultOffense)
			def = c.calculateEvoStat(matDef, matDef, firstEvo.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, matSoldier, firstEvo.DefaultFollower)
		} else {
			// 1* cards use the result card stats
			atk = c.calculateEvoStat(matAtk, matAtk, c.DefaultOffense)
			def = c.calculateEvoStat(matDef, matDef, c.DefaultDefense)
			soldier = c.calculateEvoStat(matSoldier, matSoldier, c.DefaultFollower)
		}
	}
	if atk > rarity.LimtOffense {
		atk = rarity.LimtOffense
	}
	if def > rarity.LimtDefense {
		def = rarity.LimtDefense
	}
	if soldier > rarity.LimtMaxFollower {
		soldier = rarity.LimtMaxFollower
	}
	return
}

// Archwitch If this card was used as an AW, get the AW information.
// This can be used to get Likability information.
func (c *Card) Archwitch(v *VFile) *Archwitch {
	if c.archwitch == nil {
		for _, aw := range v.Archwitches {
			if c.ID == aw.CardMasterID {
				c.archwitch = &aw
				break
			}
		}
	}
	return c.archwitch
}

// EvoAccident If this card can produce an evolution accident, get the result card.
func (c *Card) EvoAccident(cards []Card) *Card {
	return CardScan(c.TransCardID, cards)
}

// EvoAccidentOf If this card is the result of an evo accident, get the source card.
func (c *Card) EvoAccidentOf(cards []Card) *Card {
	for key, val := range cards {
		if val.TransCardID == c.ID {
			return &(cards[key])
		}
	}
	return nil
}

// Amalgamations get any amalgamations for this card (material or result)
func (c *Card) Amalgamations(v *VFile) []Amalgamation {
	ret := make([]Amalgamation, 0)
	for _, a := range v.Amalgamations {
		if c.ID == a.FusionCardID ||
			c.ID == a.Material1 ||
			c.ID == a.Material2 ||
			c.ID == a.Material3 ||
			c.ID == a.Material4 {

			ret = append(ret, a)
		}
	}
	return ret
}

// AwakensTo Gets the card this card awakens to.
func (c *Card) AwakensTo(v *VFile) *Card {
	for _, val := range v.Awakenings {
		if c.ID == val.BaseCardID {
			return CardScan(val.ResultCardID, v.Cards)
		}
	}
	return nil
}

// AwakensFrom gets the source card of this awoken card
func (c *Card) AwakensFrom(v *VFile) *Card {
	for _, val := range v.Awakenings {
		if c.ID == val.ResultCardID {
			return CardScan(val.BaseCardID, v.Cards)
		}
	}
	return nil
}

// HasAmalgamation returns true if this card has an amalgamation
// (is used as a material)
func (c *Card) HasAmalgamation(a []Amalgamation) bool {
	for _, v := range a {
		if c.ID == v.Material1 ||
			c.ID == v.Material2 ||
			c.ID == v.Material3 ||
			c.ID == v.Material4 {
			return true
		}
	}
	return false
}

// IsAmalgamation returns true if this card has an amalgamation
// (is the result of amalgamating other material)
func (c *Card) IsAmalgamation(a []Amalgamation) bool {
	for _, v := range a {
		if c.ID == v.FusionCardID {
			return true
		}
	}
	return false
}

// Skill1 of the card
func (c *Card) Skill1(v *VFile) *Skill {
	if c.skill1 == nil && c.SkillID1 > 0 {
		c.skill1 = SkillScan(c.SkillID1, v.Skills)
	}
	return c.skill1
}

// Skill2 of the card
func (c *Card) Skill2(v *VFile) *Skill {
	if c.skill2 == nil && c.SkillID2 > 0 {
		c.skill2 = SkillScan(c.SkillID2, v.Skills)
	}
	return c.skill2
}

// Skill3 of the card
func (c *Card) Skill3(v *VFile) *Skill {
	if c.skill3 == nil && c.SkillID3 > 0 {
		c.skill3 = SkillScan(c.SkillID3, v.Skills)
	}
	return c.skill3
}

// SpecialSkill1 of the card
func (c *Card) SpecialSkill1(v *VFile) *Skill {
	if c.specialSkill1 == nil && c.SpecialSkillID1 > 0 {
		c.specialSkill1 = SkillScan(c.SpecialSkillID1, v.Skills)
	}
	return c.specialSkill1
}

// ThorSkill1 of the card
func (c *Card) ThorSkill1(v *VFile) *Skill {
	if c.thorSkill1 == nil && c.ThorSkillID1 > 0 {
		c.thorSkill1 = SkillScan(c.ThorSkillID1, v.Skills)
	}
	return c.thorSkill1
}

// CardScan searches for a card by ID
func CardScan(id int, cards []Card) *Card {
	if id <= 0 {
		return nil
	}
	l := len(cards)
	i := sort.Search(l, func(i int) bool { return cards[i].ID >= id })
	if i >= 0 && i < l && cards[i].ID == id {
		return &(cards[i])
	}
	return nil
}

// CardScanCharacter searches for a card by the character ID
func CardScanCharacter(charID int, cards []Card) *Card {
	if charID > 0 {
		for k, val := range cards {
			//return the first one we find.
			if val.CardCharaID == charID {
				return &cards[k]
			}
		}
	}
	return nil
}

// CardScanImage searches for a card by the card image number
func CardScanImage(cardID string, cards []Card) *Card {
	if cardID != "" {
		i, err := strconv.Atoi(cardID)
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

// Skill1Name returns the name of the first skill
func (c *Card) Skill1Name(v *VFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

// SkillMin returns the minimum skill level info
func (c *Card) SkillMin(v *VFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.SkillMin()
}

// SkillMax returns the maximum skill level info
func (c *Card) SkillMax(v *VFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.SkillMax()
}

// SkillProcs rturns the number of times a skill can activate.
// a negative number indicates infinite procs
func (c *Card) SkillProcs(v *VFile) string {
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

// SkillTarget gets the target scope of the skill
func (c *Card) SkillTarget(v *VFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.TargetScope()
}

// SkillTargetLogic gets the target logic for the skill
func (c *Card) SkillTargetLogic(v *VFile) string {
	s := c.Skill1(v)
	if s == nil {
		return ""
	}
	return s.TargetLogic()
}

// Skill2Name name of the second skill
func (c *Card) Skill2Name(v *VFile) string {
	s := c.Skill2(v)
	if s == nil {
		return ""
	}
	return s.Name
}

// Skill3Name name of the 3rd skill
func (c *Card) Skill3Name(v *VFile) string {
	s := c.Skill3(v)
	if s == nil {
		return ""
	}
	return s.Name
}

// SpecialSkill1Name name of the 1st special skill
func (c *Card) SpecialSkill1Name(v *VFile) string {
	s := c.SpecialSkill1(v)
	if s == nil {
		return ""
	}
	if s.Name != "" {
		return s.Name
	}
	i := strconv.Itoa(s.ID)
	return i
}

// ThorSkill1Name name of the 1st thor skill
func (c *Card) ThorSkill1Name(v *VFile) string {
	s := c.ThorSkill1(v)
	if s == nil {
		return ""
	}
	return s.Name
}

// Description of the character
func (c *Card) Description(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Description
}

// Friendship quote for the character
func (c *Card) Friendship(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Friendship
}

// Login quote for the character (not used on newer cards)
func (c *Card) Login(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Login
}

// Meet quote for the character
func (c *Card) Meet(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.Meet
}

// BattleStart quote for the character
func (c *Card) BattleStart(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.BattleStart
}

// BattleEnd quote for the character
func (c *Card) BattleEnd(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.BattleEnd
}

// FriendshipMax quote for the character
func (c *Card) FriendshipMax(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipMax
}

// FriendshipEvent quote for the character
func (c *Card) FriendshipEvent(v *VFile) string {
	ch := c.Character(v)
	if ch == nil {
		return ""
	}
	return ch.FriendshipEvent
}

// CardList helper interface for looking at lists of cards
type CardList []Card

// Earliest gets the ealiest released card from a list of cards. Determined by ID
func (d CardList) Earliest() *Card {
	var min *Card
	for _, card := range d {
		if min == nil || min.ID > card.ID {
			min = &card
		}
	}
	return min
}

// Latest gets the latest released card from a list of cards. Determined by ID
func (d CardList) Latest() *Card {
	var max *Card
	for _, card := range d {
		if max == nil || max.ID < card.ID {
			max = &card
		}
	}
	return max
}

func getAmalBaseCard(card *Card, v *VFile) *Card {
	if card.IsAmalgamation(v.Amalgamations) {
		os.Stdout.WriteString(fmt.Sprintf("Checking Amalgamation base for Card: %d, Name: %s, Evo: %d\n", card.ID, card.Name, card.EvolutionRank))
		for _, amal := range card.Amalgamations(v) {
			if card.ID == amal.FusionCardID {
				// material 1
				ac := CardScan(amal.Material1, v.Cards)
				if ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation(v.Amalgamations) {
						return getAmalBaseCard(ac, v)
					}
					return ac
				}
				// material 2
				ac = CardScan(amal.Material2, v.Cards)
				if ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation(v.Amalgamations) {
						return getAmalBaseCard(ac, v)
					}
					return ac
				}
				// material 3
				ac = CardScan(amal.Material3, v.Cards)
				if ac != nil && ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation(v.Amalgamations) {
						return getAmalBaseCard(ac, v)
					}
					return ac
				}
				// material 4
				ac = CardScan(amal.Material4, v.Cards)
				if ac != nil && ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation(v.Amalgamations) {
						return getAmalBaseCard(ac, v)
					}
					return ac
				}
			}
		}
	}
	return card
}

func checkEndCards(c *Card, v *VFile) (awakening, amalCard, amalAwakening *Card) {
	awakening = c.AwakensTo(v)
	// check for Amalgamation
	if c.HasAmalgamation(v.Amalgamations) {
		amals := c.Amalgamations(v)
		for _, amal := range amals {
			// get the result card
			tamalCard := CardScan(amal.FusionCardID, v.Cards)
			if tamalCard != nil && tamalCard.ID != c.ID {
				os.Stdout.WriteString(fmt.Sprintf("Found amalgamation: %d, Name: '%s', Evo: %d\n", tamalCard.ID, tamalCard.Name, tamalCard.EvolutionRank))
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

// GetEvolutions gets the evolutions for a card including Awakening and same character(by name) amalgamations
func (c *Card) GetEvolutions(v *VFile) map[string]*Card {
	if c._allEvos == nil {
		ret := make(map[string]*Card)

		// handle cards like Chimrey and Time Traveler (enemy)
		if c.CardCharaID < 1 {
			os.Stdout.WriteString(fmt.Sprintf("No character info Card: %d, Name: %s, Evo: %d\n", c.ID, c.Name, c.EvolutionRank))
			ret["0"] = c
			c._allEvos = ret
			return ret
		}

		c2 := c
		// check if this is an awoken card
		if c2.Rarity()[0] == 'G' {
			tmp := c2.AwakensFrom(v)
			if tmp == nil {
				ch := c2.Character(v)
				if ch != nil && ch.Cards(v)[0].Name == c2.Name {
					c2 = &(ch.Cards(v)[0])
				}
				// the name changed, so we'll keep this card
			} else {
				c2 = tmp
			}
		}

		// get earliest evo
		for tmp := c2.PrevEvo(v); tmp != nil; tmp = tmp.PrevEvo(v) {
			c2 = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))
		}

		// at this point we should have the first card in the evolution path
		c2 = getAmalBaseCard(c2, v)

		// get earliest evo (again...)
		for tmp := c2.PrevEvo(v); tmp != nil; tmp = tmp.PrevEvo(v) {
			c2 = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))
		}

		os.Stdout.WriteString(fmt.Sprintf("Base Card: %d, Name: '%s', Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))

		// populate the actual evos.

		for nextEvo := c2; nextEvo != nil; nextEvo = nextEvo.NextEvo(v) {
			os.Stdout.WriteString(fmt.Sprintf("Next Evo is Card: %d, Name: '%s', Evo: %d\n", nextEvo.ID, nextEvo.Name, nextEvo.EvolutionRank))
			if nextEvo.EvolutionRank <= 0 {
				evoRank := "0"
				if nextEvo.Rarity()[0] == 'H' {
					evoRank = "H"
				}
				ret[evoRank] = nextEvo
				if nextEvo.LastEvolutionRank < 0 {
					// check for awakening
					awakening, amalCard, amalAwakening := checkEndCards(nextEvo, v)
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
			} else if nextEvo.EvolutionRank == c2.LastEvolutionRank || nextEvo.Rarity()[0] == 'H' || nextEvo.LastEvolutionRank < 0 {
				ret["H"] = nextEvo
				// check for awakening
				awakening, amalCard, amalAwakening := checkEndCards(nextEvo, v)
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
			os.Stdout.WriteString(fmt.Sprintf("(%s: %d) ", key, card.ID))
		}
		os.Stdout.WriteString("\n")
		c._allEvos = ret
		return ret
	}
	return c._allEvos
}

// Elements of the cards.
var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}

// Rarity of the cards
var Rarity = [14]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR", "LR", "HLR", "GLR"}
