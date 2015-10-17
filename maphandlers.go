package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
		http.Error(w, "Invalid map id ", http.StatusNotFound)
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
