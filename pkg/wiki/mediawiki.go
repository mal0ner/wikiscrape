package wiki

import (
	"github.com/mal0ner/wikiscrape/pkg/export"
	"github.com/mal0ner/wikiscrape/pkg/logging"
	"github.com/mal0ner/wikiscrape/pkg/manifest"
	"github.com/mal0ner/wikiscrape/pkg/scrape"
	"github.com/sirupsen/logrus"
)

type MediaWiki struct {
	Name    string
	BaseURL string
	scrape.Scraper
	export.Exporter
	manifest.Manifest
}

func NewMediaWiki(name string, baseURL string) MediaWiki {
	return MediaWiki{
		Name:     name,
		BaseURL:  baseURL,
		Scraper:  &scrape.MediaWikiScraper{BaseURL: baseURL},
		Exporter: &export.TestExporter{},
		Manifest: []string{},
	}
}

func (w *MediaWiki) ScrapeAndExport() error {
	logging.Log.WithField("name", w.Name).Info("Crawling")
	for _, path := range w.Manifest {
		page, err := w.GetPage(path)
		if err != nil {
			logging.Log.WithFields(logrus.Fields{
				"page":  path,
				"error": err.Error(),
			}).Error("Failed to get page")
			continue
		}
		w.Export(*page)
	}
	return nil
}

func (w *MediaWiki) Page(path string) error {
	page, err := w.GetPage(path)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

func (w *MediaWiki) Section(path string, heading string) error {
	page, err := w.GetSection(path, heading)
	if err != nil {
		return err
	}
	w.Export(*page) // TODO: Pointer here?
	return nil
}

// TODO: Maybe do a refactor in here and the mediawikiscraper to allow
// for the user to scrape a list of pages and only get the sections they
// are looking for?

// func (w MediaWiki) Export() error {
// 	var wg sync.WaitGroup
// 	for _, p := range w.Pages {
// 		wg.Add(1)
// 		go func(p scrape.Page) {
// 			defer wg.Done()
// 			w.Exporter.Export(p)
// 		}(p)
// 	}
// 	wg.Wait()
// 	return nil
// }

// func (w *MediaWiki) Export() error {
// 	logging.Log.WithField("name", w.Name).Info("Exporting")
// 	for _, p := range w.Pages {
// 		w.Exporter.Export(p)
// 	}
// 	return nil
// }
