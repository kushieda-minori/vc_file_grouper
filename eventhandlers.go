package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zetsuboushita.net/vc_file_grouper/vc"
)

const (
	wikiFmt = "15:04 January 2 2006"
)

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	filter := func(event *vc.Event) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		if eventType := qs.Get("eventType"); isInt(eventType) {
			eventTypeId, _ := strconv.Atoi(eventType)
			match = match && event.EventTypeId == eventTypeId
		}
		return
	}
	io.WriteString(w, "<html><head><title>All Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>_id</th><th>Event Name</th><th>Event Type</th><th>Start Date</th><th>End Date</th><th>King Series</th><th>Guild Battle</th><th>Tower Event</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for i := len(VcData.Events) - 1; i >= 0; i-- {
		e := VcData.Events[i]
		if !filter(&e) {
			continue
		}
		fmt.Fprintf(w, `<tr>
	<td><a href="/events/detail/%[1]d">%[1]d</a></td>
	<td><a href="/events/detail/%[1]d">%[2]s</a></td>
	<td>%d</td>
	<td>%s</td>
	<td>%s</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
</tr>`,
			e.Id,
			e.Name,
			e.EventTypeId,
			e.StartDatetime.Format(time.RFC3339),
			e.EndDatetime.Format(time.RFC3339),
			e.KingSeriesId,
			e.GuildBattleId,
			e.TowerEventId,
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

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", event.Name)
	if event.BannerId > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Banner"/></a><br />`, event.BannerId, url.QueryEscape(event.Name))
	}
	if event.TexIdImage > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Texture Image" /></a>br />`, event.TexIdImage, url.QueryEscape(event.Name))
	}
	if event.TexIdImage2 > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Texture Image 2" /></a><br />`, event.TexIdImage2, url.QueryEscape(event.Name))
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
	switch event.EventTypeId {
	case 1: // archwitch event
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
		rankReward := ""
		rr := event.RankRewards(VcData)
		if rr != nil {
			mid := rr.MidRewards(VcData)
			if mid != nil {
				midCaption := fmt.Sprintf("Mid Rankings<br /><small>Cutoff@ %s (JST)</small>",
					rr.MidBonusDistributionDate.Format(wikiFmt),
				)
				midrewards = genWikiAWRewards(mid, midCaption)
			}
			finalRewardList := rr.FinalRewards(VcData)
			finalrewards = genWikiAWRewards(finalRewardList, "Final Rankings")
			for _, fr := range finalRewardList {
				if fr.CardId > 0 {
					rrCard := vc.CardScan(fr.CardId, VcData.Cards)
					rankReward = rrCard.Name
					break
				}
			}
		}

		fmt.Fprintf(w, getEventTemplate(event.EventTypeId), event.EventTypeId,
			event.StartDatetime.Format(wikiFmt), // start
			event.EndDatetime.Format(wikiFmt),   // end
			eHallStart,                          // E-Hall opening
			"1",                                 // E-Hall rotation
			rankReward,                          // rank reward
			legendary,                           // legendary archwitch
			faws,                                // Fantasy Archwitch
			aws,                                 // Regular Archwitch
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			(midrewards + finalrewards), //rewards
			genWikiRankTrend(event),     // Rank trend
			"",            // sub event (Alliance Battle)
			prevEventName, //Previous event name
			nextEventName, // next event name
		)
	case 16: // alliance bingo battle
		// rewards
		//rr := event.RankRewards(VcData)
		//finalRewardList := rr.FinalRewards(VcData)
		//finalrewards := genWikiAWRewards(finalRewardList, "Ranking")

		fmt.Fprintf(w, getEventTemplate(event.EventTypeId), event.EventTypeId,
			event.StartDatetime.Format(wikiFmt),
			event.EndDatetime.Format(wikiFmt),
			"",    // RR 1
			"",    // RR 2
			"",    // Individual Card 1
			"",    // Individual Card 2
			"",    // Booster 1
			"",    // Booster 2
			"#th", // Guild Battle Number spelled out (first, second, third, etc)
			"",    // Overlap AW Event
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			"", // Ring Exchange
			"", // Rewards (combined)
			prevEventName,
			nextEventName,
		)
	case 18: // Tower Event
		tower := event.Tower(VcData)
		fmt.Fprintf(w,
			getEventTemplate(event.EventTypeId), event.EventTypeId,
			event.StartDatetime.Format(wikiFmt),
			event.EndDatetime.Format(wikiFmt),
			tower.ElementId,
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			genWikiAWRewards(tower.ArrivalRewards(VcData), "Floor Arrival Rewards"), // RR 1
			genWikiAWRewards(tower.RankRewards(VcData), "Rank Rewards"),             // RR 2
			prevEventName,
			nextEventName,
		)
	case 11: // special campaign (Abyssal AW and others)
		// may just do the THOR event seprately and leave this as just news
		fallthrough
	case 13: //Alliance Ultimate Battle
		fallthrough
	case 12: // Alliance Duel
		fallthrough
	case 10: //Alliance Battle
		fallthrough
	default:
		fmt.Fprintf(w,
			getEventTemplate(event.EventTypeId), event.EventTypeId,
			event.StartDatetime.Format(wikiFmt), // start
			event.EndDatetime.Format(wikiFmt),   // end
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			prevEventName,
			nextEventName,
		)
	}
	io.WriteString(w, "</textarea></div>")

	io.WriteString(w, "</body></html>")

}

