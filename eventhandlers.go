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
	if err != nil || eventId < 1 {
		http.Error(w, "Invalid event id "+pathParts[2], http.StatusNotFound)
		return
	}

	event := vc.EventScan(eventId, VcData.Events)

	var prevEvent, nextEvent *vc.Event = nil, nil

	for i := event.Id - 1; i > 0; i-- {
		tmp := vc.EventScan(i, VcData.Events)
		if tmp != nil && tmp.EventTypeId == event.EventTypeId {
			prevEvent = tmp
			break
		}
	}

	prevEventName := ""
	if prevEvent != nil {
		prevEventName = strings.Replace(prevEvent.Name, "【New Event】", "", -1)
	}

	for i := event.Id + 1; i <= vc.MaxEventId(VcData.Events); i++ {
		tmp := vc.EventScan(i, VcData.Events)
		if tmp != nil && tmp.EventTypeId == event.EventTypeId {
			nextEvent = tmp
			break
		}
	}

	nextEventName := ""
	if nextEvent != nil {
		nextEventName = strings.Replace(nextEvent.Name, "【New Event】", "", -1)
	}

	evntTemplate1 := `{{Event
|start jst = %s
|end jst = %s
|elementalHallOpen=%s
|elementHallRotate=%s
|image = Banner {{PAGENAME}}.png
|story = yes
|%s|Ranking Reward
|%s|Legendary Archwitch
%s%s||Amalgamation Material
||Amalgamation
||Elemental Hall
||Event 10/15x damage<br/>60/120%% Points+
||Event 10/15x damage<br/>60/120%% Points+
||Event 10/15x damage<br/>30/50%% Points+
||Event 10/15x damage<br/>30/50%% Points+
}}

%s

==Rewards==


==Ranking Trend==
{| class="article-table" style="text-align:right" border="1"
|-
!Date (JST)
!Rank 1
!Rank 50
!Rank 100
!Rank 300
!Rank 500
!Rank 1000
!Rank 2000%s
|}

%s

{{NavEvent|%s|%s}}`

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", event.Name)
	if event.BannerId > 0 {
		fmt.Fprintf(w, `<img src="/images/event/largeimage/%d/event_image_en" alt="Banner"/><br />`, event.BannerId)
	}
	if event.TexIdImage > 0 {
		fmt.Fprintf(w, `<img src="/images/event/largeimage/%d/event_image_en" alt="Texture Image" /><br />`, event.TexIdImage)
	}
	if event.TexIdImage2 > 0 {
		fmt.Fprintf(w, `<img src="/images/event/largeimage/%d/event_image_en" alt="Texture Image 2" /><br />`, event.TexIdImage2)
	}
	if prevEventName != "" {
		fmt.Fprintf(w, "<div style=\"float:left\"><a href=\"%d\">%s</a>\n</div>", prevEvent.Id, prevEventName)
	}
	if nextEventName != "" {
		fmt.Fprintf(w, "<div style=\"float:right\"><a href=\"%d\">%s</a>\n</div>", nextEvent.Id, nextEventName)
	}
	fmt.Fprintf(w, "<div style=\"clear:both;float:left\">Edit on the <a href=\"https://valkyriecrusade.wikia.com/wiki/%s?action=edit\">wikia</a>\n<br />", strings.Replace(event.Name, "【New Event】", "", -1))
	if event.MapId > 0 {
		fmt.Fprintf(w, "<a href=\"/maps/%d\">Map Information</a>\n<br />", event.MapId)
	}
	io.WriteString(w, "<textarea style=\"width:800px;height:760px\">")
	if event.EventTypeId == 1 {
		rtrend := ""
		for i := event.StartDatetime.Add(24 * time.Hour); event.EndDatetime.After(i.Add(-24 * time.Hour)); i = i.Add(24 * time.Hour) {
			rtrend += fmt.Sprintf("\n|-\n|%s\n|\n|\n|\n|\n|\n|\n|", i.Format("January _2"))
		}

		var legendary string
		var faws string
		var aws string

		for _, aw := range event.Archwitches(VcData) {
			cardMaster := vc.CardScan(aw.CardMasterId, VcData.Cards)
			if aw.IsLAW() {
				legendary = cardMaster.Name
			} else if aw.IsFAW() {
				faws += "|" + cardMaster.Name + "|Fantasy Archwitch\n"
			} else {
				aws += "|" + cardMaster.Name + "|Archwitch\n"
			}
		}

		fmt.Fprintf(w, evntTemplate1,
			event.StartDatetime.Format(wikiFmt), // start
			event.EndDatetime.Format(wikiFmt),   // end
			"",        // E-Hall opening
			"",        // E-Hall rotation
			"",        // rank reward
			legendary, // legendary archwitch
			faws,      // Fantasy Archwitch
			aws,       // Regular Archwitch
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			rtrend,        // Rank trend
			"",            // sub event (Alliance Battle)
			prevEventName, //Previous event name
			nextEventName, // next event name
		)
	} else {
		io.WriteString(w, html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)))
	}
	io.WriteString(w, "</textarea></div>")

	io.WriteString(w, "</body></html>")

}
