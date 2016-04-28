package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func imageCardSDHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/md"
	serveCardImage(vcfilepath+"/card/sd/", "/images/cardSD/", w, r)
}

func imageCardHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/md"
	serveCardImage(vcfilepath+"/card/md/", "/images/card/", w, r)
}

func imageCardThumbHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/thumb"
	serveCardImage(vcfilepath+"/card/thumb/", "/images/cardthumb/", w, r)
}

func imageCardHDHandler(w http.ResponseWriter, r *http.Request) {
	//vcfilepath+"/card/hd"
	serveCardImage(vcfilepath+"/card/hd/", "/images/cardHD/", w, r)
}

func imageHandlerFor(urlPath string, imageDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//vcfilepath+"/event"
		servImageDir(w, r, urlPath, imageDir)
	}
}

func servImageDir(w http.ResponseWriter, r *http.Request, urlPath string, root string) {
	imgname := r.URL.Path[len("/images"+urlPath):]
	for strings.HasPrefix(imgname, "/") {
		imgname = strings.TrimPrefix(imgname, "/")
	}
	for strings.HasPrefix(imgname, "../") {
		imgname = strings.TrimPrefix(imgname, "../")
	}
	if strings.Contains(imgname, "../") {
		http.Error(w, "Invalid Image location "+imgname+
			"<br />Relative path modification not allowed", http.StatusNotFound)
		return
	}
	fullpath := path.Join(vcfilepath, root, imgname)

	finfo, err := os.Stat(fullpath)
	if err != nil {
		http.Error(w, "Invalid Image location "+imgname+"<br />"+err.Error(), http.StatusNotFound)
		return
	}
	if finfo.Mode().IsRegular() {
		_, fName := filepath.Split(fullpath)
		writeout(true, fullpath, fName+".png", w, r)
		return
	} else if finfo.IsDir() {
		if !strings.HasSuffix(fullpath, "/") {
			fullpath = fullpath + "/"
		}
		io.WriteString(w, `<html>
<head>
	<link rel="stylesheet" type="text/css" href="/css/style.css">
</head>
<body class="stary-night">`)
		err := filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			b := make([]byte, 4)
			_, err = f.Read(b)
			f.Close()
			if err != nil {
				return err
			}
			if bytes.Equal(b, []byte("CODE")) {
				relPath := path[len(fullpath):]
				fmt.Fprintf(w, `<div><a href="%[1]s"><img src="%[1]s"/></a><br />%[1]s</div>`, relPath)
			}
			return nil
		})

		if err != nil {
			io.WriteString(w, err.Error()+"<br />\n")
		}
		io.WriteString(w, "</body></html>")
		return
	}
	http.Error(w, "Invalid Image location "+imgname, http.StatusNotFound)
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
			if _, err := os.Stat(fullpath[:len(fullpath)-4]); os.IsNotExist(err) {
				// base image does not exist
				http.Error(w, "Invalid Image location "+fullpath, http.StatusNotFound)
				return
			}
			fullpath = fullpath[:len(fullpath)-4]
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
		os.Stderr.WriteString("Card info not found for image " + cardId + "\n")
		if decodeOnFly {
			fileName = imgname + ext
		} else {
			fileName = imgname[:len(imgname)-4] + ext
		}
	}

	writeout(decodeOnFly, fullpath, fileName, w, r)

}

func writeout(decodeOnFly bool, fullpath string, fileName string, w http.ResponseWriter, r *http.Request) {
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
