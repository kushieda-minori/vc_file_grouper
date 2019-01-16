package main

import (
	//"bufio"
	//"github.com/bitly/go-simplejson"
	//"sort"
	"bytes"
	"encoding/json"
	"flag"
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

// VcData Main data file
var VcData *vc.VFile
var masterDataStr string
var vcfilepath string

// Main function that starts the program
func main() {

	cmdLang := flag.String("lang", "en", "The language pack to use. 'en' for English, 'zhs' for Chinese. ")
	flag.Parse()

	if cmdLang == nil {
		vc.LangPack = "en"
	} else {
		vc.LangPack = *cmdLang
	}

	if len(flag.Args()) == 0 {
		vcfilepath = "."
	} else {
		vcfilepath = flag.Args()[0]
	}

	if _, err := os.Stat(vcfilepath); os.IsNotExist(err) {
		usage()
		VcData = &vc.VFile{}
		//return
	} else {
		readMasterData(vcfilepath)
	}

	//main page
	http.HandleFunc("/", masterDataHandler)
	http.HandleFunc("/css/", cssHandler)
	//image locations
	http.HandleFunc("/images/card/", imageCardHandler)
	http.HandleFunc("/images/cardthumb/", imageCardThumbHandler)
	http.HandleFunc("/images/cardHD/", imageCardHDHandler)
	http.HandleFunc("/images/event/", imageHandlerFor("/event/", "/event/"))
	http.HandleFunc("/images/battle/", imageHandlerFor("/battle/", "/battle/"))
	http.HandleFunc("/images/garden/", imageHandlerFor("/garden/", "/garden/"))
	http.HandleFunc("/images/garden/map/", handleStructureImages)
	http.HandleFunc("/images/dungeon/", imageHandlerFor("/dungeon/", "/dungeon/"))
	http.HandleFunc("/images/alliance/", imageHandlerFor("/alliance/", "/guild/"))
	http.HandleFunc("/images/summon/", imageHandlerFor("/summon/", "/gacha/"))
	http.HandleFunc("/images/item/", imageHandlerFor("/item/", "/item/"))
	http.HandleFunc("/images/treasure/", imageHandlerFor("/treasure/", "/treasure/"))
	http.HandleFunc("/images/navi/", imageHandlerFor("/navi/", "/navi/"))

	// vc master data
	http.HandleFunc("/data/", dataHandler)
	//dynamic pages
	http.HandleFunc("/cards/", cardHandler)
	http.HandleFunc("/cards/table/", cardTableHandler)
	http.HandleFunc("/cards/csv/", cardCsvHandler)
	http.HandleFunc("/cards/glrcsv/", cardCsvGLRHandler)
	http.HandleFunc("/cards/detail/", cardDetailHandler)
	http.HandleFunc("/cards/levels/", cardLevelHandler)
	http.HandleFunc("/archwitches/", archwitchHandler)
	http.HandleFunc("/characters/", characterTableHandler)
	http.HandleFunc("/characters/detail/", characterDetailHandler)
	// http.HandleFunc("/character/csv/", characterCsvHandler)

	http.HandleFunc("/items/", itemHandler)

	http.HandleFunc("/skills/", skillTableHandler)
	http.HandleFunc("/skills/csv/", skillCsvHandler)

	http.HandleFunc("/deckbonus/", deckBonusHandler)
	http.HandleFunc("/deckbonus/WIKI/", deckBonusWikiHandler)

	http.HandleFunc("/events/", eventHandler)
	http.HandleFunc("/events/detail/", eventDetailHandler)

	http.HandleFunc("/thor/", thorHandler)

	http.HandleFunc("/maps/", mapHandler)

	http.HandleFunc("/garden/structures/", structureListHandler)
	http.HandleFunc("/garden/structures/detail/", structureDetailHandler)

	http.HandleFunc("/awakenings/", awakeningsTableHandler)
	http.HandleFunc("/awakenings/csv/", awakeningsCsvHandler)

	http.HandleFunc("/decode/", decodeHandler)

	http.HandleFunc("/raw/", rawDataHandler)
	http.HandleFunc("/raw/KEYS", rawDataKeysHandler)

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
	data := vc.VFile{}
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
<a href="/data">Set Data Location</a><br />
<br />
<a href="/cards/table">Card List as a Table</a><br />
<a href="/events">Event List</a><br />
<a href="/thor">Thor Event List</a><br />
<a href="/items">Item List</a><br />
<a href="/deckbonus">Deck Bonuses</a><br />
<a href="/maps">Map List</a><br />
<a href="/archwitches">Archwitch List</a><br />
<a href="/cards/levels">Card Levels</a><br />
<a href="/garden/structures">Garden Structures</a><br />
<a href="/characters">Character List as a Table</a><br />
<br />
Images:<br />
<a href="/images/card/?unused=1">Unused Card Images</a><br />
<a href="/images/battle/bg/">Battle Backgrounds</a><br />
<a href="/images/battle/map/">Battle Maps</a><br />
<a href="/images/event/">Event</a><br />
<a href="/images/garden/">Garden</a><br />
<a href="/images/garden/map">Garden Structures</a><br />
<a href="/images/alliance/">Alliance</a><br />
<a href="/images/dungeon/">Dungeon</a><br />
<a href="/images/summon/">Summon</a><br />
<a href="/images/item/">Items</a><br />
<a href="/images/treasure/">Sacred Relics</a><br />
<a href="/images/navi/">Navi</a><br />
<br />
<a href="/cards/csv">Card List as CSV</a><br />
<a href="/skills/csv">Skill List as CSV</a><br />
<a href="/cards/glrcsv">GLR Card List as CSV</a><br />
<br />
<a href="/awakenings">List of Awakenings</a><br />
<a href="/awakenings/csv">List of Awakenings as CSV</a><br />
<a href="/raw">Raw data</a><br />
<a href="/raw/KEYS">Raw data Keys</a><br />
<br />
<br />
<a href="/decode">Decode All Files</a><br />
<br />
<a href="/SHUTDOWN">SHUTDOWN</a><br />
</body></html>`,
		VcData.Version,
		VcData.Common.UnixTime.Unix(),
		VcData.Common.UnixTime.Format(time.RFC3339),
	)
	// io.WriteString(w, "<a href=\"/cards\">Card List</a><br />\n")
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

func rawDataKeysHandler(w http.ResponseWriter, r *http.Request) {
	c := make(map[string]interface{})
	err := json.Unmarshal([]byte(masterDataStr), &c)
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
