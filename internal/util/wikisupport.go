package util

import (
	"fmt"
	"net/url"
	"strings"
)

type wikiInfo struct {
	Name           string
	Hostname       string
	APIPath        string
	PagePathPrefix string
}

func newWikiInfo(name string, hostname string, endpoint string, pagePrefix string) *wikiInfo {
	return &wikiInfo{name, hostname, endpoint, pagePrefix}
}

// MAP SUPPORTED WIKI NAMES TO RELEVANT QUERY INFO
var supportedWikiNames = map[string]*wikiInfo{
	"wikipedia": newWikiInfo("wikipedia", "en.wikipedia.org", "/w/api.php", "/wiki/"),
	"osrs":      newWikiInfo("osrs", "oldschool.runescape.wiki", "api.php", "/w/"),
}

// MAP SUPPORTED WIKI HOSTS TO RELEVANT QUERY INFO
var supportedWikiHosts = map[string]*wikiInfo{
	"en.wikipedia.org":         newWikiInfo("wikipedia", "en.wikipedia.org", "/w/api.php", "/wiki/"),
	"oldschool.runescape.wiki": newWikiInfo("osrs", "oldschool.runescape.wiki", "api.php", "/w/"),
}

type WikiNotSupportedError struct {
	Code string
	Info string
}

func (e *WikiNotSupportedError) Error() string {
	return fmt.Sprintf("WikiNotSupportedError: [code] %s [info] %s", e.Code, e.Info)
}

func GetWikiInfoFromURL(rawURL string) (string, string, error) {
	// ParseRequestURI is much more informative than url.Parse on whether
	// or not a url is actually valid.
	//    url.Parse("banana") -> NO PROBLEM MATE!
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", "", err
	}
	if info, ok := supportedWikiHosts[parsedURL.Host]; ok {
		queryUrl, err := buildQueryURL(parsedURL.Host, info.APIPath)
		if err != nil {
			return "", "", err
		}
		pageName, err := getPageNameFromPath(parsedURL.Path, info.PagePathPrefix)
		if err != nil {
			return "", "", err
		}
		return queryUrl, pageName, nil
	}
	return "", "", &WikiNotSupportedError{
		Code: "hostnotfound",
		Info: "The provided host is not yet supported (unknown api endpoint or page prefix)",
	}
}

func GetWikiInfoFromName(name string) (string, error) {
	if info, ok := supportedWikiNames[name]; ok {
		queryURL, err := buildQueryURL(info.Hostname, info.APIPath)
		if err != nil {
			return "", err
		}
		return queryURL, nil
	}
	return "", &WikiNotSupportedError{
		Code: "namenotfound",
		Info: "The provided wiki name is not yet supported",
	}
}

func buildQueryURL(host string, endpoint string) (string, error) {
	res := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   endpoint,
	}
	return res.String(), nil
}

func getPageNameFromPath(path string, prefix string) (string, error) {
	if !strings.HasPrefix(path, prefix) {
		return "", fmt.Errorf("page path does not start with expected prefix")
	}
	return path[len(prefix):], nil
}

func GetSupportedWikis() []string {
	keys := make([]string, len(supportedWikiNames))
	i := 0
	for k := range supportedWikiNames {
		keys[i] = k
		i++
	}
	return keys
}
