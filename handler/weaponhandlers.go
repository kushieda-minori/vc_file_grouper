package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// WeaponHandler handle weapon information
func WeaponHandler(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "<html><head><title>All Weapons</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>ID</th><th>Weapon Names</th><th>Descriptions</th><th>Max Rarity</th><th>Max Rank</th><th>Rank Group</th><th>Rarity Group</th><th>Status ID</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for i := len(vc.Data.Weapons) - 1; i >= 0; i-- {
		wp := vc.Data.Weapons[i]
		fmt.Fprintf(w, `<tr>
	<td><a href="/weapons/detail/%[1]d">%[1]d</a></td>
	<td><a href="/weapons/detail/%[1]d">%[2]s</a></td>
	<td>%[3]s</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
</tr>`,
			wp.ID,
			strings.Join(wp.Names, "<br />"),
			strings.Join(wp.Descriptions, "<br />"),
			wp.MaxRarity(),
			wp.MaxRank(),
			wp.RankGroupID,
			wp.RarityGroupID,
			wp.StatusID,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

// WeaponDetailHandler show details for a single weapon
func WeaponDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "weapons/detail/id"
	if len(pathParts) < 3 {
		http.Error(w, "Invalid weapon id ", http.StatusNotFound)
		return
	}
	weaponID, err := strconv.Atoi(pathParts[2])
	if err != nil || weaponID < 1 {
		http.Error(w, "Invalid weapon id "+pathParts[2], http.StatusNotFound)
		return
	}

	weapon := vc.WeaponScan(weaponID)

	if weapon == nil {
		http.Error(w, "Invalid weapon id "+pathParts[2], http.StatusNotFound)
		return
	}

	var prevWeapon, nextWeapon *vc.Weapon = nil, nil

	prevWeaponName := ""
	prevWeapon = vc.WeaponScan(weaponID - 1)
	if prevWeapon != nil {
		prevWeaponName = prevWeapon.MaxRarityName()
	}

	nextWeaponName := ""
	nextWeapon = vc.WeaponScan(weaponID - 1)
	if nextWeapon != nil {
		nextWeaponName = nextWeapon.MaxRarityName()
	}

	weaponName := weapon.MaxRarityName()

	fmt.Fprintf(w, `<html>
<head>
	<title>%s</title>
	<style>
		table, th, td {border:1px solid black; padding: 2px;}
		table {border-collapse: collapse; padding-bottom: 10px; margin-bottom: 20px; margin-right: 15px;}
	</style>
</head>
<body>
	<h1>%[1]s</h1>
`, weaponName)

	if prevWeaponName != "" {
		fmt.Fprintf(w, "<div style=\"float:left; width: 33%%;\"><a href=\"%[1]d%[3]s\">&lt;&lt; %[2]s &lt;&lt;</a></div>\n",
			prevWeapon.ID,
			prevWeaponName,
			addQueryMark(r.URL.RawQuery),
		)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}
	fmt.Fprint(w, "<div style=\"float:left; width: 33%;text-align:center;\"><a href=\"../\">All Weapons</a></div>\n")
	if nextWeaponName != "" {
		fmt.Fprintf(w, "<div style=\"float:right; width: 33%%;text-align:right;\"><a href=\"%[1]d%[3]s\">&gt;&gt; %[2]s &gt;&gt;</a></div>\n",
			nextWeapon.ID,
			nextWeaponName,
			addQueryMark(r.URL.RawQuery),
		)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}

	qs := r.URL.Query()
	qWiki := qs.Get("wiki") != ""
	if qWiki {
		writeWeaponWiki(w, weapon)
	} else {

		fmt.Fprintf(w, "<div><a href=\"./%d?wiki=1\">Wiki View</a></div>\n", weapon.ID)
		// stats and stuff
		io.WriteString(w, "<div style=\"clear:both;\">")

		// Overview
		printWeaponConfig(w, weapon)

		io.WriteString(w, "<div style=\"clear:both;float:left;\">")
		printWeaponStatus(w, weapon)
		printWeaponRarity(w, weapon)
		io.WriteString(w, "</div>")

		printWeaponSkills(w, weapon)

		io.WriteString(w, "</div>")

		// thumbs
		io.WriteString(w, "<div style=\"clear:both;\">")
		io.WriteString(w, "<h2>Weapon Image Icons</h2>")
		printWeaponIcons(w, weapon)
		io.WriteString(w, "</div>")

		// full images
		io.WriteString(w, "<div style=\"clear:both;\">")
		io.WriteString(w, "<h2>Weapon Images</h2>")
		printWeaponImages(w, weapon)
		io.WriteString(w, "</div>")

		io.WriteString(w, "<div style=\"clear:both;\">")
		printWeaponUpgradeMaterials(w, weapon)
		printWeaponRanks(w, weapon)

		io.WriteString(w, "</div>")
	}
	io.WriteString(w, "</body></html>")

}

func printWeaponConfig(w io.Writer, weapon *vc.Weapon) {
	printHTMLTable(w,
		"float: left;",
		"Weapon configuration",
		[]string{"Config", "Value"},
		[][]interface{}{
			{"Names", strings.Join(weapon.Names, "<br />")},
			{"Descriptions", strings.Join(weapon.Descriptions, "<br />")},
			{"Max Rarity", weapon.MaxRarity()},
			{"Max Rank", weapon.MaxRank()},
			{"Rank Group", weapon.RankGroupID},
			{"Rarity Group", weapon.RarityGroupID},
			{"Status", fmt.Sprintf("%d-%s", weapon.StatusID, weapon.StatusDescription())},
		},
	)
}

func printWeaponStatus(w io.Writer, weapon *vc.Weapon) {
	//general stats
	status := weapon.Status()
	printHTMLTable(w,
		"",
		fmt.Sprintf("General Weapon Stats - %d: %s", weapon.StatusID, weapon.StatusDescription()),
		[]string{"Stat", "min", "max"},
		[][]interface{}{
			{"Attack", status.AtkMin, status.AtkMax},
			{"Defense", status.DefMin, status.DefMax},
			{"Soldiers", status.SoldiersMin, status.SoldiersMax},
		},
	)
}

func printWeaponRarity(w io.Writer, weapon *vc.Weapon) {
	// rarities
	rows := make([][]interface{}, 0)
	for _, rarity := range weapon.Rarities() {
		rows = append(rows, []interface{}{rarity.Rarity, rarity.UnlockRank})
	}

	printHTMLTable(w,
		"",
		fmt.Sprintf("Rarity Unlocks for Rarity Group: %d", weapon.RarityGroupID),
		[]string{"Rarity", "Unlocked at Rank"},
		rows,
	)
}

func printWeaponSkills(w io.Writer, weapon *vc.Weapon) {
	// skills
	rows := make([][]interface{}, 0)
	for _, skill := range weapon.SkillUnlocks() {
		rows = append(rows, []interface{}{skill.Skill().TypeName(), skill.SkillLevel, skill.UnlockRank, skill.Skill().DescriptionFormatted()})
	}

	printHTMLTable(w,
		"float: left;",
		"Skill Unlocks",
		[]string{"Skill Type", "Skill Level", "Unlocked at Rank", "Description"},
		rows,
	)
}

func printWeaponUpgradeMaterials(w io.Writer, weapon *vc.Weapon) {
	// upgrade materials
	rows := make([][]interface{}, 0)
	for _, material := range weapon.UpgradeMaterials() {
		item := material.Item()
		itemImg := fmt.Sprintf("<a href=\"/images/item/shop/%[1]d?filename=%[2]s\"><img src=\"/images/item/shop/%[1]d\"/><br/>%[3]s</a>",
			item.ItemNo,
			url.QueryEscape(vc.CleanCustomSkillNoImage(item.NameEng)),
			item.NameEng,
		)
		rows = append(rows, []interface{}{itemImg, material.Rarity, material.Exp})
	}

	printHTMLTable(w,
		"float: left;",
		"Upgrade Material",
		[]string{"Item", "Rarity Of Item", "Exp Given"},
		rows,
	)
}

func printWeaponRanks(w io.Writer, weapon *vc.Weapon) {
	// ranks
	rows := make([][]interface{}, 0)
	for _, rank := range weapon.Ranks() {
		rows = append(rows, []interface{}{rank.Rank, rank.NeedExp, rank.Gold, rank.Iron, rank.Ether, rank.Gem})
	}

	printHTMLTable(w,
		"float: left;",
		fmt.Sprintf("Ranks for Rank Group: %d", weapon.RankGroupID),
		[]string{"Rank", "Exp Needed", "Gold", "Iron", "Ether", "Gem"},
		rows,
	)
}

func printWeaponIcons(w io.Writer, weapon *vc.Weapon) {
	rlen := weapon.MaxRarity()
	name := url.QueryEscape(weapon.MaxRarityName())
	for i := 1; i <= rlen; i++ {
		fmt.Fprintf(w, `<a href="/images/weapon/thumb/wp_%05[1]d_%02[2]d?filename=%[3]s_%[2]d_icon"><img src="/images/weapon/thumb/wp_%05[1]d_%02[2]d" alt="Thumbnail"/></a>`,
			weapon.ID,
			i,
			name,
		)
	}
}

func printWeaponImages(w io.Writer, weapon *vc.Weapon) {
	rlen := weapon.MaxRarity()
	name := url.QueryEscape(weapon.MaxRarityName())
	for i := 1; i <= rlen; i++ {
		pathPart := ""
		imgName := fmt.Sprintf("wp_%05[1]d_%02[2]d", weapon.ID, i)
		if _, err := os.Stat(vc.FilePath + "weapon/hd/" + imgName); err == nil {
			pathPart = "hd"
		} else if _, err := os.Stat(vc.FilePath + "weapon/md/" + imgName); err == nil {
			pathPart = "md"
		} else {
			pathPart = "sd"
		}

		fmt.Fprintf(w, `<a href="/images/weapon/%[4]s/wp_%05[1]d_%02[2]d?filename=%[3]s_%[2]d"><img src="/images/weapon/%[4]s/wp_%05[1]d_%02[2]d" alt="image"/></a>`,
			weapon.ID,
			i,
			name,
			pathPart,
		)
	}
}

func writeWeaponWiki(w io.Writer, weapon *vc.Weapon) {
	fmt.Fprintf(w, "<div><a href=\"./%d\">Data View</a></div>\n", weapon.ID)
	fmt.Fprintf(w,
		"<div style=\"clear:both;\">Edit on <a href=\"https://valkyriecrusade.fandom.com/wiki/%s?action=edit\">fandom</a>\n<br /></div>",
		weapon.MaxRarityName(),
	)
	io.WriteString(w, "<textarea style=\"width:90%;height:760px\">")
	eventNames := weapon.EventNames()
	availability := ""
	if len(eventNames) > 0 {
		availability = "[[" + strings.Join(eventNames, "]]<br />[[") + "]]"
	}

	skills := weapon.SkillUnlocks()
	fmt.Fprintf(w, `{{Weapon
|status = %d
|rarity group = %d
|rank group = %d
<!-- descriptions for rarities -->
%s
<!-- skill defs are in the format of SkillType=rank1,rank2,rankX... and each skill is separated by !! -->
|skill defs = %s
<!-- events the weapon appeared in -->
|availability = %s
}}
`,
		weapon.StatusID,
		weapon.RarityGroupID,
		weapon.RankGroupID,
		formatWeaponWikiTemplateDescriptions(weapon.Descriptions), // rarity descriptions
		formatSkillArray(skills),
		availability, // weapon event link titles
	)
	io.WriteString(w, "</textarea>")
}

func formatWeaponWikiTemplateDescriptions(descriptions []string) string {
	ret := ""
	for i, desc := range descriptions {
		ret += fmt.Sprintf("|description %d = %s\n", i+1, desc)
	}
	return ret
}

func formatSkillArray(skills vc.WeaponSkillUnlockRankList) string {
	wsg := make(weaponSkillGroup)
	skillOrder := []int{}
	for _, skill := range skills {
		if list, ok := wsg[skill.SkillType]; ok {
			wsg[skill.SkillType] = append(list, skill.UnlockRank)
		} else {
			newlist := make(weaponSkillInfoList, 1)
			newlist[0] = skill.UnlockRank
			wsg[skill.SkillType] = newlist
			skillOrder = append(skillOrder, skill.SkillType)
		}
	}
	return wsg.String(skillOrder)
}

type weaponSkillGroup map[int]weaponSkillInfoList

type weaponSkillInfoList []int

func (wsg *weaponSkillGroup) String(keyOrder []int) string {
	if wsg == nil || len(*wsg) == 0 {
		return ""
	}
	ret := ""
	dlen := len(vc.WeaponSkillTypes)
	for i, k := range keyOrder {
		v := (*wsg)[k]
		sep := ""
		if i > 0 {
			sep = " !!"
		}
		if k < dlen {
			ret += fmt.Sprintf("%s%s=%s", sep, vc.WeaponSkillTypes[k], v.String())
		} else {
			ret += fmt.Sprintf("%s%d=%s", sep, k, v.String())
		}
		i++
	}
	return ret
}

func (wsil *weaponSkillInfoList) String() string {
	if wsil == nil || len(*wsil) == 0 {
		return ""
	}
	l := len(*wsil)
	asString := make([]string, l, l)
	// make sure the ranks are in order (although they should be)
	sort.Ints(*wsil)
	for i, v := range *wsil {
		asString[i] = strconv.Itoa(v)
	}
	return strings.Join(asString, ",")
}
