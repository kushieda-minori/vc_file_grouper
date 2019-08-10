package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	io.WriteString(w, "<th>ID</th><th>Weapon Names</th><th>Descriptions</th><th>Rank Group</th><th>Rarity Group</th><th>Status ID</th>\n")
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
</tr>`,
			wp.ID,
			strings.Join(wp.Name[:], "<br />"),
			strings.Join(wp.Description[:], "<br />"),
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

	var prevWeapon, nextWeapon *vc.Weapon = nil, nil

	prevWeaponName := ""
	prevWeapon = vc.WeaponScan(weaponID - 1)
	if prevWeapon != nil {
		prevWeaponName = prevWeapon.Name[3]
	}

	nextWeaponName := ""
	nextWeapon = vc.WeaponScan(weaponID - 1)
	if nextWeapon != nil {
		nextWeaponName = nextWeapon.Name[3]
	}

	weaponName := weapon.Name[3]

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", weaponName)

	if prevWeaponName != "" {
		fmt.Fprintf(w, "<div style=\"float:left; width: 33%%;\"><a href=\"%d\">&lt;&lt; %s &lt;&lt;</a></div>\n",
			prevWeapon.ID,
			prevWeaponName,
		)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}
	fmt.Fprint(w, "<div style=\"float:left; width: 33%;text-align:center;\"><a href=\"../\">All Weapons</a></div>\n")
	if nextWeaponName != "" {
		fmt.Fprintf(w, "<div style=\"float:right; width: 33%%;;text-align:right;\"><a href=\"%d\">&gt;&gt; %s &gt;&gt;</a></div>\n",
			nextWeapon.ID,
			nextWeaponName,
		)
	} else {
		fmt.Fprint(w, "<div style=\"float:left; width: 33%;\"></div>\n")
	}

	// stats and stuff
	io.WriteString(w, "<div style=\"clear:both;\">")
	status := weapon.Status()

	//general stats
	printHTMLTable(w,
		"General Weapon Stats",
		[]string{"Stat", "min", "max"},
		[][]interface{}{
			{"Attack", status.AtkMin, status.AtkMax},
			{"Defense", status.DefMin, status.DefMax},
			{"Soldiers", status.SoldiersMin, status.SoldiersMax},
		},
	)

	// rarities
	rows := make([][]interface{}, 0)
	for _, rarity := range weapon.Rarities() {
		rows = append(rows, []interface{}{rarity.Rarity, rarity.UnlockRank})
	}

	printHTMLTable(w,
		"Rarity Unlocks",
		[]string{"Rarity", "Unlocked at Rank"},
		rows,
	)

	// skills
	rows = make([][]interface{}, 0)
	for _, skill := range weapon.SkillUnlocks() {
		rows = append(rows, []interface{}{skill.Skill().TypeName(), skill.SkillLevel, skill.UnlockRank, skill.Skill().DescriptionFormatted()})
	}

	printHTMLTable(w,
		"Skill Unlocks",
		[]string{"Skill Type", "Skill Level", "Unlocked at Rank", "Description"},
		rows,
	)

	// upgrade materials
	rows = make([][]interface{}, 0)
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
		"Upgrade Material",
		[]string{"Item", "Apply to Rarity", "Exp Given"},
		rows,
	)

	// ranks
	rows = make([][]interface{}, 0)
	for _, rank := range weapon.Ranks() {
		rows = append(rows, []interface{}{rank.Rank, rank.NeedExp, rank.Gold, rank.Iron, rank.Ether, rank.Gem})
	}

	printHTMLTable(w,
		"Ranks",
		[]string{"Rank", "Exp Needed", "Gold", "Iron", "Ether", "Gem"},
		rows,
	)

	io.WriteString(w, "</div>")

	// thumbs
	io.WriteString(w, "<div style=\"clear:both;\">")
	io.WriteString(w, "<h2>Weapon Image Icons</h2>")
	for i := 1; i <= 4; i++ {
		fmt.Fprintf(w, `<a href="/images/weapon/thumb/wp_%05[1]d_%02[2]d?filename=%[3]s_%[2]d_icon"><img src="/images/weapon/thumb/wp_%05[1]d_%02[2]d" alt="Thumbnail"/></a>`,
			weapon.ID,
			i,
			url.QueryEscape(weapon.Name[3]),
		)
	}
	io.WriteString(w, "</div>")

	// full images
	io.WriteString(w, "<div style=\"clear:both;\">")
	io.WriteString(w, "<h2>Weapon Images</h2>")
	for i := 1; i <= 4; i++ {
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
			url.QueryEscape(weapon.Name[3]),
			pathPart,
		)
	}
	io.WriteString(w, "</div>")

	io.WriteString(w, "</body></html>")

}
