package handler

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// DeckBonusHandler show deck bonuses as a table
func DeckBonusHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html><head><title>Deck Bonuses</title>
<style>table, th, td {border: 1px solid black;};</style>
</head><body>
<div>
<a href="WIKI/">Wiki Formatted</a>
<table><thead><tr>
  <th>_id</th>
  <th>Name</th>
  <th>Description</th>
  <th>Atk/Def</th>
  <th>Value</th>
  <th>Down Grade</th>
  <th>Cond Type</th>
  <th>Cards Req.</th>
  <th>Dups?</th>
  <th>Conditions</th>
</tr></thead>
<tbody>`)

	//sort.Sort(vc.DeckBonusByCountAndName(vc.Data.DeckBonuses))

	for _, d := range vc.Data.DeckBonuses {
		fmt.Fprintf(w, `<tr>
  <td>%d</td>
  <td>%s</td>
  <td>%s</td>
  <td>%d</td>
  <td>%d</td>
  <td>%d</td>
  <td>%d</td>
  <td>%d</td>
  <td>%d</td>
  <td>%s</td>
</tr>`,
			d.ID,
			d.Name,
			d.Description,
			d.AtkDefFlg,
			d.Value,
			d.DownGrade,
			d.CondType,
			d.ReqNum,
			d.DupFlg,
			d.Conditions(vc.Data),
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

// DeckBonusWikiHandler show deck bonuses as wiki formatted
func DeckBonusWikiHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html><head><title>Deck Bonuses</title>
<style>table, th, td {border: 1px solid black;};</style>
</head><body>
<div>
<a href="../">HTML</a>
<textarea readonly="readonly" style="width:100%;height:450px">
`)

	tableHeader := `{| class="article-table" border="1"
!Name of Bonus
!Effect
!Eligible Cards
`
	tableFooter := `
|}

`

	sort.Sort(vc.DeckBonusByCountAndName(vc.Data.DeckBonuses))

	reg := regexp.MustCompile(`\[|【(.+)\]|】\n?(.*)`)

	oldReq := -1

	for _, d := range vc.Data.DeckBonuses {
		if oldReq != d.ReqNum {
			if oldReq > 1 {
				io.WriteString(w, tableFooter)
			}
			io.WriteString(w, tableHeader)
			oldReq = d.ReqNum
		}
		descMatch := reg.FindStringSubmatch(d.Description)

		dca := d.Conditions(vc.Data)
		switch d.CondType {
		case 2:
			fmt.Fprintf(w, `|-
|'''%s'''<br />"%s"
|%s
|`,
				d.Name,
				descMatch[2],
				descMatch[1],
			)
			strs := make([]string, len(dca))
			for i, dc := range dca {
				strs[i] = fmt.Sprintf("[[%s]]", dc.RefName)
			}
			strs = removeDuplicates(strs)
			sort.Strings(strs)
			io.WriteString(w, strings.Join(strs, ", "))
		case 3, 8:
			fmt.Fprintf(w, `|-
|'''%s'''
|%s
|%s`,
				d.Name,
				descMatch[1],
				descMatch[2],
			)
		}
		io.WriteString(w, "\n")
	}
	io.WriteString(w, tableFooter)
	io.WriteString(w, "</textarea></div></body></html>")
}
