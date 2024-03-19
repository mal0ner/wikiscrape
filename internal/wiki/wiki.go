package wiki

import "github.com/mal0ner/wikiscrape/internal/util"

type Wiki interface {
	ScrapeManifest(util.Manifest) error
	ScrapePage(string) error
	ScrapeSection(string, string) error
}
