package main

import (
	"io"
	"net/http"
	"os"
	"zetsuboushita.net/vc_file_grouper/vc_grouper"
)

var VcData vc_grouper.VcFile

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

	err := VcData.Read(os.Args[1])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	http.HandleFunc("/", masterDataHandler)
	http.ListenAndServe(":8080", nil)
}

func masterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	io.WriteString(w, "<html><body>\n")

	io.WriteString(w, "</body></html>")
}
