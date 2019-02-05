package vc

// CardAwaken list of possible card awakeneings and their cost from master file field "card_awaken"
type CardAwaken struct {
	// awakening id
	ID int `json:"_id"`
	// case card
	BaseCardID int `json:"base_card_id"`
	// result card
	ResultCardID int `json:"result_card_id"`
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
	// Order in the "Awoken Card List maybe?"
	Order int `json:"order"`
	// IsClosed true if unreleased (hides it from the "Awoken Card List" in the upgrade screen)
	IsClosed int `json:"is_closed"`
}

// Item needed to awaken the source card
func (ca *CardAwaken) Item(i int) *Item {
	if i < 1 || i > 5 {
		return nil
	}
	switch i {
	case 1:
		if ca.Material1Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material1Item)
	case 2:
		if ca.Material2Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material2Item)
	case 3:
		if ca.Material3Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material3Item)
	case 4:
		if ca.Material4Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material4Item)
	case 5:
		if ca.Material5Item <= 0 {
			return nil
		}
		return ItemScan(ca.Material5Item)
	}
	return nil
}
