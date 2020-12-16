package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
)

//URL root of the wiki
var URL string = "https://valkyriecrusade.fandom.com"

var client *http.Client = &http.Client{
	Jar: getCookieJar(),
}

func getCookieJar() *cookiejar.Jar {
	j, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: nil,
	})
	return j
}

//CardNameToWiki Converts a card name to a wiki name. The result is safe to use in URL paths
func CardNameToWiki(name string) string {
	return url.QueryEscape(strings.ReplaceAll(name, " ", "_"))
}

//GetLoginToken Gets a login token so a login can happen. Records the token in the Credentials
func GetLoginToken() (err error) {
	resp, err := client.Get(URL + "/api.php?action=query&meta=tokens&type=login&format=json")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(string(body))
	tr := tokenResponse{}
	err = json.Unmarshal(body, &tr)
	if err != nil {
		return
	}
	if tr.Warnings != nil {
		return errors.New("Unable to fetch token")
	}
	MyCreds.LoginToken = tr.Query.Tokens.LoginToken
	log.Printf("Login Token: %s", MyCreds.LoginToken)
	return
}

//Login uses the MyCreds to perform a login.
func Login() (err error) {
	if MyCreds.LoginToken == "" {
		err = GetLoginToken()
		if err != nil {
			return
		}
	}
	formVals := url.Values{}
	formVals.Add("lgname", MyCreds.Username)
	formVals.Add("lgpassword", MyCreds.Password)
	formVals.Add("lgtoken", MyCreds.LoginToken)
	resp, err := client.PostForm(URL+"/api.php?action=login&format=json", formVals)
	if err != nil {
		MyCreds.LoginToken = ""
		client.Jar = getCookieJar()
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		MyCreds.LoginToken = ""
		client.Jar = getCookieJar()
		return
	}
	log.Println(string(body))
	if resp.StatusCode != 200 {
		MyCreds.LoginToken = ""
		return fmt.Errorf("Invalid response. Expected HTTP 200, instead got %d", resp.StatusCode)
	}
	lr := loginResponse{}
	err = json.Unmarshal(body, &lr)
	if err != nil {
		MyCreds.LoginToken = ""
		client.Jar = getCookieJar()
		return
	}
	if strings.ToLower(lr.Login.Result) != "success" {
		MyCreds.LoginToken = ""
		client.Jar = getCookieJar()
		return errors.New(lr.Login.Reason)
	}
	return
}

//Logout Removes authentication tokens
func Logout() {
	client.Get(URL + "api.php?action=logout&token=" + url.QueryEscape(MyCreds.LoginToken))
	MyCreds.CSRFToken = ""
	MyCreds.LoginToken = ""
	client.Jar = getCookieJar()
}

//GetCardPage Gets a card page
func GetCardPage(card *vc.Card) (ret *wiki.CardPage, err error) {
	if card == nil || card.Name == "" {
		return nil, nil
	}
	resp, err := client.Get(URL + "/index.php?action=raw&title=" + CardNameToWiki(card.Name))

	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	*ret, err = wiki.ParseCardPage(string(body))
	return
}
