package vc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"zetsuboushita.net/vc_file_grouper/util"
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
	character   *CardCharacter
	archwitches ArchwitchList
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

// CardList helper interface for looking at lists of cards
type CardList []*Card

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

// EvosWithDistinctImages return the list of evos that have distinct images
func (c *Card) EvosWithDistinctImages(icons bool) []string {
	ret := make([]string, 0)
	evos := c.GetEvolutions()
	images := make(map[string][]byte, 0)
	for _, evoID := range EvoOrder {
		if evo, ok := evos[evoID]; ok {
			file := FilePath + "/card/"
			if icons {
				file += "thumb/"
			} else {
				file += "md/"
			}
			file += evo.Image()
			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Error reading image file %s", file)
				continue
			}

			if len(images) == 0 || isDistinctImage(&images, &(data)) {
				images[evoID] = data
				ret = append(ret, evoID)
			}
		}
	}
	return ret
}

func isDistinctImage(m *map[string][]byte, i *[]byte) bool {
	for _, v := range *m {
		if bytes.Equal(v, *i) {
			return false
		}
	}
	return true
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

// EvoIsFirst returns true if the Evolution of this card is the first for this card
func (c *Card) EvoIsFirst() bool {
	evos := c.GetEvolutionCards()
	return c.ID == evos[0].ID
	//return c.PrevEvo() == nil && c.RebirthsFrom() == nil && c.AwakensFrom() == nil
}

// EvoIsMidOf4 returns true if the Evolution of this card is a middle evolution for this card.
// This means it would be a 1*, 2* or 3* card.
func (c *Card) EvoIsMidOf4() bool {
	return c.LastEvolutionRank == 4 && c.EvolutionRank > 0 && c.EvolutionRank < 4
	//return c.PrevEvo() == nil && c.RebirthsFrom() == nil && c.AwakensFrom() == nil
}

// EvoIsHigh returns true if the Evolution of this card is an Awoken evolution
func (c *Card) EvoIsHigh() bool {
	s := c.Rarity()
	l := len(s)
	return l >= 2 && s[0] == 'H'
}

// EvoIsAwoken returns true if the Evolution of this card is an Awoken evolution
func (c *Card) EvoIsAwoken() bool {
	s := c.Rarity()
	l := len(s)
	return l == 3 && s[0] == 'G'
}

// EvoIsReborn returns true if the Evolution of this card is a Rebirth evolution
func (c *Card) EvoIsReborn() bool {
	s := c.Rarity()
	l := len(s)
	return l == 3 && s[0] == 'X'
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
		if c.CardCharaID <= 0 || c.EvolutionCardID <= 0 || c.EvoIsHigh() {
			return nil
		}

		var tmp *Card
		character := c.Character()
		for i, cd := range character.Cards() {
			if cd.ID == c.EvolutionCardID {
				tmp = character._cards[i]
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
				tmp = c.Character()._cards[i]
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

// FirstEvo Gets the first evolution for the card excluding pre-awakening/amalgamations
func (c *Card) FirstEvo() *Card {
	t := c
	for t.PrevEvo() != nil {
		t = t.PrevEvo()
	}
	return t
}

// LastEvo Gets the last evolution for the card excluding awakening/amalgamations
func (c *Card) LastEvo() *Card {
	t := c
	for t.NextEvo() != nil {
		t = t.NextEvo()
	}
	return t
}

// Archwitches If this card was used as an AW, get the AW information.
// This can be used to get Likability information. If a card was used
// as an Achwitch in multiple events, multiple items can be returned here.
func (c *Card) Archwitches() ArchwitchList {
	if c == nil {
		return ArchwitchList{}
	}
	if c.archwitches == nil {
		c.archwitches = make(ArchwitchList, 0)
		for _, aw := range Data.Archwitches {
			if c.ID == aw.CardMasterID && !c.archwitches.Contains(aw) {
				c.archwitches = append(c.archwitches, aw)
			}
		}
		sort.Slice(c.archwitches, func(a, b int) bool {
			awa := c.archwitches[a]
			awb := c.archwitches[b]
			if awa.KingSeriesID == awb.KingSeriesID {
				return awa.ID < awb.ID
			}
			return awa.KingSeriesID < awb.KingSeriesID
		})
	}
	log.Printf("Found %d Archwitch records for card %d:%s", len(c.archwitches), c.ID, c.Name)
	return c.archwitches
}

// ArchwitchesWithLikeabilityQuotes If this card was used as an AW, get the AW information.
// This can be used to get Likability information. If a card was used
// as an Achwitch in multiple events, multiple items can be returned here.
func (c *Card) ArchwitchesWithLikeabilityQuotes() ArchwitchList {
	aws := c.Archwitches()
	ret := make(ArchwitchList, 0)
	for _, aw := range aws {
		if len(aw.Likeability()) > 0 {
			ret = append(ret, aw)
		}
	}
	log.Printf("Found %d Archwitch records with likeability quotes for card %d:%s", len(ret), c.ID, c.Name)
	return ret
}

// EvoAccident If this card can produce an evolution accident, get the result card.
func (c *Card) EvoAccident() *Card {
	return CardScan(c.TransCardID)
}

// EvoAccidentOf If this card is the result of an evo accident, get the source card.
func (c *Card) EvoAccidentOf() *Card {
	for key, val := range Data.Cards {
		if val.TransCardID == c.ID {
			return Data.Cards[key]
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
		if val.IsClosed != 0 {
			continue
		}
		if c.ID == val.BaseCardID {
			return CardScan(val.ResultCardID)
		}
	}
	return nil
}

// AwakensFrom gets the source card of this awoken card
func (c *Card) AwakensFrom() *Card {
	for _, val := range Data.Awakenings {
		if val.IsClosed != 0 {
			continue
		}
		if c.ID == val.ResultCardID {
			return CardScan(val.BaseCardID)
		}
	}
	return nil
}

// HasRebirth Gets the card this card rebirths to.
func (c *Card) HasRebirth() bool {
	for _, val := range Data.Rebirths {
		if val.IsClosed != 0 {
			continue
		}
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
		if val.IsClosed != 0 {
			continue
		}
		if c.ID == val.BaseCardID {
			return CardScan(val.ResultCardID)
		}
	}
	return nil
}

// RebirthsFrom gets the source card of this rebirth card
func (c *Card) RebirthsFrom() *Card {
	for _, val := range Data.Rebirths {
		if val.IsClosed != 0 {
			continue
		}
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
	if c == nil {
		return nil
	}
	if c.skill1 == nil && c.SkillID1 > 0 {
		c.skill1 = SkillScan(c.SkillID1)
	}
	return c.skill1
}

// Skill2 of the card
func (c *Card) Skill2() *Skill {
	if c == nil {
		return nil
	}
	if c.skill2 == nil && c.SkillID2 > 0 {
		c.skill2 = SkillScan(c.SkillID2)
	}
	return c.skill2
}

// Skill3 of the card
func (c *Card) Skill3() *Skill {
	if c == nil {
		return nil
	}
	if c.skill3 == nil && c.SkillID3 > 0 {
		c.skill3 = SkillScan(c.SkillID3)
	}
	return c.skill3
}

// SpecialSkill1 of the card (Awoken Burst)
func (c *Card) SpecialSkill1() *Skill {
	if c == nil {
		return nil
	}
	if c.specialSkill1 == nil && c.SpecialSkillID1 > 0 {
		c.specialSkill1 = SkillScan(c.SpecialSkillID1)
	}
	return c.specialSkill1
}

// ThorSkill1 of the card
func (c *Card) ThorSkill1() *Skill {
	if c == nil {
		return nil
	}
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
		return Data.Cards[i]
	}
	return nil
}

// CardScanCharacter searches for a card by the character ID
func CardScanCharacter(charID int) *Card {
	if charID > 0 {
		for k, val := range Data.Cards {
			//return the first one we find.
			if val.CardCharaID == charID {
				return Data.Cards[k]
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
				return Data.Cards[k]
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

// Earliest gets the ealiest released card from a list of cards. Determined by ID
func (d CardList) Earliest() (min *Card) {
	for idx, card := range d {
		if min == nil || min.ID > card.ID {
			// log.Printf("'Earliest' Card: %d, Name: %s\n", card.ID, card.Name)
			min = d[idx]
		}
	}
	// if min != nil {
	// log.Printf("-Earliest Card: %d, Name: %s\n", min.ID, min.Name)
	// }
	return
}

// Latest gets the latest released card from a list of cards. Determined by ID
func (d CardList) Latest() (max *Card) {
	for idx, card := range d {
		if max == nil || max.ID < card.ID {
			// log.Printf("'Latest' Card: %d, Name: %s\n", card.ID, card.Name)
			max = d[idx]
		}
	}
	// if max != nil {
	// log.Printf("-Latest Card: %d, Name: %s\n", max.ID, max.Name)
	// }
	return
}

//Copy returns a copy of this list. Useful for local sorting
func (d CardList) Copy() CardList {
	ret := make(CardList, len(d), len(d))
	copy(ret, d)
	return ret
}

//MinimumEvolutionRank gets the lowest evolution rank in the set
func (d CardList) MinimumEvolutionRank() (min string) {
	var minRare *CardRarity
	for idx, card := range d {
		if minRare == nil || minRare.ID > card.CardRareID {
			// log.Printf("'Latest' Card: %d, Name: %s\n", card.ID, card.Name)
			minRare = d[idx].CardRarity()
			min = d[idx].Rarity()
		}
	}

	return min
}

func getAmalBaseCard(card *Card) *Card {
	if card.IsAmalgamation() {
		log.Printf("Checking Amalgamation base for Card: %d, Name: %s, Evo: %d\n", card.ID, card.Name, card.EvolutionRank)
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
				log.Printf("Found amalgamation: %d, Name: '%s', Evo: %d\n", tamalCard.ID, tamalCard.Name, tamalCard.EvolutionRank)
				if tamalCard.Name == c.Name {
					amalCard = tamalCard
					// check for amal awakening
					amalAwakening = amalCard.AwakensTo()
					rebirthAmal = amalCard.RebirthsTo()
					if rebirthAmal == nil && amalAwakening != nil {
						rebirthAmal = amalAwakening.RebirthsTo()
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
		if c.EvoIsAwoken() {
			if isIcon {
				return fileName + "_G"
			}
			return fileName + "_H"
		}
		if c.EvoIsHigh() {
			return fileName + "_H"
		}
		return fileName
	}
	if !isIcon {
		if thisKey[0] == 'G' {
			if _, ok := evos["H"]; ok {
				evoImages := c.EvosWithDistinctImages(isIcon)
				if !util.Contains(evoImages, thisKey) {
					return fileName + "_H"
				}
			} else {
				return fileName + "_H"
			}
		} else if thisKey == "A" {
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
			log.Printf("No character info Card: %d, Name: %s, Evo: %d\n", c.ID, c.Name, c.EvolutionRank)
			ret["0"] = c
			c._allEvos = ret
			return ret
		}

		c2 := c
		// check if this is a rebirth card
		if c2.EvoIsReborn() {
			log.Printf("Card %d:%s is reborn, finding it's source", c2.ID, c2.Name)
			tmp := c2.RebirthsFrom()
			if tmp == nil {
				ch := c2.Character()
				if ch != nil && ch.Cards()[0].Name == c2.Name {
					c2 = ch.Cards()[0]
					log.Printf("Found card %d:%s", c2.ID, c2.Name)
				}
				// the name changed, so we'll keep this card
			} else {
				c2 = tmp
				log.Printf("Found card %d:%s", c2.ID, c2.Name)
			}
		}
		// check if this is an awoken card
		if c2.EvoIsAwoken() {
			log.Printf("Card %d:%s is awoken, finding it's source", c2.ID, c2.Name)
			tmp := c2.AwakensFrom()
			if tmp == nil {
				ch := c2.Character()
				if ch != nil && ch.Cards()[0].Name == c2.Name {
					c2 = ch.Cards()[0]
					log.Printf("Found card %d:%s", c2.ID, c2.Name)
				}
				// the name changed, so we'll keep this card
			} else {
				c2 = tmp
				log.Printf("Found card %d:%s", c2.ID, c2.Name)
			}
		}

		// get earliest evo
		for tmp := c2.PrevEvo(); tmp != nil; tmp = tmp.PrevEvo() {
			c2 = tmp
			log.Printf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank)
		}

		// at this point we should have the first card in the evolution path
		c2 = getAmalBaseCard(c2)

		// get earliest evo (again...)
		for tmp := c2.PrevEvo(); tmp != nil; tmp = tmp.PrevEvo() {
			c2 = tmp
			log.Printf("Looking for earliest Evo for Card: %d, Name: %s, Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank)
		}

		log.Printf("Base Card: %d, Name: '%s', Evo: %d\n", c2.ID, c2.Name, c2.EvolutionRank)

		// assigns evolutions found
		assignLastEvos := func(awakening, amalCard, amalAwakening, rebirth, rebirthAmal *Card) {
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

		// populate the actual evos.
		for nextEvo := c2; nextEvo != nil; nextEvo = nextEvo.NextEvo() {
			log.Printf("Next Evo is Card: %d, Name: '%s', Evo: %d\n", nextEvo.ID, nextEvo.Name, nextEvo.EvolutionRank)
			if nextEvo.EvolutionRank <= 0 {
				evoRank := "0"
				if nextEvo.Rarity()[0] == 'H' {
					evoRank = "H"
				}
				ret[evoRank] = nextEvo
				if nextEvo.LastEvolutionRank < 0 {
					// check for awakening
					awakening, amalCard, amalAwakening, rebirth, rebirthAmal := checkEndCards(nextEvo)
					assignLastEvos(awakening, amalCard, amalAwakening, rebirth, rebirthAmal)
				}
			} else if nextEvo.EvoIsReborn() {
				// for some reason we hit a X during Evo traversal. Probably a X originating
				// from amalgamation

				// check for awakening/rebirth
				_, rebirthAmal, _, _, _ := checkEndCards(nextEvo)
				assignLastEvos(nil, nil, nil, nextEvo, rebirthAmal)
			} else if nextEvo.EvoIsAwoken() {
				// for some reason we hit a G during Evo traversal. Probably a G originating
				// from amalgamation

				// check for awakening/rebirth
				_, amalCard, amalAwakening, rebirth, rebirthAmal := checkEndCards(nextEvo)
				assignLastEvos(nextEvo, amalCard, amalAwakening, rebirth, rebirthAmal)
			} else if nextEvo.EvolutionRank == c2.LastEvolutionRank || nextEvo.EvoIsHigh() || nextEvo.LastEvolutionRank < 0 {
				ret["H"] = nextEvo
				// check for awakening
				awakening, amalCard, amalAwakening, rebirth, rebirthAmal := checkEndCards(nextEvo)
				assignLastEvos(awakening, amalCard, amalAwakening, rebirth, rebirthAmal)
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

		log.Printf("Found Evos: ")
		for key, card := range ret {
			card._allEvos = ret
			log.Printf("(%s: %d) ", key, card.ID)
		}
		log.Printf("\n")
		c._allEvos = ret
		return ret
	}
	return c._allEvos
}

// GetEvolutionCards same as GetEvolutions, but only returns the cards
func (c *Card) GetEvolutionCards() CardList {
	evos := c.GetEvolutions()
	cards := make(CardList, 0, len(evos))
	//log.Printf("Card: %d, Name: %s, Evos: %d\n", c.ID, c.Name, len(evos))
	for _, ek := range EvoOrder {
		evo := evos[ek]
		if evo == nil {
			continue
		}
		cards = append(cards, evo)
	}
	//log.Printf("Cards: %d\n", len(cards))
	return cards
}

//CardsByName gets the cards by name. If multiple cards have the same name, they are groupped together
func CardsByName() map[string]CardList {
	ret := make(map[string]CardList, 0)

	for _, card := range Data.Cards {
		if _, ok := ret[card.Name]; !ok {
			ret[card.Name] = make(CardList, 0)
		}
		ret[card.Name] = append(ret[card.Name], card)
	}

	return ret
}

//CardsByNameByLowestID gets the CardsByName then sorts them by lowest ID first.
func CardsByNameByLowestID(asc bool) []CardList {
	byName := CardsByName()
	ret := make([]CardList, 0, len(byName))
	for _, cl := range byName {
		ret = append(ret, cl)
	}
	if asc {
		// newest card last by original release order
		sort.Slice(ret, func(a, b int) bool {
			return ret[a].Earliest().ID < ret[b].Earliest().ID
		})
	} else {
		// newest cards first by original release order
		sort.Slice(ret, func(a, b int) bool {
			return ret[a].Earliest().ID > ret[b].Earliest().ID
		})
	}
	return ret
}

// EvoOrder order of evolutions in the map.
var EvoOrder = [10]string{"0", "1", "2", "3", "H", "A", "G", "GA", "X", "XA"}

// Elements of the cards.
var Elements = [5]string{"Light", "Passion", "Cool", "Dark", "Special"}

// Rarity of the cards
var Rarity = [17]string{"N", "R", "SR", "HN", "HR", "HSR", "X", "UR", "HUR", "GSR", "GUR", "LR", "HLR", "GLR", "XSR", "XUR", "XLR"}
