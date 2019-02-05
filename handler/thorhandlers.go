package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// ThorHandler handle Thor events
func ThorHandler(w http.ResponseWriter, r *http.Request) {
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
		ThorTableHandler(w, r)
		return
	}

	thorID, err := strconv.Atoi(pathParts[1])
	if err != nil || thorID < 1 || thorID > len(vc.Data.ThorEvents) {
		http.Error(w, "Invalid Thor Event id "+pathParts[1], http.StatusNotFound)
		return
	}
	t := vc.ThorEventScan(thorID, vc.Data.ThorEvents)

	if t == nil {
		http.Error(w, "Invalid Thor Event id "+pathParts[1], http.StatusNotFound)
		return
	}

	if len(pathParts) >= 3 && "WIKI" == pathParts[2] {
		ThorDetailWikiHandler(w, r, t)
		return
	}

	ThorDetailHandler(w, r, t)
}

// ThorDetailWikiHandler does nothing
func ThorDetailWikiHandler(w http.ResponseWriter, r *http.Request, t *vc.ThorEvent) {
}

// ThorDetailHandler does nothing
func ThorDetailHandler(w http.ResponseWriter, r *http.Request, t *vc.ThorEvent) {
}

// ThorTableHandler shows thor events as a table
func ThorTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Thor Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>ID</th>"+
		"<th>Start</th>"+
		"<th>End</th>"+
		"<th>Rank Start</th>"+
		"<th>Rank End</th>"+
		"<th>Reward Distribution</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")

	for _, t := range vc.Data.ThorEvents {
		fmt.Fprintf(w, "<tr><td><a href=\"/thor/%[1]d\">%[1]d</a></td>"+
			"<td>%[2]s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>",
			t.ID,
			t.PublicStartDatetime.Format(time.RFC3339),
			t.PublicEndDatetime.Format(time.RFC3339),
			t.RankingStartDatetime.Format(time.RFC3339),
			t.RankingEndDatetime.Format(time.RFC3339),
			t.RankingRewardDestributionStartDatetime.Format(time.RFC3339),
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
