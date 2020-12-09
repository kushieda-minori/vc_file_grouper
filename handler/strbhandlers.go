package handler

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"vc_file_grouper/vc"
)

// StrbHandler handle STRB files
func StrbHandler(w http.ResponseWriter, r *http.Request) {
	fpath := r.URL.Path
	var pathLen int
	if fpath[len(fpath)-1] == '/' {
		pathLen = len(fpath) - 1
	} else {
		pathLen = len(fpath)
	}

	pathParts := strings.Split(fpath[1:pathLen], "/")
	// "strb/id/TYPE"
	if len(pathParts) < 3 {
		StrbTableHandler(w, r)
		return
	}

	lpathParts := len(pathParts)

	// validate that the file is a valid strb file
	strbFile := filepath.Join(pathParts[1 : lpathParts-1]...)
	if strings.HasPrefix(strbFile, ".") || !strings.HasSuffix(strings.ToLower(strbFile), ".strb") {
		fmt.Fprintf(w, "Illegal file path: %s", strbFile)
		return
	}
	fullPath := filepath.Join(vc.FilePath, strbFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Fprintf(w, "File does not exist: %s", strbFile)
		return
	}

	ftype := pathParts[lpathParts-1]
	if ftype == "html" {
		contents, err := vc.ReadStringFileFilter(fullPath, false)
		if err != nil {
			fmt.Fprintf(w, "Error reading file %s: %s", strbFile, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, "<html><head><title>Strb Events</title>\n")
		io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
		io.WriteString(w, "</head><body>\n")
		io.WriteString(w, "<div>\n")
		io.WriteString(w, "<table><thead><tr>\n")
		io.WriteString(w, "<th>line #</th><th>String</th>")
		io.WriteString(w, "</tr></thead>\n")
		io.WriteString(w, "<tbody>\n")
		for i, line := range contents {
			fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td></tr>\n", i+1, html.EscapeString(line))
		}
		io.WriteString(w, "</tbody></table></div></body></html>")
	} else if ftype == "txt" {
		contents, err := vc.ReadStringFileFilter(fullPath, false)
		if err != nil {
			fmt.Fprintf(w, "Error reading file %s: %s", strbFile, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		for _, line := range contents {
			fmt.Fprintf(w, "%s\n\n", line)
		}
	} else {
		fmt.Fprintf(w, "Unsupported file type for conversion: %s", ftype)
		return
	}

}

func checkStrbName(info os.FileInfo) bool {
	name := info.Name()
	return strings.HasSuffix(strings.ToLower(name), ".strb")
}

// StrbTableHandler shows strb events as a table
func StrbTableHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Strb Events</title>\n")
	io.WriteString(w, "<style>table, th, td {border: 1px solid black;};</style>")
	io.WriteString(w, "</head><body>\n")
	io.WriteString(w, "<div>\n")
	io.WriteString(w, "<table><thead><tr>\n")
	io.WriteString(w, "<th>File</th><th>&nbsp;</th><th>&nbsp;</th>")
	io.WriteString(w, "</tr></thead>\n")
	io.WriteString(w, "<tbody>\n")

	fullpath := vc.FilePath

	err := filepath.Walk(fullpath, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !checkStrbName(info) {
			return nil
		}
		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		b := make([]byte, 4)
		_, err = f.Read(b)
		f.Close()
		if err != nil {
			return err
		}
		if bytes.Equal(b, []byte("STRB")) {
			relPath, _ := filepath.Rel(fullpath, fpath)
			relPath = filepath.ToSlash(relPath)
			fmt.Fprintf(w, "<tr>"+
				"<td>%[1]s</td>"+
				"<td><a href=\"/strb/%[2]s/html\">html</a></td>"+
				"<td><a href=\"/strb/%[2]s/txt\">txt</a></td>",
				html.EscapeString(relPath),
				url.QueryEscape(relPath),
			)
		} else {
			log.Printf("STRB file is not encoded: %s", fullpath)
		}
		return nil
	})
	if err != nil {
		io.WriteString(w, err.Error()+"<br />\n")
	}
	io.WriteString(w, "</tbody></table></div></body></html>")
}
