package handler

import (
	"archive/zip"
	"fmt"
	"html"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"vc_file_grouper/util"
	"vc_file_grouper/vc"
)

// MasterDataHandler Main index page
func MasterDataHandler(w http.ResponseWriter, r *http.Request) {

	// File header
	fmt.Fprintf(w, `<html><body>
<p>Version: %d,&nbsp;&nbsp;&nbsp;&nbsp;Timestamp: %d,&nbsp;&nbsp;&nbsp;&nbsp;JST: %s</p>
<a href="/config/dataLoc">Configure Data Location</a><br />
<a href="/config/setBotCreds">Set bot username and password</a>.
If not set, you won't be able to automate updates to the wiki.
Please create a "bot key" for your account by using the <a href="https://valkyriecrusade.fandom.com/wiki/Special:BotPasswords" target="_blank">Special:BotPasswords</a> page<br />
<br />
<a href="/wikibot">Use the WikiBot</a><br />
<br />
<a href="/cards/table">Card List as a Table</a><br />
<a href="/weapons">Weapon List</a><br />
<a href="/events">Event List</a><br />
<a href="/events/towerScenario/">Tower Scenarios</a><br />
<a href="/events/dungeonScenario/">DRV Scenarios</a><br />
<a href="/events/weaponScenario/">Weapon Scenario</a><br />
<a href="/items">Item List</a><br />
<a href="/deckbonus">Deck Bonuses</a><br />
<a href="/maps">Map List</a><br />
<a href="/archwitches">Archwitch List</a><br />
<a href="/cards/levels">Card Levels</a><br />
<a href="/garden/structures">Garden Structures</a><br />
<a href="/characters">Character List as a Table</a><br />
<a href="/thor">Thor Event List</a><br />
<br />
Formatted Data:<br />
<a href="/cards/csv">Card List as CSV</a><br />
<a href="/skills/csv">Skill List as CSV</a><br />
<a href="/cards/glrcsv">GLR Card List as CSV</a> <a href="/cards/glrjson"> as JSON</a><br />
<br />
<a href="/strb/">Binary String files</a><br />
<br />
Images:<br />
<a href="/images/card/?unused=1">Unused Card Images</a><br />
<a href="/images/battle/bg/">Battle Backgrounds</a><br />
<a href="/images/battle/map/">Battle Maps</a><br />
<a href="/images/event/">Event</a><br />
<a href="/images/garden/">Garden</a><br />
<a href="/images/garden/map">Garden Structures</a><br />
<a href="/images/alliance/">Alliance</a><br />
<a href="/images/dungeon/">Dungeon</a><br />
<a href="/images/summon/">Summon</a><br />
<a href="/images/item/">Items</a><br />
<a href="/images/treasure/">Sacred Relics</a><br />
<a href="/images/navi/">Navi</a><br />
<a href="/images/weapon/">All Weapon Images</a><br />
<a href="/images/weaponevent/">Weapon Event Images</a><br />
<br />
<a href="/awakenings">List of Awakenings</a><br />
<a href="/awakenings/csv">List of Awakenings as CSV</a><br />
<a href="/raw">Raw data</a><br />
<a href="/raw/KEYS">Raw data Keys</a><br />
<br />
<a href="/decode">Decode All Files</a><br />
<br />
<a href="/zipData">Decode All Files and store them in a Zip archive</a><br />
<br />
<a href="/downloadMaps">Download Maps</a><br />
<br />
<a href="/SHUTDOWN">SHUTDOWN</a><br />
</body></html>`,
		vc.Data.Version,
		vc.Data.Common.UnixTime.Unix(),
		vc.Data.Common.UnixTime.Format(time.RFC3339),
	)
	// io.WriteString(w, "<a href=\"/cards\">Card List</a><br />\n")
}

