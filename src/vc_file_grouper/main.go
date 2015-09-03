package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", fship)
	http.ListenAndServe(":8080", nil)
}

func fship(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><")
	fmt.Fprintf(w, "</body></html>")

}
