package handler

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// ConfigHandler configures the path for the main VC data file
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Master Data</title></head><body>\n")

	// check form value and update if valid
	newpath := r.FormValue("path")
	if newpath != "" {
		if _, err := os.Stat(newpath); os.IsNotExist(err) {
			io.WriteString(w, "<div>Invalid new path specified</div>")
		} else {
			if err = vc.ReadMasterData(newpath); err != nil {
				fmt.Fprintf(w, "<div>%s</div>", err.Error())
			} else {
				io.WriteString(w, "<div>Success</div>")
				vc.FilePath = newpath
			}
		}
	}
	// write out the form
	fmt.Fprintf(w, `<form method="post">
<label for="f_path">Data Path</label>
<input id="f_path" name="path" value="%s" style="width:300px"/>
<button type="submit">Submit</button>
<p><a href="/">back</a></p>
</form>`,
		html.EscapeString(vc.FilePath),
	)
	io.WriteString(w, "</body></html>")
}
