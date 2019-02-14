package vc

import (
	"fmt"
	"log"
	"strings"
)

// Stats stats for a card
type Stats struct {
	Attack   int
	Defense  int
	Soldiers int
}

func (s Stats) String() string {
	return fmt.Sprintf("%d, %d, %d", s.Attack, s.Defense, s.Soldiers)
}

// Add adds this stat set to another and returns the result.
func (s Stats) Add(o Stats) Stats {
	return Stats{
		Attack:   s.Attack + o.Attack,
		Defense:  s.Defense + o.Defense,
		Soldiers: s.Soldiers + o.Soldiers,
	}
}

// Subtract Subtracts another stat set from this one and returns the result.
func (s Stats) Subtract(o Stats) Stats {
	return Stats{
		Attack:   s.Attack - o.Attack,
		Defense:  s.Defense - o.Defense,
		Soldiers: s.Soldiers - o.Soldiers,
	}
}

// Subtract Subtracts another stat set from this one and returns the result.
func (s Stats) Multiply(m float64) Stats {
	return Stats{
		Attack:   int(float64(s.Attack) * m),
		Defense:  int(float64(s.Defense) * m),
		Soldiers: int(float64(s.Soldiers) * m),
	}
}

func (s *Stats) ensureMaxCap(rarity *CardRarity) {
	if s.Attack > rarity.LimtOffense {
		s.Attack = rarity.LimtOffense
	}
	if s.Defense > rarity.LimtDefense {
		s.Defense = rarity.LimtDefense
	}
	if s.Soldiers > rarity.LimtMaxFollower {
		s.Soldiers = rarity.LimtMaxFollower
	}
}

func (s *Stats) applyAmal(material Stats) {
	*s = material.Multiply(0.08).Add(*s)
}

func (s *Stats) applyAmalLvl1(material Stats) {
	*s = material.Multiply(0.03).Add(*s)
}

// calculateEvoAccidentStat calculated the stats if an evo accident happens
func (s *Stats) calculateEvoAccidentStat(materialStat, baseStat Stats) {
	*s = materialStat.Multiply(0.15).Multiply(2).Add(baseStat)
}

func (s *Stats) calculateTransferRate(transferRate float64, baseStat Stats) {
	*s = baseStat.Multiply(transferRate)
}

func (s *Stats) applyTransferRate(transferRate float64, material1Stat, material2Stat Stats) {
	*s = material1Stat.Multiply(transferRate).Add(material2Stat.Multiply(transferRate)).Add(*s)
}

