package structout

import (
	"vc_file_grouper/vc"
)

// CardStatInfo card info that can be output to JSON
type CardStatInfo struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Element             string `json:"element"`
	Rarity              string `json:"rarity"`
	BaseAtk             int    `json:"baseAtk"`
	BaseDef             int    `json:"baseDef"`
	BaseSol             int    `json:"baseSol"`
	MaxAtk              int    `json:"maxAtk"`
	MaxDef              int    `json:"maxDef"`
	MaxSol              int    `json:"maxSol"`
	MaxLevel            int    `json:"maxLevel"`
	MaxRarityAtk        int    `json:"maxRarityAtk"`
	MaxRarityDef        int    `json:"maxRarityDef"`
	MaxRaritySol        int    `json:"maxRaritySol"`
	PerfectSoldierCount int    `json:"perfectSoldierCount"`
	RebirthCardID       *int   `json:"rebirthCardId"`
	IsClosed            bool   `json:"isClosed"`
}

// ToCardStatInfo converts the full card info to the shortened form for export
func ToCardStatInfo(c *vc.Card) CardStatInfo {
	var rebirthID *int = nil
	if rb := c.RebirthsTo(); rb != nil {
		rebirthID = &rb.ID
	}
	return CardStatInfo{
		ID:                  c.ID,
		Name:                c.Name,
		Element:             c.Element(),
		Rarity:              c.Rarity(),
		BaseAtk:             c.DefaultOffense,
		BaseDef:             c.DefaultDefense,
		BaseSol:             c.DefaultFollower,
		MaxAtk:              c.MaxOffense,
		MaxDef:              c.MaxDefense,
		MaxSol:              c.MaxFollower,
		MaxLevel:            c.CardRarity().MaxCardLevel,
		MaxRarityAtk:        c.CardRarity().LimtOffense,
		MaxRarityDef:        c.CardRarity().LimtDefense,
		MaxRaritySol:        c.CardRarity().LimtMaxFollower,
		PerfectSoldierCount: c.EvoPerfect().Soldiers,
		RebirthCardID:       rebirthID,
		IsClosed:            c.IsClosed == 1,
	}
}
