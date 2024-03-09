// package wikiscrape/scrape handles interface specifications and concrete wiki-specific
// implementations for the scraping and parsing of the content from pages served by various
// Wiki frameworks.
//
// Current supported API backends:
//   - MediaWiki
package scrape

type Page struct {
	Title    string
	Sections []Section
}

type Section struct {
	Heading string
	Index   int
	Content string
}

type Response interface {
	ParseSections() ([]Section, error)
	ParseSection(heading string) (Section, error)
}

type Scraper interface {
	GetPage(path string) (*Page, error)
	GetSection(path string, heading string) (*Page, error)
}
