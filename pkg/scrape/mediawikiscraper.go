package scrape

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/mal0ner/wikiscrape/pkg/manifest"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MediaWikiScraper struct {
	BaseURL string
}

type mediaWikiPageResponse struct {
	Parse struct {
		Title string `json:"title"`
		Text  struct {
			Value string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
	Error *MediaWikiAPIError `json:"error"`
}

type MediaWikiAPIError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

func (e *MediaWikiAPIError) Error() string {
	return fmt.Sprintf("MediaWiki API error: [code] %s [info] %s", e.Code, e.Info)
}

func (s *MediaWikiScraper) pageQuery(path string) (string, error) {
	path, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("format", "json")
	params.Set("page", path)
	return s.BaseURL + "?" + params.Encode(), nil
}

func (s *MediaWikiScraper) fetchPage(path string) (mediaWikiPageResponse, error) {
	var result mediaWikiPageResponse
	url, err := s.pageQuery(path)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}
	res, err := http.Get(url)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}

	if result.Error != nil {
		return mediaWikiPageResponse{}, result.Error
	}

	return result, nil
}

func (s *MediaWikiScraper) getPage(path string) (Page, error) {
	response, err := s.fetchPage(path)
	if err != nil {
		return Page{}, err
	}
	sections, err := response.parseSections()
	if err != nil {
		return Page{}, err
	}
	return Page{
		Title:    response.Parse.Title,
		Url:      "test",
		Sections: sections,
	}, nil
}

func (response *mediaWikiPageResponse) parseSections() ([]Section, error) {
	var sections []Section
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Parse.Text.Value))
	if err != nil {
		return []Section{}, err
	}
	// Create first section manually because its header is not included in the ".mw-parser-output"
	// seems kinda sketch TODO: FIX???
	var introBuilder strings.Builder
	doc.Find("*").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Is("h2") && s.Find("span.mw-headline").Length() > 0 {
			return false
		}
		if s.Is("p") {
			introBuilder.WriteString(s.Text())
		}
		return true
	})
	intro := Section{
		Heading: "Introduction",
		Index:   0,
		Content: introBuilder.String(),
	}
	sections = append(sections, intro)
	// get remaining sections
	doc.Find("h2 > span.mw-headline").Each(func(i int, h2 *goquery.Selection) {
		title := h2.Text()
		var contentBuilder strings.Builder
		h2.Parent().NextFilteredUntil("p", "h2 > span.mw-headline").Each(func(i int, s *goquery.Selection) {
			contentBuilder.WriteString(s.Text())
		})
		sections = append(sections, Section{
			Heading: title,
			Index:   i,
			Content: contentBuilder.String(),
		})
	})
	return sections, err
}

func (s *MediaWikiScraper) Scrape(manifest manifest.Manifest) ([]Page, error) {
	var pages []Page
	for _, path := range manifest {
		page, err := s.getPage(path)
		if err != nil {
			return []Page{}, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}
