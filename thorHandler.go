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

func thorHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "thor/id/WIKI"
	if len(pathParts) < 2 {
		thorTableHandler(w, r)
		return
	}

	thorId, err := strconv.Atoi(pathParts[1])
	if err != nil || thorId < 1 || thorId > len(VcData.ThorEvents) {
		http.Error(w, "Invalid Thor Event id "+pathParts[1], http.StatusNotFound)
		return
	}
	t := vc.ThorEventScan(thorId, VcData.ThorEvents)

	if t == nil {
		http.Error(w, "Invalid Thor Event id "+pathParts[1], http.StatusNotFound)
		return
	}

	if len(pathParts) >= 3 && "WIKI" == pathParts[2] {
		thorDetailWikiHandler(w, r, t)
		return
	}

	thorDetailHandler(w, r, t)
}

func thorDetailWikiHandler(w http.ResponseWriter, r *http.Request, t *vc.ThorEvent) {
}

func thorDetailHandler(w http.ResponseWriter, r *http.Request, t *vc.ThorEvent) {
}

func thorTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Thor Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>Id</th>"+
		"<th>Start</th>"+
		"<th>End</th>"+
		"<th>Rank Start</th>"+
		"<th>Rank End</th>"+
		"<th>Reward Distribution</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")

	for _, t := range VcData.ThorEvents {
		fmt.Fprintf(w, "<tr><td><a href=\"/thor/%[1]d\">%[1]d</a></td>"+
			"<td>%[2]s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>",
			t.Id,
			t.PublicStartDatetime.Format(time.RFC3339),
			t.PublicEndDatetime.Format(time.RFC3339),
			t.RankingStartDatetime.Format(time.RFC3339),
			t.RankingEndDatetime.Format(time.RFC3339),
			t.RankingRewardDestributionStartDatetime.Format(time.RFC3339),
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
