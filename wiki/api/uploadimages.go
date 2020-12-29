package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"vc_file_grouper/vc"
	"vc_file_grouper/wiki"
)

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
