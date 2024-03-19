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
func NewMediaWiki(name string, baseURL string) Wiki {
	return &MediaWiki{
		Name:     name,
		BaseURL:  baseURL,
		Scraper:  &scrape.MediaWikiScraper{BaseURL: baseURL},
		Exporter: &export.TestExporter{},
	}
}

// ScrapeAndExport loops over a util.Manifest ([]string) list of page
// names, scraping and then exporting each page sequentially.
func (wiki *MediaWiki) ScrapeManifest(man util.Manifest) error {
	for _, path := range man {
		page, err := wiki.GetPage(path)
		if err != nil {
			continue
		}
		wiki.Export(page)
	}
	return nil
}

// Page provides a convenient wrapper around the wiki's
// scraper and exporter. Fetches, parses, and exports a single
// page given its name.
func (wiki *MediaWiki) ScrapePage(path string) error {
	page, err := wiki.GetPage(path)
	if err != nil {
		return err
	}
	wiki.Export(page)
	return nil
}

// Section provides a convenient wrapper around the wiki's
// scraper and exporter. Fetches, parses, and exports a single
// section of a single page given its name.
func (wiki *MediaWiki) ScrapeSection(path string, heading string) error {
	page, err := wiki.GetSection(path, heading)
	if err != nil {
		return err
	}
	wiki.Export(page)
	return nil
}