func ZipDataHandler(w http.ResponseWriter, r *http.Request) {
	// Get a Buffer to Write To
	w.Header().Set("Content-Disposition", "attachment; filename=\"Valkyrie Crusade Fan Archive - Final - 2021-05-30.zip\"")
	w.Header().Set("Content-Type", "application/zip")

	// Create a new zip archive.
	z := zip.NewWriter(w)

	var err error

	err = zipCards(z)
	if err != nil {
		log.Printf("Card zip error: " + err.Error() + "\n")
		http.Error(w, "Card zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipWeapons(z)
	if err != nil {
		log.Printf("Weapon zip error: " + err.Error() + "\n")
		http.Error(w, "Weapon zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipItems(z)
	if err != nil {
		log.Printf("Item zip error: " + err.Error() + "\n")
		http.Error(w, "Item zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipStructures(z)
	if err != nil {
		log.Printf("Structure zip error: " + err.Error() + "\n")
		http.Error(w, "Structure zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipTreasure(z)
	if err != nil {
		log.Printf("Treasure zip error: " + err.Error() + "\n")
		http.Error(w, "Treasure zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipAudio(z)
	if err != nil {
		log.Printf("Audio zip error: " + err.Error() + "\n")
		http.Error(w, "Audio zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipAlliance(z)
	if err != nil {
		log.Printf("Alliance zip error: " + err.Error() + "\n")
		http.Error(w, "Alliance zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipEventStory(z)
	if err != nil {
		log.Printf("Story zip error: " + err.Error() + "\n")
		http.Error(w, "Story zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipNavi(z)
	if err != nil {
		log.Printf("Navi zip error: " + err.Error() + "\n")
		http.Error(w, "Navi zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//err = zipGardenSprites(z)
	if err != nil {
		log.Printf("Navi zip error: " + err.Error() + "\n")
		http.Error(w, "Navi zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//err = zipBattleImages(z) // dungeon, weapon event, battle backgrounds, battle maps
	if err != nil {
		log.Printf("Navi zip error: " + err.Error() + "\n")
		http.Error(w, "Navi zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = zipApkImages(z)
	if err != nil {
		log.Printf("ApkImages zip error: " + err.Error() + "\n")
		http.Error(w, "ApkImages zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = addFileToZip(z, "README.txt", nil, []byte(`Valkryie Crusade Fan Archive

This file was generated by Kushieda using the file grouper code located at
https://github.com/kushieda-minori/vc_file_grouper

I'm super happy to have been part of the community and glad that I got to
meet many of you.

I tried to make this archive as complete as possible, but as always, something
feels missing. Perhaps it's just that we can't play the game anymore.

Hope you all have fun in the future!`))

	if err != nil {
		log.Printf("Readme zip error: " + err.Error() + "\n")
		http.Error(w, "Readme zip error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = z.Close()
	if err != nil {
		log.Printf(err.Error() + "\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func zipCards(z *zip.Writer) error {
	seen := make([]string, 0)
	cardImageNames := make([]string, 0)
	for _, cardEvos := range vc.Data.Cards.CardsByName() {
		firstEvo := cardEvos.EarliestOpen()
		if firstEvo == nil {
			firstEvo = cardEvos.Earliest()
		}
		cardPathNameBase := ""
		if firstEvo.IsClosed > 0 {
			cardPathNameBase = "Cards/Inactive or Unreleased/" + firstEvo.MainRarity() + "/" + firstEvo.Element() + "/" + firstEvo.Name + "/"
		} else {
			cardPathNameBase = "Cards/" + firstEvo.MainRarity() + "/" + firstEvo.Element() + "/" + firstEvo.Name + "/"
		}
		character := firstEvo.Character()
		if character != nil && character.HasQuotes() {
			characterQuotes := "Description: " + character.Description + "\n"
			characterQuotes += "\nLogin: " + character.Login + "\n"
			characterQuotes += "\nMeet: " + character.Meet + "\n"
			characterQuotes += "\nFriendship: " + character.Friendship + "\n"
			characterQuotes += "\nFriendshipMax: " + character.FriendshipMax + "\n"
			characterQuotes += "\nFriendshipEvent: " + character.FriendshipEvent + "\n"
			characterQuotes += "\nBattleStart: " + character.BattleStart + "\n"
			characterQuotes += "\nBattleEnd: " + character.BattleEnd + "\n"
			characterQuotes += "\nRebirth: " + character.Rebirth + "\n"
			err := addFileToZip(z, cardPathNameBase+firstEvo.Name+" Quotes.txt", nil, []byte(characterQuotes))
			if err != nil {
				return err
			}
		}
		evolutions := firstEvo.GetEvolutions()
		for _, iconEvos := range firstEvo.EvosWithDistinctImages(true) {
			evo := evolutions[iconEvos]
			imgName, data, fsInfo, err := evo.GetImageData(true)
			if err != nil {
				return err
			}

			outputName := cardPathNameBase + imgName
			if util.Contains(seen, outputName) {
				continue
			}
			seen = append(seen, outputName)

			err = addFileToZip(z, outputName, &fsInfo, data)
			if err != nil {
				return err
			}
		}
		for _, iconEvos := range firstEvo.EvosWithDistinctImages(false) {
			evo := evolutions[iconEvos]
			imgName, data, fsInfo, err := evo.GetImageData(false)
			if err != nil {
				return err
			}

			outputName := cardPathNameBase + imgName
			if util.Contains(seen, outputName) {
				continue
			}
			seen = append(seen, outputName)

			err = addFileToZip(z, outputName, &fsInfo, data)
			if err != nil {
				return err
			}
		}

		for _, c := range evolutions {
			cardImageNames = append(cardImageNames, c.Image())
		}
	}

	filepath.Walk(vc.FilePath+"/card", func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !util.Contains(cardImageNames, info.Name()) {
			relPath, err := filepath.Rel(vc.FilePath+"/card", p)
			if err != nil {
				return err
			}
			relPath = filepath.ToSlash(relPath)
			var b []byte
			b, e = vc.Decode(p)
			if e != nil && strings.HasSuffix(e.Error(), "is not encoded") {
				b, e = ioutil.ReadFile(p)
			}
			if e != nil {
				return
			}

			uci := vc.CardScanImage(strings.TrimPrefix(info.Name(), "cd_"))
			if uci != nil {
				relPath = strings.TrimSuffix(relPath, info.Name())
				relPath += uci.Rarity() + " - " + uci.Name + " - " + info.Name()
			}

			if !strings.HasSuffix(strings.ToLower(relPath), ".png") {
				relPath += ".png"
			}
			e = addFileToZip(z, "Cards/Unused Images/"+relPath, &info, b)
			if e != nil {
				return
			}
		}
		return nil
	})

	return nil
}

func zipWeapons(z *zip.Writer) error {
	for _, weapon := range vc.Data.Weapons {
		wName := strings.Replace(weapon.MaxRarityName(), " (Weapon)", "", -1)
		pathNameBase := "Weapons/" + wName + "/"
		quotes := "Weapon Type: " + weapon.StatusDescription() + "\n"

		descriptions := weapon.Descriptions
		dLen := len(descriptions)
		rarities := weapon.Rarities()
		for _, rarity := range rarities {
			r := rarity.Rarity
			if r <= dLen {
				quotes += "\nDescription " + strconv.Itoa(r) + ": " + strings.ReplaceAll(descriptions[r-1], "\n", " ") + "\n"
			}
		}
		err := addFileToZip(z, pathNameBase+"Quotes.txt", nil, []byte(quotes))
		if err != nil {
			return err
		}
		for imageName, data := range weapon.GetImageData(true) {
			err = addFileToZip(z, pathNameBase+imageName, nil, data)
			if err != nil {
				return err
			}
		}
		for imageName, data := range weapon.GetImageData(false) {
			err = addFileToZip(z, pathNameBase+imageName, nil, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func zipItems(z *zip.Writer) error {
	seen := make([]string, 0)
	for _, item := range vc.Data.Items {
		pathNameBase := "Items/"
		if group, ok := vc.ItemGroups[item.GroupID]; ok {
			pathNameBase += group + "/"
		} else {
			pathNameBase += "Other-" + strconv.Itoa(item.GroupID) + "/"
		}
		imageName, data, fsInfo, err := item.GetImageData()
		if err != nil {
			return err
		}
		if len(imageName) > 0 {
			outputName := pathNameBase + imageName
			if util.Contains(seen, outputName) {
				continue
			}
			seen = append(seen, outputName)
			err = addFileToZip(z, outputName, &fsInfo, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func zipTreasure(z *zip.Writer) error {
	return filepath.Walk(vc.FilePath+"/treasure/", func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(vc.FilePath+"/treasure/", p)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		lrelPath := strings.ToLower(relPath)
		if strings.HasPrefix(lrelPath, "menu") {
			return
		}

		fsInfo, _ := os.Stat(p)

		pathName := relPath + ".png"

		var b []byte
		b, e = vc.Decode(p)
		if e != nil {
			return
		}

		e = addFileToZip(z, "Sacred Treasure/"+pathName, &fsInfo, b)

		return
	})
}

func zipStructures(z *zip.Writer) error {
	seen := make([]string, 0)
	for _, structure := range vc.Data.Structures {
		group := vc.ShopGroup[structure.ShopGroupDecoID]
		if strings.HasPrefix(structure.Name, "Honorable Plaque") || strings.HasPrefix(structure.Name, "Shield of Honor") {
			group = "Special/Honorable Plaque"
		}
		pathNameBase := "Structures/" + group + "/"
		binImages, err := structure.GetImageData()
		if err != nil {
			return err
		}
		//log.Printf("Found %d images for structure %s", len(binImages), structure.Name)
		// if structure.IsResource() || structure.IsBank() {
		// 	pathNameBase += "Resources/"
		// }
		if len(binImages) > 1 {
			pathNameBase += strings.ReplaceAll(structure.Name, "/", "-") + "/"
		}
		for i, binImg := range binImages {
			if len(binImg.Data) > 0 {
				outputName := pathNameBase + fmt.Sprintf("%s_%02d.png", strings.ReplaceAll(structure.Name, "/", "-"), i+1)
				if util.Contains(seen, outputName) {
					continue
				}
				seen = append(seen, outputName)
				err = addFileToZip(z, outputName, nil, binImg.Data)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func zipAudio(z *zip.Writer) error {
	return filepath.Walk(vc.FilePath, func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		lp := strings.ToLower(p)
		if strings.HasSuffix(lp, ".ogg") ||
			strings.HasSuffix(lp, ".wav") ||
			strings.HasSuffix(lp, ".flac") ||
			strings.HasSuffix(lp, ".mp3") {
			relPath, err := filepath.Rel(vc.FilePath, p)
			if err != nil {
				return err
			}
			relPath = filepath.ToSlash(relPath)
			var b []byte
			b, e = ioutil.ReadFile(p)
			if e != nil {
				return
			}
			fsInfo, _ := os.Stat(p)
			e = addFileToZip(z, "Audio/"+relPath, &fsInfo, b)
			if e != nil {
				return
			}
		}
		return
	})
}

func zipAlliance(z *zip.Writer) error {
	return filepath.Walk(vc.FilePath+"/guild/texture/", func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(vc.FilePath+"/guild/texture/", p)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		lrelPath := strings.ToLower(relPath)
		if strings.HasPrefix(lrelPath, "menu") {
			return
		}

		fsInfo, _ := os.Stat(p)

		pathName := ""
		if strings.HasPrefix(lrelPath, "sym_a") {
			pathName = "Alliance/Background Symbol/" + relPath
		} else if strings.HasPrefix(lrelPath, "sym_b") {
			pathName = "Alliance/Main Symbol/" + relPath
		} else if strings.HasPrefix(lrelPath, "sym_c") {
			pathName = "Alliance/Decoration Symbol/" + relPath
		} else if strings.HasPrefix(lrelPath, "stamp_1") {
			pathName = "Alliance/Stamp/Oracle" + strings.Replace(relPath, "stamp_1", "_", -1)
		} else if strings.HasPrefix(lrelPath, "stamp_2") {
			pathName = "Alliance/Stamp/Alchemist" + strings.Replace(relPath, "stamp_2", "_", -1)
		} else if strings.HasPrefix(lrelPath, "stamp_3") {
			pathName = "Alliance/Stamp/Pixie" + strings.Replace(relPath, "stamp_3", "_", -1)
		} else if strings.HasPrefix(lrelPath, "stamp_4") {
			pathName = "Alliance/Stamp/Hades" + strings.Replace(relPath, "stamp_4", "_", -1)
		} else if strings.HasPrefix(lrelPath, "stamp_5") {
			pathName = "Alliance/Stamp/Circo" + strings.Replace(relPath, "stamp_5", "_", -1)
		} else if strings.HasPrefix(lrelPath, "stamp_6") {
			if strings.HasSuffix(lrelPath, "9") {
				pathName = "Alliance/Stamp/Calamity" + strings.Replace(relPath, "stamp_6", "_", -1)
			} else {
				pathName = "Alliance/Stamp/Fenrir and Skoll" + strings.Replace(relPath, "stamp_6", "_", -1)
			}
		} else if strings.HasPrefix(lrelPath, "stamp_7") {
			if strings.HasSuffix(lrelPath, "7") || strings.HasSuffix(lrelPath, "8") || strings.HasSuffix(lrelPath, "9") {
				pathName = "Alliance/Stamp/Demon Ministers" + strings.Replace(relPath, "stamp_7", "_", -1)
			} else {
				pathName = "Alliance/Stamp/Calamity" + strings.Replace(relPath, "stamp_7", "_", -1)
			}
		} else if strings.HasPrefix(lrelPath, "stamp_8") {
			pathName = "Alliance/Stamp/Demon Ministers" + strings.Replace(relPath, "stamp_8", "_", -1)
		}

		if pathName != "" {
			var b []byte
			b, e = vc.Decode(p)
			if e != nil && strings.HasSuffix(e.Error(), "is not encoded") {
				b, e = ioutil.ReadFile(p)
			} else {
				pathName += ".png"
			}
			if e != nil {
				return
			}
			e = addFileToZip(z, pathName, &fsInfo, b)
		}

		return
	})
}

func zipNavi(z *zip.Writer) error {
	return filepath.Walk(vc.FilePath+"/navi/", func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(vc.FilePath+"/navi/", p)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		lrelPath := strings.ToLower(relPath)
		if strings.HasPrefix(lrelPath, "menu") {
			return
		}

		fsInfo, _ := os.Stat(p)

		pathName := "Navi-Sprites/" + strings.TrimPrefix(relPath, "flash/") + ".png"

		if pathName != "" {
			var b []byte
			b, e = vc.Decode(p)
			if e != nil {
				if strings.HasSuffix(e.Error(), "is not encoded") {
					return nil
				}
				//b, e = ioutil.ReadFile(p)
				return
			}
			e = addFileToZip(z, pathName, &fsInfo, b)
		}

		return
	})
}

func zipEventStory(z *zip.Writer) (err error) {
	var txt string

	for _, m := range vc.Data.Maps {
		if !m.HasStory() || m.CleanedEventName() == "" {
			continue
		}

		txt = fmt.Sprintf("<html><head><title>%[1]s : %[2]s</title><style>table,th,td{border:1px solid black;}</style></head></body><h1>%[1]s : %[2]s</h1>\n\n", m.EventName(), m.Name)
		if m.StartMsg != "" {
			txt += "<dl><dt>Introduction:</dt><dd>" + m.StartMsg + "</dd></dl>"
		}
		txt += `<table style="border-collapse:collapse;">`
		for _, a := range m.Areas() {
			if a.HasStory() {
				story := ""
				if a.Story != "" {
					story += fmt.Sprintf("<dl><dt>Prologue:</dt><dd>%s</dd></dl>\n", html.EscapeString(strings.ReplaceAll(a.Story, "\n", " ")))
					if a.Start != "" || a.End != "" || a.BossStart != "" || a.BossEnd != "" {
						story += "<hr/>\n\n"
					}
				}
				if a.Start != "" || a.End != "" {
					story += "<dl><dt>Guide Dialogue:</dt>"
					if a.Start != "" {
						story += fmt.Sprintf("\n<dd>%s</dd>\n", html.EscapeString(strings.ReplaceAll(a.Start, "\n", " ")))
						if a.End != "" {
							story += "<br />&nbsp;<br />"
						}
					}
					if a.End != "" {
						story += fmt.Sprintf("\n<dd>%s</dd>\n", html.EscapeString(strings.ReplaceAll(a.End, "\n", " ")))
					} else {
						story += "\n"
					}
					if a.BossStart != "" || a.BossEnd != "" {
						story += "<hr/>\n\n"
					}
				}

				if a.BossStart != "" || a.BossEnd != "" {
					story += "<dl><dt>Boss Dialogue:</dt>"
					if a.BossStart != "" {
						story += fmt.Sprintf("\n<dd>%s</dd>", html.EscapeString(strings.ReplaceAll(a.BossStart, "\n", " ")))
						if a.BossEnd != "" {
							story += "<br />&nbsp;<br />"
						}
					}
					if a.BossEnd != "" {
						story += fmt.Sprintf("\n<dd>%s</dd>", html.EscapeString(strings.ReplaceAll(a.BossEnd, "\n", " ")))
					} else {
						story += "\n"
					}
				}
				txt += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", a.Name, story)
			}
		}
		txt += "</table>\n</body>\n</html>\n"

		err = addFileToZip(z, "Event Stories/Archwitch/"+m.PublicStartDatetime.Format("2006-01-02")+" - "+m.CleanedEventName()+".html", nil, []byte(txt))
		if err != nil {
			return
		}
	}

	for _, t := range vc.Data.Towers {
		txt, err = t.ScenarioHtml()
		if err != nil {
			return
		}
		if txt == "" {
			continue
		}
		err = addFileToZip(z, "Event Stories/Tower/"+t.PublicStartDatetime.Format("2006-01-02")+" - "+t.CleanedEventName()+".html", nil, []byte(txt))
		if err != nil {
			return
		}
	}

	for _, d := range vc.Data.Dungeons {
		if d.ScenarioID > 0 {
			txt, err = d.ScenarioHtml()
			if err != nil {
				return
			}
			err = addFileToZip(z, "Event Stories/DRV/"+d.PublicStartDatetime.Format("2006-01-02")+" - "+d.CleanedEventName()+".html", nil, []byte(txt))
			if err != nil {
				return
			}
		}
	}

	for _, w := range vc.Data.WeaponEvents {
		if w.ScenarioID > 0 {
			txt, err = w.ScenarioHtml()
			if err != nil {
				return
			}
			err = addFileToZip(z, "Event Stories/Weapon/"+w.PublicStartDatetime.Format("2006-01-02")+" - "+w.CleanedEventName()+".html", nil, []byte(txt))
			if err != nil {
				return
			}
		}
	}

	return
}

func zipApkImages(z *zip.Writer) error {
	return filepath.Walk(vc.FilePath+"/apk_images/", func(p string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(p), ".png") || strings.HasSuffix(strings.ToLower(p), "template.png") {
			return nil
		}

		relPath, err := filepath.Rel(vc.FilePath, p)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		fsInfo, _ := os.Stat(p)

		if relPath != "" {
			var b []byte
			b, e = ioutil.ReadFile(p)
			if e != nil {
				return
			}
			e = addFileToZip(z, relPath, &fsInfo, b)
		}

		return
	})
}

func addFileToZip(w *zip.Writer, filePathAndName string, fsInfo *fs.FileInfo, data []byte) (err error) {
	var f io.Writer
	if fsInfo == nil {
		f, err = w.Create(filePathAndName)
		if err != nil {
			return
		}
	} else {
		var fih *zip.FileHeader
		fih, err = zip.FileInfoHeader(*fsInfo)
		if err != nil {
			return
		}
		fih.Name = filePathAndName
		fih.Method = zip.Deflate
		f, err = w.CreateHeader(fih)
		if err != nil {
			return
		}
	}
	_, err = f.Write(data)
	if err != nil {
		return
	}
	return nil
}

type mapDlResult struct {
	Success   bool
	Map       *vc.Map
	Timestamp int
}

//DownloadAwMapsHandler
func DownloadAwMapsHandler(w http.ResponseWriter, r *http.Request) {
	lMaps := len(vc.Data.Maps)
	numJobs := 30
	maps := make(chan *vc.Map, lMaps)
	results := make(chan mapDlResult, lMaps)

	// set up my worker pool
	for i := 0; i < numJobs; i++ {
		go findAndDownloadAwMap(i, maps, results)
	}
	queued := 0
	for i := range vc.Data.Maps {
		m := &(vc.Data.Maps[i])
		if !m.PublicStartDatetime.IsZero() {
			log.Printf("Searching for maps for event: %s", m.CleanedEventName())
			maps <- m
			queued++
		}
	}
	close(maps)

	found := 0
	completed := 0
	for r := range results {
		completed++
		if r.Success {
			found++
		} else {
			log.Printf("Failed to DL Map %d: %s %s", r.Map.ID, r.Map.EventName(), r.Map.Name)
			fmt.Fprintf(w, "Failed to DL Map %d: %s %s", r.Map.ID, r.Map.EventName(), r.Map.Name)
		}
		if completed == queued {
			close(results)
		}
		log.Printf("completed %d (found %d) of %d map files: map %d : %d", completed, found, queued, r.Map.ID, r.Timestamp)
		fmt.Fprintf(w, "completed %d (found %d) of %d map files: map %d : %d", completed, found, queued, r.Map.ID, r.Timestamp)
	}
	fmt.Fprintf(w, "Found %d of %d map files", found, completed)
}

type mapDl struct {
	Timestamp int
	DestName  string
}

func findAndDownloadAwMap(wkId int, maps chan *vc.Map, results chan mapDlResult) {
	defer log.Printf("Shutting down findAndDownloadAwMap worker process %d", wkId)
	for m := range maps {
		eventName := m.CleanedEventName()
		fileName := fmt.Sprintf("AreaMap_002_%02d.%s.%s", m.ID, eventName, m.CleanedName())
		fileLoc := filepath.Join(vc.FilePath, "battle", "map", fileName)
		if s, err := os.Stat(fileLoc); err == nil {
			log.Printf("File already exists on disk: %s", fileName)
			results <- mapDlResult{Success: true, Map: m, Timestamp: int(s.ModTime().Unix())}
			continue
		}

		const numJobs = 10
		timestamps := make(chan mapDl, numJobs*2)
		done := make(chan int)
		cancel := make(chan bool)
		// set up my worker pool
		for i := 0; i < numJobs; i++ {
			go downloadWorker(wkId, m, timestamps, done, cancel)
		}
		timestamp := int(m.PublicStartDatetime.Unix()) + 3600
		end := timestamp - (60 * 60 * 24) - 3600
		//found := make(chan bool)
		for i := timestamp; i > end; i-- {
			select {
			case found := <-done:
				cancel <- true
				results <- mapDlResult{Success: found > 0, Map: m, Timestamp: found}
				close(timestamps)
				return
			default:
				timestamps <- mapDl{Timestamp: i, DestName: fileLoc}
			}
			// if downloadAwMap(i, fileLoc) {
			// 	results <- true
			// 	return
			// }
		}
		close(timestamps)
		found := <-done
		cancel <- true
		results <- mapDlResult{Success: found > 0, Map: m, Timestamp: found}
	}
}

func downloadWorker(parentWkId int, m *vc.Map, timestamps chan mapDl, done chan int, cancel chan bool) {
	defer log.Printf("Shutting down downloadWorker worker process under parent %d for map %d", parentWkId, m.ID)
	for dl := range timestamps {
		select {
		case <-cancel:
			cancel <- true
			return
		default:
			if downloadAwMap(dl.Timestamp, dl.DestName) {
				cancel <- true
				done <- dl.Timestamp
				close(done)
				return
			}
		}
	}
	done <- 0
}

func downloadAwMap(timestamp int, fileName string) bool {
	url := fmt.Sprintf("http://webview.valkyriecrusade.nubee.com/download/BattleMap.zip/AreaMap_002.%d", timestamp)
	// if timestamp%50 == 0 {
	// 	log.Printf("***Trying to DL %s", url)
	// }
	var resp *http.Response
	var err error
	retries := 0 // retry on timeouts
	for ok := true; ok; ok = resp.StatusCode == 408 && retries <= 10 {
		resp, err = http.Get(url)
		if err != nil {
			log.Printf("Get Err for map %d: %s", timestamp, err.Error())
			return false
		}
		defer resp.Body.Close()
		retries++
		if resp.StatusCode == 408 && retries <= 10 {
			log.Printf("download %d failed. Retry: %d", timestamp, retries)
			time.Sleep(100 * time.Millisecond)
		}
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var b []byte
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Read Err for map %d: %s", timestamp, err.Error())
			return false
		}
		err = ioutil.WriteFile(fileName, b, 0755)
		if err != nil {
			log.Printf("Write Err for map %d: %s", timestamp, err.Error())
			return false
		}
		t := time.Unix(int64(timestamp), 0)
		os.Chtimes(fileName, t, t)
		//log.Printf("***Found for map %d", timestamp)
		return true
	}
	if resp.StatusCode != 403 && resp.StatusCode != 404 {
		log.Printf("Http Status Err for map %d: %d %s", timestamp, resp.StatusCode, resp.Status)
	}
	return false
}
