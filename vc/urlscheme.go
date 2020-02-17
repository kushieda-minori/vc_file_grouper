package vc

import "sort"

//URLScheme URL Scheme
type URLScheme struct {
	ID      int    `json:"_id"`
	Ios     string `json:"ios"`
	Android string `json:"android"`
}

// URLSchemeScan search for a URLScheme by ID
func URLSchemeScan(id int) *URLScheme {
	if id > 0 {
		l := len(Data.URLSchemes)
		i := sort.Search(l, func(i int) bool { return Data.URLSchemes[i].ID >= id })
		if i >= 0 && i < l && Data.URLSchemes[i].ID == id {
			return &(Data.URLSchemes[i])
		}
	}
	return nil
}
