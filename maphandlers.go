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

func mapHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "maps/id"
	if len(pathParts) < 2 {
		mapTableHandler(w, r)
		return
	}
	mapId, err := strconv.Atoi(pathParts[1])
	if err != nil || mapId < 1 || mapId > len(VcData.Maps) {
		http.Error(w, "Invalid map id "+pathParts[1], http.StatusNotFound)
		return
	}

	m := vc.MapScan(mapId, VcData.Maps)

	fmt.Fprintf(w, "<html><head><title>Map %s</title>\n", m.Name)
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	fmt.Fprintf(w, "<h1>%s</h1>\n%s", m.Name, m.StartMsg)
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>No</th><th>Name</th><th>Long Name</th><th>Start</th><th>End</th><th>Story</th><th>Boss Start</th><th>Boss End</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, e := range m.Areas(&VcData) {
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

func mapTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Maps</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>Id</th><th>Name</th><th>Name Jp</th><th>Start</th><th>End</th><th>Archwitch Series</th><th>Archwitch</th><th>Elemental Hall</th><th>Flags</th><th>Beginner</th><th>Navi</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")

	for _, m := range VcData.Maps {
		fmt.Fprintf(w, "<tr><td><a href=\"/maps/%[1]d\">%[1]d</a></td><td><a href=\"/maps/%[1]d\">%[2]s</a></td><td>%[3]s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td>",
			m.Id,
			m.Name,
			m.NameJp,
			m.PublicStartDatetime.Format(time.RFC3339),
			m.PublicEndDatetime.Format(time.RFC3339),
			m.KingSeriesId,
			m.KingId,
			m.ElementalhallId,
			m.Flags,
			m.ForBeginner,
			m.NaviId,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
