package main

import (
	"bufio"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"os"
)

func usage() {
	os.Stderr.WriteString("You must pass the location of the files.\n" +
		"Usage: " + os.Args[0] + " /path/to/com.nubee.valkyriecrusade/files\n")
}

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		usage()
		return
	}

	http.HandleFunc("/", masterDataHandler)
	http.ListenAndServe(":8080", nil)
}

func masterDataHandler(w http.ResponseWriter, r *http.Request) {
	// 0. decode all the files (or maybe just decode on demand?)
	// 1. load in json data from masterdata.dat
	// 2. list main keys as navigation points

	// File header
	io.WriteString(w, "<html><body>\n")

	// Step 0.
	filename := os.Args[1] + "/response/master_all.dat"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// here we can check if the non converted file exists and convert it
		fmt.Fprintf(w, "no such file or directory: %s", filename)
		io.WriteString(w, "</body></html>")
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(w, "Error opening: %s", filename)
		io.WriteString(w, "</body></html>")
		return
	}

	masterData, _ := simplejson.NewFromReader(bufio.NewReader(f))

	masterDataMap, _ := masterData.Map()

	for k, _ := range masterDataMap {
		io.WriteString(w, k)
		io.WriteString(w, "<br />\n")
	}

	io.WriteString(w, "</body></html>")
}
