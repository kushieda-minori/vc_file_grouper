package handler

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
)

// ConfigDataLocHandler configures the path for the main VC data file
func ConfigDataLocHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Master Data</title></head><body>\n")

	// check form value and update if valid
	newpath := r.FormValue("path")
	if newpath != "" {
		newpath = filepath.Clean(newpath)
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

//ConfigBotCredsHandler configure the bot Username and Password
func ConfigBotCredsHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Update Bot User Info</title></head><body>\n")

	// check form value and update if valid
	username := r.FormValue("username")
	if username != "" {
		wiki.MyCreds.Username = username
		password := r.FormValue("password")
		if password != "" {
			wiki.MyCreds.Password = password
			io.WriteString(w, "<div>Success</div>")
		} else {
			io.WriteString(w, "<div>Password can not be blank</div>")
		}
	}

	// write out the form
	fmt.Fprintf(w, `<form method="post">
	<div>
		Did you remember to create a <a href="https://valkyriecrusade.fandom.com/wiki/Special:BotPasswords" target="_blank">special bot credential</a>?
	</div>
	<label for="f_username">Username</label>
	<input id="f_username" name="username" value="%s" style="width:300px"/><br/>
	<label for="f_password">Password</label>
	<input id="f_password" type="password" name="password" value="" style="width:300px"/><br/>
<button type="submit">Submit</button>
<p><a href="/">back</a></p>
</form>`,
		html.EscapeString(wiki.MyCreds.Username),
	)
	io.WriteString(w, "</body></html>")
}
