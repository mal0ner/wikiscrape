package scrape

import (
	"net/http"

	"github.com/mal0ner/wikiscrape/pkg/manifest"
)

type Page struct {
	Title    string
	Url      string
	Sections []Section
}

type Section struct {
	Heading string
	Index   int
	Content string
}

type Parser interface {
	ParsePageResponse(res *http.Response) (Page, error)
	ParseSectionResponse(res *http.Response) (Section, error)
	ParseContent(content string) string
}

type Scraper interface {
	Scrape(manifest manifest.Manifest) []Page
	getPage(path string) (Page, error)
}
