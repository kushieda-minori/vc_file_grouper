package vc

import (
	"fmt"
	"sort"
	"strings"
)

// DeckBonus Unit Bonuses from master file field "deck_bonus"
// these match with the strings in MsgDeckBonusName_en.strb
// and MsgDeckBonusDesc_en.strb
type DeckBonus struct {
	ID          int    `json:"_id"`         // bonus id
	AtkDefFlg   int    `json:"atk_def_flg"` // Affects ATK or DEF
	ValueType   int    `json:"value_type"`  // ?
	Value       int    `json:"value"`       // amount of the modifier
	DownGrade   int    `json:"down_grade"`  // ?
	CondType    int    `json:"cond_type"`   // deck condition
	ReqNum      int    `json:"req_num"`     // number of cards required
	DupFlg      int    `json:"dup_flg"`     // allows duplicates
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Conditions that trigger the bonus
func (d *DeckBonus) Conditions(v *VFile) DeckBonusCondArray {
	ret := make([]DeckBonusCond, 0)
	for _, val := range v.DeckBonusConditions {
		if val.DeckBonusID == d.ID {
			switch d.CondType {
			case 2:
				c := CardScanCharacter(val.RefID, v.Cards)
				if c == nil || c.Name == "" {
					continue
				} else {
					val.RefName = c.Name
				}
			case 3:
				switch val.RefID {
				case 1:
					val.RefName = "Light"
				case 2:
					val.RefName = "Passion"
				case 3:
					val.RefName = "Cool"
				case 4:
					val.RefName = "Dark"
				default:
					val.RefName = fmt.Sprintf("Unknown Element (%d)", val.RefID)
				}
			case 8:
				switch val.RefID {
				case 1:
					val.RefName = "N"
				case 2:
					val.RefName = "R"
				case 3:
					val.RefName = "SR"
				case 4:
					val.RefName = "HN"
				case 5:
					val.RefName = "HR"
				case 6:
					val.RefName = "HSR"
				case 8:
					val.RefName = "UR"
				case 9:
					val.RefName = "HUR"
				case 10:
					val.RefName = "GSR"
				case 11:
					val.RefName = "GUR"
				default:
					val.RefName = fmt.Sprintf("Unknown Rarity (%d)", val.RefID)
				}
			default:
				val.RefName = fmt.Sprintf("Unknown Type (%d)", val.CondTypeID)
			}
			ret = append(ret, val)
		}
	}
	return ret
}

// DeckBonusCond Deck Bonus Conditions from masfter file field "deck_bonus_cond"
type DeckBonusCond struct {
	ID          int `json:"_id"`           // deck condition id
	DeckBonusID int `json:"deck_bonus_id"` // deck bonus id
	Group       int `json:"group"`         // group
	CondTypeID  int `json:"cond_type_id"`  // type id
	/* reference to the card character id, (type 2)
	* the element, (type 3)
	* or the rarity (type 8)
	* Elements: 2=Passion, 3=Cool, 1=Light, 4=Dark
	* Rarity: 4=H/HN, 5=R/HR, 6=SR/HSR, 9=UR/HUR, 11=GSR/GUR
	 */
	RefID   int    `json:"ref_id"`
	RefName string `json:"-"` //lookup of card name
}

// DeckBonusByCountAndName sort interface to sort by number of items
// under the deck bonus and by name
type DeckBonusByCountAndName []DeckBonus

func (d DeckBonusByCountAndName) Len() int {
	return len(d)
}

func (d DeckBonusByCountAndName) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DeckBonusByCountAndName) Less(i, j int) bool {
	if d[i].ReqNum < d[j].ReqNum {
		return true
	}
	if d[i].ReqNum == d[j].ReqNum {
		if d[i].CondType > d[j].CondType {
			return true
		}
		if d[i].CondType == d[j].CondType {
			if d[i].AtkDefFlg != d[j].AtkDefFlg && d[i].AtkDefFlg == 1 {
				return true
			}
			if d[i].AtkDefFlg == d[j].AtkDefFlg {
				if d[i].Value < d[j].Value {
					return true
				}
				if d[i].Value == d[j].Value && strings.Compare(d[i].Name, d[j].Name) < 0 {
					return true
				}
			}
		}
	}
	return false
}

func (d *DeckBonusCond) String() string {
	return d.RefName
}

// DeckBonusCondArray helper to work on arrays of deck bonus conditions
type DeckBonusCondArray []DeckBonusCond

func (da DeckBonusCondArray) String() string {
	if len(da) < 1 {
		return "[]"
	}
	strs := make([]string, len(da))
	for i, d := range da {
		strs[i] = d.String()
	}
	if da[0].CondTypeID == 2 {
		sort.Strings(strs)
	}
	return "[ " + strings.Join(strs, ", ") + " ]"
}
