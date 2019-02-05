package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// ItemHandler shows item details as a table
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>All Items</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>_id</th><th>Item Name</th><th>Image</th><th>Description</th><th>Group</th><th>End Date</th><th>Max Own</th><th>Limited Item</th><th>Is Deleted</th>\n")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")
	for i := len(vc.Data.Items) - 1; i >= 0; i-- {
		e := vc.Data.Items[i]
		fmt.Fprintf(w, "<tr>"+
			"<td>%d</td>"+
			"<td>%s<br />%s</td>"+
			"<td><a href=\"/images/item/shop/%[5]d?filename=%[4]s\"><img src=\"/images/item/shop/%[5]d\"/></a></td>"+
			"<td><p>Description: %s</p><p>Shop Description: %s</p><p>Sub Item Description: %s</p><p>Use: %s</p></td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"</tr>",
			e.ID,
			vc.CleanCustomSkillImage(e.Name),
			vc.CleanCustomSkillImage(e.NameEng),
			url.QueryEscape(vc.CleanCustomSkillNoImage(e.NameEng)),
			e.ItemNo,
			e.Description,
			e.DescriptionInShop,
			e.DescriptionSub,
			e.MsgUse,
			e.GroupID,
			e.EndDatetime.Format(time.RFC3339),
			e.MaxCount,
			e.LimitedItemFlg,
			e.IsDelete,
		)
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