func maxStats(c *Card) Stats {
	return Stats{
		Attack:   c.MaxOffense,
		Defense:  c.MaxDefense,
		Soldiers: c.MaxFollower,
	}
}
func baseStats(c *Card) Stats {
	return Stats{
		Attack:   c.DefaultOffense,
		Defense:  c.DefaultDefense,
		Soldiers: c.DefaultFollower,
	}
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
func (c *Card) calculateEvoStats(material1Stat, material2Stat, resultMax Stats) (ret Stats) {
	if c.EvolutionRank == c.LastEvolutionRank && c.EvolutionRank == 1 {
		// single evo gets no final stage bonus
		ret = resultMax
	} else if c.EvolutionRank == c.LastEvolutionRank {
		// 4* evo bonus for final stage
		if c.CardCharaID == 250 || c.CardCharaID == 315 {
			// queen of ice, strategist
			ret.calculateTransferRate(1.209, resultMax)
		} else {
			// all other N-SR (UR and LR? ATM, there are no UR 4*)
			ret.calculateTransferRate(1.1, resultMax)
		}
	} else {
		//4* evo bonus for intermediate stage
		if c.CardCharaID == 250 || c.CardCharaID == 315 {
			// queen of ice, strategist
			ret.calculateTransferRate(1.155, resultMax)
		} else {
			// all other N-SR (UR and LR? ATM, there are no UR 4*)
			ret.calculateTransferRate(1.05, resultMax)
		}
	}
	ret.applyTransferRate(0.15, material1Stat, material2Stat)

	return
}

// calculateLrAmalStat calculates amalgamation where one of the material is not max stat.
// this is generaly used for GUR+LR=LR or GLR+LR=GLR
func calculateLrAmalStat(sourceMatStat, lrMatStat, resultMax int) (ret int) {
	ret = resultMax
	ret += int(float64(sourceMatStat)*0.08) + int(float64(lrMatStat)*0.03)
	return
}

// calculateAwakeningStat calculated the stats after awakening
func (c *Card) calculateAwakeningStat(materialStatGain Stats, atLevel1 bool) (stats Stats) {
	// Awakening calculation thanks to Elle (https://docs.google.com/spreadsheets/d/1CT41xSuHyibfDSHQOyON4DkCPRlwW2WCy_ag6Pb76z4/edit?usp=sharing)
	if atLevel1 {
		stats = baseStats(c).Add(materialStatGain)
	} else {
		stats = maxStats(c).Add(materialStatGain)
	}
	return
}

// AmalgamationStandard calculated the stats if the materials have been evo'd using
// standard processes (4* 5 card evos).
func (c *Card) AmalgamationStandard() (stats Stats) {
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
	stats = maxStats(c)
	for _, mat := range mats {
		log.Printf("Calculating Standard Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.EvoStandard()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// AmalgamationStandardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 5 card evos).
func (c *Card) AmalgamationStandardLvl1() (stats Stats) {
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
	stats = baseStats(c)
	for _, mat := range mats {
		log.Printf("Calculating Standard Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.EvoStandard()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// Amalgamation6Card calculated the stats if the materials have been evo'd using
// standard processes (4* 6 card evos).
func (c *Card) Amalgamation6Card() (stats Stats) {
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
	stats = maxStats(c)
	for _, mat := range mats {
		log.Printf("Calculating 6-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.Evo6Card()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// Amalgamation6CardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 6 card evos).
func (c *Card) Amalgamation6CardLvl1() (stats Stats) {
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
	stats = baseStats(c)
	for _, mat := range mats {
		log.Printf("Calculating 6-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.Evo6Card()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// Amalgamation9Card calculated the stats if the materials have been evo'd using
// standard processes (4* 9 card evos).
func (c *Card) Amalgamation9Card() (stats Stats) {
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
	stats = maxStats(c)
	for _, mat := range mats {
		log.Printf("Calculating 9-card Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.Evo9Card()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// Amalgamation9CardLvl1 calculated the stats at level 1 if the materials have been evo'd using
// standard processes (4* 9 card evos).
func (c *Card) Amalgamation9CardLvl1() (stats Stats) {
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
	stats = baseStats(c)
	for _, mat := range mats {
		log.Printf("Calculating 9-card Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.Evo9Card()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// AmalgamationPerfect calculates the amalgamation stats if the material cards
// have all been evo'd perfectly (4* 16-card evos)
func (c *Card) AmalgamationPerfect() (stats Stats) {
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
	stats = maxStats(c)
	for _, mat := range mats {
		log.Printf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.EvoPerfect()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// AmalgamationPerfectLvl1 calculates the amalgamation stats at level 1 if the material cards
// have all been evo'd perfectly (4* 16-card evos)
func (c *Card) AmalgamationPerfectLvl1() (stats Stats) {
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
	stats = baseStats(c)
	for _, mat := range mats {
		log.Printf("Calculating Perfect Evo for lvl1 %s Amal Mat: %s\n", c.Name, mat.Name)
		matStats := mat.EvoPerfect()
		stats.applyAmal(matStats)
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// AmalgamationLRStaticLvl1 calculates the amalgamation stats if the material cards
// have all been evo'd perfectly (4* 16-card evos) with the exception of LR materials
// that have unchanging stats (i.e. 9999 / 9999)
func (c *Card) AmalgamationLRStaticLvl1() (stats Stats) {
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
	stats = maxStats(c)
	for _, mat := range mats {
		log.Printf("Calculating Perfect Evo for %s Amal Mat: %s\n", c.Name, mat.Name)
		if strings.HasSuffix(mat.Rarity(), "LR") && mat.Element() == "Special" {
			matStats := mat.AmalgamationLRStaticLvl1()
			// no 5% bonus for max level
			stats.applyAmalLvl1(matStats)
		} else {
			matStats := mat.AmalgamationLRStaticLvl1()
			stats.applyAmal(matStats)
		}
	}
	stats.ensureMaxCap(c.CardRarity())
	return
}

// EvoStandard calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 5-card evo
func (c *Card) EvoStandard() (stats Stats) {
	log.Printf("Calculating Standard Evo for %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			stats = maxStats(c)
			log.Printf("Using collected stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
			// calculate the amalgamation stats here
			// stats.Attack, stats.Defense, stats.Soldiers = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.EvoStandardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, false)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.EvoStandardLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, false)
					log.Printf("Awakening standard card %d (%s) -> %d (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.EvoStandard()
						stats.calculateEvoAccidentStat(matStats, maxStats(c))
						log.Printf("Using Evo Accident stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = maxStats(c)
						log.Printf("Using base Max stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matStats := materialCard.EvoStandard()
		firstEvo := c.GetEvolutions()["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), maxStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), maxStats(c))
		}
		log.Printf("Using Evo stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	}
	stats.ensureMaxCap(rarity)
	log.Printf("Final stats for Standard %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	return
}

// EvoStandardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 5-card evo
func (c *Card) EvoStandardLvl1() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			stats = baseStats(c)
			log.Printf("Using collected stats for StandardLvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
			// calculate the amalgamation stats here
			// stats = c.AmalgamationStandardLvl1()
			// log.Printf("Using Amalgamation stats for StandardLvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.EvoStandardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, true)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.EvoStandardLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, true)
					log.Printf("Awakening StandardLvl1 card %d (%s) -> %d lvl1 (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.EvoStandard()
						stats.calculateEvoAccidentStat(matStats, maxStats(c))
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = baseStats(c)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matStats := materialCard.EvoStandard()
		firstEvo := c.GetEvolutions()["0"]
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), baseStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), baseStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// EvoMixed calculates the standard evolution for cards with a evo[0] amalgamation it calculates
// the evo[1] using 1 amalgamated, and one as a "drop"
func (c *Card) EvoMixed() (stats Stats) {
	log.Printf("Calculating Mixed Evo for %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			stats = maxStats(c)
			log.Printf("Using collected stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
			// calculate the amalgamation stats here
			// stats = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				log.Printf("%d:%s is a Rebirth evo. Calculating the rebirth stats", c.ID, c.Name)
				// if this is an rebirth, calculate the max...
				matStats := mat.EvoMixedLvl1()
				log.Printf("Material Stats at lvl1: %s", matStats)
				matStats = matStats.Subtract(baseStats(mat))
				log.Printf("Material Stats Gains: %s, Rebirth Base Stats: %s / %s", matStats, baseStats(c), maxStats(c))
				stats = c.calculateAwakeningStat(matStats, false)
				log.Printf("Rebirth Mixed card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.EvoMixedLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, false)
					log.Printf("Awakening Mixed card %d (%s) -> %d (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.EvoStandard()
						stats.calculateEvoAccidentStat(matStats, maxStats(c))
						log.Printf("Using Evo Accident stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = maxStats(c)
						log.Printf("Using base Max stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var matStats Stats
		if c.EvolutionRank == 1 && c.PossibleMixedEvo() && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matStats = materialCard.AmalgamationStandard()
		} else {
			matStats = materialCard.EvoStandard()
		}
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), maxStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), maxStats(c))
		}
		log.Printf("Using Evo stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	}
	stats.ensureMaxCap(rarity)
	log.Printf("Final stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	return
}

// EvoMixedLvl1 calculates the standard evolution for cards with a evo[0] amalgamation it calculates
// the evo[1] using 1 amalgamated, and one as a "drop"
func (c *Card) EvoMixedLvl1() (stats Stats) {
	log.Printf("Calculating Mixed Evo for Lvl1 %s: %d\n", c.Name, c.EvolutionRank)
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// assume we collected this card directly (drop)?
			stats = baseStats(c)
			log.Printf("Using collected stats for Mixed Lvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
			// calculate the amalgamation stats here
			// stats = c.AmalgamationStandard()
			// log.Printf("Using Amalgamation stats for Mixed %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.EvoMixedLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, true)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.EvoMixedLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, false)
					log.Printf("Awakening Mixed card Lvl1 %d (%s) -> %d (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.EvoStandard()
						stats.calculateEvoAccidentStat(matStats, baseStats(c))
						log.Printf("Using Evo Accident stats for Mixed Lvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = baseStats(c)
						log.Printf("Using base Max stats for Mixed Lvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var matStats Stats
		if c.EvolutionRank == 1 && c.PossibleMixedEvo() && materialCard.ID == firstEvo.ID {
			// perform the mixed evo here
			matStats = materialCard.AmalgamationStandard()
		} else {
			matStats = materialCard.EvoStandard()
		}
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), baseStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, maxStats(firstEvo), baseStats(c))
		}
		log.Printf("Using Evo stats for Mixed Lvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	}
	stats.ensureMaxCap(rarity)
	log.Printf("Final stats for Mixed Lvl1 %s: %d (%s)\n", c.Name, c.EvolutionRank, stats)
	return
}

// Evo6Card calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 6-card evo
func (c *Card) Evo6Card() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.Amalgamation6Card()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.Evo6CardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, false)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.Evo6CardLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, false)
					log.Printf("Awakening 6 card %d (%s) -> %d (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.Evo6Card()
						stats.calculateEvoAccidentStat(matStats, maxStats(c))
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = maxStats(c)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		mat1Stats := materialCard.Evo6Card()
		var mat2Stats Stats
		if c.EvolutionRank == 4 {
			mat2Stats = c.GetEvolutions()["1"].Evo6Card()
		} else {
			mat2Stats = firstEvo.Evo6Card()
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, maxStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, maxStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// Evo6CardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 6-card evo
func (c *Card) Evo6CardLvl1() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.Amalgamation6CardLvl1()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.Evo6CardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, true)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.Evo6CardLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, true)
					log.Printf("Awakening 6 card %d (%s) -> %d lvl1 (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.Evo6Card()
						stats.calculateEvoAccidentStat(matStats, baseStats(c))
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = baseStats(c)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		mat1Stats := materialCard.Evo6Card()
		var mat2Stats Stats
		if c.EvolutionRank == 4 {
			mat2Stats = c.GetEvolutions()["1"].Evo6Card()
		} else {
			mat2Stats = firstEvo.Evo6Card()
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, baseStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, baseStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// Evo9Card calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 9-card evo
// this one should only be used for <= SR
func (c *Card) Evo9Card() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.Amalgamation9Card()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.Evo9CardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, false)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.Evo9CardLvl1().Subtract(baseStats(mat))
					stats = c.calculateAwakeningStat(matStats, false)
					log.Printf("Awakening 9 card %d (%s) -> %d (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.Evo9Card()
						stats.calculateEvoAccidentStat(matStats, maxStats(c))
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = maxStats(c)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var mat1Stats, mat2Stats Stats
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Stats = materialCard.EvoPerfect()
			mat2Stats = mat1Stats
		} else if c.EvolutionRank == 3 {
			mat1Stats = c.GetEvolutions()["2"].EvoStandard()
			mat2Stats = c.GetEvolutions()["1"].EvoStandard()
		} else {
			// this would be the materials to get 4*
			mat1Stats = c.GetEvolutions()["2"].Evo9Card()
			mat2Stats = c.GetEvolutions()["3"].Evo9Card()
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, maxStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, maxStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// Evo9CardLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 9-card evo
// this one should only be used for <= SR
func (c *Card) Evo9CardLvl1() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.Amalgamation9CardLvl1()
		} else {
			mat := c.RebirthsFrom()
			if mat != nil {
				// if this is an rebirth, calculate the max...
				matStats := mat.Evo9CardLvl1().Subtract(baseStats(mat))
				stats = c.calculateAwakeningStat(matStats, true)
				log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
					mat.ID, matStats, c.ID, stats)
			} else {
				mat := c.AwakensFrom()
				if mat != nil {
					// if this is an awakwening, calculate the max...
					matStats := mat.Evo9CardLvl1()
					stats = c.calculateAwakeningStat(matStats, true)
					log.Printf("Awakening 9 card %d (%s) -> %d lvl1 (%s)\n",
						mat.ID, matStats, c.ID, stats)
				} else {
					// check for Evo Accident
					mat = c.EvoAccidentOf()
					if mat != nil {
						// calculate the transfered stats of the 2 material cards
						// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
						matStats := mat.Evo9Card()
						stats.calculateEvoAccidentStat(matStats, baseStats(c))
					} else {
						// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
						stats = baseStats(c)
					}
				}
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		firstEvo := c.GetEvolutions()["0"]
		var mat1Stats, mat2Stats Stats
		if c.EvolutionRank <= 2 {
			// pretend it's a perfect evo. this gets tricky after evo 2
			mat1Stats = materialCard.EvoPerfect()
			mat2Stats = mat1Stats
		} else if c.EvolutionRank == 3 {
			mat1Stats = c.GetEvolutions()["2"].Evo9Card()
			mat2Stats = c.GetEvolutions()["1"].Evo9Card()
		} else {
			// this would be the materials to get 4*
			mat1Stats = c.GetEvolutions()["2"].Evo9Card()
			mat2Stats = c.GetEvolutions()["3"].Evo9Card()
		}

		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, baseStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(mat1Stats, mat2Stats, baseStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// EvoPerfect calculates the standard evolution stat but at max level. If this is a 4* card,
// calculates for 16-card evo
func (c *Card) EvoPerfect() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.AmalgamationPerfect()
		} else if c.RebirthsFrom() != nil {
			mat := c.RebirthsFrom()
			// if this is an rebirth, calculate the max...
			matStats := mat.EvoPerfectLvl1()
			log.Printf("Material Stats at lvl1: %s", matStats)
			matStats = matStats.Subtract(baseStats(mat))
			log.Printf("Material Stats Gains: %s, Rebirth Base Stats: %s / %s", matStats, baseStats(c), maxStats(c))
			stats = c.calculateAwakeningStat(matStats, false)
			log.Printf("Rebirth Perfect card %d (%s) -> %d (%s)\n",
				mat.ID, matStats, c.ID, stats)
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matStats := mat.EvoPerfectLvl1().Subtract(baseStats(mat))
			stats = c.calculateAwakeningStat(matStats, false)
			log.Printf("Awakening Perfect card %d (%s) -> %d (%s)\n",
				mat.ID, matStats, c.ID, stats)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf()
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matStats := turnOver.EvoPerfect()
				stats.calculateEvoAccidentStat(matStats, maxStats(c))
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				stats = maxStats(c)
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matStats := materialCard.EvoPerfect()
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions()["0"]
			stats = c.calculateEvoStats(matStats, matStats, maxStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, matStats, maxStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}

// EvoPerfectLvl1 calculates the standard evolution stat but at level 1. If this is a 4* card,
// calculates for 16-card evo
func (c *Card) EvoPerfectLvl1() (stats Stats) {
	materialCard := c.PrevEvo()
	rarity := c.CardRarity()
	if materialCard == nil {
		// log.Printf("No previous evo found for card %v\n", c.ID)
		// check for amalgamation
		if c.IsAmalgamation() {
			// calculate the amalgamation stats here
			stats = c.AmalgamationPerfectLvl1()
		} else if c.RebirthsFrom() != nil {
			mat := c.RebirthsFrom()
			// if this is an rebirth, calculate the max...
			matStats := mat.EvoPerfectLvl1().Subtract(baseStats(mat))
			stats = c.calculateAwakeningStat(matStats, true)
			log.Printf("Rebirth standard card %d (%s) -> %d (%s)\n",
				mat.ID, matStats, c.ID, stats)
		} else if c.AwakensFrom() != nil {
			// if this is an awakwening, calculate the max...
			mat := c.AwakensFrom()
			matStats := mat.EvoPerfectLvl1().Subtract(baseStats(mat))
			stats = c.calculateAwakeningStat(matStats, true)
			log.Printf("Awakening Perfect card %d (%s) -> %d lvl1 (%s)\n",
				mat.ID, matStats, c.ID, stats)
		} else {
			// check for Evo Accident
			turnOver := c.EvoAccidentOf()
			if turnOver != nil {
				// calculate the transfered stats of the 2 material cards
				// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk)
				matStats := turnOver.EvoPerfect()
				stats.calculateEvoAccidentStat(matStats, baseStats(c))
			} else {
				// if not an amalgamation or awoken card, use the default MAX (this should be evo 0*)
				stats = baseStats(c)
			}
		}
	} else {
		//TODO: For LR cards do we want to use level 1 evos for the first 3?
		matStats := materialCard.EvoPerfect()
		// calculate the transfered stats of the 2 material cards
		// ret = (0.15 * previous evo max atk) + (0.15 * [0*] max atk) + (newCardMax * bonus)
		if c.LastEvolutionRank == 4 {
			// 4* cards do not use the result card stats
			firstEvo := c.GetEvolutions()["0"]
			stats = c.calculateEvoStats(matStats, matStats, baseStats(firstEvo))
		} else {
			// 1* cards use the result card stats
			stats = c.calculateEvoStats(matStats, matStats, baseStats(c))
		}
	}
	stats.ensureMaxCap(rarity)
	return
}
