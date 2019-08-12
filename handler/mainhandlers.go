package handler

import (
	"fmt"
	"net/http"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// MasterDataHandler Main index page
func MasterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	fmt.Fprintf(w, `<html><body>
<p>Version: %d,&nbsp;&nbsp;&nbsp;&nbsp;Timestamp: %d,&nbsp;&nbsp;&nbsp;&nbsp;JST: %s</p>
<a href="/config">Configure Data Location</a><br />
<br />
<a href="/cards/table">Card List as a Table</a><br />
<a href="/weapons">Weapon List</a><br />
<a href="/events">Event List</a><br />
<a href="/events/dungeonScenario/">Dungeon Scenario</a><br />
<a href="/items">Item List</a><br />
<a href="/deckbonus">Deck Bonuses</a><br />
<a href="/maps">Map List</a><br />
<a href="/archwitches">Archwitch List</a><br />
<a href="/cards/levels">Card Levels</a><br />
<a href="/garden/structures">Garden Structures</a><br />
<a href="/characters">Character List as a Table</a><br />
<a href="/thor">Thor Event List</a><br />
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
<a href="/images/weapon/">All Weapon Images</a><br />
<a href="/images/weaponevent/">Weapon Event Images</a><br />
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
		vc.Data.Version,
		vc.Data.Common.UnixTime.Unix(),
		vc.Data.Common.UnixTime.Format(time.RFC3339),
	)
	// io.WriteString(w, "<a href=\"/cards\">Card List</a><br />\n")
}
