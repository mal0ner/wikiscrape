package scrape

import (
	"bytes"
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

// Error returns a formatted MediaWiki error including code and
// additional information
func (e *MediaWikiAPIError) Error() string {
	return fmt.Sprintf("MediaWiki API error: [code] %s [info] %s", e.Code, e.Info)
}

// pageQuery builds an encoded Media Wiki page request url given the path of a page.
//
// In this case,
// 'path' refers to the segment of a Media Wiki url which holds the unique page
// name/title. This can be found using the url.ParseRequestURI assuming pre-existing
// knowledge of the way in which the wiki prefixes their page names in urls.
// The ParseRequestURI method returns:
//
//	Input: "https://en.wikipedia.org/wiki/Bear"
//	Output:
//	  url.Scheme: "https"
//	  url.Host:   "en.wikipedia.org"
//	  url.Path:   "wiki/Bear"
//
// With a known page prefix: "wiki/", mapped to by the host name, we can simply
// strip this from the path and receive the page name.
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

// fetchPage makes a http request to the MediaWiki API endpoint
// for the page specified by the path, then unmarshals the response.
// Returns a mediaWikiPageResponse.
// Can return a MediaWikiAPIError if (for example):
//   - The page does not exist
//   - The user is denied read access to the page
//   - The user has been rate-limited and should try again
func (s *MediaWikiScraper) fetchPage(path string) (*mediaWikiPageResponse, error) {
	var result mediaWikiPageResponse
	url, err := s.pageQuery(path)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	// body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(&buf).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &result, nil
}

// GetPage first fetches the page specified by path before parsing
// its content.
// Returns a Page containing the content parsed from the HTML response.
//
// Can error when:
//   - page fetch fails.
//   - section parsing fails
func (s *MediaWikiScraper) GetPage(path string) (*Page, error) {
	response, err := s.fetchPage(path)
	if err != nil {
		return nil, err
	}
	sections, err := response.ParseSections()
	if err != nil {
		return nil, err
	}
	return &Page{
		Title:    response.Parse.Title,
		Sections: sections,
	}, nil
}

// GetSection searches for a section of a page by heading, and returns it if found.
func (s *MediaWikiScraper) GetSection(path string, heading string) (*Page, error) {
	response, err := s.fetchPage(path)
	if err != nil {
		return nil, err
	}
	section, err := response.ParseSection(heading)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title:    response.Parse.Title,
		Sections: []*Section{section},
	}, nil
}

// ParseSections parses raw HTML from mediaWikiPageResponse.
// Returns an array of Sections containing headlines and body text.
//
// Can error when:
//   - The content in the response is not valid HTML
//
// TODO: Add table parsing support
func (response *mediaWikiPageResponse) ParseSections() ([]*Section, error) {
	var sections []*Section
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(response.Parse.Text.Value),
	)
	if err != nil {
		return nil, err
	}
	// Create first section manually because its header is not included in the
	// ".mw-parser-output"
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
	sections = append(sections, &intro)
	// get remaining sections
	doc.Find("h2 > span.mw-headline").Each(func(i int, h2 *goquery.Selection) {
		title := h2.Text()
		var contentBuilder strings.Builder
		h2.Parent().NextFilteredUntil("p", "h2 > span.mw-headline").Each(func(i int, s *goquery.Selection) {
			contentBuilder.WriteString(s.Text())
		})
		sections = append(sections, &Section{
			Heading: title,
			Index:   i + 1,
			Content: contentBuilder.String(),
		})
	})
	return sections, err
}

// ParseSection parses the raw HTML of a mediaWikiPageResponse and searches for a section
// that contains the heading specified by the function argument.
// Returns a Section containing the heading and its corresponding body text.
//
// Can error when:
//   - The content in the response is not valid HTML
//   - The heading is not found.
func (response *mediaWikiPageResponse) ParseSection(heading string) (*Section, error) {
	var section Section
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Parse.Text.Value))
	if err != nil {
		return nil, err
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
		return nil, &MediaWikiAPIError{
			Code: "sectionnotfound",
			Info: fmt.Sprintf("The section with heading %s was not found on page %s", heading, response.Parse.Title),
		}
	}
	return &section, nil
}
