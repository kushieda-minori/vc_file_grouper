package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// DungeonScenarioHandler outputs the raw JSON
func DungeonScenarioHandler(w http.ResponseWriter, r *http.Request) {
	// string file location vcRoot/scenario/MsgScenarioString_<lang>.strb
	lines, err := vc.ReadStringFileFilter(vc.FilePath+"/scenario/MsgScenarioString_"+vc.LangPack+".strb", false)
	io.WriteString(w, "<html><head><title>DRV Scenario</title></head><body>\n")
	if err != nil {
		// write out our error...
		io.WriteString(w, "<pre>")
		fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
		io.WriteString(w, "</pre>")
		io.WriteString(w, "</body></html>")
		return
	}

	io.WriteString(w, "<pre>")
	scenario := 0
	for _, line := range lines {
		if line == "" {
			io.WriteString(w, "\n")
		}
		if strings.HasPrefix(line, "Chapter") {
			if strings.HasPrefix(line, "Chapter 1") {
				scenario++
				fmt.Fprintf(w, "</pre><h1>Scenario %d</h1><pre>", scenario)
			}
			io.WriteString(w, "\n")
			lines := filterStoryLine(line)
			for _, l := range lines {
				io.WriteString(w, l)
				io.WriteString(w, "\n")
			}
		} else {
			lines := filterStoryLine(line)
			if len(lines) > 0 {
				io.WriteString(w, ";[[SPEAKER]]\n")
				for _, l := range lines {
					io.WriteString(w, ":")
					io.WriteString(w, l)
					io.WriteString(w, "\n")
				}
			}
		}
	}
	io.WriteString(w, "</pre></body></html>")
}

func filterStoryLine(line string) []string {
	line = strings.ReplaceAll(line, "\n", " ")
	line = strings.ReplaceAll(line, "  ", " ")
	line = strings.ReplaceAll(line, "<i><break>", "\n")
	line = strings.TrimSpace(line)
	lines := strings.Split(line, "\n")

	ll := len(lines)
	ret := make([]string, 0, ll)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			ret = append(ret, l)
		}
	}
	return ret
}
