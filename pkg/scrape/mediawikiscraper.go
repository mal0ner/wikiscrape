package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mal0ner/wikiscrape/internal/util"
	logger "github.com/mal0ner/wikiscrape/pkg/logging"
	"github.com/mal0ner/wikiscrape/pkg/manifest"
)

// TODO: CLEAN UP LOGGING IN THIS FILE ITS DISGUSTING
// INFO: Maybe embed the parser struct inside?
// This seems to be ok for now
type MediaWikiScraper struct {
	BaseURL string
	MediaWikiParser
	util.MediaWikiQueryBuilder
}

type MediaWikiParser struct{}

// TODO: http request error handling: page does not exist, see postman
// do we just add in an error struct alongside the parse one and check if it's null?
type mediaWikiPageResponse struct {
	Parse struct {
		Title    string `json:"title"`
		Sections []struct {
			Heading   string `json:"line"`
			Index     string `json:"index"`
			FromTitle string `json:"fromtitle"`
		} `json:"sections"`
	} `json:"parse"`
	Error *MediaWikiAPIError `json:"error"`
}

type mediaWikiSectionResponse struct {
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

func (*MediaWikiParser) ParsePageResponse(res *http.Response) (mediaWikiPageResponse, error) {
	var pageResponse mediaWikiPageResponse

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}

	err = json.Unmarshal(body, &pageResponse)
	if err != nil {
		return mediaWikiPageResponse{}, err
	}

	if pageResponse.Error != nil {
		return mediaWikiPageResponse{}, pageResponse.Error
	}

	return pageResponse, nil
}

func (*MediaWikiParser) ParseSectionResponse(res *http.Response) (mediaWikiSectionResponse, error) {
	var sectionResponse mediaWikiSectionResponse

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mediaWikiSectionResponse{}, err
	}
	err = json.Unmarshal(body, &sectionResponse)
	if err != nil {
		return mediaWikiSectionResponse{}, err
	}
	if pageResponse.Error != nil {
		return mediaWikiSectionResponse{}, sectionResponse.Error
	}
	return sectionResponse, nil
}

func (scraper *MediaWikiScraper) getPage(path string) (Page, error) {
	logger.Log.WithField("page", path).Debug("scraping")
	query, err := scraper.BuildPageQuery(path)
	if err != nil {
		return Page{}, err
	}
	reqURL := scraper.BaseURL + "?" + query

	logger.Log.WithField("url", reqURL).Debug("requesting page sections")
	res, err := http.Get(reqURL)
	if err != nil {
		return Page{}, err
	}
	pageResponse, err := scraper.ParsePageResponse(res)
	if err != nil {
		logger.Log.Warn("FAILED TO PARSE PAGE RESPONSE")
		return Page{}, err
	}
	page := Page{
		Title:    pageResponse.Parse.Title,
		Url:      reqURL,
		Sections: []Section{},
	}

	for i, s := range pageResponse.Parse.Sections {
		logger.Log.WithField("index", i).Debug("request section")
		query, err := scraper.BuildSectionQuery(path, i)
		if err != nil {
			logger.Log.Warn("FAILED TO BUILD SECTION QUERY")
			return Page{}, err
		}
		reqURL := scraper.BaseURL + "?" + query
		res, err := http.Get(reqURL)
		if err != nil {
			logger.Log.Warn("FAILED TO GET PAGE SECTION")
			return Page{}, err
		}
		sectionResponse, err := scraper.ParseSectionResponse(res)
		if err != nil {
			logger.Log.Warn("FAILED TO PARSE SECTION RESPONSE")
			return Page{}, err
		}
		content, err := sectionResponse.parseContent()
		if err != nil {
			logger.Log.Warn("FAILED TO PARSE SECTION CONTENT")
			return Page{}, err
		}
		result := Section{
			Heading: s.Heading,
			Content: content,
			Index:   i,
		}
		page.Sections = append(page.Sections, result)
	}
	return page, nil
}

func (response mediaWikiSectionResponse) parseContent() (string, error) {
	var content string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Parse.Text.Value))
	if err != nil {
		return "", err
	}
	doc.Find("p").Each(func(_ int, selection *goquery.Selection) {
		content = content + selection.Text()
	})
	return content, nil
}

func (scraper *MediaWikiScraper) Scrape(manifest manifest.Manifest) ([]Page, error) {
	var pages []Page
	for _, path := range manifest {
		page, err := scraper.getPage(path)
		if err != nil {
			return []Page{}, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}
