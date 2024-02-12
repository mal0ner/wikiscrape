package main

import (
	"fmt"
	"os"

	logger "github.com/mal0ner/wikiscrape/pkg/logging"
	"github.com/mal0ner/wikiscrape/pkg/manifest"
	"github.com/mal0ner/wikiscrape/pkg/wiki"
)

func main() {
	logger.Log.SetLevel(5)
	paths, err := manifest.ReadFrom("manifest.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w := wiki.NewMediaWiki("Old School Runescape", "https://oldschool.runescape.wiki/api.php")
	w.Manifest = paths

	err = w.Crawl()
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Export()
}
