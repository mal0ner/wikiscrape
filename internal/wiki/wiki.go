package wiki

import (
	"github.com/mal0ner/wikiscrape/internal/manifest"
)

type Wiki interface {
	ScrapeAndExport(manifest.Manifest) error
	Page(string) error
	Section(string, string) error
}
