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
	"strconv"
	"time"
	"zetsuboushita.net/vc_file_grouper/vc"
)

var VcData *vc.VcFile
var masterDataStr string
var vcfilepath string

// Main function that starts the program
func main() {
	if len(os.Args) == 1 {
		vcfilepath = "."
	} else {
		vcfilepath = os.Args[1]
	}

	if _, err := os.Stat(vcfilepath); os.IsNotExist(err) {
		usage()
		VcData = &vc.VcFile{}
		//return
	} else {
		readMasterData(vcfilepath)
	}

	//main page
	http.HandleFunc("/", masterDataHandler)
	//image locations
	http.HandleFunc("/images/card/", imageCardHandler)
	http.HandleFunc("/images/cardthumb/", imageCardThumbHandler)
	http.HandleFunc("/images/cardHD/", imageCardHDHandler)
	http.HandleFunc("/images/event/", imageEventHandler)
	http.HandleFunc("/images/battle/bg/", imageBattleBGHandler)
	http.HandleFunc("/images/battle/map/", imageBattleMapHandler)

	// vc master data
	http.HandleFunc("/data/", dataHandler)
	//dynamic pages
	http.HandleFunc("/cards/", cardHandler)
	http.HandleFunc("/cards/table/", cardTableHandler)
	http.HandleFunc("/cards/csv/", cardCsvHandler)
	http.HandleFunc("/cards/detail/", cardDetailHandler)
	http.HandleFunc("/archwitches/", archwitchHandler)
	// http.HandleFunc("/character/", characterTableHandler)
	// http.HandleFunc("/character/csv/", characterCsvHandler)

	http.HandleFunc("/skills/", skillTableHandler)
	http.HandleFunc("/skills/csv/", skillCsvHandler)

	http.HandleFunc("/deckbonus/", deckBonusHandler)
	http.HandleFunc("/deckbonus/WIKI/", deckBonusWikiHandler)

	http.HandleFunc("/events/", eventHandler)
	http.HandleFunc("/events/detail/", eventDetailHandler)

	http.HandleFunc("/maps/", mapHandler)

	http.HandleFunc("/awakenings/", awakeningsTableHandler)
	http.HandleFunc("/awakenings/csv/", awakeningsCsvHandler)

	http.HandleFunc("/decode/", decodeHandler)

	http.HandleFunc("/raw/", rawDataHandler)

	http.HandleFunc("/SHUTDOWN/", func(w http.ResponseWriter, r *http.Request) { os.Exit(0) })

	os.Stdout.WriteString("Listening on port 8585. Connect to http://localhost:8585/\nPress <CTRL>+C to stop or close the terminal.\n")
	err := http.ListenAndServe("localhost:8585", nil)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
}

// Prints useage to the console
func usage() {
	os.Stderr.WriteString("You must pass the location of the files.\n" +
		"Usage: " + os.Args[0] + " /path/to/com.nubee.valkyriecrusade/files\n")
}

func readMasterData(files string) error {
	data := vc.VcFile{}
	b, err := data.Read(files)
	if err != nil {
		if VcData == nil {
			VcData = &data
		}
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}
	VcData = &data
	masterDataStr = string(b)
	return nil
}

//Main index page
func masterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	fmt.Fprintf(w, `<html><body>
<p>Version: %d,&nbsp;&nbsp;&nbsp;&nbsp;Timestamp: %d,&nbsp;&nbsp;&nbsp;&nbsp;JST: %s</p>
<a href="/data" >Set Data Location</a><br />
<br />
<a href="/decode" >Decode All Files</a><br />
<br />
<a href="/cards/table" >Card List as a Table</a><br />
<a href="/deckbonus" >Deck Bonuses</a><br />
<a href="/events" >Event List</a><br />
<a href="/archwitches" >Archwitch List</a><br />
<br />
<a href="/images/event/">Event Images</a>
<br />
<a href="/cards/csv" >Card List as CSV</a><br />
<a href="/skills/csv" >Skill List as CSV</a><br />
<br />
<a href="/awakenings" >List of Awakenings</a><br />
<a href="/awakenings/csv" >List of Awakenings as CSV</a><br />
<a href="/raw" >Raw data</a><br />
<br />
<br />
<a href="/SHUTDOWN" >SHUTDOWN</a><br />
</body></html>`,
		VcData.Version,
		VcData.Common.UnixTime.Unix(),
		VcData.Common.UnixTime.Format(time.RFC3339),
	)
	// io.WriteString(w, "<a href=\"/cards\" >Card List</a><br />\n")
}

func rawDataHandler(w http.ResponseWriter, r *http.Request) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(masterDataStr), "", "\t")
	if err != nil {
		// File header
		io.WriteString(w, "<html><body>\n")

		io.WriteString(w, "<pre>")
		fmt.Fprintf(w, " : ERROR: %s<br />\n", err.Error())
		io.WriteString(w, "</pre>")
		io.WriteString(w, "</body></html>")
		return
	}
	w.Header().Set("Content-Disposition", "filename="+"vcData-raw-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".json")
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, string(prettyJSON.Bytes()))
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>File Decode</title></head><body>\nDecodng files<br />\n")
	err := filepath.Walk(vcfilepath, func(path string, info os.FileInfo, err error) error {
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

func removeDuplicates(a []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}
