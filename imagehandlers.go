package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func imageCardHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/md"
	serveCardImage(vcfilepath+"/card/md/", "/images/card/", w, r)
}
func imageCardThumbHandler(w http.ResponseWriter, r *http.Request) {
	//(vcfilepath+"/card/thumb"
	serveCardImage(vcfilepath+"/card/thumb/", "/images/cardthumb/", w, r)
}
func imageCardHDHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/hd/"
	serveCardImage(vcfilepath+"/card/hd/", "/images/cardHD/", w, r)
}
func imageEventHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/event/largeimage"
	//imgname := r.URL.Path[len("/images/event/"):]
}

func serveCardImage(imagePath string, urlprefix string, w http.ResponseWriter, r *http.Request) {
	imgname := r.URL.Path[len(urlprefix):]
	if imgname == "" || imgname == "/" || strings.HasPrefix(imgname, "../") {
		// trying to read the entire directory or break out of the dir with a ../
		http.Error(w, "Invalid Image location", http.StatusForbidden)
		return
	}

	fullpath := imagePath + imgname

	var cardId, fileName string
	decodeOnFly := false
	if strings.HasSuffix(strings.ToLower(imgname), ".png") {
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			// png file is not cached
			if _, err := os.Stat(fullpath[len(fullpath)-4:]); os.IsNotExist(err) {
				// base image does not exist
				http.Error(w, "Invalid Image location "+fullpath, http.StatusNotFound)
				return
			}
			fullpath = fullpath[len(fullpath)-4:]
			decodeOnFly = true
		}
		cardId = imgname[3 : len(imgname)-4]
	} else {
		if _, err := os.Stat(fullpath + ".png"); os.IsNotExist(err) {
			// png file is not cached
			if _, err := os.Stat(fullpath); os.IsNotExist(err) {
				// base image does not exist
				http.Error(w, "Invalid Image location "+fullpath, http.StatusNotFound)
				return
			}
			decodeOnFly = true
		} else {
			fullpath = fullpath + ".png"
		}
		cardId = imgname[3:]
	}

	card := vc.CardScanImage(cardId, VcData.Cards)
	ext := ".png"
	if strings.Contains(fullpath, "/thumb/") {
		ext = "_icon.png"
	}
	if card != nil {
		fileName = card.Name
		if strings.HasPrefix(card.Rarity(), "G") {
			fileName = fileName + "_G" + ext
		} else if card.EvolutionRank == 0 {
			fileName = fileName + ext
		} else if card.EvolutionRank == card.LastEvolutionRank || card.EvolutionCardId <= 0 {
			fileName = fileName + "_H" + ext
		} else {
			fileName = fileName + "_" + strconv.Itoa(card.EvolutionRank) + ext
		}
	} else {
		os.Stderr.WriteString("Card infor not found for image " + cardId + "\n")
		if decodeOnFly {
			fileName = imgname + ext
		} else {
			fileName = imgname[:len(imgname)-4] + ext
		}
	}

	var b []byte
	var err error
	if decodeOnFly {
		// decode the file
		b, err = vc.Decode(fullpath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// read the entire png file to a byte array
		b, err = ioutil.ReadFile(fullpath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "image/png")

	var buff bytes.Buffer
	buff.Write(b)
	buff.WriteTo(w)

}

func nthPos(s string, pat string, n int) int {
	l := len(s)
	i := 0
	for ; n-1 > 0 && i+1 < l; n = n - 1 {
		i = strings.Index(s[i:], pat)
		i = i + 1
	}
	return i
}
