package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// RawDataHandler outputs the raw JSON
func RawDataHandler(w http.ResponseWriter, r *http.Request) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(vc.MasterDataStr), "", "\t")
	if err != nil {
		// File header
		io.WriteString(w, "<html><body>\n")

		io.WriteString(w, "<pre>")
		fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
		io.WriteString(w, "</pre>")
		io.WriteString(w, "</body></html>")
		return
	}
	w.Header().Set("Content-Disposition", "filename="+"vcData-raw-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".json")
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(prettyJSON.Bytes()))
}

// RawDataKeysHandler outputs all keys in the main JSON object
func RawDataKeysHandler(w http.ResponseWriter, r *http.Request) {
	c := make(map[string]interface{})
	err := json.Unmarshal([]byte(vc.MasterDataStr), &c)
	if err != nil {
		// File header
		io.WriteString(w, "<html><body>\n")

		io.WriteString(w, "<pre>")
		fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
		io.WriteString(w, "</pre>")
		io.WriteString(w, "</body></html>")
		return
	}
	w.Header().Set("Content-Disposition", "filename="+"vcData-raw-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".json")
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, "[\n")

	mk := make([]string, len(c))
	i := 0
	for s := range c {
		mk[i] = s
		i++
	}
	sort.Strings(mk)

	for _, s := range mk {
		fmt.Fprintf(w, "\t%s,\n", s)
	}
	io.WriteString(w, "]")
}

// DecodeHandler decodes all files
func DecodeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>File Decode</title></head><body>\nDecodng files<br />\n")
	err := filepath.Walk(vc.FilePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileEncoded, err := vc.IsFileEncoded(path)
		if err != nil {
			return err
		}
		if fileEncoded {
			fmt.Fprintf(w, "Decoding: %s ", path)
			nf, _, err := vc.DecodeAndSave(path)
			if err != nil {
				fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
				return err
			}
			fmt.Fprintf(w, " : %s<br />\n", nf)
		}
		return nil
	})
	if err == nil {
		io.WriteString(w, "Decode complete<br />\n")
	} else {
		io.WriteString(w, err.Error()+"<br />\n")
	}
	io.WriteString(w, "</body></html>")
}
