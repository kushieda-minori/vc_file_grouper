package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"zetsuboushita.net/vc_file_grouper/bot"
	"zetsuboushita.net/vc_file_grouper/handler"
	"zetsuboushita.net/vc_file_grouper/vc"
)

// Main function that starts the program
func main() {

	cmdLang := flag.String("lang", "en", "The language pack to use. 'en' for English, 'zhs' for Chinese. ")
	cmdHelp := flag.Bool("help", false, "Show the help message")
	cmdDbg := flag.Bool("debug", false, "Outputs error message to the standard error console")
	flag.Parse()

	if *cmdHelp {
		usage()
		return
	}

	if !(*cmdDbg) {
		log.SetOutput(ioutil.Discard)
	}

	if cmdLang == nil {
		vc.LangPack = "en"
	} else {
		vc.LangPack = *cmdLang
	}

	if len(flag.Args()) == 0 {
		vc.FilePath = "."
	} else {
		vc.FilePath = flag.Args()[0]
		if len(flag.Args()) > 1 {
			bot.DbFileLocation = flag.Args()[1]
		}
	}

	if _, err := os.Stat(vc.FilePath); os.IsNotExist(err) {
		usage()
		vc.Data = &vc.VFile{}
		//return
	} else {
		vc.ReadMasterData(vc.FilePath)
	}
	if bot.DbFileLocation != "" {
		if err := bot.LoadDb(); err != nil {
			log.Printf("Error loading Bot DB: " + err.Error())
		}

	}

	//main page
	http.HandleFunc("/", handler.MasterDataHandler)
	http.HandleFunc("/css/", cssHandler)
	//image locations
	http.HandleFunc("/images/card/", handler.ImageCardHandler)
	http.HandleFunc("/images/cardthumb/", handler.ImageCardThumbHandler)
	http.HandleFunc("/images/cardHD/", handler.ImageCardHDHandler)
	http.HandleFunc("/images/event/", handler.ImageHandlerFor("/event/", "/event/"))
	http.HandleFunc("/images/battle/", handler.ImageHandlerFor("/battle/", "/battle/"))
	http.HandleFunc("/images/garden/", handler.ImageHandlerFor("/garden/", "/garden/"))
	http.HandleFunc("/images/garden/map/", handler.StructureImagesHandler)
	http.HandleFunc("/images/dungeon/", handler.ImageHandlerFor("/dungeon/", "/dungeon/"))
	http.HandleFunc("/images/alliance/", handler.ImageHandlerFor("/alliance/", "/guild/"))
	http.HandleFunc("/images/summon/", handler.ImageHandlerFor("/summon/", "/gacha/"))
	http.HandleFunc("/images/item/", handler.ImageHandlerFor("/item/", "/item/"))
	http.HandleFunc("/images/treasure/", handler.ImageHandlerFor("/treasure/", "/treasure/"))
	http.HandleFunc("/images/navi/", handler.ImageHandlerFor("/navi/", "/navi/"))

	// vc master data
	http.HandleFunc("/config/", handler.ConfigHandler)
	//dynamic pages
	http.HandleFunc("/cards/", handler.CardHandler)
	http.HandleFunc("/cards/table/", handler.CardTableHandler)
	http.HandleFunc("/cards/csv/", handler.CardCsvHandler)
	http.HandleFunc("/cards/glrcsv/", handler.CardCsvGLRHandler)
	http.HandleFunc("/cards/detail/", handler.CardDetailHandler)
	http.HandleFunc("/cards/levels/", handler.CardLevelHandler)
	http.HandleFunc("/archwitches/", handler.ArchwitchHandler)
	http.HandleFunc("/characters/", handler.CharacterTableHandler)
	http.HandleFunc("/characters/detail/", handler.CharacterDetailHandler)
	// http.HandleFunc("/character/csv/", handler.CharacterCsvHandler)

	http.HandleFunc("/items/", handler.ItemHandler)

	http.HandleFunc("/skills/", handler.SkillTableHandler)
	http.HandleFunc("/skills/csv/", handler.SkillCsvHandler)

	http.HandleFunc("/deckbonus/", handler.DeckBonusHandler)
	http.HandleFunc("/deckbonus/WIKI/", handler.DeckBonusWikiHandler)

	http.HandleFunc("/events/", handler.EventHandler)
	http.HandleFunc("/events/detail/", handler.EventDetailHandler)

	http.HandleFunc("/thor/", handler.ThorHandler)

	http.HandleFunc("/maps/", handler.MapHandler)

	http.HandleFunc("/garden/structures/", handler.StructureListHandler)
	http.HandleFunc("/garden/structures/detail/", handler.StructureDetailHandler)

	http.HandleFunc("/awakenings/", handler.AwakeningsTableHandler)
	http.HandleFunc("/awakenings/csv/", handler.AwakeningsCsvHandler)

	http.HandleFunc("/decode/", handler.DecodeHandler)

	http.HandleFunc("/bot/", handler.BotHandler)
	http.HandleFunc("/bot/config", handler.BotConfigHandler)
	http.HandleFunc("/bot/update", handler.BotUpdateHandler)

	http.HandleFunc("/raw/", handler.RawDataHandler)
	http.HandleFunc("/raw/KEYS", handler.RawDataKeysHandler)

	http.HandleFunc("/SHUTDOWN/", func(w http.ResponseWriter, r *http.Request) { os.Exit(0) })

	os.Stdout.WriteString("Listening on port 8585. Connect to http://localhost:8585/\nPress <CTRL>+C to stop or close the terminal.\n")
	err := http.ListenAndServe("localhost:8585", nil)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
}

// Prints useage to the console
func usage() {
	os.Stdout.WriteString(fmt.Sprintf("To use this program you can specify the following command options:\n"+
		"-help\n\tShow this help message\n"+
		"-lang\n\tSelect a language pack to use. 'en' is the default\n"+
		"-debug\n\tOutputs error message to the standard error console\n"+
		"file1\n\tlocation of the VC master data file\n"+
		"file2\n\tlocation of the VC bot data file\n\n"+
		"example usages:\n\t%[1]s -help\n"+
		"\t%[1]s -lang %[2]s\n"+
		"\t%[1]s \"%[3]s\"\n"+
		"\t%[1]s \"%[3]s\" \"%[4]s\"\n"+
		"\t%[1]s -lang %[2]s \"%[3]s\" \"%[4]s\"\n",
		os.Args[0],
		"zhs",
		"/path/to/vc/data/file",
		"/path/to/bot/db/file",
	))
}
