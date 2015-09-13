package vc_grouper

// From the master data file, this is the "cardlevel" field
type CardLevel struct {
	// card level
	Id int `json:"_id"`
	// experiance needed to be at this level
	Exp int `json:"exp"`
}
