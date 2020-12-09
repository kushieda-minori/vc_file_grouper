package handler

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vc_file_grouper/vc"
)

// MapHandler handle Map requests
func MapHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "maps/id/WIKI"
	if len(pathParts) < 2 {
		MapTableHandler(w, r)
		return
	}

	mapID, err := strconv.Atoi(pathParts[1])
	if err != nil || mapID < 1 || mapID > len(vc.Data.Maps) {
		http.Error(w, "Invalid map id "+pathParts[1], http.StatusNotFound)
		return
	}
	m := vc.MapScan(mapID, vc.Data.Maps)

	if m == nil {
		http.Error(w, "Invalid map id "+pathParts[1], http.StatusNotFound)
		return
	}

	if len(pathParts) >= 3 && "WIKI" == pathParts[2] {
		MapDetailWikiHandler(w, r, m)
		return
	}

	MapDetailHandler(w, r, m)
}

// MapDetailHandler show details for a single map
func MapDetailHandler(w http.ResponseWriter, r *http.Request, m *vc.Map) {
	fmt.Fprintf(w, "<html><head><title>Map %s</title>\n", m.Name)
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	fmt.Fprintf(w, "<h1>%s</h1>\n%s", m.Name, m.StartMsg)
	fmt.Fprintf(w, "<p><a href=\"/maps/%d/WIKI\">Wiki Formatted</a></p>", m.ID)
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>No</th><th>Name</th><th>Long Name</th><th>Start</th><th>End</th><th>Story</th><th>Boss Start</th><th>Boss End</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, e := range m.Areas(vc.Data) {
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			e.AreaNo,
			e.Name,
			e.LongName,
			e.Start,
			e.End,
			e.Story,
			e.BossStart,
			e.BossEnd,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

// MapDetailWikiHandler show map details but in wiki format
func MapDetailWikiHandler(w http.ResponseWriter, r *http.Request, m *vc.Map) {
	fmt.Fprintf(w, "<html><head><title>Map %s</title>\n", m.Name)
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	fmt.Fprintf(w, `<h1>%s</h1>
<p><a href="../%d/WIKI">prev</a> &nbsp; <a href="../%d/WIKI">next</a></p>
<p>%s</p>`,
		m.Name,
		m.ID-1,
		m.ID+1,
		m.StartMsg,
	)
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<textarea style=\"width:800px;height:760px\">\n")
	fmt.Fprintf(w, `{{#tag:gallery|
Banner {{#titleparts:{{PAGENAME}}|1}}.png
AreaMap %s.png
BattleBG %d.png
|type="slider"
|widths="680"
|position="left"
|captionposition="within"
|captionalign="center"
|captionsize="small"
|bordersize="none"
|bordercolor="transparent"
|hideaddbutton="true"
|spacing="small"
}}

{| border="0" cellpadding="1" cellspacing="1" class="article-table wikitable" style="width:680px;" 
|-
! scope="col" style="width:120px;" |Area
! scope="col"|Dialogue
`,
		m.Name,
		0,
	)

	if m.StartMsg != "" {
		fmt.Fprintf(w, "|-\n| align=\"center\" |%s\n|%s\n", m.Name, html.EscapeString(strings.ReplaceAll(m.StartMsg, "\n", " ")))
	}

	for _, e := range m.Areas(vc.Data) {
		if e.Story != "" || e.Start != "" || e.End != "" || e.BossStart != "" || e.BossEnd != "" {
			fmt.Fprintf(w, "|-\n| align=\"center\" |%s\n|\n", e.LongName)

			if e.Story != "" {
				fmt.Fprintf(w, "; Prologue\n: %s\n", html.EscapeString(strings.ReplaceAll(e.Story, "\n", " ")))
				if e.Start != "" || e.End != "" || e.BossStart != "" || e.BossEnd != "" {
					io.WriteString(w, "----\n\n")
				}
			}

			if e.Start != "" || e.End != "" {
				io.WriteString(w, "; Guide Dialogue")
				if e.Start != "" {
					fmt.Fprintf(w, "\n: ''%s''",
						html.EscapeString(strings.ReplaceAll(e.Start, "\n", " ")))
					if e.End != "" {
						io.WriteString(w, "<br />&amp;nbsp;<br />")
					}
				}
				if e.End != "" {
					fmt.Fprintf(w, "\n: ''%s''\n",
						html.EscapeString(strings.ReplaceAll(e.End, "\n", " ")))
				} else {
					io.WriteString(w, "\n")
				}
				if e.BossStart != "" || e.BossEnd != "" {
					io.WriteString(w, "----\n\n")
				}
			}

			if e.BossStart != "" || e.BossEnd != "" {
				fmt.Fprintf(w, "; Boss Dialogue")
				if e.BossStart != "" {
					fmt.Fprintf(w, "\n: %s",
						html.EscapeString(strings.ReplaceAll(e.BossStart, "\n", " ")))
					if e.BossEnd != "" {
						io.WriteString(w, "<br />&amp;nbsp;<br />")
					}
				}
				if e.BossEnd != "" {
					fmt.Fprintf(w, "\n: %s\n",
						html.EscapeString(strings.ReplaceAll(e.BossEnd, "\n", " ")))
				} else {
					io.WriteString(w, "\n")
				}
			}
		}
	}
	io.WriteString(w, "|}\n[[Category:Story]]\n</textarea></div></body></html>")
}

// MapTableHandler show maps as a table
func MapTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Maps</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>ID</th><th>Name</th><th>Name Jp</th><th>Start</th><th>End</th><th>Archwitch Series</th><th>Archwitch</th><th>Elemental Hall</th><th>Flags</th><th>Beginner</th><th>Navi</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")

	for i := len(vc.Data.Maps) - 1; i >= 0; i-- {
		m := vc.Data.Maps[i]
		fmt.Fprintf(w, "<tr><td><a href=\"/maps/%[1]d\">%[1]d</a></td><td><a href=\"/maps/%[1]d\">%[2]s</a></td><td>%[3]s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td>",
			m.ID,
			m.Name,
			m.NameJp,
			m.PublicStartDatetime.Format(time.RFC3339),
			m.PublicEndDatetime.Format(time.RFC3339),
			m.KingSeriesID,
			m.KingID,
			m.ElementalhallID,
			m.Flags,
			m.ForBeginner,
			m.NaviID,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
