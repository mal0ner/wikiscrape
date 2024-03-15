// package scrape handles interface specifications and concrete wiki-specific
// implementations for the scraping and parsing of the content from pages served by various
// Wiki frameworks.
//
// Currently supported API backends: see 'wikiscrape list backends'
package scrape

// Page represents a wiki/backend agnostic container for storing the content
// of a wiki page.
type Page struct {
	Title    string
	Sections []*Section
}

// Section represents a wiki/backend agnostic container for storing the contents
// of a single section of a wiki page.
type Section struct {
	Heading string
	Index   int
	Content string
}

// Response denotes the methods one should implement on the API
// response struct for a specific wiki.
type Response interface {
	ParseSections() ([]Section, error)
	ParseSection(heading string) (Section, error)
}

// Scraper denotes the methods one should implement on the scraper
// struct for a specific wiki in order to handle requests.
type Scraper interface {
	GetPage(path string) (*Page, error)
	GetSection(path string, heading string) (*Page, error)
}
