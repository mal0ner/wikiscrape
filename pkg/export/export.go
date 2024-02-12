package export

import "github.com/mal0ner/wikiscrape/pkg/scrape"

type Exporter interface {
	Export(page scrape.Page)
}
