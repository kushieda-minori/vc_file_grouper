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
	fmt.Fprintf(w, "<html><body><img src='http://davidhehenberger.com/wp-content/uploads/2013/10/fuck-it-ship-it1.jpg'/><br/>Fuck it, ship it!!</body></html>")
}
