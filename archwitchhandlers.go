package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func archwitchHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Archwitches</title></head><body>\n")
	io.WriteString(w, "<table><thead><tr><th>Series ID</th><th>Reward Card Name</th><th>Event Start</th><th>Event End</th><th>Recieve Limit</th><th>Is Beginner</th></tr></thead><tbody>\n")
	for _, series := range VcData.ArchwitchSeries {
		rewardCard := vc.CardScan(series.RewardCardId, VcData.Cards)
		fmt.Fprintf(w,
			"<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>",
			series.Id,
			imageLink(rewardCard),
			series.Description,
			series.PublicStartDatetime.Format(time.RFC3339),
			series.PublicEndDatetime.Format(time.RFC3339),
			series.ReceiveLimitDatetime.Format(time.RFC3339),
			series.IsBeginnerKing,
		)
		io.WriteString(w, "\n<tr><td></td><td></td><td colspan=5><table border=1>")
		io.WriteString(w, "<thead><tr><th>ID</th><th>Card Master / servants</th><th>Status Group</th><th>Public</th><th>Rarity</th><th>RareIntensity</th><th>Battle Time</th><th>Exp</th><th>Max Friendship</th><th>Skill 1</th><th>Skill 2</th><th>Weather</th><th>Model</th><th>Chain Ratio 2</th><th>Likability</th></tr></thead><tbody>")
		for _, aw := range series.Archwitches(VcData) {
			cardMaster := vc.CardScan(aw.CardMasterId, VcData.Cards)
			skill1 := vc.SkillScan(aw.SkillId1, VcData.Skills)
			skill2 := vc.SkillScan(aw.SkillId2, VcData.Skills)
			//servant1 := vc.CardScanCharacter(aw.ServantId1, VcData.Cards)
			//servant2 := vc.CardScanCharacter(aw.ServantId2, VcData.Cards)
			fmt.Fprintf(w,
				"<tr><td>%d</td><td>%s",
				aw.Id,
				imageLink(cardMaster),
			)
			if aw.ServantId1 > 0 {
				fmt.Fprintf(w,
					"<br />%s (%d)<br />%s (%d)",
					imageLink(nil),
					aw.ServantId1,
					imageLink(nil),
					aw.ServantId2,
				)
			}
			io.WriteString(w, "</td>")
			fmt.Fprintf(w, "<td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td>%d</td><td>%s</td>",
				aw.StatusGroupId,
				aw.PublicFlg,
				aw.RareFlg,
				aw.RareIntensity,
				aw.BattleTime,
				aw.Exp,
				aw.MaxFriendship,
				printSkill(skill1),
				printSkill(skill2),
				aw.WeatherId,
				aw.ModelName,
				aw.ChainRatio2,
				printFriendship(&aw),
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
		card.Id,
		card.Id,
		card.Name,
	)
}

func printSkill(skill *vc.Skill) string {
	if skill == nil {
		return ""
	}
	return "<b>" + skill.Name + "</b><br />" + skill.FireMin()
}

func printFriendship(aw *vc.Archwitch) string {
	s := "<ol>"
	for _, af := range aw.Likeability(VcData) {
		s = s + fmt.Sprintf("<li>%d%% chance to the next heart</li>\n", af.UpRate)
	}
	return s + "</ol>"
}
