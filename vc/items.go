package vc

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

//GetImageData gets the image data for the item
func (i Item) GetImageData() (imageName string, data []byte, fsInfo fs.FileInfo, err error) {
	if i.ItemNo < 1 || i.NameEng == "" {
		return
	}
	var path string = filepath.Join(FilePath, "item", "shop")
	fileName := fmt.Sprintf("%d", i.ItemNo)
	fullpath := filepath.Join(path, fileName)
	if fsInfo, err = os.Stat(fullpath); os.IsNotExist(err) {
		log.Printf("Unable to find Item %s : %s", i.NameEng, fileName)
		return
	}
	data, err = Decode(fullpath)
	if err != nil {
		return
	}
	imageName = strings.ReplaceAll(CleanCustomSkillNoImage(i.NameEng), "/", "-") + ".png"
	return
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

var ItemGroups = map[int]string{
	1:  "Resource",
	2:  "Resource",
	3:  "Resource",
	4:  "Other",
	5:  "Enhancement",
	6:  "Recovery",
	7:  "Recovery",
	8:  "Deck Enhancement",
	9:  "Arcana",
	10: "Arcana",
	11: "Arcana",
	12: "Arcana",
	13: "Arcana",
	14: "Arcana",
	15: "Arcana",
	16: "Arcana",
	17: "Tickets",
	18: "Exchange",
	19: "Exchange",
	20: "Event",
	22: "Recovery",
	23: "Event",
	24: "Deck Enhancement",
	25: "Event",

	29: "Awakening and Rebirth",
	30: "Arcana",
	31: "Recovery",
	32: "Exchange",

	36: "Recovery",
	37: "Resource",
	38: "Custom Skill Recipe",
	39: "Custom Skill",
	40: "Other",

	42: "Other",
	43: "Exchange",
	44: "Other",

	47: "Recovery",
	48: "Recovery",

	50: "Recovery",
	51: "Exchange",
	52: "Memory Core",
	53: "Other",
	54: "Other",
	55: "Enhancement",
	56: "Enhancement",
	57: "Enhancement",
}
