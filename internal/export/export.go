// Package export handles interface specification for
// wiki-agnostic exporters for page instances created by the
// scrape package
package export

import "github.com/mal0ner/wikiscrape/internal/scrape"

type Exporter interface {
	Export(page scrape.Page)
}
