package handler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"../vc"
)

// ImageCardSDHandler show SD card images
func ImageCardSDHandler(w http.ResponseWriter, r *http.Request) {
	//vc.FilePath+"/card/sd"
	serveCardImage(filepath.Join(vc.FilePath, "card", "sd"), "/images/cardSD/", w, r)
}

// ImageCardHandler show MD card images
func ImageCardHandler(w http.ResponseWriter, r *http.Request) {
	//vc.FilePath+"/card/md"
	serveCardImage(filepath.Join(vc.FilePath, "card", "md"), "/images/card/", w, r)
}

// ImageCardThumbHandler show thumbnail card images
func ImageCardThumbHandler(w http.ResponseWriter, r *http.Request) {
	//vc.FilePath+"/card/thumb"
	serveCardImage(filepath.Join(vc.FilePath, "card", "thumb"), "/images/cardthumb/", w, r)
}

// ImageCardHDHandler show HD card images
func ImageCardHDHandler(w http.ResponseWriter, r *http.Request) {
	//vc.FilePath+"/card/hd"
	serveCardImage(filepath.Join(vc.FilePath, "card", "hd"), "/images/cardHD/", w, r)
}

// ImageHandlerFor handles images under a specified path
func ImageHandlerFor(urlPath string, imageDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//vc.FilePath+"/event"
		servImageDir(w, r, urlPath, imageDir)
	}
}

type fileFilterFunc func(os.FileInfo) bool

func servImageDir(w http.ResponseWriter, r *http.Request, urlPath string, root string, filters ...fileFilterFunc) {
	imgname := filepath.FromSlash(r.URL.Path)
	urlRoot := path.Join("/images", urlPath)
	urlRoot = filepath.FromSlash(urlRoot)
	imgname, _ = filepath.Rel(urlRoot, imgname)
	imgname = filepath.ToSlash(imgname)
	queryValues := r.URL.Query()

	var forceFileName string
	_, qok := queryValues["filename"]
	if qok {
		forceFileName = queryValues["filename"][0]
	}

	for strings.HasPrefix(imgname, "/") {
		imgname = strings.TrimPrefix(imgname, "/")
	}
	for strings.HasPrefix(imgname, "../") {
		imgname = strings.TrimPrefix(imgname, "../")
	}
	if strings.Contains(imgname, "../") {
		http.Error(w, "Invalid Image location "+imgname+
			"\nRelative path modification not allowed", http.StatusNotFound)
		return
	}
	fullpath := path.Join(vc.FilePath, root, imgname)

	finfo, err := os.Stat(fullpath)
	if err != nil {
		http.Error(w, "Invalid Image location "+imgname+"\n"+err.Error(), http.StatusNotFound)
		return
	}
	if finfo.Mode().IsRegular() {
		var fName string
		if forceFileName == "" {
			_, fName = filepath.Split(fullpath)
		} else {
			fName = forceFileName
		}
		if !strings.HasSuffix(strings.ToLower(fName), ".png") {
			fName += ".png"
		}
		log.Printf("sending image {%s} as {%s}", fullpath, fName)
		writeout(true, fullpath, fName, w, r)
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
		err := filepath.Walk(fullpath, func(p string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filters != nil {
				for _, filter := range filters {
					if !filter(info) {
						return nil
					}
				}
			}
			fileEncoded, err := vc.IsFileEncoded(p)
			if err != nil {
				return err
			}
			if fileEncoded {
				relPath, _ := filepath.Rel(fullpath, p)
				relPath = filepath.ToSlash(relPath)
				fmt.Fprintf(w, `<div class="image"><a href="%[1]s"><img src="%[1]s"/></a><br>`, relPath)
				if "/card/" == urlPath {
					fmt.Fprintf(w, `<a href="../cardthumb/%[1]s"><img src="../cardthumb/%[1]s" /></a><br />`, relPath)
				}
				fmt.Fprintf(w, `%[1]s</div>`, relPath)
			} else {
				log.Printf("Image is not encoded: %s", fullpath)
			}
			return nil
		})

		if err != nil {
			io.WriteString(w, err.Error()+"<br />\n")
		}
		io.WriteString(w, "</body></html>")
		return
	} else {
		log.Printf("Unknown file mode: %v", finfo.Mode())
	}
	http.Error(w, "Invalid Image location "+imgname, http.StatusNotFound)
}

func checkImageName(info os.FileInfo) bool {
	imageName := info.Name()
	for _, card := range vc.Data.Cards {
		if card.Image() == imageName {
			return false
		}
	}
	return true
}

func serveCardImage(imagePath string, urlprefix string, w http.ResponseWriter, r *http.Request) {
	imgname := r.URL.Path[len(urlprefix):]
	qs := r.URL.Query()
	if imgname == "" || imgname == "/" || strings.HasPrefix(imgname, "../") {
		if len(qs) > 0 {
			if unused := qs.Get("unused"); unused != "" {
				relPath, _ := filepath.Rel(vc.FilePath, imagePath)
				servImageDir(w, r, strings.TrimPrefix(urlprefix, "/images"), relPath, checkImageName)
				return
			}
		}
		// trying to read the entire directory or break out of the dir with a ../
		http.Error(w, "Invalid Image location", http.StatusForbidden)
		return
	}

	fullpath := filepath.Join(imagePath, imgname)

	var cardID, fileName string
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
		cardID = imgname[3 : len(imgname)-4]
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
		cardID = imgname[3:]
	}

	card := vc.CardScanImage(cardID)
	ext := ".png"
	isIcon := false
	if strings.Contains(fullpath, filepath.FromSlash("/thumb/")) {
		ext = "_icon.png"
		isIcon = true
	}
	if card != nil {
		fileName = card.GetEvoImageName(isIcon) + ext
	} else {
		//log.Printf("Card info not found for image " + cardID + "\n")
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
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", "image/png")

	var buff bytes.Buffer
	buff.Write(b)
	buff.WriteTo(w)
}
