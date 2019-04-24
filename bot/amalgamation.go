package bot

import (
	"zetsuboushita.net/vc_file_grouper/vc"
)

// AmalgamationItem // single item used in amalgamation
type AmalgamationItem struct {
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

// AmalgamationItems slice of amalgamations
type AmalgamationItems []AmalgamationItem

// AmalgamationRecipe recipe for amalgamations
type AmalgamationRecipe struct {
	recipeID  int
	Materials AmalgamationItems `json:"materials,omitempty"`
	Result    AmalgamationItem  `json:"result,omitempty"`
}

// AmalgamationRecipes recipies
type AmalgamationRecipes []AmalgamationRecipe

// Amagamations amalgamations
type Amagamations struct {
	AsMaterial AmalgamationRecipes `json:"asMaterial,omitempty"`
	AsResult   AmalgamationRecipes `json:"asResult,omitempty"`
}

func cleanHXAmal(rarity string) string {
	if rarity == "HX" {
		return "X"
	}
	return rarity
}

func newRecipe(a vc.Amalgamation) AmalgamationRecipe {
	mats := a.MaterialsOnly()
	l := len(mats)
	materials := make(AmalgamationItems, 0, l-1)
	for _, mat := range mats {
		materials = append(materials, AmalgamationItem{
			Name:   mat.Name,
			Rarity: cleanHXAmal(mat.Rarity()),
		})
	}
	res := a.Result()
	return AmalgamationRecipe{
		recipeID:  a.ID,
		Materials: materials,
		Result: AmalgamationItem{
			Name:   res.Name,
			Rarity: cleanHXAmal(res.Rarity()),
		},
	}
}

func (r AmalgamationRecipes) contains(n AmalgamationRecipe) bool {
	for _, i := range r {
		if i.recipeID == n.recipeID {
			return true
		}
	}
	return false
}

func (r *AmalgamationRecipes) addNewRecipe(n AmalgamationRecipe) {
	if !r.contains(n) {
		*r = append(*r, n)
	}
}
