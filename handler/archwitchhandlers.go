package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"vc_file_grouper/vc"
)

// ArchwitchHandler displays archwitch data as a table.
func ArchwitchHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Archwitches</title></head><body>\n")
	io.WriteString(w, "<table><thead><tr><th>Series ID</th><th>Reward Card Name</th><th>Event Start</th><th>Event End</th><th>Recieve Limit</th><th>Is Beginner</th></tr></thead><tbody>\n")
	for i := len(vc.Data.ArchwitchSeries) - 1; i >= 0; i-- {
		series := vc.Data.ArchwitchSeries[i]
		rewardCard := vc.CardScan(series.RewardCardID)
		fmt.Fprintf(w,
			"<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>",
			series.ID,
			imageLink(rewardCard),
			series.Description,
			series.PublicStartDatetime.Format(time.RFC3339),
			series.PublicEndDatetime.Format(time.RFC3339),
			series.ReceiveLimitDatetime.Format(time.RFC3339),
			series.IsBeginnerKing,
		)
		io.WriteString(w, "\n<tr><td></td><td></td><td colspan=5><table border=1>")
		io.WriteString(w, "<thead><tr><th>ID</th><th>Card Master / servants</th><th>Skill 1</th><th>Skill 2</th><th>Status Group</th><th>Public</th><th>Rarity</th><th>RareIntensity</th><th>Battle Time</th><th>Exp</th><th>Max Friendship</th><th>Weather</th><th>Model</th><th>Chain Ratio 2</th><th>Likability</th></tr></thead><tbody>")
		for _, aw := range series.Archwitches() {
			cardMaster := vc.CardScan(aw.CardMasterID)
			skill1 := vc.SkillScan(aw.SkillID1)
			skill2 := vc.SkillScan(aw.SkillID2)
			//servant1 := vc.CardScanCharacter(aw.ServantID1)
			//servant2 := vc.CardScanCharacter(aw.ServantID2)
			fmt.Fprintf(w,
				"<tr><td>%d</td><td>%s",
				aw.ID,
				imageLink(cardMaster),
			)
			if aw.ServantID1 > 0 {
				fmt.Fprintf(w,
					"<br />%s (%d)<br />%s (%d)",
					imageLink(nil),
					aw.ServantID1,
					imageLink(nil),
					aw.ServantID2,
				)
			}
			io.WriteString(w, "</td>")
			fmt.Fprintf(w, "<td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%d</td><td>%s</td>",
				printSkill(skill1),
				printSkill(skill2),
				aw.StatusGroupID,
				aw.PublicFlg,
				aw.RareFlg,
				aw.RareIntensity,
				aw.BattleTime,
				aw.Exp,
				aw.MaxFriendship,
				aw.WeatherID,
				aw.ModelName,
				aw.ChainRatio2,
				printFriendship(aw),
			)
		}
		io.WriteString(w, "</tr></tbody></table>\n")
	}
	io.WriteString(w, "</tbody></table>\n")
	io.WriteString(w, "</body></html>")
}

func imageLink(card *vc.Card) string {
	if card == nil {
		return ""
	}
	return fmt.Sprintf("<img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">(%d) %s</a>",
		card.Image(),
		card.ID,
		card.ID,
		card.Name,
	)
}

func printSkill(skill *vc.Skill) string {
	if skill == nil {
		return ""
	}
	return fmt.Sprintf("<b>%s</b><br />%s<br /> Procs: %d<br /> Chance: %d%% - %d%%", skill.Name, skill.FireMin(), skill.Activations(), skill.DefaultRatio, skill.MaxRatio)
}

func printFriendship(aw *vc.Archwitch) string {
	s := "<ol>"
	for _, af := range aw.Likeability() {
		s = s + fmt.Sprintf("<li>%d%%: \"%s\"</li>\n", af.UpRate, af.Likability)
	}
	return s + "</ol>"
}
