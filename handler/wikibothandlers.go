package handler

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki/api"
)

// WikibotHandler shows cards in order
func WikibotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Wikibot tasks</title></head><body>\n")
	io.WriteString(w, `
	<a href="/">Home</a>
<ul>
	<li><a href="/wikibot/testLogin">Test Your Login</a></li>
	<li><a href="/wikibot/testCardFetch">Test Fetch Oracle (R).</a></li>
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

//TestCardFetchHandler tests fetching a single card page
func TestCardFetchHandler(w http.ResponseWriter, r *http.Request) {
	/*
		 313 Oracle R - HR
		2582 Sulis SR GSR with Amal
		1879 Oracle Ascendant UR - GUR
		3493 Summer Oracle UR - XUR
		8408 New Year Alchemist LR - XLR
		9490 Cheerleader Pixie VR - GVR
	*/
	card := vc.CardScan(9490)
	cardPage, rawPagebody, err := api.GetCardPage(card)
	fmt.Fprintf(w, "<html><head><title>Wikibot Test page updates</title><style type=\"text/css\">%s</style></head><body>\n",
		`
div.flex {
	display:flex;
	max-width:100%;
	overflow:auto;
}
div.flex > div {
	margin: 2px;
	min-width: 575px;
	max-width: 49%;
	overflow: auto;
}
pre {
	padding:5px;
	border:solid black 1px;
	width: 98%;
	overflow: auto;
}
`,
	)
	if err != nil {
		fmt.Fprintf(w, "<h1>%s</h1>", err.Error())
	} else {
		fmt.Fprintf(w, "<h1>%s</h1>", card.Name)
		io.WriteString(w, `<div class="flex">`)
		fmt.Fprintf(w, "<div>Wiki Version<pre>%s</pre></div>", html.EscapeString(rawPagebody))
		cardPage.CardInfo.UpdateBaseData(*card)
		cardPage.CardInfo.UpdateSkills(card.GetEvolutions())
		cardPage.CardInfo.UpdateExchangeInfo(card.GetEvolutions())
		cardPage.CardInfo.UpdateAwakenRebirthInfo(card.GetEvolutions())
		cardPage.CardInfo.UpdateQuotes(card)
		fmt.Fprintf(w, "<div>Adjusted Version<pre>%s</pre></div>", html.EscapeString(cardPage.String()))
		io.WriteString(w, `</div>`)
	}
	io.WriteString(w, "<br /><a href=\"/wikibot\">Wikibot home</a><br /><a href=\"/\">Home</a>")
	io.WriteString(w, "</body></html>")
}

//StartMassUpdateCardsHandler starts the mass update wizard.
func StartMassUpdateCardsHandler(w http.ResponseWriter, r *http.Request) {
}
