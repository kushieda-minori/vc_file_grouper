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
	ID                        int    `json:"_id"`                                      // card id
	CardNo                    int    `json:"card_no"`                                  // card number, matches to the image file
	CardCharaID               int    `json:"card_chara_id"`                            // card character id
	CardRareID                int    `json:"card_rare_id"`                             // rarity of the card
	CardTypeID                int    `json:"card_type_id"`                             // type of the card (Passion, Cool, Light, Dark)
	DeckCost                  int    `json:"deck_cost"`                                // unit cost
	LastEvolutionRank         int    `json:"last_evolution_rank"`                      // number of evolution statges available to the card
	EvolutionRank             int    `json:"evolution_rank"`                           // this card current evolution stage
	EvolutionCardID           int    `json:"evolution_card_id"`                        // id of the card that this card evolves into, -1 for no evolution
	TransCardID               int    `json:"trans_card_id"`                            // id of a possible turnover accident
	FollowerKindID            int    `json:"follower_kind_id"`                         // cost of the followers?
	DefaultFollower           int    `json:"default_follower"`                         // base soldiers
	MaxFollower               int    `json:"max_follower"`                             // max soldiers if evolved minimally
	DefaultOffense            int    `json:"default_offense"`                          // base ATK
	MaxOffense                int    `json:"max_offense"`                              // max ATK if evolved minimally
	DefaultDefense            int    `json:"default_defense"`                          // base DEF
	MaxDefense                int    `json:"max_defense"`                              // max DEF if evolved minimally
	SkillID1                  int    `json:"skill_id_1"`                               // First Skill
	SkillID2                  int    `json:"skill_id_2"`                               // second Skill
	SkillID3                  int    `json:"skill_id_3"`                               // third Skill (LR)
	SpecialSkillID1           int    `json:"special_skill_id_1"`                       // Awakened Burst type (GSR,GUR,GLR)
	ThorSkillID1              int    `json:"thor_skill_id_1"`                          // Temporary Thor skills used for AAW
	CustomSkillCost           int    `json:"custom_skill_cost_1"`                      // initial skill cost
	CustomSkillCostIncPattern int    `json:"custom_skill_cost_increment_pattern_id_1"` // ?
	MedalRate                 int    `json:"medal_rate"`                               // amount of medals can be traded for
	Price                     int    `json:"price"`                                    // amount of gold can be traded for
	StunRate                  int    `json:"stun_rate"`                                // ?
	IsClosed                  int    `json:"is_closed"`                                // is closed
	Name                      string `json:"name"`                                     // name from the strings file

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
	CardExpCoefficient     int `json:"card_exp_coefficient"`
	EvolutionCoefficient   int `json:"evolution_coefficient"`
	GuildbattleCoefficient int `json:"guildbattle_coefficient"`
	Order                  int `json:"order"`
	// Signature lowercase rarity name "n" "hn" etc
	Signature            string `json:"signature"`
	CardLevelCoefficient int    `json:"card_level_coefficient"`
	FragmentSlot         int    `json:"fragment_slot"`
	LimtOffense          int    `json:"limt_offense"`
	LimtDefense          int    `json:"limt_defense"`
	LimtMaxFollower      int    `json:"limt_max_follower"`
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

// MainRarity gets the main rarity of this card instead of the exact evo rarity
// (i.e. GUR => UR, HSR => SR)
func (c *Card) MainRarity() string {
	s := c.Rarity()
	l := len(s)
	switch l {
	case 1:
		// N, X, R
		return s
	case 2:
		// HN, HX, HR, SR, UR, LR
		return strings.TrimPrefix(s, "H")
	case 3:
		// HSR, GSR, HUR, GUR, HLR, GLR, XSR, XUR, XLR
		return strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(s, "H"), "G"), "X")
	default:
		// not a known rarity!
		return s
	}
}

// CardRarity with full rarity information
func (c *Card) CardRarity() *CardRarity {
	return CardRarityScan(c.CardRareID)
}

