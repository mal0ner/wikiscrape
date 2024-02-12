package util

import (
	"net/url"
	"strconv"

	logger "github.com/mal0ner/wikiscrape/pkg/logging"
	"github.com/sirupsen/logrus"
)

// TODO: Add options for the query params here?
type MediaWikiQueryBuilder struct{}

// Build the parameter values for a mediaWiki API page request
func (MediaWikiQueryBuilder) BuildPageQuery(path string) (string, error) {
	path, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("prop", "sections")
	params.Set("page", path)
	params.Set("format", "json")
	logger.Log.WithFields(logrus.Fields{
		"page": path,
	}).Debug("generated page sections query")
	return params.Encode(), nil
}

// Build the parameter values for a mediaWiki API page request
func (MediaWikiQueryBuilder) BuildSectionQuery(path string, section int) (string, error) {
	path, err := url.QueryUnescape(path)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("prop", "text")
	params.Set("page", path)
	params.Set("format", "json")
	params.Set("section", strconv.Itoa(section))
	logger.Log.WithFields(logrus.Fields{
		"page":    path,
		"section": section,
	}).Debug("generated page section query")
	return params.Encode(), nil
}
