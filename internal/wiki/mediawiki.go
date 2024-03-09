package wiki

import (
	"github.com/mal0ner/wikiscrape/internal/export"
	"github.com/mal0ner/wikiscrape/internal/logging"
	"github.com/mal0ner/wikiscrape/internal/manifest"
	"github.com/mal0ner/wikiscrape/internal/scrape"
	"github.com/sirupsen/logrus"
)

type MediaWiki struct {
	Name    string
	BaseURL string
	scrape.Scraper
	export.Exporter
	manifest.Manifest
}

// TODO: Decide on best approach for reflecting command-line args for
// scraper and export initialization. Do we embed them in the wiki
// struct like this and just pass them into the constructor via
// an argument, or pass an instance of scrape.Scraper / export.Exporter
// into the ScrapeAndExport methods manually?
// Passing them into the NewMediaWiki might be the easiest way for other
// people to interact with this as a library??
func NewMediaWiki(name string, baseURL string) MediaWiki {
	return MediaWiki{
		Name:     name,
		BaseURL:  baseURL,
		Scraper:  &scrape.MediaWikiScraper{BaseURL: baseURL},
		Exporter: &export.TestExporter{},
	}
}

func (w *MediaWiki) ScrapeAndExport(man manifest.Manifest) error {
	// Exporting pages individually as they are gathered seems to be the
	// best approach here, aggregating thousands of pages into a slice
	// of structs only to export them immediately after seems like a waste
	// bc extra memory consumption. We do lose the ability to
	// export concurrently though
	logging.Log.WithField("name", w.Name).Info("Crawling")
	for _, path := range man {
		page, err := w.GetPage(path)
		if err != nil {
			logging.Log.WithFields(logrus.Fields{
				"page":  path,
				"error": err.Error(),
			}).Error("Failed to get page")
			continue
		}
		w.Export(*page) // TODO: Cleanup pointers
	}
	return nil
}

// TODO: Rename this and maybe make the scraper methods private??
func (w *MediaWiki) Page(path string) error {
	page, err := w.GetPage(path)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

// TODO: Rename this and maybe make the scraper methods private??
func (w *MediaWiki) Section(path string, heading string) error {
	page, err := w.GetSection(path, heading)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}
