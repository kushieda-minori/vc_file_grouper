package vc

import (
	"sort"
	"strings"
)

//"items" is the list of Archwitches
type Item struct {
	Id                int       `json:"_id"`
	GroupId           int       `json:"group_id"`
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

func ItemScan(id int, items []Item) *Item {
	if id > 0 {
		l := len(items)
		i := sort.Search(l, func(i int) bool { return items[i].Id >= id })
		if i >= 0 && i < l && items[i].Id == id {
			return &(items[i])
		}
	}
	return nil
}

func CleanCustomSkillNoImage(name string) string {
	ret := name
	ret = strings.Replace(ret, "<img=24>", "PASSION", -1)
	ret = strings.Replace(ret, "<img=25>", "COOL", -1)
	ret = strings.Replace(ret, "<img=26>", "DARK", -1)
	ret = strings.Replace(ret, "<img=27>", "LIGHT", -1)
	return ret
}

func CleanCustomSkillImage(name string) string {
	ret := name
	ret = strings.Replace(ret, "<img=24>", "{{Passion}}", -1)
	ret = strings.Replace(ret, "<img=25>", "{{Cool}}", -1)
	ret = strings.Replace(ret, "<img=26>", "{{Dark}}", -1)
	ret = strings.Replace(ret, "<img=27>", "{{Light}}", -1)
	return ret
}
