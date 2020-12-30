package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"vc_file_grouper/vc"
)

//UploadNewCardUniqueImages uploads images that don't yet exist
func UploadNewCardUniqueImages(card *vc.Card) (err error) {
	if card == nil || card.Name == "" {
		return
	}

	if MyCreds.LoginToken == "" {
		err = Login()
		if err != nil {
			return
		}
	}

	// upload full card images
	err = uploadImages(card, false)
	if err != nil {
		return
	}

	// upload thumbnails
	err = uploadImages(card, true)

	return
}

func uploadImages(card *vc.Card, thumbs bool) (err error) {
	evos := card.GetEvolutions()
	for _, evoID := range card.EvosWithDistinctImages(thumbs) {
		evo := evos[evoID]
		var name string
		var data []byte
		name, data, err = evo.GetImageData(thumbs)
		if err != nil {
			return
		}

		var formData bytes.Buffer
		w := multipart.NewWriter(&formData)
		err = w.WriteField("filename", name)
		if err != nil {
			return
		}
		err = w.WriteField("token", MyCreds.CSRFToken)
		if err != nil {
			return
		}
		err = createMultiPartFormFile(w, data, "file", name)
		if err != nil {
			return
		}
		err = w.Close()
		if err != nil {
			return
		}

		contentType := w.FormDataContentType()

		//log.Println(string(formData))

		// query := fmt.Sprintf("/api.php?action=upload&format=json&filename=%s&token=%s",
		// 	url.QueryEscape(name),
		// 	url.QueryEscape(MyCreds.CSRFToken),
		// )

		var resp *http.Response
		resp, err = client.Post(URL+"/api.php?action=upload&format=json&ignorewarnings=true", contentType, bytes.NewReader(formData.Bytes()))
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

func createMultiPartFormFile(w *multipart.Writer, data []byte, key, fileName string) (err error) {
	var fw io.Writer
	if fw, err = w.CreateFormFile(key, fileName); err != nil {
		return // error
	}
	_, err = fw.Write(data)
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
				return // error
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return // error
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return // error
		}
	}
	return formData.Bytes(), w.FormDataContentType(), nil
}
