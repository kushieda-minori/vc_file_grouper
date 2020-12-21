package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
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
		err = errors.New("Unable to fetch token")
	}
	if tokenType == "login" {
		token = tr.Query.Tokens.LoginToken
	} else if tokenType == "csrf" {
		token = tr.Query.Tokens.CSRFToken
	}
	log.Printf("Token: %s", token)
	return
}

//Login uses the MyCreds to perform a login.
func Login() (err error) {
	if MyCreds.Username == "" || MyCreds.Password == "" {
		return errors.New("User information not setup")
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
		return fmt.Errorf("Invalid response. Expected HTTP 200, instead got %d", resp.StatusCode)
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

//Logout Removes authentication tokens
func Logout() {
	client.Get(URL + "api.php?action=logout&token=" + url.QueryEscape(MyCreds.LoginToken))
	MyCreds.CSRFToken = ""
	MyCreds.LoginToken = ""
	client.Jar = getCookieJar()
}

//GetCardPage Gets a card page
func GetCardPage(card *vc.Card) (ret *wiki.CardPage, raw string, err error) {
	if card == nil || card.Name == "" {
		return
	}
	resp, err := client.Get(URL + "/index.php?action=raw&title=" + CardNameToWiki(card.Name))

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

	ret = &wiki.CardPage{}
	raw = string(body)
	*ret, err = wiki.ParseCardPage(raw)
	return
}

//UploadNewCardUniqueImages uploads images that don't yet exist
func UploadNewCardUniqueImages(card *vc.Card) (ret *wiki.CardPage, err error) {
	if card == nil || card.Name == "" {
		return
	}

	if MyCreds.LoginToken == "" {
		err = Login()
		if err != nil {
			return
		}
	}

	err = uploadImages(card, false)
	if err != nil {
		return
	}
	err = uploadImages(card, true)

	return
}

func uploadImages(card *vc.Card, thumbs bool) (err error) {
	evos := card.GetEvolutions()
	for _, evoID := range card.EvosWithDistinctImages(thumbs) {
		evo := evos[evoID]
		var name string
		var data []byte
		name, data, err = evo.GetImageData(true)
		if err != nil {
			return
		}

		var contentType string
		data, contentType, err = createMultipartForm(map[string]io.Reader{
			"filename": strings.NewReader(name),
			"token":    strings.NewReader(MyCreds.CSRFToken),
			"file":     bytes.NewReader(data),
		})
		if err != nil {
			return
		}

		// query := fmt.Sprintf("/api.php?action=upload&format=json&filename=%s&token=%s",
		// 	url.QueryEscape(name),
		// 	url.QueryEscape(MyCreds.CSRFToken),
		// )

		var resp *http.Response
		resp, err = client.Post(URL+"/api.php?action=upload&format=json", contentType, bytes.NewReader(data))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		log.Println(string(body))
	}
	return
}

func createMultipartForm(values map[string]io.Reader) (form []byte, contentType string, err error) {
	var formData bytes.Buffer
	w := multipart.NewWriter(&formData)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add a file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return
		}
	}
	return formData.Bytes(), w.FormDataContentType(), nil
}
