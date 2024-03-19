package util

import (
	"fmt"
	"net/url"
	"strings"
)

// wikiInfo represents the basic information needed
// to add support for querying a wiki's API.
type wikiInfo struct {
	Name           string
	APIPath        string
	PagePathPrefix string
	Backend        string
}

// QueryData represents all the information the
// internal/wiki and internal/scraper packages need
// to make queries to their respective APIs.
type QueryData struct {
	Info *wikiInfo
	Page string
}

// newWikiInfo initializes a new wikiInfo object, which represents the basic information
// needed to add support for querying a wiki's API.
func newWikiInfo(name string, apiPath string, pagePrefix string, backend string) *wikiInfo {
	return &wikiInfo{name, apiPath, pagePrefix, backend}
}

// Map supported wiki names to relevant query info
var wikiNameInfo = map[string]*wikiInfo{
	"wikipedia": newWikiInfo("Wikipedia", "https://en.wikipedia.org/w/api.php", "/wiki/", "mediawiki"),
	"osrs":      newWikiInfo("Old School Runescape", "https://oldschool.runescape.wiki/api.php", "/w/", "mediawiki"),
}

// Map supported wiki hosts to relevant query info
var wikiHostInfo = map[string]*wikiInfo{
	"en.wikipedia.org":         newWikiInfo("Wikipedia", "https://en.wikipedia.org/w/api.php", "/wiki/", "mediawiki"),
	"oldschool.runescape.wiki": newWikiInfo("Old School Runescape", "https://oldschool.runescape.wiki/api.php", "/w/", "mediawiki"),
}

var supportedBackends = []string{
	"mediawiki",
}

// Custom error designed indicate to the user that the
// wiki they are trying to escape is not explicitly supported
// by the program.
type WikiNotSupportedError struct {
	Code string
	Info string
}

// Error returns a formatted WikiNotSupportedError including code
// and additional information.
func (e *WikiNotSupportedError) Error() string {
	return fmt.Sprintf("WikiNotSupportedError: [code] %s [info] %s", e.Code, e.Info)
}

// GetQueryDataFromURL accepts a raw wiki url, parses it's host name, checks for explicit
// support (existing wikiInfo entry in the wikiHostInfo map) and returns a QueryData
// object which provides all the necessary information to make a query to the wiki's api.
func GetQueryDataFromURL(rawURL string) (*QueryData, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	if info, ok := wikiHostInfo[parsedURL.Host]; ok {
		pageName, err := getPageNameFromPath(parsedURL.Path, info.PagePathPrefix)
		if err != nil {
			return nil, err
		}
		return &QueryData{
			Page: pageName,
			Info: info,
		}, nil
	}
	return nil, &WikiNotSupportedError{
		Code: "hostnotfound",
		Info: "The provided host is not yet supported (unknown api endpoint or page prefix)",
	}
}

// GetQueryDataFromName accepts a page name and wikiname, checks for explicit support
// (existing wikiInfo entry in the wikiNameInfo map) and returns a QueryData object
// which provides all the necessary information to make a query to the wiki's api.
func GetQueryDataFromName(pageName string, wikiName string) (*QueryData, error) {
	if info, ok := wikiNameInfo[wikiName]; ok {
		return &QueryData{
			Page: pageName,
			Info: info,
		}, nil
	}
	return nil, &WikiNotSupportedError{
		Code: "namenotfound",
		Info: "The provided name is not yet supported (unknown api endpoint or page prefix)",
	}
}

// getPageNameFromPath strips a prefix from the beginning of a string. In this
// use case, it is designed to take a url.URL.path from a parsed wiki URL and
// remove the page prefix so as to return the full page name. This function
func getPageNameFromPath(path string, prefix string) (string, error) {
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("page path does not start with expected prefix")
	}
	return path[len(prefix):], nil
}

// GetSupportedWikis returns a list of the names of all wikis supported by
// wikiscrape
func GetSupportedWikis() []string {
	keys := make([]string, len(wikiNameInfo))
	i := 0
	for k := range wikiNameInfo {
		keys[i] = k
		i++
	}
	return keys
}

// GetSupportedBackends returns a list of the names of all wiki backends
// supported by wikiscrape.
func GetSupportedBackends() []string {
	return supportedBackends
}

// GetWikiInfoStrings returns a formatted slice of strings containing relevant
// information about the wikis supported by wikiscrape. This function is designed for
// use with the List command.
func GetWikiInfoStrings(backendFilter string) []string {
	items := make([]string, len(wikiNameInfo))
	i := 0
	for k, v := range wikiNameInfo {
		if v.Backend == backendFilter || backendFilter == "" {
			items[i] = fmt.Sprintf("%s: [alias: %s, backend: %s]", v.Name, k, v.Backend)
		}
		i++
	}
	return items
}

// GetWikiInfoFromHost takes a URL host segment and returns its corresponding wikiInfo.
// Fails if the wiki is not supported by wikiscrape.
func GetWikiInfoFromHost(host string) (*wikiInfo, error) {
	info, ok := wikiHostInfo[host]
	if !ok {
		return nil, &WikiNotSupportedError{
			Code: "hostnotfound",
			Info: "The provided host if not yet supported (uknown api endpoint or page prefix)",
		}
	}
	return info, nil
}

// GetWikiInfoFromName takes a wiki name and returns its corresponding wikiInfo.
// Fails if the wiki is not supported by wikiscrape.
func GetWikiInfoFromName(name string) (*wikiInfo, error) {
	info, ok := wikiNameInfo[name]
	if !ok {
		return nil, &WikiNotSupportedError{
			Code: "namenotfound",
			Info: "The provided name is not yet supported (unknown api endpoint or page prefix)",
		}
	}
	return info, nil
}
