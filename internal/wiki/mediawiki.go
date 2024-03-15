package wiki

import (
	"github.com/mal0ner/wikiscrape/internal/export"
	"github.com/mal0ner/wikiscrape/internal/scrape"
	"github.com/mal0ner/wikiscrape/internal/util"
)

type MediaWiki struct {
	Name    string
	BaseURL string
	scrape.Scraper
	export.Exporter
	util.Manifest
}

func NewMediaWiki(name string, baseURL string) MediaWiki {
	return MediaWiki{
		Name:     name,
		BaseURL:  baseURL,
		Scraper:  &scrape.MediaWikiScraper{BaseURL: baseURL},
		Exporter: &export.TestExporter{},
	}
}

// TODO: Add Pointer receiver here
func (w MediaWiki) ScrapeAndExport(man util.Manifest) error {
	// Exporting pages individually as they are gathered seems to be the
	// best approach here, aggregating thousands of pages into a slice
	// of structs only to export them immediately after seems like a waste
	// bc extra memory consumption. We do lose the ability to
	// export concurrently though
	for _, path := range man {
		page, err := w.GetPage(path)
		if err != nil {
			continue
		}
		w.Export(*page) // TODO: Cleanup pointers
	}
	return nil
}

// TODO: Rename this and maybe make the scraper methods private??
// TODO: Add Pointer receiver here
func (w MediaWiki) Page(path string) error {
	page, err := w.GetPage(path)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

// TODO: Rename this and maybe make the scraper methods private??
// TODO: Add Pointer receiver here
func (w MediaWiki) Section(path string, heading string) error {
	page, err := w.GetSection(path, heading)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}
