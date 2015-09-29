package main

import (
	"fmt"
	"io"
	"net/http"
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
}
