package scrape

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

type Response interface {
	ParseSections() ([]Section, error)
	ParseSection(heading string) (Section, error)
}

type Scraper interface {
	GetPage(path string) (*Page, error)
	GetSection(path string, heading string) (*Page, error)
}
