package wiki

type Wiki interface {
	Crawl() error
	Export() error
}
