package export

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/pkg/scrape"
)

type TestExporter struct{}

func (te *TestExporter) Export(page scrape.Page) {
	fmt.Println("Title: " + page.Title)
	for _, s := range page.Sections {
		fmt.Println("Section: " + s.Heading + "\n")
		fmt.Println(s.Content)
	}
}
