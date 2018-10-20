package vc

// CardLevel From the master data file, this is the "cardlevel" field
type CardLevel struct {
	// card level
	ID int `json:"_id"`
	// experiance needed to be at this level
	Exp int `json:"exp"`
}
