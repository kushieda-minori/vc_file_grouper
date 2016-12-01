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
	io.WriteString(w, "<th>_id</th><th>Event Name</th><th>Event Type</th><th>Start Date</th><th>End Date</th><th>King Series</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for _, e := range VcData.Events {
		fmt.Fprintf(w, `<tr>
	<td><a href="/events/detail/%[1]d">%[1]d</a></td>
	<td><a href="/events/detail/%[1]d">%[2]s</a></td>
	<td>%d</td>
	<td>%s</td>
	<td>%s</td>
	<td>%d</td>
</tr>`,
			e.Id,
			e.Name,
			e.EventTypeId,
			e.StartDatetime.Format(time.RFC3339),
			e.EndDatetime.Format(time.RFC3339),
			e.KingSeriesId,
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
||Event 10/15x damage<br />60/120%% Points+
||Event 10/15x damage<br />60/120%% Points+
||Event 10/15x damage<br />30/50%% Points+
||Event 10/15x damage<br />30/50%% Points+
}}

%s

==Rewards==
%s
{{clr}}

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

		eventMap := event.Map(VcData)
		var eHallStart string
		if eventMap == nil || eventMap.ElementalhallStart.IsZero() || event.EndDatetime.Before(eventMap.ElementalhallStart.Time) {
			eHallStart = ""
		} else {
			eHallStart = eventMap.ElementalhallStart.Format(wikiFmt)
		}

		midrewards := ""
		finalrewards := ""
		rr := event.RankRewards(VcData)
		if rr != nil {
			mid := rr.MidRewards(VcData)
			if mid != nil {
				midCaption := fmt.Sprintf("Mid Rankings<br /><small> Cutoff@ %s (JST)</small>",
					rr.MidBonusDistributionDate.Format(wikiFmt),
				)
				midrewards = genWikiRewards(mid, midCaption)
			}
			finalrewards = genWikiRewards(rr.FinalRewards(VcData), "Final Rankings")
		}

		fmt.Fprintf(w, evntTemplate1,
			event.StartDatetime.Format(wikiFmt), // start
			event.EndDatetime.Format(wikiFmt),   // end
			eHallStart,                          // E-Hall opening
			"1",                                 // E-Hall rotation
			"",                                  // rank reward
			legendary,                           // legendary archwitch
			faws,                                // Fantasy Archwitch
			aws,                                 // Regular Archwitch
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			(midrewards + finalrewards), //rewards
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

func genWikiRewards(rewards []vc.RankRewardSheet, caption string) string {
	ret := `
{| border="1" cellpadding="1" cellspacing="1" class="mw-collapsible mw-collapsed article-table" style="float:left"
|-
! colspan="2"|` + caption + `
|-
! style="text-align:right"|Rank
! Reward
`
	rrange := `|-
|style="text-align:right"|%d~%d
|%s
`

	prevRangeStart, prevRangeEnd := 0, 0
	rewardList := ""
	count := len(rewards) - 1
	for k, reward := range rewards {
		if k >= count {
			// get the last reward
			if reward.RankFrom != prevRangeStart || reward.RankTo != prevRangeEnd {
				// if the last reward range is a single item, write out the previous range
				ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
				rewardList = ""
			}
			rewardList += getWikiReward(reward, rewardList != "")
			ret += fmt.Sprintf(rrange, reward.RankFrom, reward.RankTo, rewardList)
		} else {
			if reward.RankFrom != prevRangeStart || reward.RankTo != prevRangeEnd {
				if prevRangeStart > 0 {
					ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
				}
				prevRangeStart, prevRangeEnd = reward.RankFrom, reward.RankTo
				rewardList = ""
			}
			newline := rewardList != ""
			rewardList += getWikiReward(reward, newline)
		}
	}

	return ret + "|}\n"
}

func getWikiReward(reward vc.RankRewardSheet, newline bool) string {
	rlist := "%s x%d"
	if newline {
		rlist = "<br />" + rlist
	}

	var r string
	if reward.CardId > 0 {
		card := vc.CardScan(reward.CardId, VcData.Cards)
		if card == nil {
			r = "{{Card Icon|Unknown Card Id}}"
		} else {
			r = fmt.Sprintf("{{Card Icon|%s}}", card.Name)
		}
	} else if reward.ItemId > 0 {
		item := vc.ItemScan(reward.ItemId, VcData.Items)
		if item == nil {
			r = fmt.Sprintf("__UNKNOWN_ITEM_ID:%d__", reward.ItemId)
		} else if item.GroupId == 17 {
			// tickets
			r = fmt.Sprintf("{{Ticket|%s}}", cleanTicketName(item.NameEng))
		} else if item.GroupId == 30 ||
			(item.GroupId >= 10 &&
				item.GroupId <= 16) {
			// Arcana
			r = fmt.Sprintf("{{Arcana|%s}}", cleanArcanaName(item.NameEng))
		} else if (item.GroupId >= 5 && item.GroupId <= 7) || item.GroupId == 31 || item.GroupId == 22 {
			// sword, shoe, key, rod, potion
			r = fmt.Sprintf("{{Valkyrie|%s}}", cleanItemName(item.NameEng))
		} else if item.GroupId == 18 {
			r = fmt.Sprintf("[[File:%[1]s.png|28px|link=Items#%[1]s]] [[Items#%[1]s|%[1]s]]", item.NameEng)
		} else {
			r = fmt.Sprintf("__UNKNOWN_GROUP:_%d_%s__", item.GroupId, item.NameEng)
		}
	} else {
		r = "Unknown Reward Type"
	}

	return fmt.Sprintf(rlist, r, reward.Num)
}

func cleanTicketName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "ticket", "", -1)
	ret = strings.Replace(ret, "summon", "", -1)
	ret = strings.Replace(ret, "guaranteed", "", -1)
	ret = strings.Replace(ret, "★★★", "3star", -1)
	ret = strings.Replace(ret, "★★", "2star", -1)
	ret = strings.Replace(ret, "★", "1star", -1)
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanItemName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "valkyrie", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanArcanaName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "arcana", "", -1)
	ret = strings.Replace(ret, "%", "", -1)
	ret = strings.Replace(ret, "+", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.TrimSpace(ret)
	return ret
}
