package vc

// Amalgamation List of possible Fusions (Amalgamations) from master file field "fusion_list"
type Amalgamation struct {
	// internal id
	ID int `json:"_id"`
	// card 1
	Material1 int `json:"material_1"`
	// card 2
	Material2 int `json:"material_2"`
	// card 3
	Material3 int `json:"material_3"`
	// card 4
	Material4 int `json:"material_4"`
	// resulting card
	FusionCardID int `json:"fusion_card_id"`
}

// MaterialCount number of materials used in an amalgamation
func (a *Amalgamation) MaterialCount() int {
	if a.Material4 > 0 {
		return 4
	}
	if a.Material3 > 0 {
		return 3
	}
	return 2
}

// Materials material used in the amalgamation including the result
func (a *Amalgamation) Materials() []*Card {
	return append(a.MaterialsOnly(), a.Result())
}

// Result result card from the amalgamation
func (a *Amalgamation) Result() *Card {
	return CardScan(a.FusionCardID)
}

// MaterialsOnly material used in the amalgamation excluding the result
func (a *Amalgamation) MaterialsOnly() []*Card {
	ret := make([]*Card, 0)
	ret = append(ret, CardScan(a.Material1))
	ret = append(ret, CardScan(a.Material2))
	if a.Material3 > 0 {
		ret = append(ret, CardScan(a.Material3))
	}
	if a.Material4 > 0 {
		ret = append(ret, CardScan(a.Material4))
	}
	return ret
}

// ByMaterialCount sorting interface for sorting amalgamations
// by the number of materials
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
