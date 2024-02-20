package scrape

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/mal0ner/wikiscrape/internal/util"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Wraps methods for retrieving and parsing pages on
// a MediaWiki based website.
type MediaWikiScraper struct {
	BaseURL string
}

// Representation of the json response returned
// by making a request for the raw HTML of a MediaWiki page
type mediaWikiPageResponse struct {
	Parse struct {
		Title string `json:"title"`
		Text  struct {
			Value string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
	Error *MediaWikiAPIError `json:"error"`
}

// MediaWiki API error format, for a list of error codes
// and their associated information see:
// http://tinyurl.com/mwerrorcodes
type MediaWikiAPIError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

// Returns formatted MediaWiki error including code and
// additional information
func (e *MediaWikiAPIError) Error() string {
	return fmt.Sprintf("MediaWiki API error: [code] %s [info] %s", e.Code, e.Info)
}

// Build MediaWiki page request url with escaped parameters
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

// Makes a http request to the MediaWiki API endpoint
// for the page specified by the path, then parses and
// returns the response as a mediaWikiPageResponse
// Can return a MediaWikiAPIError if (for example):
//   - The page does not exist
//   - The user is denied read access to the page
//   - The user has been rate-limited and should try again
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

// Fetches a page specified by path, parses its sections and returns
// a GetPage{} instance
func (s *MediaWikiScraper) GetPage(path string) (*Page, error) {
	response, err := s.fetchPage(path)
	if err != nil {
		return &Page{}, err
	}
	sections, err := response.ParseSections()
	if err != nil {
		return &Page{}, err
	}
	return &Page{
		Title:    response.Parse.Title,
		Sections: sections,
	}, nil
}

// Searches for a specific section of the specified page given a heading.
// parses and returns the section if found, otherwise returns a mediaWikiAPIError
func (s *MediaWikiScraper) GetSection(path string, heading string) (*Page, error) {
	response, err := s.fetchPage(path)
	if err != nil {
		return &Page{}, err
	}
	section, err := response.ParseSection(heading)
	if err != nil {
		return &Page{}, err
	}
	return &Page{
		Title:    response.Parse.Title,
		Sections: []Section{section},
	}, nil
}

// Parses raw HTML from mediaWikiPageResponse into an array of
// Section{}, including headings and body text
// TODO: Add table parsing support
func (response *mediaWikiPageResponse) ParseSections() ([]Section, error) {
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

func (response *mediaWikiPageResponse) ParseSection(heading string) (Section, error) {
	var section Section
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Parse.Text.Value))
	if err != nil {
		return Section{}, err
	}
	found := false
	doc.Find("h2 > span.mw-headline").EachWithBreak(func(i int, h2 *goquery.Selection) bool {
		if util.TrimLower(h2.Text()) == util.TrimLower(heading) {
			section.Heading = h2.Text()
			section.Index = i
			var contentBuilder strings.Builder
			h2.Parent().NextFilteredUntil("p", "h2 > span.mw-headline").Each(func(i int, s *goquery.Selection) {
				contentBuilder.WriteString(s.Text())
			})
			section.Content = contentBuilder.String()
			found = true
			return false
		}
		return true
	})

	if !found {
		return Section{}, &MediaWikiAPIError{
			Code: "sectionnotfound",
			Info: fmt.Sprintf("The section with heading %s was not found on page %s", heading, response.Parse.Title),
		}
	}
	return section, nil
}
