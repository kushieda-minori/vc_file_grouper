package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"sort"
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
	qEventType := qs.Get("eventType")
	eventTypeID, _ := strconv.Atoi(qEventType)
	qSearch := strings.ToLower(strings.TrimSpace(qs.Get("search")))
	qIsActive := strings.ToLower(strings.TrimSpace(qs.Get("isActive")))

	filter := func(event *vc.Event) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		if eventTypeID > 0 {
			match = match && event.EventTypeID == eventTypeID
		}
		if qSearch != "" {
			match = match && (strings.Contains(strings.ToLower(event.Description), qSearch) || strings.Contains(strings.ToLower(event.Name), qSearch))
		}
		if qIsActive != "" {
			match = match && event.StartDatetime.Before(time.Now()) && event.EndDatetime.After(time.Now())
		}
		return
	}

	eventTypeKeys := make([]int, 0, len(vc.EventType))
	for i := range vc.EventType {
		eventTypeKeys = append(eventTypeKeys, i)
	}
	sort.Ints(eventTypeKeys)

	io.WriteString(w, "<html><head><title>All Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<form name=\"searchForm\">\n")
	fmt.Fprintf(w, "<label for=\"f_search\">Contains Text:</label><input id=\"f_search\" name=\"search\" value=\"%s\" />\n", qs.Get("search"))
	io.WriteString(w, "<label for=\"f_eventType\">Event Type:</label><select id=\"f_eventType\" name=\"eventType\">\n")
	io.WriteString(w, "<option value=\"\"></option>\n")
	for _, key := range eventTypeKeys {
		sKey := strconv.Itoa(key)
		selected := ""
		if sKey == qEventType {
			selected = " selected=\"selected\""
		}
		fmt.Fprintf(w, "<option value=\"%d\"%[3]s>%[1]d: %s</option>\n", key, vc.EventType[key], selected)
	}
	io.WriteString(w, "</select>\n")
	fmt.Fprintf(w, "<label for=\"f_isActive\">Is Active:</label><input id=\"f_isActive\" name=\"isActive\" type=\"checkbox\" value=\"checked\" %s/>", qIsActive)
	io.WriteString(w, "<button type=\"submit\">Submit</button>\n</form>\n")
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
			e.ID,
			e.Name,
			e.EventTypeID,
			e.StartDatetime.Format(time.RFC3339),
			e.EndDatetime.Format(time.RFC3339),
			e.KingSeriesID,
			e.GuildBattleID,
			e.TowerEventID,
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
	eventID, err := strconv.Atoi(pathParts[2])
	if err != nil || eventID < 1 {
		http.Error(w, "Invalid event id "+pathParts[2], http.StatusNotFound)
		return
	}

	event := vc.EventScan(eventID, VcData.Events)

	var prevEvent, nextEvent *vc.Event = nil, nil

	for i := event.ID - 1; i > 0; i-- {
		tmp := vc.EventScan(i, VcData.Events)
		if tmp != nil && tmp.EventTypeID == event.EventTypeID && !strings.Contains(tmp.Name, "Rune Boss") && !strings.Contains(tmp.Name, " 2x ") {
			prevEvent = tmp
			break
		}
	}

	prevEventName := ""
	if prevEvent != nil {
		prevEventName = strings.Replace(prevEvent.Name, "【New Event】", "", -1)
	}

	for i := event.ID + 1; i <= vc.MaxEventID(VcData.Events); i++ {
		tmp := vc.EventScan(i, VcData.Events)
		if tmp != nil && tmp.EventTypeID == event.EventTypeID && !strings.Contains(tmp.Name, "Rune Boss") && !strings.Contains(tmp.Name, " 2x ") {
			nextEvent = tmp
			break
		}
	}

	nextEventName := ""
	if nextEvent != nil {
		nextEventName = strings.Replace(nextEvent.Name, "【New Event】", "", -1)
	}

	fmt.Fprintf(w, "<html><head><title>%s</title></head><body><h1>%[1]s</h1>\n", event.Name)
	if event.BannerID > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Banner"/></a><br />`, event.BannerID, url.QueryEscape(event.Name))
	}
	if event.TexIDImage > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Texture Image" /></a>br />`, event.TexIDImage, url.QueryEscape(event.Name))
	}
	if event.TexIDImage2 > 0 {
		fmt.Fprintf(w, `<a href="/images/event/largeimage/%[1]d/event_image_en?filename=Banner_%[2]s.png"><img src="/images/event/largeimage/%[1]d/event_image_en" alt="Texture Image 2" /></a><br />`, event.TexIDImage2, url.QueryEscape(event.Name))
	}
	if prevEventName != "" {
		fmt.Fprintf(w, "<div style=\"float:left\"><a href=\"%d\">%s</a>\n</div>", prevEvent.ID, prevEventName)
	}
	if nextEventName != "" {
		fmt.Fprintf(w, "<div style=\"float:right\"><a href=\"%d\">%s</a>\n</div>", nextEvent.ID, nextEventName)
	}
	fmt.Fprintf(w, "<div style=\"clear:both;float:left\">Edit on the <a href=\"https://valkyriecrusade.fandom.com/wiki/%s?action=edit\">fandom</a>\n<br />", strings.Replace(event.Name, "【New Event】", "", -1))
	if event.MapID > 0 {
		fmt.Fprintf(w, "<a href=\"/maps/%d\">Map Information</a>\n<br />", event.MapID)
	}
	io.WriteString(w, "<textarea style=\"width:800px;height:760px\">")
	switch event.EventTypeID {
	case 1: // archwitch event
		var legendary string
		var faws string
		var aws string

		for _, aw := range event.Archwitches(VcData) {
			cardMaster := vc.CardScan(aw.CardMasterID, VcData.Cards)
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
		midRewardTime := time.Time{}
		finalrewards := ""
		rankReward := ""
		rr := event.RankRewards(VcData)
		if rr != nil {
			mid := rr.MidRewards(VcData)
			if mid != nil {
				midRewardTime = rr.MidBonusDistributionDate.Time
				midCaption := fmt.Sprintf("Mid Rankings<br /><small>Cutoff@ %s (JST)</small>",
					midRewardTime.Format(wikiFmt),
				)
				midrewards = genWikiAWRewards(mid, midCaption, "Rank")
			}
			finalRewardList := rr.FinalRewards(VcData)
			finalrewards = genWikiAWRewards(finalRewardList, "Final Rankings", "Rank")
			for _, fr := range finalRewardList {
				if fr.CardID > 0 {
					rrCard := vc.CardScan(fr.CardID, VcData.Cards)
					rankReward = rrCard.Name
					break
				}
			}
		}

		var ranks = []int{1, 100, 200, 300, 500, 1000, 2000}

		fmt.Fprintf(w, getEventTemplate(event.EventTypeID), event.EventTypeID,
			event.StartDatetime.Format(wikiFmt), // start
			event.EndDatetime.Format(wikiFmt),   // end
			eHallStart,                          // E-Hall opening
			"1",                                 // E-Hall rotation
			rankReward,                          // rank reward
			legendary,                           // legendary archwitch
			faws,                                // Fantasy Archwitch
			aws,                                 // Regular Archwitch
			html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
			(midrewards + finalrewards),                                    //rewards
			genWikiRankTrend(event, eventMap, midRewardTime, ranks, false), // Rank trend
			"",            // sub event (Alliance Battle)
			prevEventName, //Previous event name
			nextEventName, // next event name
		)
	case 16: // alliance bingo battle
		// rewards
		//rr := event.RankRewards(VcData)
		//finalRewardList := rr.FinalRewards(VcData)
		//finalrewards := genWikiAWRewards(finalRewardList, "Ranking")
		gb := event.GuildBattle(VcData)
		bb := gb.BingoBattle(VcData)

		// aws := bb.Archwitches(VcData)
		// aw := ""
		//os.Stderr.WriteString(fmt.Sprintf("found %d archwitches on guild battle %d king series id %d\n", len(aws), bb.ID, bb.KingSeriesID))
		// if len(aws) > 0 {
		// 	king := aws[0]
		// 	kingCard := vc.CardScan(king.CardMasterID, VcData.Cards)
		// 	aw = kingCard.Name
		// 	if len(aws) > 1 {
		// 		// append extra AW cards
		// 		for i := 1; i < len(aws); i++ {
		// 			king = aws[i]
		// 			kingCard = vc.CardScan(king.CardMasterID, VcData.Cards)
		// 			aw += " |Archwitch Panel Encounter\n| " + kingCard.Name
		// 		}
		// 	}
		// }

		rankRewards := genWikiAWRewards(gb.RankRewards(VcData), "Ranking", "Rank") +
			genWikiAWRewards(gb.IndividualRewards(VcData), "Point Reward", "Points")

		fmt.Fprintf(w, getEventTemplate(event.EventTypeID), event.EventTypeID,
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
			genWikiExchange(bb.ExchangeRewards(VcData)), // Ring Exchange
			rankRewards,                                 // Rewards (combined)
			prevEventName,
			nextEventName,
		)
	case 18: // Tower Event
		tower := event.Tower(VcData)
		if tower == nil {
			fmt.Fprintf(w, "Unable to find tower event")
		} else {
			element := tower.ElementID - 1
			towerShield := vc.Elements[element]
			var ranks = []int{1, 100, 300, 500, 1000, 2000, 3000, 5000}
			fmt.Fprintf(w,
				getEventTemplate(event.EventTypeID), event.EventTypeID,
				event.StartDatetime.Format(wikiFmt),
				event.EndDatetime.Format(wikiFmt),
				towerShield,
				html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
				genWikiAWRewards(tower.ArrivalRewards(VcData), "Floor Arrival Rewards", "Floor"), // RR 1
				genWikiAWRewards(tower.RankRewards(VcData), "Rank Rewards", "Rank"),              // RR 2
				genWikiRankTrend(event, nil, time.Unix(0, 0), ranks, true),                       // rank trend
				prevEventName,
				nextEventName,
			)
		}
	case 19: // Demon Realm Voyage
		realm := event.DemonRealm(VcData)
		if realm == nil {
			fmt.Fprintf(w, "Unable to find demon realm event")
		} else {
			element := realm.ElementID - 1
			shield := vc.Elements[element]
			var ranks = []int{1, 100, 300, 500, 1000, 2000, 3000, 5000}
			fmt.Fprintf(w,
				getEventTemplate(event.EventTypeID), event.EventTypeID,
				event.StartDatetime.Format(wikiFmt),
				event.EndDatetime.Format(wikiFmt),
				shield,
				html.EscapeString(strings.Replace(event.Description, "\n", "\n\n", -1)),
				genWikiAWRewards(realm.ArrivalRewards(VcData), "Point Rewards", "Floor"), // RR 1
				genWikiAWRewards(realm.RankRewards(VcData), "Rank Rewards", "Rank"),      // RR 2
				genWikiRankTrend(event, nil, time.Unix(0, 0), ranks, true),               // rank trend
				prevEventName,
				nextEventName,
			)
		}
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
			getEventTemplate(event.EventTypeID), event.EventTypeID,
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

func genWikiAWRewards(rewards []vc.RankRewardSheet, caption string, rankTitle string) string {
	ret := `
{| border="1" cellpadding="1" cellspacing="1" class="mw-collapsible mw-collapsed article-table" style="float:left"
|-
! colspan="2"|` + caption + `
|-
! style="text-align:right"|` + rankTitle + `
! Reward
`
	rrange := `|-
|style="text-align:right"|%d~%d
|%s
`
	prange := `|-
|style="text-align:right"|%d
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
				if rewardList != "" {
					ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
					rewardList = ""
				}
			} else if reward.Point != prevPoint {
				if rewardList != "" {
					ret += fmt.Sprintf(prange, prevPoint, rewardList)
					rewardList = ""
				}
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
					if rewardList != "" {

						ret += fmt.Sprintf(rrange, prevRangeStart, prevRangeEnd, rewardList)
					}
				} else if reward.Point > 0 && reward.Point != prevPoint {
					if rewardList != "" {
						ret += fmt.Sprintf(prange, prevPoint, rewardList)
					}
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
	if reward.CardID > 0 {
		card := vc.CardScan(reward.CardID, VcData.Cards)
		if card == nil {
			r = "{{Card Icon|Unknown Card ID}}"
		} else {
			r = fmt.Sprintf("{{Card Icon|%s}}", card.Name)
		}
	} else if reward.ItemID > 0 {
		item := vc.ItemScan(reward.ItemID, VcData.Items)
		if item == nil {
			r = fmt.Sprintf("__UNKNOWN_ITEM_ID:%d__", reward.ItemID)
		} else {
			r = getWikiItem(item)
		}
	} else if reward.Cash > 0 {
		r = "{{icon|jewel}}"
	} else {
		r = "Unknown Reward Type"
	}

	return fmt.Sprintf(rlist, r, reward.Num)
}

func getWikiItem(item *vc.Item) (r string) {
	if item.GroupID == 17 {
		// tickets
		r = fmt.Sprintf("{{Ticket|%s}}", cleanTicketName(item.NameEng))
	} else if item.GroupID == 30 ||
		(item.GroupID >= 9 &&
			item.GroupID <= 16) {
		// Arcana
		r = fmt.Sprintf("{{Arcana|%s}}", cleanArcanaName(item.NameEng))
	} else if (item.GroupID >= 5 && item.GroupID <= 7) || item.GroupID == 31 || item.GroupID == 22 || item.GroupID == 47 || item.GroupID == 48 {
		// sword, shoe, key, rod, potion, crystal, feather
		r = fmt.Sprintf("{{Valkyrie|%s}}", cleanItemName(item.NameEng))
	} else if item.GroupID == 18 || item.GroupID == 19 || item.GroupID == 43 {
		switch item.ID {
		case 29:
			r = "{{MaidenTicket}}"
		case 138:
			r = "{{AWCore}}"
		default:
			// exchange items
			r = fmt.Sprintf("[[File:%[1]s.png|28px|link=Items#%[1]s]] [[Items#%[1]s|%[1]s]]", item.NameEng)
		}
	} else if item.GroupID == 29 {
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
		if itemName == "" {
			// check rebirth items
			element := ""
			if strings.Contains(item.NameEng, "Stellar") {
				element = "Light"
			} else if strings.Contains(item.NameEng, "Flame") {
				element = "Passion"
			} else if strings.Contains(item.NameEng, "Frost") {
				element = "Cool"
			} else if strings.Contains(item.NameEng, "Crescent") {
				element = "Dark"
			}
			itype := ""
			if strings.Contains(item.NameEng, "Bud") {
				itype = "Bud"
			} else if strings.Contains(item.NameEng, "Bloom") {
				itype = "Bloom"
			} else if strings.Contains(item.NameEng, "Flora") {
				itype = "Flora"
			} else if strings.Contains(item.NameEng, "Secret Elixir") {
				itype = "Secret Elixir"
			} else if strings.Contains(item.NameEng, "Medicinal Herb") {
				itype = "Medicinal Herb"
			} else if strings.Contains(item.NameEng, "Zera") {
				itype = "Zera"
			}
			r = fmt.Sprintf("{{Flower|%s|%s}}", element, itype)
		} else {
			r = fmt.Sprintf("{{Stone|%s}}", itemName)
		}
	} else if item.GroupID == 32 {
		// ABB Ring
		r = fmt.Sprintf("{{Icon|%s}}", item.NameEng)
	} else if item.GroupID == 38 {
		// Custom Skill Recipies
		r = fmt.Sprintf("{{Skill Recipe|%s}}", cleanCustomSkillRecipe(item.NameEng))
	} else if item.GroupID == 39 {
		// Custom Skill items
		r = fmt.Sprintf("[[File:%[1]s.png|28px|link=Custom Skills#Skill_Materials]] [[Custom Skills#Skill_Materials|%[2]s]]",
			vc.CleanCustomSkillNoImage(item.NameEng),
			vc.CleanCustomSkillNoImage(item.NameEng),
		)
	} else {
		r = fmt.Sprintf("__UNKNOWN_GROUP:_%d_%s__", item.GroupID, item.NameEng)
	}
	return
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
||Event 10/15x damage<br />100/200%% Points+
||Event 10/15x damage<br />100/200%% Points+
||Event 10/15x damage<br />50/100%% Points+
||Event 10/15x damage<br />50/100%% Points+
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
| Mirror Maiden (LR) |Ring Exchange
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
To exchange Rings for prizes, go to '''Menu > Items > Tickets / Medals''' and use them.

%s

==Rewards==
%s

{{clr}}
==Local ABB Times==
{{AUBLocalTime|start jst = %[2]s|consecutive=1}}

{{NavEvent|%[15]s|%s}}
`
	case 18: // Tower Events
		return `{{Event|eventType = %d
|start jst = %s
|end jst = %s
|towerShield=%s
|image = Banner {{PAGENAME}}.png
||Ranking Reward<br />Amalgamation
||Floor Reward
||Floor Reward
||Amalgamation
||Amalgamation
||Amalgamation
||Fantasy Archwitch
||Archwitch
||Floor Reward<br/>Amalgamation Material
||Floor Reward<br/>Amalgamation Material
||Floor Reward<br/>Amalgamation Material
||Ranking Reward<br />Amalgamation Material
||Floor Reward<br />Amalgamation Material
||Ranking Reward<br />Amalgamation Material
||Event ATK and DEF 10x <br /> KO Gauge 100%% <br /> Pass 180%% / 460%% UP
||Event ATK and DEF 10x <br /> KO Gauge 100%% <br /> Pass 180%% / 460%% UP
||Event ATK and DEF 10x <br /> KO Gauge 100%% UP
}}

%s

==Rewards==
%s%s
{{clr}}

==Final Ranking==
%s

{{clr}}
{{NavEvent|%s|%s}}
`
	case 19: // Demon Realm
		return `{{Event|eventType = %d
|start jst = %s
|end jst = %s
|towerShield=%s
|story=yes
|image = Banner {{PAGENAME}}.png
||Ranking Reward<br>Amalgmation
||Point Rewards
||Point Rewards
||Point Rewards
||Point Rewards
||Fantasy Archwitch
||Archwitch
||ATK • DEF 10x<br>Soldiers +50%% / 100%%<br>Demon Core +20%% / 50%%<br>Pts +30%% / 100%%
||ATK • DEF 10x<br>Soldiers +50%% / 100%%<br>Demon Core +20%% / 50%%<br>Pts +30%% / 100%%
||ATK • DEF 10x<br>Soldiers +50%% / 100%%
}}

%s

==Demon Core Exchange==
{| class="article-table mw-collapsible mw-collapsed" border="1" cellpadding="1" cellspacing="1"
|-
! scope="col" |Prize
! scope="col" |Cost
! scope="col" |Limit

|-
| prize || [[File:Demon Core.png|30px|link=]] xCost || limit
|}

==Rewards==
%s%s
{{clr}}

==Final Ranking==
%s

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

func genWikiRankTrend(event *vc.Event, eventMap *vc.Map, midRewardTime time.Time, ranks []int, finalDayOnly bool) (rtrend string) {
	rtrend = `{| class="article-table rank-trend" border="1"
|-
!Date (JST)`
	for _, rank := range ranks {
		// headers
		rtrend += fmt.Sprintf(" !! Rank %d", rank)
	}
	if finalDayOnly {
		i := event.EndDatetime
		rtrend += fmt.Sprintf("\n|-\n|%s", i.Format("Jan _2"))
		for r := 0; r < len(ranks); r++ {
			rtrend += "\n|"
		}
	} else {
		for i := event.StartDatetime.Add(24 * time.Hour); event.EndDatetime.After(i.Add(-24 * time.Hour)); i = i.Add(24 * time.Hour) {
			ehElement := ""
			if eventMap != nil && !eventMap.ElementalhallStart.IsZero() && !i.Before(eventMap.ElementalhallStart.Time) {
				ehElement = fmt.Sprintf("{{ {{subst:#invoke:ElementalHall|elementForDate|%s|1 }} }} ", i.Format("January _2, 2006"))
			}
			pre := ""
			post := ""
			if !midRewardTime.IsZero() {
				midRewardDate := midRewardTime.Truncate(time.Duration(24) * time.Hour).Unix()
				iDate := i.Truncate(time.Duration(24) * time.Hour).Unix()
				if iDate == midRewardDate {
					pre = "{{tooltip|"
					post = "|Mid Ranking}}"
				}
			}
			rtrend += fmt.Sprintf("\n|-\n|%s%s%s%s", ehElement, pre, i.Format("Jan _2"), post)
			for r := 0; r < len(ranks); r++ {
				rtrend += "\n|"
			}
		}
	}
	rtrend += "\n|}"
	return
}

func genWikiExchange(exchanges []vc.GuildBingoExchangeReward) (ret string) {
	ret = `{| class="article-table mw-collapsible mw-collapsed sortable" border="1" cellpadding="1" cellspacing="1" style="width:400px;"
|-
! scope="col" |Prize
! scope="col" data-sort-type="number"|Cost
`
	atLeatOneLimitedItem := false
	for _, exchange := range exchanges {
		if exchange.ExchangeLimit > 0 {
			atLeatOneLimitedItem = true
			break
		}
	}

	if atLeatOneLimitedItem {
		ret += `! scope="col" data-sort-type="number"|Limit` + "\n"
	}
	for _, exchange := range exchanges {
		itemSortCode := fmt.Sprintf("%02d", (10 - exchange.RewardType))
		switch exchange.RewardType {
		case 1: // card
			card := vc.CardScan(exchange.RewardID, VcData.Cards)
			itemSortCode += fmt.Sprintf("%02d-%s", card.CardRareID, strings.Replace(card.Name, " ", "_", -1))
			ret += fmt.Sprintf("\n|-\n|data-sort-value=\"%s\"| {{Card Icon|%s}} ||data-sort-value=%d| x%[3]d",
				itemSortCode,
				card.Name,
				exchange.RequireNum,
			)
		case 2: //item
			item := vc.ItemScan(exchange.RewardID, VcData.Items)
			if item == nil {
				ret += fmt.Sprintf("\n|-\n|data-sort-value=\"%s\"|__UNKNOWN_ITEM_ID:%d__ ||data-sort-value=%d| x%[3]d",
					itemSortCode,
					exchange.RewardID,
					exchange.RequireNum,
				)
			} else {
				itemSortCode += fmt.Sprintf("%03d-%s",
					item.GroupID,
					strings.Replace(item.NameEng, " ", "_", -1),
				)
				ret += fmt.Sprintf("\n|-\n|data-sort-value=\"%s\"| %s ||data-sort-value=%d| x%[3]d",
					itemSortCode,
					getWikiItem(item),
					exchange.RequireNum,
				)
			}
		default: //unknown!
			ret += fmt.Sprintf("\n|-\n|Unknown reward type %d for reward id %d ||data-sort-value=%d| x%[3]d",
				exchange.RewardType,
				exchange.RewardID,
				exchange.RequireNum,
			)
		}
		if atLeatOneLimitedItem {
			var limit int
			var slimit string
			if exchange.ExchangeLimit > 0 {
				limit = exchange.ExchangeLimit
				slimit = strconv.Itoa(limit)
			} else {
				limit = int(int32(^uint32(0) >> 1)) // max int
				slimit = "Infinite"
			}
			ret += fmt.Sprintf("||data-sort-value=%d| %s", limit, slimit)
		}
	}
	ret += "\n|}"
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
