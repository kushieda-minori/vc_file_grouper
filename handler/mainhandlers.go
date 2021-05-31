package handler

import (
	"archive/zip"
	"fmt"
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"vcData-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".zip\"")
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

	err = zipAudio(z)
	if err != nil {
		log.Printf("Structure zip error: " + err.Error() + "\n")
		http.Error(w, "Structure zip error: "+err.Error(), http.StatusInternalServerError)
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
			err = addFileToZip(z, cardPathNameBase+imgName, &fsInfo, data)
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
			err = addFileToZip(z, cardPathNameBase+imgName, &fsInfo, data)
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
	itemsSeen := make([]string, 0)
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
			if util.Contains(itemsSeen, outputName) {
				continue
			}
			itemsSeen = append(itemsSeen, outputName)
			err = addFileToZip(z, outputName, &fsInfo, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
