package main

import (
	//"bufio"
	//"github.com/bitly/go-simplejson"
	//"sort"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"zetsuboushita.net/vc_file_grouper/vc"
)

var VcData vc.VcFile
var masterDataStr string

// Main function that starts the program
func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		usage()
		return
	}

	b, err := VcData.Read(os.Args[1])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}
	masterDataStr = string(b)

	//main page
	http.HandleFunc("/", masterDataHandler)
	//image locations
	http.Handle("/cardimages/", http.StripPrefix("/cardimages/", http.FileServer(http.Dir(os.Args[1]+"/card/md"))))
	http.Handle("/cardthumbs/", http.StripPrefix("/cardthumbs/", http.FileServer(http.Dir(os.Args[1]+"/card/thumb"))))
	http.Handle("/cardimagesHD/", http.StripPrefix("/cardimagesHD/", http.FileServer(http.Dir(os.Args[1]+"/../hd/"))))
	http.Handle("/eventimages/", http.StripPrefix("/eventimages/", http.FileServer(http.Dir(os.Args[1]+"/event/largeimage"))))

	//dynamic pages
	http.HandleFunc("/cards/", cardHandler)
	http.HandleFunc("/cards/table/", cardTableHandler)
	http.HandleFunc("/cards/csv/", cardCsvHandler)
	// http.HandleFunc("/cards/raw/", cardRawHandler)
	http.HandleFunc("/cards/detail/", cardDetailHandler)
	http.HandleFunc("/character/", characterTableHandler)
	http.HandleFunc("/character/csv/", characterCsvHandler)
	http.HandleFunc("/skills/", skillTableHandler)
	http.HandleFunc("/skills/csv/", skillCsvHandler)
	http.HandleFunc("/awakenings/", awakeningsTableHandler)
	http.HandleFunc("/awakenings/csv/", awakeningsCsvHandler)
	http.HandleFunc("/decode/", decodeHandler)
	http.HandleFunc("/raw/", rawDataHandler)

	os.Stdout.WriteString("Listening on port 8585\n")
	err = http.ListenAndServe(":8585", nil)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
}

// Prints useage to the console
func usage() {
	os.Stderr.WriteString("You must pass the location of the files.\n" +
		"Usage: " + os.Args[0] + " /path/to/com.nubee.valkyriecrusade/files\n")
}

//Main index page
func masterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	io.WriteString(w, "<html><body>\n")
	io.WriteString(w, "<a href=\"/decode\" >Decode All Files</a><br />\n")
	io.WriteString(w, "<a href=\"/cards\" >Card List</a><br />\n")
	io.WriteString(w, "<a href=\"/cards/table\" >Card List as a Table</a><br />\n")
	io.WriteString(w, "<a href=\"/cards/csv\" >Card List as CSV</a><br />\n")
	io.WriteString(w, "<a href=\"/skills/csv\" >Skill List as CSV</a><br />\n")
	io.WriteString(w, "<a href=\"/awakenings\" >List of Awakenings</a><br />\n")
	io.WriteString(w, "<a href=\"/awakenings/csv\" >List of Awakenings as CSV</a><br />\n")
	io.WriteString(w, "<a href=\"/raw\" >Raw data</a><br />\n")

	io.WriteString(w, "</body></html>")
}

func rawDataHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	io.WriteString(w, "<html><body>\n")

	io.WriteString(w, "<pre>")
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(masterDataStr), "", "\t")
	if err != nil {
		fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
		return
	}

	io.WriteString(w, string(prettyJSON.Bytes()))
	io.WriteString(w, "</pre>")
	// // read the decoded master data file
	// masterData, _ := simplejson.NewJson(bytes(masterDataStr))

	// *masterDataMap, _ = masterData.Map()

	// // sort the keys
	// mk := make([]string, len(*masterDataMap))
	// i := 0
	// var k string
	// for k, _ = range *masterDataMap {
	// 	mk[i] = k
	// 	i++
	// }
	// sort.Strings(mk)

	// // print all the keys in order
	// for _, k = range mk {
	// 	io.WriteString(w, "<p><a href=\"")
	// 	io.WriteString(w, k)
	// 	io.WriteString(w, "\">")
	// 	io.WriteString(w, k)
	// 	io.WriteString(w, "</a><div style=\"padding-left:10px\">")
	// 	// val, _ := masterData.Get(k).EncodePretty()
	// 	// io.WriteString(w, strings.Replace(string(val), "\n", "<br />", -1))
	// 	io.WriteString(w, "</div></p>\n")
	// }

	io.WriteString(w, "</body></html>")
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>File Decode</title></head><body>\nDecodng files<br />\n")
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		b := make([]byte, 4)
		_, err = f.Read(b)
		f.Close()
		if err != nil {
			return err
		}
		if bytes.Equal(b, []byte("CODE")) {
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
