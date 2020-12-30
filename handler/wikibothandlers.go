package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
	"vc_file_grouper/wiki/api"
)

// WikibotHandler shows cards in order
func WikibotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><head><title>Wikibot tasks</title></head><body>\n")
	io.WriteString(w, `
	<a href="/">Home</a>
<ul>
	<li><a href="/wikibot/testLogin">Test Your Login</a></li>
	<li><a href="/wikibot/testCardFetch">Test Fetch and compare.</a></li>
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
		1879 Oracle Ascendant UR - GUR
		2582 Sulis SR GSR with Amal
		3493 Summer Oracle UR - XUR
		3934 Dark Succubus - XUR Random Skill
		8408 New Year Alchemist LR - XLR
		9490 Cheerleader Pixie VR - GVR
		9517 Christmas Lum Lum - XSR - ABB (skill expire)
	*/
	card := vc.CardScan(3934)
	writeCardReviewForm(w, card, 1, 1, "", "checked", "", nil)
}

var botCardList vc.CardList = nil
var botLogRoot = ""

//StartMassUpdateCardsHandler starts the mass update wizard.
func StartMassUpdateCardsHandler(w http.ResponseWriter, r *http.Request) {
	if botCardList == nil {
		// initialilze the list
		tmp := vc.CardsByNameByLowestID(true)
		botCardList = make(vc.CardList, 0)
		for _, cl := range tmp {
			botCardList = append(botCardList, cl.Earliest())
		}
		botCardList = botCardList.Filter(func(c vc.Card) bool {
			return c.CardCharaID > 0 && c.IsClosed == 0 && c.Name != "" && c.SkillID1 > 0
		})
	}

	var card *vc.Card
	currentID := 0
	lenCardList := len(botCardList)
	if pos := r.FormValue("pos"); pos != "" {
		posID, err := strconv.Atoi(pos)
		if err != nil {
			io.WriteString(w, "Requested position is invalid: "+err.Error())
			return
		}
		if posID < 0 || posID >= lenCardList {
			io.WriteString(w, "Requested position is invalid: "+pos)
			return
		}
		currentID = posID
		card = botCardList[currentID]
	} else {
		card = botCardList[0]
	}

	var err error
	if r.Method == "POST" {
		if currentID == 0 || botLogRoot == "" {
			botLogRoot = "changelog/"
			//Mon Jan 2 15:04:05 -0700 MST 2006
			botLogRoot += time.Now().Format("20060102-150405")
			os.MkdirAll(botLogRoot, 0700)
		}
		origPage := r.FormValue("orig")
		if origPage == "" {
			io.WriteString(w, "Original page not provided")
			return
		}
		fixedPage := r.FormValue("data")
		if fixedPage == "" {
			io.WriteString(w, "Can not update with a blank page using this BOT")
			return
		}
		origCardPage := wiki.CardPage{
			PageName: api.CardNameToWiki(card.Name),
		}
		err = origCardPage.Parse(origPage)
		if err != nil {
			io.WriteString(w, "Error parsing the original page for comparisson: "+err.Error())
			return
		}
		newCardPage := &wiki.CardPage{
			PageName: api.CardNameToWiki(card.Name),
		}
		err = newCardPage.Parse(fixedPage)
		if err != nil {
			io.WriteString(w, "Error parsing the updated page for comparisson: "+err.Error())
			return
		}

		diff := origCardPage.CardInfo.Differences(newCardPage.CardInfo)
		if len(diff) == 0 {
			log.Printf("Card %s had no updates, so nothing will be saved to the wiki", card.Name)
		} else {
			log.Printf("*****Card %s has updates, will be saved to the wiki", card.Name)
			json, _ := json.MarshalIndent(diff, "", "\t")
			fName := fmt.Sprintf("%s/%06d.%s.diff.json", botLogRoot, currentID, card.Name)
			err = ioutil.WriteFile(fName, json, 0700)
			// only save pages that actually have changes to page content.
			//
			if r.FormValue("dryrun") != "checked" {
				err = api.EditCardPage(newCardPage, r.FormValue("summary"))
			}
		}

		if err == nil {
			if currentID+1 == lenCardList {
				// if we are at the last item, go back to the bot menu
				http.Redirect(w, r, `/wikibot/`, http.StatusSeeOther)
			} else {
				// bring up the next card
				http.Redirect(w, r,
					fmt.Sprintf(`?pos=%d&auto=%s&dryrun=%s&summary=%s`, currentID+1, r.FormValue("auto"), r.FormValue("dryrun"), url.QueryEscape(r.FormValue("summary"))),
					http.StatusSeeOther,
				)
			}
			return
		}
	}

	writeCardReviewForm(w, card, currentID, lenCardList, r.FormValue("auto"), r.FormValue("dryrun"), r.FormValue("summary"), err)
}

func writeCardReviewForm(w io.Writer, card *vc.Card, currentID, listLength int, isAuto string, isDryRun string, summary string, err error) {
	fmt.Fprintf(w, `<html>
	<head>
		<title>Wikibot updates %d of %d</title>
		<style type="text/css">
			div.flex {
				display:flex;
				max-width:100%%;
				overflow:auto;
			}
			div.flex > div {
				margin: 2px;
				min-width: 575px;
				max-width: 49%%;
				overflow: auto;
			}
			pre,textarea {
				padding:5px;
				border:solid black 1px;
				width: 100%%;
				height: 600px;
				overflow: auto;
			}
			textarea {
				white-space: pre;
				overflow-wrap: normal;
				overflow-x: scroll;
			}
			div.nav span, div.nav div {
				margin-left: 15px;
			}
		</style>
		<script type="text/javascript">
		var vc = vc || {};
		vc.pageTimer = null;
		vc.onPageLoad = function() {
			var autoF = document.getElementById("f_auto");
			if (autoF && autoF.checked) {
				vc.pageTimer = setTimeout(vc.submit, 3 * 1000)
			}
		}
		vc.submit = function() {
			// disable the timer in case the user pressed submit so we don't double submit.
			if (vc.pageTimer) { clearTimeout(vc.pageTimer); }
			var f = document.getElementById("cardChanges");
			f.submit();
			return false;
		}
		</script>
	</head>
	<body onload="vc.onPageLoad()">
`,
		currentID+1,
		listLength,
	)
	if err != nil {
		fmt.Fprintf(w, "<h1>%s</h1>", err.Error())
		io.WriteString(w, `<br /><a href="/wikibot">Wikibot home</a><br /><a href="/">Home</a>
		</body>
	</html>
	`)
		return
	}
	cardPage, rawPagebody, err := api.GetCardPage(card)
	if err != nil {
		fmt.Fprintf(w, "<h1>%s: %s</h1>", card.Name, err.Error())
	} else {
		fmt.Fprintf(w, "<h1>%s</h1>\n", card.Name)
		io.WriteString(w, `<form id="cardChanges" action="./" name="cardChanges" method="post" onsubmit="return vc.submit();">`)
		fmt.Fprintf(w, `<input type="hidden" name="pos" value="%d" />`, currentID)
		fmt.Fprintf(w, `<div><label for="f_summary">Bot Edit Summary:<input id="f_summary" name="summary" type="text" value="%s"/></label></div>`, html.EscapeString(summary))
		io.WriteString(w, `<div class="nav"><span><a href="/wikibot">Cancel</a></span>`)
		if currentID < listLength {
			fmt.Fprintf(w, `<span><a href="?pos=%d">Skip with no update</a></span>`, currentID+1)
			io.WriteString(w, `<span><button name="s" type="submit">Submit and move Next</button></span>`)
			fmt.Fprintf(w, `<span><label for="f_auto">Auto Advance:<input id="f_auto" name="auto" type="checkbox" value="checked" %s onchange="clearTimeout(vc.pageTimer)"/></label><small>(every 3 seconds)</small></span>`, isAuto)
			fmt.Fprintf(w, `<span><label for="f_dryrun">Dry Run:<input id="f_dryrun" name="dryrun" type="checkbox" value="checked" %s /></label><small>(no actual edits)</small></span>`, isDryRun)
		} else {
			io.WriteString(w, `<button name="submit" type="submit">Submit and End</button>`)
		}
		io.WriteString(w, `</div>`)
		io.WriteString(w, `<div class="flex">`)
		fmt.Fprintf(w, `<div>Wiki Version<textarea readonly="readonly" name="orig">%s</textarea></div>`, html.EscapeString(rawPagebody))
		cardPage.CardInfo.UpdateBaseData(card)
		cardPage.CardInfo.UpdateSkills(card.GetEvolutions())
		cardPage.CardInfo.UpdateExchangeInfo(card.GetEvolutions())
		//cardPage.CardInfo.UpdateEvoStats(card.GetEvolutions())
		cardPage.CardInfo.UpdateAwakenRebirthInfo(card.GetEvolutions())
		cardPage.CardInfo.UpdateQuotes(card)
		fmt.Fprintf(w, `<div>Adjusted Version<textarea name="data">%s</textarea></div>`, html.EscapeString(cardPage.String()))
		io.WriteString(w, `</div>`)
		io.WriteString(w, `</form>`)
	}
	io.WriteString(w, `<br /><a href="/wikibot">Wikibot home</a><br /><a href="/">Home</a>
	</body>
</html>
`)
}
