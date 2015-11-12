package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zetsuboushita.net/vc_file_grouper/vc"
)

const (
	wikiFmt = "15:04 January 2 2006"
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>_id</th><th>Event Name</th><th>Event Type</th><th>Start Date</th><th>End Date</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, e := range VcData.Events {
		fmt.Fprintf(w, "<tr><td><a href=\"/events/detail/%[1]d\">%[1]d</a></td><td><a href=\"/events/detail/%[1]d\">%[2]s</a></td><td>%d</td><td>%s</td><td>%s</td></tr>",
			e.Id,
			e.Name,
			e.EventTypeId,
			e.StartDatetime.Format(time.RFC3339),
			e.EndDatetime.Format(time.RFC3339),
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}

func eventDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "events/detail/id"
	if len(pathParts) < 3 {
		http.Error(w, "Invalid event id ", http.StatusNotFound)
		return
	}
	eventId, err := strconv.Atoi(pathParts[2])
	if err != nil || eventId < 1 || eventId > len(VcData.Events) {
		http.Error(w, "Invalid event id "+pathParts[2], http.StatusNotFound)
		return
	}

	event := vc.EventScan(eventId, VcData.Events)

	evntTemplate := `{{Event
|start jst = %s
|end jst = %s
|image = Banner {{PAGENAME}}.png
|story = yes
|%s|Ranking Reward
|%s|Legendary Archwitch
|%s|Fantasy Archwitch
|%s|Archwitch
}}

%s

==Rewards==
{{Rewards8
|ranking ur = %[3]s
|ranking sr = %[5]s
|ranking r = %[8]s
|progress point = 
|progress point sr = 
}}

==Ranking Trend==
{| class="article-table" style="text-align:right" border="1"
|-
!Date (JST)
!Rank 1
!Rank 50
!Rank 100
!Rank 300
!Rank 1000
!Rank 3000%s
|}
%s

{{NavEvent|%s|%s}}`

	rtrend := ""
	for i := event.StartDatetime.Add(24 * time.Hour); event.EndDatetime.After(i.Add(-24 * time.Hour)); i = i.Add(24 * time.Hour) {
		rtrend += fmt.Sprintf("\n|-\n|%s\n|\n|\n|\n|\n|\n|", i.Format("January _2"))
	}

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", event.Name)
	fmt.Fprintf(w, "<div style=\"float:left\">Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br />", event.Name)
	fmt.Fprintf(w, "<a href=\"/maps/%d\">Map Information</a>\n<br />", event.MapId)
	io.WriteString(w, "<textarea style=\"width:800px;height:760px\">")
	fmt.Fprintf(w, evntTemplate,
		event.StartDatetime.Format(wikiFmt), // start
		event.EndDatetime.Format(wikiFmt),   // end
		"", // rank reward
		"", // legendary archwitch
		"", // Fantasy Archwitch
		"", // Regular Archwitch
		html.EscapeString(event.Description),
		"",     // R reward
		rtrend, // Rank trend
		"",     // sub event (AUB)
		"",     //Previous event name
		"",     // next event name
	)
	io.WriteString(w, "</textarea></div>")

	io.WriteString(w, "</body></html>")

}
