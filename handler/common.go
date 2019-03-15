package handler

import (
	"fmt"
	"net/http"
	"strings"

	"zetsuboushita.net/vc_file_grouper/util"
	"zetsuboushita.net/vc_file_grouper/vc"
)

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
	ret = strings.ReplaceAll(ret, "ticket", "")
	ret = strings.ReplaceAll(ret, "summon", "")
	ret = strings.ReplaceAll(ret, "guaranteed", "")
	ret = strings.ReplaceAll(ret, "★★★", "3star")
	ret = strings.ReplaceAll(ret, "★★", "2star")
	ret = strings.ReplaceAll(ret, "★", "1star")
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanItemName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.ReplaceAll(ret, "valkyrie", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanArcanaName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.ReplaceAll(ret, "arcana's", "")
	ret = strings.ReplaceAll(ret, "arcana", "")
	ret = strings.ReplaceAll(ret, "%", "")
	ret = strings.ReplaceAll(ret, "+", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.Replace(ret, "forced", "", 1)
	ret = strings.Replace(ret, "strongdef", "def", 1)
	ret = strings.TrimSpace(ret)
	return ret
}

func isChecked(values []string, e string) string {
	if util.Contains(values, e) {
		return "checked"
	}
	return ""
}

func printHTMLTable(w http.ResponseWriter, headers []string, bodyRows [][]interface{}) {
	fmt.Fprintf(w, "<table>")
	printHTMLTableHeader(w, headers...)

	fmt.Fprintf(w, "<tbody>\n")
	for _, row := range bodyRows {
		printHTMLTableRow(w, row)
	}
	fmt.Fprintf(w, "\n</tbody>")
	fmt.Fprintf(w, "\n</table>")
}

func printHTMLTableHeader(w http.ResponseWriter, headers ...string) {
	fmt.Fprintf(w, "<thead>\n<tr>\n")
	for _, header := range headers {
		fmt.Fprintf(w, "<th>%s</th>", header)
	}
	fmt.Fprintf(w, "\n</tr>\n</thead>\n")
}
func printHTMLTableRow(w http.ResponseWriter, columns ...interface{}) {
	fmt.Fprintf(w, "<tr>\n")
	for _, col := range columns {
		fmt.Fprintf(w, "<td>%v</td>", col)
	}
	fmt.Fprintf(w, "\n</tr>\n")
}
