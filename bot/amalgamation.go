package bot

import (
	"zetsuboushita.net/vc_file_grouper/vc"
)

// AmalgamationItem // single item used in amalgamation
type AmalgamationItem struct {
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

// AmalgamationRecipe recipe for amalgamations
type AmalgamationRecipe struct {
	Materials []AmalgamationItem `json:"materials"`
	Result    AmalgamationItem   `json:"result"`
}

func newRecipe(a *vc.Amalgamation) AmalgamationRecipe {
	mats := a.Materials()
	l := len(mats)
	materials := make([]AmalgamationItem, 0, l-1)
	res := mats[l-1]
	for _, mat := range mats[0 : l-2] {
		materials = append(materials, AmalgamationItem{
			Name:   mat.Name,
			Rarity: mat.Rarity(),
		})
	}
	return AmalgamationRecipe{
		Materials: materials,
		Result: AmalgamationItem{
			Name:   res.Name,
			Rarity: res.Rarity(),
		},
	}
}
