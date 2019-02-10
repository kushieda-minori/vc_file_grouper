package vc

import (
	"log"
	"strings"
)

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
		log.Printf("Calculating Standard Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating Standard Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating 6-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating 6-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating 9-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating 9-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating Perfect Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
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
		log.Printf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
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
	log.Printf("Calculating Standard Evo for %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			log.Printf("Using collected stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoStandardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoStandardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					log.Printf("Awakening standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
						log.Printf("Using Evo Accident stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						atk = c.MaxOffense
						def = c.MaxDefense
						soldier = c.MaxFollower
						log.Printf("Using base Max stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
		log.Printf("Using Evo stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
	log.Printf("Final stats for Standard %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
	return
}

// EvoStandardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 5-card evo
func (c *Card) EvoStandardLvl1() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			log.Printf("Using collected stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandardLvl1()
			// log.Printf("Using Amalgamation stats for StandardLvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoStandardLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoStandardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					log.Printf("Awakening StandardLvl1 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
	log.Printf("Calculating Mixed Evo for %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.MaxOffense, c.MaxDefense, c.MaxFollower
			log.Printf("Using collected stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoMixedLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, false)
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoMixedLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					log.Printf("Awakening Mixed card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
						log.Printf("Using Evo Accident stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						atk = c.MaxOffense
						def = c.MaxDefense
						soldier = c.MaxFollower
						log.Printf("Using base Max stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
		log.Printf("Using Evo stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
	log.Printf("Final stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
	return
}

// EvoMixedLvl1 calculates the standard evolution for cards with a evo[0] amalgamation it calculates
// the evo[1] using 1 amalgamated, and one as a "drop"
func (c *Card) EvoMixedLvl1() (atk, def, soldier int) {
	log.Printf("Calculating Mixed Evo for Lvl1 %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			atk, def, soldier = c.DefaultOffense, c.DefaultDefense, c.DefaultFollower
			log.Printf("Using collected stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
			// calculate the amalgamation stats here
			// atk, def, soldier = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Mixed %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				awakenMat := mat.AwakensFrom()
				matAtk, matDef, matSoldier := awakenMat.EvoMixedLvl1()
				aatk, adef, asoldier := mat.calculateAwakeningStat(awakenMat, matAtk, matDef, matSoldier, true)
				atk, def, soldier = c.calculateAwakeningStat(mat, aatk, adef, asoldier, true)
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.EvoMixedLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					log.Printf("Awakening Mixed card Lvl1 %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
						log.Printf("Using Evo Accident stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						atk = c.DefaultOffense
						def = c.DefaultDefense
						soldier = c.DefaultFollower
						log.Printf("Using base Max stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
		log.Printf("Using Evo stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
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
	log.Printf("Final stats for Mixed Lvl1 %s: %d (%d, %d, %d)\n", c.Name, c.EvolutionRank, atk, def, soldier)
	return
}

// Evo6Card calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 6-card evo
func (c *Card) Evo6Card() (atk, def, soldier int) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo6CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					log.Printf("Awakening 6 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo6CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					log.Printf("Awakening 6 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo9CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
					log.Printf("Awakening 9 card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
				log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
					mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matAtk, matDef, matSoldier := mat.Evo9CardLvl1()
					atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
					log.Printf("Awakening 9 card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
						mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
			log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1()
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, false)
			log.Printf("Awakening Perfect card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
		// log.Printf("No previous evo found for card %v\n", c.ID)
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
			log.Printf("Rebirth standard card %d (%d, %d, %d) -> %d (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matAtk, matDef, matSoldier := mat.EvoPerfectLvl1()
			atk, def, soldier = c.calculateAwakeningStat(mat, matAtk, matDef, matSoldier, true)
			log.Printf("Awakening Perfect card %d (%d, %d, %d) -> %d lvl1 (%d, %d, %d)\n",
				mat.ID, matAtk, matDef, matSoldier, c.ID, atk, def, soldier)
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
