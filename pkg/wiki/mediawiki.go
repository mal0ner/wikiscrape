package wiki

import (
	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/mal0ner/wikiscrape/pkg/export"
	"github.com/mal0ner/wikiscrape/pkg/logging"
	"github.com/mal0ner/wikiscrape/pkg/manifest"
	"github.com/mal0ner/wikiscrape/pkg/scrape"
)

type MediaWiki struct {
	Name     string
	BaseURL  string
	Pages    []scrape.Page
	Scraper  scrape.MediaWikiScraper
	Exporter export.Exporter
	Manifest manifest.Manifest
}

func NewMediaWiki(name string, baseURL string) MediaWiki {
	return MediaWiki{
		Name:    name,
		BaseURL: baseURL,
		Pages:   []scrape.Page{},
		Scraper: scrape.MediaWikiScraper{
			BaseURL:               baseURL,
			MediaWikiParser:       scrape.MediaWikiParser{},
			MediaWikiQueryBuilder: util.MediaWikiQueryBuilder{},
		},
		Exporter: &export.TestExporter{},
	}
}

func (w *MediaWiki) Crawl() error {
	logging.Log.WithField("name", w.Name).Info("Crawling")
	pages, err := w.Scraper.Scrape(w.Manifest)
	if err != nil {
		return err
	}
	w.Pages = pages
	return nil
}

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

func (w *MediaWiki) Export() error {
	logging.Log.WithField("name", w.Name).Info("Exporting")
	for _, p := range w.Pages {
		w.Exporter.Export(p)
	}
	return nil
}
