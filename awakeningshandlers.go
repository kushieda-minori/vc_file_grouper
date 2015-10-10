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
	io.WriteString(w, "<table><thead><tr><th>From Card</th><th>To Card</th></tr></thead><tbody>\n")
	for _, value := range VcData.Awakenings {
		baseCard := vc.CardScan(value.BaseCardId, VcData.Cards)
		resultCard := vc.CardScan(value.ResultCardId, VcData.Cards)
		fmt.Fprintf(w,
			"<tr><td><img src=\"/cardthumbs/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td><td><img src=\"/cardthumbs/%s\"/><br /><a href=\"/cards/detail/%d\">%s</a></td></tr>",
			baseCard.Image(),
			baseCard.Id,
			baseCard.Name,
			resultCard.Image(),
			resultCard.Id,
			resultCard.Name)
	}
	io.WriteString(w, "</tbody>\n")
	io.WriteString(w, "</body></html>")
}

func awakeningsCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-awaken-"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{
		"Id",
		"Name",
		"BaseCardId",
		"ResultCardId",
		"Percent",
		"Material1Item ",
		"Material1Count",
		"Material2Item ",
		"Material2Count",
		"Material3Item ",
		"Material3Count",
		"Material4Item ",
		"Material4Count",
		"Order",
		"IsClosed",
	})
	for _, value := range VcData.Awakenings {
		baseCard := vc.CardScan(value.BaseCardId, VcData.Cards)
		cw.Write([]string{
			strconv.Itoa(value.Id),
			baseCard.Name,
			strconv.Itoa(value.BaseCardId),
			strconv.Itoa(value.ResultCardId),
			strconv.Itoa(value.Percent),
			strconv.Itoa(value.Material1Item),
			strconv.Itoa(value.Material1Count),
			strconv.Itoa(value.Material2Item),
			strconv.Itoa(value.Material2Count),
			strconv.Itoa(value.Material3Item),
			strconv.Itoa(value.Material3Count),
			strconv.Itoa(value.Material4Item),
			strconv.Itoa(value.Material4Count),
			strconv.Itoa(value.Order),
			strconv.Itoa(value.IsClosed),
		})
	}
	cw.Flush()
}
