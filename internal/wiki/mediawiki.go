package wiki

import (
	"github.com/mal0ner/wikiscrape/internal/export"
	"github.com/mal0ner/wikiscrape/internal/scrape"
	"github.com/mal0ner/wikiscrape/internal/util"
)

// MediaWiki represents a "Wiki" interface compliant implementation
// with fields and methods designed to add support for querying from
// the Media Wiki API and parsing / manipulating the pages.
type MediaWiki struct {
	Name    string
	BaseURL string
	scrape.Scraper
	export.Exporter
	util.Manifest
}

// NewMediaWiki instantiates a new Media Wiki with a provided name and
// base url, as well as sensible defaults for the scraper and exporter.
func NewMediaWiki(name string, baseURL string) MediaWiki {
	return MediaWiki{
		Name:     name,
		BaseURL:  baseURL,
		Scraper:  &scrape.MediaWikiScraper{BaseURL: baseURL},
		Exporter: &export.TestExporter{},
	}
}

// ScrapeAndExport loops over a util.Manifest ([]string) list of page
// names, scraping and then exporting each page sequentially.
func (w MediaWiki) ScrapeAndExport(man util.Manifest) error {
	// TODO: Add Pointer receiver here
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

// Page provides a convenient wrapper around the wiki's
// scraper and exporter. Fetches, parses, and exports a single
// page given its name.
func (w MediaWiki) Page(path string) error {
	// TODO: Add Pointer receiver here
	page, err := w.GetPage(path)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

// Section provides a convenient wrapper around the wiki's
// scraper and exporter. Fetches, parses, and exports a single
// section of a single page given its name.
func (w MediaWiki) Section(path string, heading string) error {
	// TODO: Add Pointer receiver here
	page, err := w.GetSection(path, heading)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

// TODO: Maybe fix these above methods? Especially the names.
