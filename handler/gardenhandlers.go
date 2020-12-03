package handler

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"../vc"
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

	structure := vc.StructureScan(structureID)
	if structure == nil {
		http.Error(w, "Structure not found with id "+pathParts[3], http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "<html><head><title>Structure %d: %s</title>"+
		"<style>\n"+
		".nav {display:flex;flex-direction:row;flex-wrap:wrap;}\n"+
		".nav a {padding: 5px;}\n"+
		".images{display:flex;flex-direction:row;flex-wrap:wrap;}\n"+
		".image-wrapper{align-self:flex-end;text-align:center;padding:2px;border:1px solid black;}\n"+
		"</style>"+
		"</head><body>\n<h1>%[2]s</h1>%s<hr />\n",
		structure.ID,
		structure.Name,
		structure.Description,
	)
	pc := structure.PurchaseCosts()
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
				p.Gem,
				p.Cash,
			)
		}
		io.WriteString(w, "\n|}\n</textarea>")
	}
	if structure.IsResource() {
		printResource(w, structure)
	} else if structure.IsBank() {
		printBank(w, structure)
	} else if structure.MaxLv > 1 {
		printGenericStructureLevel(w, structure)
	} else {
		io.WriteString(w, "No details...")
	}

	gardenBin := filepath.Join(vc.FilePath, "garden", "map_01.bin")
	texIds := structure.TextureIDs()
	images, err := vc.GetBinFileImages(gardenBin, texIds...)

	log.Printf("Found %d images for structure %d: %s", len(images), structure.ID, structure.Name)

	io.WriteString(w, "<div class=\"images\">")
	for idx, image := range images {
		io.WriteString(
			w,
			inlineImageTag(
				fmt.Sprintf("%s_%d.png", structure.Name, idx+1),
				&image.Data,
			),
		)
	}

	io.WriteString(w, "</div>\n</body></html>")
}

// StructureImagesHandler show structure images
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

	gardenBin := filepath.Join(vc.FilePath, "garden", "map_01.bin")
	log.Printf("reading garden image file\n")
	images, err := vc.ReadBinFileImages(gardenBin)
	limages := len(images)
	log.Printf("read garden image file\n")

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
				p = image.Name[0:strings.Index(image.Name, "_")]
			}
			f, err := z.Create(filepath.Join(p, image.Name))
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
		io.WriteString(w,
			inlineImageTag(
				image.Name,
				&image.Data,
			),
		)
	}

	io.WriteString(w, "</div></body></html>")
}

func inlineImageTag(imageName string, data *[]byte) string {
	return fmt.Sprintf(
		"<div class=\"image-wrapper\"><a download=\"%[1]s\" href=\"data:image/png;name=%[1]s;charset=utf-8;base64, %[2]s\"><img src=\"data:image/png;name=%[1]s;charset=utf-8;base64, %[2]s\" /><br/>%[1]s</a></div>\n",
		imageName,
		base64.StdEncoding.EncodeToString(*data),
	)
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
!Jewel Cost
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
		structure.MaxQty(),
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
	levels := structure.Levels()
	expTot, goldTot, ethTot, ironTot, gemTot, jewelTot := 0, 0, 0, 0, 0, 0
	for _, l := range levels {
		buildTime := time.Duration(l.Time) * time.Second
		fmt.Fprintf(w, `
|-
| %d || Level %d%s%s || %d || %d/min || %d || %d || %d || %d || %s || %s || %d`,
			l.Level,
			l.LevelCap,
			castleReq,
			areaReq,
			l.Resource.Income,
			l.Resource.Rate(),
			l.Coin,
			l.Ether,
			l.Iron,
			l.Cash,
			buildTime,
			l.Resource.FillTime(),
			l.Exp,
		)
		goldTot += l.Coin
		ethTot += l.Ether
		ironTot += l.Iron
		gemTot += l.Gem
		jewelTot += l.Cash
		expTot += l.Exp
	}
	// summary line
	fmt.Fprintf(w, `
|-
!Total
!colspan=3|
!%d !!%d !!%d !!%d !! !! !!%d`,
		goldTot,
		ethTot,
		ironTot,
		jewelTot,
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
		structure.MaxQty(),
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
	levels := structure.Levels()
	expTot, goldTot, ethTot, ironTot, gemTot, jewelTot := 0, 0, 0, 0, 0, 0
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
		gemTot += l.Gem
		jewelTot += l.Cash
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

func printGenericStructureLevel(w http.ResponseWriter, structure *vc.Structure) {
	levels := structure.Levels()
	if levels == nil || len(levels) == 0 {
		return
	}

	expTot, goldTot, ethTot, ironTot, gemTot, jewelTot, maxEffectParams := 0, 0, 0, 0, 0, 0, 0
	buildTot := time.Duration(0)
	hasSpecialEffect := false
	for _, l := range levels {
		buildTime := time.Duration(l.Time) * time.Second
		buildTot = buildTot + buildTime
		goldTot += l.Coin
		ethTot += l.Ether
		ironTot += l.Iron
		expTot += l.Exp
		gemTot += l.Gem
		jewelTot += l.Cash
		hasSpecialEffect = hasSpecialEffect || l.SpecialEffect != nil
		effect := l.SpecialEffect
		if effect != nil {
			for i, p := range []int{effect.Param1, effect.Param2, effect.Param3, effect.Param4} {
				if p > 0 {
					maxEffectParams = int(math.Max(float64(i), float64(maxEffectParams)))
				}
			}
		}
	}

	lvlHeader := `{| class="mw-collapsible mw-collapsed article-table" style="min-width:677px"
|-
!Lvl !!Requirement`
	if hasSpecialEffect {
		lvlHeader += ` !!Effect`
	}
	if goldTot > 0 {
		lvlHeader += ` !!Gold Cost`
	}
	if ethTot > 0 {
		lvlHeader += ` !!Ether Cost`
	}
	if ironTot > 0 {
		lvlHeader += ` !!Iron Cost`
	}
	if gemTot > 0 {
		lvlHeader += ` !!Gem Cost`
	}
	if jewelTot > 0 {
		lvlHeader += ` !!Jewel Cost`
	}
	lvlHeader += ` !!Build Time !!Exp`

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
		structure.MaxQty(),
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
	for _, l := range levels {
		if l.Level > structure.MaxLv {
			break
		}
		buildTime := time.Duration(l.Time) * time.Second

		resources := ""
		if hasSpecialEffect {
			effect := l.SpecialEffect
			for i, p := range []int{effect.Param1, effect.Param2, effect.Param3, effect.Param4} {
				if i <= maxEffectParams {
					resources += fmt.Sprintf(" || %d", p)
				}
			}
		}
		if goldTot > 0 {
			resources += fmt.Sprintf(" || %d",
				l.Coin,
			)
		}
		if ethTot > 0 {
			resources += fmt.Sprintf(" || %d",
				l.Ether,
			)
		}
		if ironTot > 0 {
			resources += fmt.Sprintf(" || %d",
				l.Iron,
			)
		}
		if gemTot > 0 {
			resources += fmt.Sprintf(" || %d",
				l.Gem,
			)
		}
		if jewelTot > 0 {
			resources += fmt.Sprintf(" || %d",
				l.Cash,
			)
		}

		fmt.Fprintf(w, `
|-
| %d || Level %d%s%s%s || %s || %d`,
			l.Level,
			l.LevelCap,
			castleReq,
			areaReq,
			resources,
			buildTime,
			l.Exp,
		)
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
