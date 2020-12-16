package handler

import (
	"fmt"
	"io"
	"net/http"
	"vc_file_grouper/wiki/api"
)

// WikibotHandler shows cards in order
func WikibotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Wikibot tasks</title></head><body>\n")
	io.WriteString(w, `
	<a href="/">Home</a>
<ul>
	<li><a href="/wikibot/testLogin">Test Your Login</a></li>
	<li><a href="/wikibot/startMassUpdate">Start a mass update.</a></li>
</ul>
`)
	io.WriteString(w, "</body></html>")
}

//LogoutHandler logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	api.Logout()
}

//TestLoginHandler tests the wiki login and reports back
func TestLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := api.Login()
	io.WriteString(w, "<html><head><title>Wikibot test login</title></head><body>\n")
	if err != nil {
		fmt.Fprintf(w, "<h1>%s</h1>", err.Error())
		io.WriteString(w, "<a href=\"/config/setBotCreds\">Update Login info</a>")
	} else {
		io.WriteString(w, "<h1>Success</h1>")
	}
	io.WriteString(w, "<br /><a href=\"/wikibot\">Wikibot home</a><br /><a href=\"/\">Home</a>")
	io.WriteString(w, "</body></html>")
}

//StartMassUpdateCardsHandler starts the mass update wizard.
func StartMassUpdateCardsHandler(w http.ResponseWriter, r *http.Request) {
}
