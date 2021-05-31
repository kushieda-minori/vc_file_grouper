package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

//Login uses the MyCreds to perform a login.
func Login() (err error) {
	if MyCreds.Username == "" || MyCreds.Password == "" {
		return errors.New("user information not setup")
	}
	if MyCreds.LoginToken == "" {
		MyCreds.LoginToken, err = getToken("login")
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
		MyCreds.CSRFToken = ""
		client.Jar = getCookieJar()
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		MyCreds.LoginToken = ""
		MyCreds.CSRFToken = ""
		client.Jar = getCookieJar()
		return
	}
	log.Println(string(body))
	if resp.StatusCode != 200 {
		MyCreds.LoginToken = ""
		MyCreds.CSRFToken = ""
		return fmt.Errorf("invalid response. Expected HTTP 200, instead got %d", resp.StatusCode)
	}
	lr := loginResponse{}
	err = json.Unmarshal(body, &lr)
	if err != nil {
		MyCreds.LoginToken = ""
		MyCreds.CSRFToken = ""
		client.Jar = getCookieJar()
		return
	}
	if strings.ToLower(lr.Login.Result) != "success" {
		MyCreds.LoginToken = ""
		MyCreds.CSRFToken = ""
		client.Jar = getCookieJar()
		return errors.New(lr.Login.Reason)
	}

	MyCreds.CSRFToken, err = getToken("csrf")
	if err != nil {
		MyCreds.LoginToken = ""
		MyCreds.CSRFToken = ""
		client.Jar = getCookieJar()
		return
	}
	return
}

//getToken Gets a token so a login can happen. Records the token in the Credentials
func getToken(tokenType string) (token string, err error) {
	resp, err := client.Get(URL + "/api.php?action=query&meta=tokens&type=" + tokenType + "&format=json")
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
		err = errors.New("unable to fetch token")
	}
	if tokenType == "login" {
		token = tr.Query.Tokens.LoginToken
	} else if tokenType == "csrf" {
		token = tr.Query.Tokens.CSRFToken
	}
	log.Printf("Token: %s", token)
	return
}