func genWikiAWRewards(rewards []vc.RankRewardSheet, caption string) string {
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
	prange := `|-
|style="text-align:right"|%d F
|%s
`

	prevRangeStart, prevRangeEnd, prevPoint := 0, 0, 0
	rewardList := ""
	count := len(rewards) - 1
	for k, reward := range rewards {
		if k >= count {
			// get the last reward
			if reward.Point <= 0 && (reward.RankFrom != prevRangeStart || reward.RankTo != prevRangeEnd) {
				// if the last reward range is a single item, write out the previous range
				ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
				rewardList = ""
			} else if reward.Point != prevPoint {
				ret += fmt.Sprintf(prange, prevPoint, rewardList)
				rewardList = ""
			}
			rewardList += getWikiAWRewards(reward, rewardList != "")
			if reward.Point > 0 {
				ret += fmt.Sprintf(prange, reward.Point, rewardList)
			} else {
				ret += fmt.Sprintf(rrange, reward.RankFrom, reward.RankTo, rewardList)
			}
		} else {
			if reward.RankFrom != prevRangeStart || reward.RankTo != prevRangeEnd || reward.Point != prevPoint {
				if prevRangeStart > 0 {
					ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
				} else if reward.Point > 0 && reward.Point != prevPoint {
					ret += fmt.Sprintf(prange, prevPoint, rewardList)
				}
				prevRangeStart, prevRangeEnd, prevPoint = reward.RankFrom, reward.RankTo, reward.Point
				rewardList = ""
			}
			newline := rewardList != ""
			rewardList += getWikiAWRewards(reward, newline)
		}
	}

	return ret + "|}\n"
}

func getWikiAWRewards(reward vc.RankRewardSheet, newline bool) string {
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
			(item.GroupId >= 9 &&
				item.GroupId <= 16) {
			// Arcana
			r = fmt.Sprintf("{{Arcana|%s}}", cleanArcanaName(item.NameEng))
		} else if (item.GroupId >= 5 && item.GroupId <= 7) || item.GroupId == 31 || item.GroupId == 22 {
			// sword, shoe, key, rod, potion
			r = fmt.Sprintf("{{Valkyrie|%s}}", cleanItemName(item.NameEng))
		} else if item.GroupId == 18 || item.GroupId == 19 || item.GroupId == 43 || item.GroupId == 47 {
			switch item.Id {
			case 29:
				r = "{{MaidenTicket}}"
			case 138:
				r = "{{AWCore}}"
			default:
				// exchange items
				r = fmt.Sprintf("[[File:%[1]s.png|28px|link=Items#%[1]s]] [[Items#%[1]s|%[1]s]]", item.NameEng)
			}
		} else if item.GroupId == 29 {
			itemName := ""
			if strings.Contains(item.Name, "LIGHT") {
				itemName = "Light"
			} else if strings.Contains(item.Name, "PASSION") {
				itemName = "Passion"
			} else if strings.Contains(item.Name, "COOL") {
				itemName = "Cool"
			} else if strings.Contains(item.Name, "DARK") {
				itemName = "Dark"
			}
			if strings.Contains(item.NameEng, "Crystal") {
				itemName += "C"
			} else if strings.Contains(item.NameEng, "Orb") {
				itemName += "O"
			} else if strings.Contains(item.NameEng, "(L)") {
				itemName += "L"
			} else if strings.Contains(item.NameEng, "(M)") {
				itemName += "M"
			} else if strings.Contains(item.NameEng, "(S)") {
				itemName += "S"
			}
			r = fmt.Sprintf("{{Stone|%s}}", itemName)
		} else if item.GroupId == 38 {
			// Custom Skill Recipies
			r = fmt.Sprintf("{{Skill Recipe|%s}}", cleanCustomSkillRecipe(item.NameEng))
		} else if item.GroupId == 39 {
			// Custom Skill items
			r = fmt.Sprintf("[[File:%[1]s.png|28px|link=Custom Skills#Skill_Materials]] [[Custom Skills#Skill_Materials|%[2]s]]",
				vc.CleanCustomSkillNoImage(item.NameEng),
				vc.CleanCustomSkillNoImage(item.NameEng),
			)
		} else {
			r = fmt.Sprintf("__UNKNOWN_GROUP:_%d_%s__", item.GroupId, item.NameEng)
		}
	} else if reward.Cash > 0 {
		r = "{{icon|jewel}}"
	} else {
		r = "Unknown Reward Type"
	}

	return fmt.Sprintf(rlist, r, reward.Num)
}

