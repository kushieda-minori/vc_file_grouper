package vc

import (
	"sort"
	"strings"
)

// Item "items" is the list of Archwitches
type Item struct {
	ID                int       `json:"_id"`
	GroupID           int       `json:"group_id"`
	Name              string    `json:"name"`
	Ratio             int       `json:"ratio"`
	Day               int       `json:"day"`
	Order             int       `json:"order"`
	UseSituation      int       `json:"use_situation"`
	ZeroDisp          int       `json:"zero_disp"`
	ItemNo            int       `json:"item_no"`
	ItemNoSub         int       `json:"item_no_sub"`
	Tab               int       `json:"tab"`
	MaxCount          int       `json:"max_count"`
	UseButtonAct      int       `json:"use_button_act"`
	ChangeSceneAct    int       `json:"change_scene_act"`
	EndDatetime       Timestamp `json:"end_datetime"`
	LimitedItemFlg    int       `json:"limited_item_flg"`
	Present           int       `json:"present"`
	MaxBuyCount       int       `json:"max_buy_count"`
	IsDelete          int       `json:"is_delete"`
	DailyMaxBuyCount  int       `json:"daily_max_buy_count"`
	ArcanaType        int       `json:"arcana_type"`
	Description       string    `json:"-"` // MsgShopItemDesc_en.strb
	DescriptionInShop string    `json:"-"` // MsgShopItemDescInShop_en.strb
	DescriptionSub    string    `json:"-"` // MsgShopItemDescSub_en.strb
	NameEng           string    `json:"-"` // MsgShopItemName_en.strb
	MsgUse            string    `json:"-"` // MsgShopItemUseResult_en.strb
}

// ItemScan searches for an item by ID
func ItemScan(id int) *Item {
	if id > 0 {
		l := len(Data.Items)
		i := sort.Search(l, func(i int) bool { return Data.Items[i].ID >= id })
		if i >= 0 && i < l && Data.Items[i].ID == id {
			return &(Data.Items[i])
		}
	}
	return nil
}

// CleanCustomSkillNoImage cleans VC images and replaces with the element name
func CleanCustomSkillNoImage(name string) string {
	ret := name
	ret = strings.ReplaceAll(ret, "<img=24>", "PASSION")
	ret = strings.ReplaceAll(ret, "<img=25>", "COOL")
	ret = strings.ReplaceAll(ret, "<img=26>", "DARK")
	ret = strings.ReplaceAll(ret, "<img=27>", "LIGHT")
	return ret
}

// CleanCustomSkillImage cleans the VC images and replaces with Wiki tempaltes
func CleanCustomSkillImage(name string) string {
	ret := name
	ret = strings.ReplaceAll(ret, "<img=24>", "{{Passion}}")
	ret = strings.ReplaceAll(ret, "<img=25>", "{{Cool}}")
	ret = strings.ReplaceAll(ret, "<img=26>", "{{Dark}}")
	ret = strings.ReplaceAll(ret, "<img=27>", "{{Light}}")
	return ret
}
