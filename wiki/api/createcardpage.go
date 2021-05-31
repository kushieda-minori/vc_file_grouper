package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
)

//CreateCardPage Creates a new card page
func CreateCardPage(c *vc.Card, editSummary string) (err error) {
	// verify basic page information is available
	if c == nil {
		err = errors.New("page is nil")
		return
	}
	if c.Name == "" {
		err = errors.New("page name can not be blank")
		return
	}

	// verify we are logged into the wiki API
	if MyCreds.LoginToken == "" {
		err = Login()
		if err != nil {
			return
		}
	}

	//
	cp := wiki.CardPage{
		PageName: CardNameToWiki(c.Name),
	}
	if c.IsClosed != 0 {
		cp.PageHeader = "{{Unreleased}}"
	}
	cp.CardInfo.UpdateAll(c, "")

	alamgamations := wiki.GetAmalgamations(c)

	if len(alamgamations) > 0 {
		cp.PageFooter = "==''[[Amalgamation]]''==\n\n" + alamgamations.String()
	}

	pageName, _ := url.QueryUnescape(cp.PageName)
	formVals := url.Values{}
	formVals.Add("token", MyCreds.CSRFToken)
	formVals.Add("bot", "true")
	formVals.Add("createonly", "true")
	formVals.Add("recreate", "true")
	formVals.Add("title", pageName)
	formVals.Add("summary", editSummary)
	formVals.Add("text", cp.String())
	resp, err := client.PostForm(URL+"/api.php?action=edit&format=json", formVals)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	er := editResponse{}
	err = json.Unmarshal(body, &er)
	if err != nil {
		return
	}
	if er.Error != nil {
		data, err := er.Error.MarshalJSON()
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}
	return nil
}