func getEventTemplate(eventType int) string {
	switch eventType {
	case 1: // AW Events
		return `{{Event|eventType = %d
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
%s


%s

{{NavEvent|%s|%s}}
`
	case 16: // ABB Events
		return `{{event|eventType = %d
|image=Banner_{{PAGENAME}}.png
|start jst=%s
|end jst=%s
| %s |Rank Reward
| %s |Rank Reward
| %s |Individual Point Reward<br />Ring Exchange
| %s |Individual Point Reward<br />Ring Exchange
| Mirror Maiden (UR) |Ring Exchange
| Mirror Maiden (SR) |Ring Exchange
| Mirror Maiden (R) |Ring Exchange
| Slime Queen |Ring Exchange
| Mirror Maiden Shard | Ring Exchange
| %s |Alliance Battle Point Booster<br>+60%%/150%%
| %s |Alliance Battle Point Booster<br>+20%%/50%%
}}
:''The %s [[Alliance Bingo Battle]] was held during the [[%s]] event.''

%s

==Ring Exchange==
%s

==Rewards==
%s

==Local ABB Times==
{{AUBLocalTime|start jst = %[2]s|consecutive=1}}

{{NavEvent|%[15]s|%s}}
`
	case 18: // Tower Events
		return `{{Event|eventType = %d
|start jst = %s
|end jst = %s
|elementHallRotate=%s
|image = Banner {{PAGENAME}}.png
}}

%s

==Rewards==
%s%s
{{clr}}

{{NavEvent|%s|%s}}
`
	default: // Default event handler
		return `{{Event|eventType = %d
|start jst = %s
|end jst = %s
|image = Banner {{PAGENAME}}.png
}}

%s

{{clr}}

{{NavEvent|%s|%s}}
`
	}
}

func genWikiRankTrend(event *vc.Event) (rtrend string) {
	rtrend = `{| class="article-table" style="text-align:right" border="1"
|-
!Date (JST) !! Rank 1 !! Rank 50 !! Rank 100 !! Rank 300 !! Rank 500 !! Rank 1000 !! Rank 2000`
	for i := event.StartDatetime.Add(24 * time.Hour); event.EndDatetime.After(i.Add(-24 * time.Hour)); i = i.Add(24 * time.Hour) {
		rtrend += fmt.Sprintf("\n|-\n|%s\n|\n|\n|\n|\n|\n|\n|", i.Format("January _2"))
	}
	rtrend += "\n|}"
	return
}

func cleanCustomSkillRecipe(name string) string {
	ret := ""
	lower := strings.ToLower(vc.CleanCustomSkillNoImage(name))
	if strings.Contains(lower, "all enemies") {
		ret += "aoe"
	}
	if strings.Contains(lower, "stop") {
		ret += "+ts "
	} else if ret != "" {
		ret += " "
	}
	if strings.Contains(lower, "fixed") {
		ret += "fixed "
	}
	if strings.Contains(lower, "proportional") {
		ret += "proportional "
	}
	if strings.Contains(lower, "awoken burst") {
		ret += "awoken burst "
	} else if strings.Contains(lower, "recover") {
		ret += "recover "
	} else if strings.Contains(lower, "own atk up") {
		ret += "own atk up "
	}
	if strings.Contains(lower, "passion") {
		ret += "passion "
	}
	if strings.Contains(lower, "cool") {
		ret += "cool "
	}
	if strings.Contains(lower, "light") {
		ret += "light "
	}
	if strings.Contains(lower, "dark") {
		ret += "dark "
	}
	if strings.Contains(lower, "1") {
		ret += "1"
	}
	if strings.Contains(lower, "2") {
		ret += "2"
	}
	if strings.Contains(lower, "3") {
		ret += "3"
	}
	return ret
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
	ret = strings.Replace(ret, "arcana's", "", -1)
	ret = strings.Replace(ret, "arcana", "", -1)
	ret = strings.Replace(ret, "%", "", -1)
	ret = strings.Replace(ret, "+", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.Replace(ret, "forced", "", 1)
	ret = strings.Replace(ret, "strongdef", "def", 1)
	ret = strings.TrimSpace(ret)
	return ret
}
