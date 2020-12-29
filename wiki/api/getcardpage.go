package api

import (
	"fmt"
	"io/ioutil"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
)

//GetCardPage Gets a card page
func GetCardPage(card *vc.Card) (ret *wiki.CardPage, raw string, err error) {
	if card == nil || card.Name == "" {
		return
	}

	pageName := CardNameToWiki(card.Name)
	resp, err := client.Get(URL + "/index.php?action=raw&title=" + pageName)

	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("Invalid HTTP Status returned - %d: %s", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ret = &wiki.CardPage{
		PageName: pageName,
	}
	raw = string(body)
	err = ret.Parse(raw)
	return
}