//CardRarityScan scans for a card rarity by id
func CardRarityScan(id int) *CardRarity {
	if id >= 0 {
		for idx, cr := range Data.CardRarities {
			if cr.ID == id {
				return &(Data.CardRarities[idx])
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

// IsRetired returns true if this card is no longer available because a newer card
// of the same character was released.
func (c *Card) IsRetired() bool {
	oldIDx := sort.SearchInts(retiredCards, c.ID)
	return oldIDx >= 0 && oldIDx < len(retiredCards) && retiredCards[oldIDx] == c.ID
}

// Character information of the card
func (c *Card) Character() *CardCharacter {
	if c.character == nil && c.CardCharaID > 0 {
		for k, val := range Data.CardCharacters {
			if val.ID == c.CardCharaID {
				c.character = &(Data.CardCharacters[k])
				break
			}
		}
	}
	return c.character
}

// NextEvo is the next evolution of this card, or nil if no further evolutions are possible.
// Amalgamations, Awakenings, and Rebirths may still be possible.
func (c *Card) NextEvo() *Card {
	if c.ID == c.EvolutionCardID {
		// bad data
		return nil
	}
	if c.nextEvo == nil {
		if c.CardCharaID <= 0 || c.EvolutionCardID <= 0 || c.Rarity()[0] == 'H' {
			return nil
		}

		var tmp *Card
		character := c.Character()
		for i, cd := range character.Cards() {
			if cd.ID == c.EvolutionCardID {
				tmp = &(character._cards[i])
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
func (c *Card) PrevEvo() *Card {
	if c.prevEvo == nil {
		// no charcter ID or already lowest evo rank
		if c.CardCharaID <= 0 || c.EvolutionRank <= 0 {
			return nil
		}

		var tmp *Card
		for i, cd := range c.Character().Cards() {
			if c.ID == cd.EvolutionCardID {
				tmp = &(c.Character()._cards[i])
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

// FirstEvo Gets the first evolution for the card
func (c *Card) FirstEvo() *Card {
	t := c
	for t.PrevEvo() != nil {
		t = t.PrevEvo()
	}
	return t
}

// LastEvo Gets the last evolution for the card
func (c *Card) LastEvo() *Card {
	t := c
	for t.NextEvo() != nil {
		t = t.NextEvo()
	}
	return t
}

// PossibleMixedEvo Checks if this card has a possible mixed evo:
// Amalgamation at evo[0] with a standard evo after it. This may issue
// false positives for cards that can only be obtained through
// amalgamation at evo[0] and there is no drop/RR
func (c *Card) PossibleMixedEvo() bool {
	if len(c._allEvos) == 1 && c.IsAmalgamation() {
		// if this is an amalgamation only card, we want to find the previous
		// cards and see if any of those is a "possible mix evos"
		for _, amal := range c.Amalgamations() {
			if c.ID == amal.FusionCardID {
				// check if each material is a possible mixed evo
				for _, amalCard := range amal.MaterialsOnly() {
					if amalCard.PossibleMixedEvo() {
						return true
					}
				}
			}
		}
		return false
	}
	firstEvo := c.GetEvolutions()["0"]
	secondEvo := c.GetEvolutions()["1"]
	if secondEvo == nil {
		secondEvo = c.GetEvolutions()["H"]
	}
	return firstEvo != nil && secondEvo != nil &&
		firstEvo.IsAmalgamation() &&
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
func (c *Card) AmalgamationStandard() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoStandard()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.EvoStandard()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) AmalgamationStandardLvl1() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoStandardLvl1()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.EvoStandard()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) Amalgamation6Card() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo6Card()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating 6-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.Evo6Card()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) Amalgamation6CardLvl1() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo6CardLvl1()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating 6-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.Evo6Card()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) Amalgamation9Card() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo9Card()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating 9-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.Evo9Card()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) Amalgamation9CardLvl1() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.Evo9CardLvl1()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating 9-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.Evo9Card()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) AmalgamationPerfect() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoPerfect()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.EvoPerfect()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) AmalgamationPerfectLvl1() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		return c.EvoPerfectLvl1()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.DefaultOffense
	def = c.DefaultDefense
	soldier = c.DefaultFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name))
		matAtk, matDef, matSoldier := mat.EvoPerfect()
		atk += int(float64(matAtk) * 0.08)
		def += int(float64(matDef) * 0.08)
		soldier += int(float64(matSoldier) * 0.08)
	}
	rarity := c.CardRarity()
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
func (c *Card) AmalgamationLRStaticLvl1() (atk, def, soldier int) {
	amalgs := c.Amalgamations()
	var myAmal *Amalgamation
	for _, a := range amalgs {
		if c.ID == a.FusionCardID {
			myAmal = &a
			break
		}
	}
	if myAmal == nil {
		if strings.HasSuffix(c.Rarity(), "LR") && c.Element() == "Special" {
			return c.EvoPerfectLvl1()
		}
		return c.EvoPerfect()
	}
	mats := myAmal.MaterialsOnly()
	atk = c.MaxOffense
	def = c.MaxDefense
	soldier = c.MaxFollower
	for _, mat := range mats {
		os.Stdout.WriteString(fmt.Sprintf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name))
		if strings.HasSuffix(mat.Rarity(), "LR") && mat.Element() == "Special" {
			matAtk, matDef, matSoldier := mat.AmalgamationLRStaticLvl1()
			// no 5% bonus for max level
			atk += int(float64(matAtk) * 0.03)
			def += int(float64(matDef) * 0.03)
			soldier += int(float64(matSoldier) * 0.03)
		} else {
			matAtk, matDef, matSoldier := mat.AmalgamationLRStaticLvl1()
			atk += int(float64(matAtk) * 0.08)
			def += int(float64(matDef) * 0.08)
			soldier += int(float64(matSoldier) * 0.08)
		}
	}
	rarity := c.CardRarity()
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
func (c *Card) EvoStandard() (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Standard Evo for %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoStandardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoStandardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					os.Stdout.WriteString(fmt.Sprintf("Awakening standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.EvoStandard()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoStandard()
		firstEvo := c.GetEvolutions()["0"]
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
func (c *Card) EvoStandardLvl1() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandardLvl1()
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoStandardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoStandardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					os.Stdout.WriteString(fmt.Sprintf("Awakening StandardLvl1 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.EvoStandard()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matAtk, matDef, matSoldier := materialCard.EvoStandard()
		firstEvo := c.GetEvolutions()["0"]
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
func (c *Card) EvoMixed() (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Mixed Evo for %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoMixedLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoMixedLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					os.Stdout.WriteString(fmt.Sprintf("Awakening Mixed card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.EvoStandard()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var matAtk, matDef, matSoldier int
		if c.EvolutionRank == 1 && c.PossibleMixedEvo() && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matAtk, matDef, matSoldier = materialCard.AmalgamationStandard()
		} else {
			matAtk, matDef, matSoldier = materialCard.EvoStandard()
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
func (c *Card) EvoMixedLvl1() (atk, def, soldier int) {
	os.Stdout.WriteString(fmt.Sprintf("Calculating Mixed Evo for Lvl1 %s: %d\n", c.Name, c.EvolutionRank))
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			os.Stdout.WriteString(fmt.Sprintf("Using collected stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// os.Stdout.WriteString(fmt.Sprintf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier))
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoMixedLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoMixedLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					os.Stdout.WriteString(fmt.Sprintf("Awakening Mixed card Lvl1 %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.EvoStandard()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var matAtk, matDef, matSoldier int
		if c.EvolutionRank == 1 && c.PossibleMixedEvo() && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matAtk, matDef, matSoldier = materialCard.AmalgamationStandard()
		} else {
			matAtk, matDef, matSoldier = materialCard.EvoStandard()
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
func (c *Card) Evo6Card() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation6Card()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.Evo6CardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo6CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					os.Stdout.WriteString(fmt.Sprintf("Awakening 6 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.Evo6Card()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		mat1Atk, mat1Def, mat1Soldier := materialCard.Evo6Card()
		var mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank == 4 {
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["1"].Evo6Card()
		} else {
			mat2Atk, mat2Def, mat2Soldier = firstEvo.Evo6Card()
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
func (c *Card) Evo6CardLvl1() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation6CardLvl1()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.Evo6CardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo6CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					os.Stdout.WriteString(fmt.Sprintf("Awakening 6 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.Evo6Card()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		mat1Atk, mat1Def, mat1Soldier := materialCard.Evo6Card()
		var mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank == 4 {
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["1"].Evo6Card()
		} else {
			mat2Atk, mat2Def, mat2Soldier = firstEvo.Evo6Card()
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
func (c *Card) Evo9Card() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation9Card()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.Evo9CardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo9CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					os.Stdout.WriteString(fmt.Sprintf("Awakening 9 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.Evo9Card()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var mat1Atk, mat1Def, mat1Soldier, mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Atk, mat1Def, mat1Soldier = materialCard.EvoPerfect()
			mat2Atk, mat2Def, mat2Soldier = mat1Atk, mat1Def, mat1Soldier
		} else if c.EvolutionRank == 3 {
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions()["2"].EvoStandard()
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["1"].EvoStandard()
		} else {
			// this would be the materials to get 4*
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions()["2"].Evo9Card()
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["3"].Evo9Card()
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
func (c *Card) Evo9CardLvl1() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.Amalgamation9CardLvl1()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.Evo9CardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo9CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					os.Stdout.WriteString(fmt.Sprintf("Awakening 9 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matAtk, matDef, matSoldier := mat.Evo9Card()
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
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var mat1Atk, mat1Def, mat1Soldier, mat2Atk, mat2Def, mat2Soldier int
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Atk, mat1Def, mat1Soldier = materialCard.EvoPerfect()
			mat2Atk, mat2Def, mat2Soldier = mat1Atk, mat1Def, mat1Soldier
		} else if c.EvolutionRank == 3 {
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions()["2"].Evo9Card()
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["1"].Evo9Card()
		} else {
			// this would be the materials to get 4*
			mat1Atk, mat1Def, mat1Soldier = c.GetEvolutions()["2"].Evo9Card()
			mat2Atk, mat2Def, mat2Soldier = c.GetEvolutions()["3"].Evo9Card()
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
func (c *Card) EvoPerfect() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.AmalgamationPerfect()
		} else if c.RebirthsFrom() != nil {
			mat := c.RebirthsFrom()
			// if this is an rebirth, calculate the max...
			awakenMat := mat.AwakensFrom()
			matAtk, matDef, matSoldier := awakenMat.EvoPerfectLvl1()
			aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
			atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
			os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1()
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
			os.Stdout.WriteString(fmt.Sprintf("Awakening Perfect card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf()
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matAtk, matDef, matSoldier := turnOver.EvoPerfect()
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
		matAtk, matDef, matSoldier := materialCard.EvoPerfect()
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions()["0"]
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
func (c *Card) EvoPerfectLvl1() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// os.Stderr.WriteString(fmt.Sprintf("No previous evo found for card %v\n", c.ID))
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			atk, def, soldier = c.AmalgamationPerfectLvl1()
		} else if c.RebirthsFrom() != nil {
			mat := c.RebirthsFrom()
			// if this is an rebirth, calculate the max...
			awakenMat := mat.AwakensFrom()
			matAtk, matDef, matSoldier := awakenMat.EvoPerfectLvl1()
			aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
			atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
			os.Stdout.WriteString(fmt.Sprintf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1()
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
			os.Stdout.WriteString(fmt.Sprintf("Awakening Perfect card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier))
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf()
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matAtk, matDef, matSoldier := turnOver.EvoPerfect()
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
		matAtk, matDef, matSoldier := materialCard.EvoPerfect()
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions()["0"]
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
func (c *Card) Archwitch() *Archwitch {
	if c.archwitch == nil {
		for _, aw := range Data.Archwitches {
			if c.ID == aw.CardMasterID {
				c.archwitch = &aw
				break
			}
		}
	}
	return c.archwitch
}

// EvoAccident If this card can produce an evolution accident, get the result card.
func (c *Card) EvoAccident() *Card {
	return CardScan(c.TransCardID)
}

// EvoAccidentOf If this card is the result of an evo accident, get the source card.
func (c *Card) EvoAccidentOf() *Card {
	for key, val := range Data.Cards {
		if val.TransCardID == c.ID {
			return &(Data.Cards[key])
		}
	}
	return nil
}

// Amalgamations get any amalgamations for this card (material or result)
func (c *Card) Amalgamations() []Amalgamation {
	ret := make([]Amalgamation, 0)
	for _, a := range Data.Amalgamations {
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

// AwakensTo Gets the card this card awakens to. Call LastEvo first if
// you want the awoken card and aren't sure if this is the direct material.
func (c *Card) AwakensTo() *Card {
	for _, val := range Data.Awakenings {
		if c.ID == val.BaseCardID {
			return CardScan(val.ResultCardID)
		}
	}
	return nil
}

// AwakensFrom gets the source card of this awoken card
func (c *Card) AwakensFrom() *Card {
	for _, val := range Data.Awakenings {
		if c.ID == val.ResultCardID {
			return CardScan(val.BaseCardID)
		}
	}
	return nil
}

// HasRebirth Gets the card this card rebirths to.
func (c *Card) HasRebirth() bool {
	for _, val := range Data.Rebirths {
		if c.ID == val.BaseCardID {
			return true
		}
	}
	return false
}

// RebirthsTo Gets the card this card rebirths to. Call LastEvo().AwakensTo()
// first if you want the rebith card and aren't sure if this is the direct
// material.
func (c *Card) RebirthsTo() *Card {
	for _, val := range Data.Rebirths {
		if c.ID == val.BaseCardID {
			return CardScan(val.ResultCardID)
		}
	}
	return nil
}

// RebirthsFrom gets the source card of this rebirth card
func (c *Card) RebirthsFrom() *Card {
	for _, val := range Data.Rebirths {
		if c.ID == val.ResultCardID {
			return CardScan(val.BaseCardID)
		}
	}
	return nil
}

// HasAmalgamation returns true if this card has an amalgamation
// (is used as a material)
func (c *Card) HasAmalgamation() bool {
	for _, a := range Data.Amalgamations {
		if c.ID == a.Material1 ||
			c.ID == a.Material2 ||
			c.ID == a.Material3 ||
			c.ID == a.Material4 {
			return true
		}
	}
	return false
}

// IsAmalgamation returns true if this card has an amalgamation
// (is the result of amalgamating other material)
func (c *Card) IsAmalgamation() bool {
	for _, a := range Data.Amalgamations {
		if c.ID == a.FusionCardID {
			return true
		}
	}
	return false
}

// Skill1 of the card
func (c *Card) Skill1() *Skill {
	if c.skill1 == nil && c.SkillID1 > 0 {
		c.skill1 = SkillScan(c.SkillID1)
	}
	return c.skill1
}

// Skill2 of the card
func (c *Card) Skill2() *Skill {
	if c.skill2 == nil && c.SkillID2 > 0 {
		c.skill2 = SkillScan(c.SkillID2)
	}
	return c.skill2
}

// Skill3 of the card
func (c *Card) Skill3() *Skill {
	if c.skill3 == nil && c.SkillID3 > 0 {
		c.skill3 = SkillScan(c.SkillID3)
	}
	return c.skill3
}

// SpecialSkill1 of the card
func (c *Card) SpecialSkill1() *Skill {
	if c.specialSkill1 == nil && c.SpecialSkillID1 > 0 {
		c.specialSkill1 = SkillScan(c.SpecialSkillID1)
	}
	return c.specialSkill1
}

// ThorSkill1 of the card
func (c *Card) ThorSkill1() *Skill {
	if c.thorSkill1 == nil && c.ThorSkillID1 > 0 {
		c.thorSkill1 = SkillScan(c.ThorSkillID1)
	}
	return c.thorSkill1
}

// CardScan searches for a card by ID
func CardScan(id int) *Card {
	if id <= 0 {
		return nil
	}
	l := len(Data.Cards)
	i := sort.Search(l, func(i int) bool { return Data.Cards[i].ID >= id })
	if i >= 0 && i < l && Data.Cards[i].ID == id {
		return &(Data.Cards[i])
	}
	return nil
}

// CardScanCharacter searches for a card by the character ID
func CardScanCharacter(charID int) *Card {
	if charID > 0 {
		for k, val := range Data.Cards {
			//return the first one we find.
			if val.CardCharaID == charID {
				return &(Data.Cards[k])
			}
		}
	}
	return nil
}

// CardScanImage searches for a card by the card image number
func CardScanImage(cardID string) *Card {
	if cardID != "" {
		i, err := strconv.Atoi(cardID)
		if err != nil {
			return nil
		}
		for k, val := range Data.Cards {
			if val.CardNo == i {
				return &(Data.Cards[k])
			}
		}
	}
	return nil
}

// Skill1Name returns the name of the first skill
func (c *Card) Skill1Name() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	return s.Name
}

// SkillMin returns the minimum skill level info
func (c *Card) SkillMin() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	return s.SkillMin()
}

// SkillMax returns the maximum skill level info
func (c *Card) SkillMax() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	return s.SkillMax()
}

// SkillProcs rturns the number of times a skill can activate.
// a negative number indicates infinite procs
func (c *Card) SkillProcs() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	// battle start skills seem to have random Max Count values. Force it to 1
	// since they can only proc once anyway
	if strings.Contains(strings.ToLower(c.SkillMin()), "battle start") {
		return "1"
	}
	return strconv.Itoa(s.MaxCount)
}

// SkillTarget gets the target scope of the skill
func (c *Card) SkillTarget() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	return s.TargetScope()
}

// SkillTargetLogic gets the target logic for the skill
func (c *Card) SkillTargetLogic() string {
	s := c.Skill1()
	if s == nil {
		return ""
	}
	return s.TargetLogic()
}

// Skill2Name name of the second skill
func (c *Card) Skill2Name() string {
	s := c.Skill2()
	if s == nil {
		return ""
	}
	return s.Name
}

// Skill3Name name of the 3rd skill
func (c *Card) Skill3Name() string {
	s := c.Skill3()
	if s == nil {
		return ""
	}
	return s.Name
}

// SpecialSkill1Name name of the 1st special skill
func (c *Card) SpecialSkill1Name() string {
	s := c.SpecialSkill1()
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
func (c *Card) ThorSkill1Name() string {
	s := c.ThorSkill1()
	if s == nil {
		return ""
	}
	return s.Name
}

// Description of the character
func (c *Card) Description() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.Description
}

// Friendship quote for the character
func (c *Card) Friendship() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.Friendship
}

// Login quote for the character (not used on newer cards)
func (c *Card) Login() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.Login
}

// Meet quote for the character
func (c *Card) Meet() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.Meet
}

// BattleStart quote for the character
func (c *Card) BattleStart() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.BattleStart
}

// BattleEnd quote for the character
func (c *Card) BattleEnd() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.BattleEnd
}

// FriendshipMax quote for the character
func (c *Card) FriendshipMax() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.FriendshipMax
}

// FriendshipEvent quote for the character
func (c *Card) FriendshipEvent() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.FriendshipEvent
}

// RebirthEvent quote for the character
func (c *Card) RebirthEvent() string {
	ch := c.Character()
	if ch == nil {
		return ""
	}
	return ch.Rebirth
}

// CardList helper interface for looking at lists of cards
type CardList []Card

// Earliest gets the ealiest released card from a list of cards. Determined by ID
func (d CardList) Earliest() (min *Card) {
	for idx, card := range d {
		if min == nil || min.ID > card.ID {
			// os.Stdout.WriteString(fmt.Sprintf("'Earliest' Card: %d, Name: %s\n", card.ID, card.Name))
			min = &(d[idx])
		}
	}
	// if min != nil {
	// os.Stdout.WriteString(fmt.Sprintf("-Earliest Card: %d, Name: %s\n", min.ID, min.Name))
	// }
	return
}

// Latest gets the latest released card from a list of cards. Determined by ID
func (d CardList) Latest() (max *Card) {
	for idx, card := range d {
		if max == nil || max.ID < card.ID {
			// os.Stdout.WriteString(fmt.Sprintf("'Latest' Card: %d, Name: %s\n", card.ID, card.Name))
			max = &(d[idx])
		}
	}
	// if max != nil {
	// os.Stdout.WriteString(fmt.Sprintf("-Latest Card: %d, Name: %s\n", max.ID, max.Name))
	// }
	return
}

func getAmalBaseCard(card *Card) *Card {
	if card.IsAmalgamation() {
		os.Stdout.WriteString(fmt.Sprintf("Checking Amalgamation base for Card: %d, Name: %s, Evo: %d\n", card.ID, card.Name, card.EvolutionRank))
		for _, amal := range card.Amalgamations() {
			if card.ID == amal.FusionCardID {
				// material 1
				ac := CardScan(amal.Material1)
				if ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation() {
						return getAmalBaseCard(ac)
					}
					return ac
				}
				// material 2
				ac = CardScan(amal.Material2)
				if ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation() {
						return getAmalBaseCard(ac)
					}
					return ac
				}
				// material 3
				ac = CardScan(amal.Material3)
				if ac != nil && ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation() {
						return getAmalBaseCard(ac)
					}
					return ac
				}
				// material 4
				ac = CardScan(amal.Material4)
				if ac != nil && ac.ID != card.ID && ac.Name == card.Name {
					if ac.IsAmalgamation() {
						return getAmalBaseCard(ac)
					}
					return ac
				}
			}
		}
	}
	return card
}

func checkEndCards(c *Card) (awakening, amalCard, amalAwakening, rebirth, rebirthAmal *Card) {
	awakening = c.AwakensTo()
	rebirth = c.RebirthsTo()
	if rebirth == nil && awakening != nil {
		rebirth = awakening.RebirthsTo()
	}
	// check for Amalgamation
	if c.HasAmalgamation() {
		amals := c.Amalgamations()
		for _, amal := range amals {
			// get the result card
			tamalCard := CardScan(amal.FusionCardID)
			if tamalCard != nil && tamalCard.ID != c.ID {
				os.Stdout.WriteString(fmt.Sprintf("Found amalgamation: %d, Name: '%s', Evo: %d\n", tamalCard.ID, tamalCard.Name, tamalCard.EvolutionRank))
				if tamalCard.Name == c.Name {
					amalCard = tamalCard
					// check for amal awakening
					amalAwakening = amalCard.AwakensTo()
					if amalAwakening != nil {
						rebirthAmal = amalAwakening.RebirthsTo()
					} else {
						rebirthAmal = amalCard.RebirthsTo()
					}
					return // awakening, amalCard, amalAwakening, rebirth, rebirthAmal
				}
			}
		}
	}
	return // awakening, amalCard, amalAwakening, rebirth, rebirthAmal
}

// GetEvoImageName gets the nice name of the image for this card's evolution for use on the wiki
func (c *Card) GetEvoImageName(isIcon bool) string {
	evos := c.GetEvolutions()
	thisKey := ""
	for k, e := range evos {
		if e.ID == c.ID {
			thisKey = k
			break
		}
	}
	fileName := c.Name
	if fileName == "" {
		fileName = c.Character().FirstEvoCard().Image()
	}
	if thisKey == "0" {
		if c.Rarity()[0] == 'G' {
			if isIcon {
				return fileName + "_G"
			}
			return fileName + "_H"
		}
		if c.Rarity()[0] == 'H' {
			return fileName + "_H"
		}
		return fileName
	}
	if !isIcon {
		if len(evos) == 1 && thisKey == "G" {
			return fileName + "_H"
		}
		if thisKey == "A" {
			return fileName + "_H"
		}
	}
	if thisKey == "" {
		return fileName
	}
	return fileName + "_" + thisKey
}

// GetEvolutions gets the evolutions for a card including Awakening and same character(by name) amalgamations
func (c *Card) GetEvolutions() map[string]*Card {
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
		// check if this is an rebirth card
		if len(c2.Rarity()) > 2 && c2.Rarity()[0] == 'X' {
			tmp := c2.RebirthsFrom()
			if tmp == nil {
				ch := c2.Character()
				if ch != nil && ch.Cards()[0].Name == c2.Name {
					c2 = &(ch.Cards()[0])
				}
				// the name changed, so we'll keep this card
			} else {
				c2 = tmp
			}
		}
		// check if this is an awoken card
		if c2.Rarity()[0] == 'G' {
			tmp := c2.AwakensFrom()
			if tmp == nil {
				ch := c2.Character()
				if ch != nil && ch.Cards()[0].Name == c2.Name {
					c2 = &(ch.Cards()[0])
				}
				// the name changed, so we'll keep this card
			} else {
				c2 = tmp
			}
		}

		// get earliest evo
		for tmp := c2.PrevEvo(); tmp != nil; tmp = tmp.PrevEvo() {
			c2 = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))
		}

		// at this point we should have the first card in the evolution path
		c2 = getAmalBaseCard(c2)

		// get earliest evo (again...)
		for tmp := c2.PrevEvo(); tmp != nil; tmp = tmp.PrevEvo() {
			c2 = tmp
			os.Stdout.WriteString(fmt.Sprintf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))
		}

		os.Stdout.WriteString(fmt.Sprintf("Base Card: %d, Name: '%s', Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank))

		// populate the actual evos.

		for nextEvo := c2; nextEvo != nil; nextEvo = nextEvo.NextEvo() {
			os.Stdout.WriteString(fmt.Sprintf("Next Evo is Card: %d, Name: '%s', Evo: %d\n", nextEvo.ID, nextEvo.Name, nextEvo.EvolutionRank))
			if nextEvo.EvolutionRank <= 0 {
				evoRank := "0"
				if nextEvo.Rarity()[0] == 'H' {
					evoRank = "H"
				}
				ret[evoRank] = nextEvo
				if nextEvo.LastEvolutionRank < 0 {
					// check for awakening
					awakening, amalCard, amalAwakening, rebirth, rebirthAmal := checkEndCards(nextEvo)
					if awakening != nil {
						ret["G"] = awakening
					}
					if amalCard != nil {
						ret["A"] = amalCard
					}
					if amalAwakening != nil {
						ret["GA"] = amalAwakening
					}
					if rebirth != nil {
						ret["X"] = rebirth
						if rebirthAmal != nil {
							ret["XA"] = rebirthAmal
						}
					} else if rebirthAmal != nil {
						ret["X"] = rebirthAmal
					}
				}
			} else if nextEvo.Rarity()[0] == 'G' {
				// for some reason we hit a G during Evo traversal. Probably a G originating
				// from amalgamation
				ret["G"] = nextEvo
			} else if nextEvo.EvolutionRank == c2.LastEvolutionRank || nextEvo.Rarity()[0] == 'H' || nextEvo.LastEvolutionRank < 0 {
				ret["H"] = nextEvo
				// check for awakening
				awakening, amalCard, amalAwakening, rebirth, rebirthAmal := checkEndCards(nextEvo)
				if awakening != nil {
					ret["G"] = awakening
				}
				if amalCard != nil {
					ret["A"] = amalCard
				}
				if amalAwakening != nil {
					ret["GA"] = amalAwakening
				}
				if rebirth != nil {
					ret["X"] = rebirth
					if rebirthAmal != nil {
						ret["XA"] = rebirthAmal
					}
				} else if rebirthAmal != nil {
					ret["X"] = rebirthAmal
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

// GetEvolutionCards same as GetEvolutions, but only returns the cards
func (c *Card) GetEvolutionCards() CardList {
	evos := c.GetEvolutions()
	cards := make([]Card, 0, len(evos))
	//os.Stdout.WriteString(fmt.Sprintf("Card: %d, Name: %s, Evos: %d\n", c.ID, c.Name, len(evos)))
	for _, c := range evos {
		if c == nil {
			continue
		}
		cards = append(cards, *c)
	}
	//os.Stdout.WriteString(fmt.Sprintf("Cards: %d\n", len(cards)))
	return CardList(cards)
}

// EvoOrder order of evolutions in the map.
var EvoOrder = [10]string{"0", "1", "2", "3", "H", "A", "G", "GA", "X", "XA"}

// Elements of the cards.
var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}

// Rarity of the cards
var Rarity = [17]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR", "LR", "HLR", "GLR", "XSR", "XUR", "XLR"}
