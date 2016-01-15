package main

import (
	"fmt"
	"io"
	"net/http"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func archwitchHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Archwitches</title></head><body>\n")
	io.WriteString(w, "<table><thead><tr><th>Archwitch ID</th><th>Card Name</th><th>Max Friendship</th><th>Series</th><th>Battle Time (Minues)</th><th>Exp</th><th>Chain Ratio</th><th>Skill 1</th><th>Skill 2</th></tr></thead><tbody>\n")
	for _, value := range VcData.Awakenings {
		baseCard := vc.CardScan(value.BaseCardId, VcData.Cards)
		resultCard := vc.CardScan(value.ResultCardId, VcData.Cards)
		fmt.Fprintf(w,
			"<tr><td><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td><td><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td><td>%d%%</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			baseCard.Image(),
			baseCard.Id,
			baseCard.Name,
			resultCard.Image(),
			resultCard.Id,
			resultCard.Name,
			value.Percent,
			value.Material1Count,
			value.Material2Count,
			value.Material3Count,
			value.Material4Count,
		)
	}
	io.WriteString(w, "</tbody>\n")
	io.WriteString(w, "</body></html>")
}
