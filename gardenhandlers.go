package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func structureListHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Structures</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, `<th>_id</th>
<th>Name</th>
<th>Description</th>
<th>Max Level</th>
<th>Type</th>
<th>Event</th>
<th>Base Own</th>
<th>Stockable</th>
<th>Shop Group Deco Id</th>
<th>Enabled</th>
`)
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, s := range VcData.Structures {
		fmt.Fprintf(w, `<tr>
	<td><a href="/garden/structures/detail/%[1]d">%[1]d</a></td>
	<td>%[2]s</td>
	<td>%s</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
</tr>`,
			s.Id,
			s.Name,
			s.Description,
			s.MaxLv,
			s.StructureTypeId,
			s.EventId,
			s.BaseNum,
			s.Stockable,
			s.ShopGroupDecoId,
			s.Enable,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

func structureDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "garden/structures/detail/id"
	if len(pathParts) < 4 {
		http.Error(w, "Invalid structure id ", http.StatusNotFound)
		return
	}
	structureId, err := strconv.Atoi(pathParts[3])
	if err != nil || structureId < 1 {
		http.Error(w, "Invalid structure id "+pathParts[3], http.StatusNotFound)
		return
	}

	structure := vc.StructureScan(structureId, VcData)
	if structure == nil {
		http.Error(w, "Structure not found with id "+pathParts[3], http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "<html><head><title>Structure %d: %s</title></head><body>\n",
		structure.Id,
		structure.Name,
	)
	pc := structure.PurchaseCosts(VcData)
	if len(pc) > 0 {
		io.WriteString(w, "\nPurchase Costs<br/><textarea rows=\"25\" cols=\"80\">")
		io.WriteString(w, `{| class="article-table"
	!Num !! Gold !! Ether !! Iron !! Gem !! Jewels`)
		for _, p := range pc {
			fmt.Fprintf(w, `
|-
|%d || %d || %d || %d || %d || %d
`,
				p.Num,
				p.Coin,
				p.Iron,
				p.Ether,
				p.Elixir,
				p.Cash,
			)
		}
		io.WriteString(w, "\n|}\n</textarea>")
	}
	if structure.IsResource() {
		printResource(w, structure)
	} else if structure.IsBank() {
		printResourceStorage(w, structure)
	} else {
		io.WriteString(w, "Not ready yet.")
	}

	io.WriteString(w, "</body></html>")
}

func printResource(w http.ResponseWriter, structure *vc.Structure) {
	lvlHeader := `{| class="mw-collapsible mw-collapsed article-table" style="min-width:677px"
|-
!Lvl
!Requirement
!Stock
!Rate
!Gold Cost
!Ether Cost
!Iron Cost
!Build Time
!Stock Fill Time
!Exp`

	io.WriteString(w, "\n<br />Levels<br/><textarea rows=\"25\" cols=\"80\">")
	fmt.Fprintf(w, `=== %s ===
[[File:%[1]s.png|thumb|right]]
%[2]s

Size: %dx%d

Max quantity: %d
`,
		structure.Name,
		structure.Description,
		structure.SizeX,
		structure.SizeY,
		structure.MaxQty(VcData),
	)
	castleReq := ""
	if structure.UnlockCastleLv > 0 {
		if structure.UnlockCastleId == 7 { // main castle
			castleReq = fmt.Sprintf("<br />Castle lvl %d", structure.UnlockCastleLv)
		} else if structure.UnlockCastleId == 66 { // ward
			castleReq = fmt.Sprintf("<br />Ward lvl %d", structure.UnlockCastleLv)
		} else {
			castleReq = ""
		}
	}
	areaReq := ""
	if structure.UnlockAreaId > 0 {
		areaReq = fmt.Sprintf("<br />Clear Area %s", VcData.Areas[structure.UnlockAreaId].Name)
	}
	io.WriteString(w, lvlHeader)
	levels := structure.Levels(VcData)
	expTot, goldTot, ethTot, ironTot := 0, 0, 0, 0
	for _, l := range levels {
		buildTime := time.Duration(l.Time) * time.Second
		fmt.Fprintf(w, `
|-
| %d || Level %d%s%s || %d || %d/min || %d || %d || %d || %s || %s || %d`,
			l.Level,
			l.LevelCap,
			castleReq,
			areaReq,
			l.Resource.Income,
			l.Resource.Rate(),
			l.Coin,
			l.Ether,
			l.Iron,
			buildTime,
			l.Resource.FillTime(),
			l.Exp,
		)
		goldTot += l.Coin
		ethTot += l.Ether
		ironTot += l.Iron
		expTot += l.Exp
	}
	// summary line
	fmt.Fprintf(w, `
|-
!Total
!colspan=3|
!%d !!%d !!%d !! !! !!%d`,
		goldTot,
		ethTot,
		ironTot,
		expTot,
	)
	io.WriteString(w, "\n|}\n</textarea>")
}

func printResourceStorage(w http.ResponseWriter, structure *vc.Structure) {
	lvlHeader := `{| class="mw-collapsible mw-collapsed article-table" style="min-width:677px"
|-
!Lvl
!Requirement
!Stock
!Gold Cost
!Ether Cost
!Iron Cost
!Build Time
!Exp`

	io.WriteString(w, "\n<br />Levels<br/><textarea rows=\"25\" cols=\"80\">")
	fmt.Fprintf(w, `=== %s ===
[[File:%[1]s.png|thumb|right]]
%[2]s

Size: %dx%d

Max quantity: %d
`,
		structure.Name,
		structure.Description,
		structure.SizeX,
		structure.SizeY,
		structure.MaxQty(VcData),
	)
	castleReq := ""
	if structure.UnlockCastleLv > 0 {
		if structure.UnlockCastleId == 7 { // main castle
			castleReq = fmt.Sprintf("<br />Castle lvl %d", structure.UnlockCastleLv)
		} else if structure.UnlockCastleId == 66 { // ward
			castleReq = fmt.Sprintf("<br />Ward lvl %d", structure.UnlockCastleLv)
		} else {
			castleReq = ""
		}
	}
	areaReq := ""
	if structure.UnlockAreaId > 0 {
		areaReq = fmt.Sprintf("<br />Clear Area %s", VcData.Areas[structure.UnlockAreaId].Name)
	}
	io.WriteString(w, lvlHeader)
	levels := structure.Levels(VcData)
	expTot, goldTot, ethTot, ironTot := 0, 0, 0, 0
	for _, l := range levels {
		buildTime := time.Duration(l.Time) * time.Second
		fmt.Fprintf(w, `
|-
| %d || Level %d%s%s || %d || %d || %d || %d || %s || %d`,
			l.Level,
			l.LevelCap,
			castleReq,
			areaReq,
			l.Bank.Value,
			l.Coin,
			l.Ether,
			l.Iron,
			buildTime,
			l.Exp,
		)
		goldTot += l.Coin
		ethTot += l.Ether
		ironTot += l.Iron
		expTot += l.Exp
	}
	// summary line
	fmt.Fprintf(w, `
|-
!Total
!colspan=2|
!%d !!%d !!%d !! !!%d`,
		goldTot,
		ethTot,
		ironTot,
		expTot,
	)
	io.WriteString(w, "\n|}\n</textarea>")
}
