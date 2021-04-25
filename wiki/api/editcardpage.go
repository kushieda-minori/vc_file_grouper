package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"vc_file_grouper/wiki"
)

//EditCardPage Edits a card page
func EditCardPage(cp *wiki.CardPage, editSummary string) (err error) {
	// verify basic page information is available
	if cp == nil {
		err = errors.New("Page is nil")
		return
	}
	if cp.PageName == "" {
		err = errors.New("Page name can not be blank")
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
	pageName, _ := url.QueryUnescape(cp.PageName)
	formVals := url.Values{}
	formVals.Add("token", MyCreds.CSRFToken)
	formVals.Add("bot", "true")
	formVals.Add("nocreate", "true")
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
