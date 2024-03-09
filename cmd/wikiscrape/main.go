package main

import (
	"fmt"
	"os"

	logger "github.com/mal0ner/wikiscrape/internal/logging"
	"github.com/mal0ner/wikiscrape/internal/manifest"
	"github.com/mal0ner/wikiscrape/internal/wiki"
)

func main() {
	// lol
	logger.Log.SetLevel(5)
	paths, err := manifest.ReadFrom("manifest.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w := wiki.NewMediaWiki("Old School Runescape", "https://oldschool.runescape.wiki/api.php")
	w.Manifest = paths

	// err = w.ScrapeAndExport()
	err = w.Section("Zulrah", "drops")
	if err != nil {
		fmt.Println(err.Error())
	}
	// w.Export()
}
