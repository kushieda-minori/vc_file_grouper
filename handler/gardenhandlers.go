package handler

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// StructureListHandler show structures as a list
func StructureListHandler(w http.ResponseWriter, r *http.Request) {
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
<th>Shop Group Deco ID</th>
<th>Enabled</th>
`)
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, s := range vc.Data.Structures {
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
			s.ID,
			s.Name,
			s.Description,
			s.MaxLv,
			s.StructureTypeID,
			s.EventID,
			s.BaseNum,
			s.Stockable,
			s.ShopGroupDecoID,
			s.Enable,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

// StructureDetailHandler show details for a single structure
func StructureDetailHandler(w http.ResponseWriter, r *http.Request) {
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
	structureID, err := strconv.Atoi(pathParts[3])
	if err != nil || structureID < 1 {
		http.Error(w, "Invalid structure id "+pathParts[3], http.StatusNotFound)
		return
	}

	structure := vc.StructureScan(structureID, vc.Data)
	if structure == nil {
		http.Error(w, "Structure not found with id "+pathParts[3], http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "<html><head><title>Structure %d: %s</title></head><body>\n",
		structure.ID,
		structure.Name,
	)
	pc := structure.PurchaseCosts(vc.Data)
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
		printBank(w, structure)
	} else {
		io.WriteString(w, "Not ready yet.")
	}

	io.WriteString(w, "</body></html>")
}

func StructureImagesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "images/garden/map/zip"
	if len(pathParts) < 3 {
		http.Error(w, "Invalid garden path ", http.StatusNotFound)
		return
	}

	gardenBin := VcFilePath + "/garden/map_01.bin"
	os.Stdout.WriteString("reading garden image file\n")
	images, err := vc.ReadBinFileImages(gardenBin)
	limages := len(images)
	os.Stdout.WriteString("read garden image file\n")

	if err != nil {
		http.Error(w, "Error "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(pathParts) >= 4 && pathParts[3] == "zip" {
		// Create a buffer to write our archive to.
		buf := new(bytes.Buffer)

		withDir := false
		if len(pathParts) >= 5 && pathParts[4] == "withDirs" {
			withDir = true
		}

		// Create a new zip archive.
		z := zip.NewWriter(buf)
		// Register a custom Deflate compressor.
		// z.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		// 	return flate.NewWriter(out, flate.BestCompression)
		// })

		// Add some files to the archive.
		for i := 0; i < len(images); i++ {
			image := images[i]
			p := ""
			if withDir {
				p = image.Name[0:strings.Index(image.Name, "_")] + "/"
			}
			f, err := z.Create(p + image.Name)
			if err != nil {
				http.Error(w, "Error "+err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = f.Write(image.Data)
			if err != nil {
				http.Error(w, "Error "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Make sure to check the error on Close.
		err := z.Close()
		if err != nil {
			http.Error(w, "Error "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=\"garden_map01.zip\"")
		w.Header().Set("Content-Type", "application/zip")

		buf.WriteTo(w)
		return
	}

	if len(pathParts) >= 4 && pathParts[3] == "byName" {
		// sort the images by name
		sort.Slice(images, func(i, j int) bool {
			first := images[i]
			second := images[j]

			return first.Name < second.Name
		})
	} else {
		// sort the images by ID
		sort.Slice(images, func(i, j int) bool {
			first := images[i]
			second := images[j]

			return first.ID < second.ID
		})

	}

	io.WriteString(w, "<html><head><title>Structure Images</title>\n")
	io.WriteString(w, "<style>\n"+
		".nav {display:flex;flex-direction:row;flex-wrap:wrap;}\n"+
		".nav a {padding: 5px;}\n"+
		".images{display:flex;flex-direction:row;flex-wrap:wrap;}\n"+
		".image-wrapper{align-self:flex-end;text-align:center;padding:2px;border:1px solid black;}\n"+
		"</style>")
	io.WriteString(w, "</head><body><div class=\"nav\">\n")
	io.WriteString(w, "<a href=\"./\">Sort By ID</a>\n")
	io.WriteString(w, "<a href=\"byName\">Sort By Name</a>\n")
	io.WriteString(w, "<a href=\"zip\">Download All As Zip</a>\n")
	io.WriteString(w, "<a href=\"zip/withDirs\">Download All As Zip With Directories</a>\n")
	io.WriteString(w, "</div>\n<div class=\"images\">")
	for i := 0; i < limages; i++ {
		image := images[i]
		fmt.Fprintf(w,
			"<div class=\"image-wrapper\"><a download=\"%[1]s\" href=\"data:image/png;name=%[1]s;charset=utf-8;base64, %[2]s\"><img src=\"data:image/png;name=%[1]s;charset=utf-8;base64, %[2]s\" /><br/>%[1]s</a></div>\n",
			image.Name,
			base64.StdEncoding.EncodeToString(image.Data),
		)
	}

	io.WriteString(w, "</div></body></html>")
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
		structure.MaxQty(vc.Data),
	)
	castleReq := ""
	if structure.UnlockCastleLv > 0 {
		if structure.UnlockCastleID == 7 { // main castle
			castleReq = fmt.Sprintf("<br />Castle lvl %d", structure.UnlockCastleLv)
		} else if structure.UnlockCastleID == 66 { // ward
			castleReq = fmt.Sprintf("<br />Ward lvl %d", structure.UnlockCastleLv)
		} else {
			castleReq = ""
		}
	}
	areaReq := ""
	if structure.UnlockAreaID > 0 {
		areaReq = fmt.Sprintf("<br />Clear Area %s", vc.Data.Areas[structure.UnlockAreaID].Name)
	}
	io.WriteString(w, lvlHeader)
	levels := structure.Levels(vc.Data)
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

func printBank(w http.ResponseWriter, structure *vc.Structure) {
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
		structure.MaxQty(vc.Data),
	)
	castleReq := ""
	if structure.UnlockCastleLv > 0 {
		if structure.UnlockCastleID == 7 { // main castle
			castleReq = fmt.Sprintf("<br />Castle lvl %d", structure.UnlockCastleLv)
		} else if structure.UnlockCastleID == 66 { // ward
			castleReq = fmt.Sprintf("<br />Ward lvl %d", structure.UnlockCastleLv)
		} else {
			castleReq = ""
		}
	}
	areaReq := ""
	if structure.UnlockAreaID > 0 {
		areaReq = fmt.Sprintf("<br />Clear Area %s", vc.Data.Areas[structure.UnlockAreaID].Name)
	}
	io.WriteString(w, lvlHeader)
	levels := structure.Levels(vc.Data)
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
