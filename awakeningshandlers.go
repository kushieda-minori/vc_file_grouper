package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

func awakeningsTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Awakenings</title></head><body>\n")
	io.WriteString(w, "<table><thead><tr><th>From Card</th><th>To Card</th><th>Chance</th><th>Crystals</th><th>Orb</th><th>Large</th><th>Medium</th><th>Small</th></tr></thead><tbody>\n")
	for _, value := range VcData.Awakenings {
		baseCard := vc.CardScan(value.BaseCardID, VcData.Cards)
		resultCard := vc.CardScan(value.ResultCardID, VcData.Cards)
		fmt.Fprintf(w,
			"<tr><td><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td><td><img src=\"/images/cardthumb/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td><td>%d%%</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			baseCard.Image(),
			baseCard.ID,
			baseCard.Name,
			resultCard.Image(),
			resultCard.ID,
			resultCard.Name,
			value.Percent,
			value.Material5Count,
			value.Material1Count,
			value.Material2Count,
			value.Material3Count,
			value.Material4Count,
		)
	}
	io.WriteString(w, "</tbody></table>\n")
	io.WriteString(w, "</body></html>")
}

func awakeningsCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-awaken-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{
		"ID",
		"Name",
		"BaseCardID",
		"ResultCardID",
		"Percent",
		"Material1Item ",
		"Material1Count",
		"Material2Item ",
		"Material2Count",
		"Material3Item ",
		"Material3Count",
		"Material4Item ",
		"Material4Count",
		"Material5Item ",
		"Material5Count",
		"Order",
		"IsClosed",
	})
	for _, value := range VcData.Awakenings {
		baseCard := vc.CardScan(value.BaseCardID, VcData.Cards)
		cw.Write([]string{
			strconv.Itoa(value.ID),
			baseCard.Name,
			strconv.Itoa(value.BaseCardID),
			strconv.Itoa(value.ResultCardID),
			strconv.Itoa(value.Percent),
			strconv.Itoa(value.Material1Item),
			strconv.Itoa(value.Material1Count),
			strconv.Itoa(value.Material2Item),
			strconv.Itoa(value.Material2Count),
			strconv.Itoa(value.Material3Item),
			strconv.Itoa(value.Material3Count),
			strconv.Itoa(value.Material4Item),
			strconv.Itoa(value.Material4Count),
			strconv.Itoa(value.Material5Item),
			strconv.Itoa(value.Material5Count),
			strconv.Itoa(value.Order),
			strconv.Itoa(value.IsClosed),
		})
	}
	cw.Flush()
}
