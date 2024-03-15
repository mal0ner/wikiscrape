package wiki

import "github.com/mal0ner/wikiscrape/internal/util"

type Wiki interface {
	ScrapeAndExport(util.Manifest) error
	Page(string) error
	Section(string, string) error
}
