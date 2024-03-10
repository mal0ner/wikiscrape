package util

import (
	"fmt"
	"net/url"
	"strings"
)

type wikiInfo struct {
	APIPath        string
	PagePathPrefix string
	Backend        string
}

type queryData struct {
	Info *wikiInfo
	Page string
}

func newWikiInfo(apiPath string, pagePrefix string, backend string) *wikiInfo {
	return &wikiInfo{apiPath, pagePrefix, backend}
}

// MAP SUPPORTED WIKI NAMES TO RELEVANT QUERY INFO
var supportedWikiNames = map[string]*wikiInfo{
	"wikipedia": newWikiInfo("https://en.wikipedia.org/w/api.php", "/wiki/", "mediawiki"),
	"osrs":      newWikiInfo("https://oldschool.runescape.wiki/api.php", "/w/", "mediawiki"),
}

// MAP SUPPORTED WIKI HOSTS TO RELEVANT QUERY INFO
var supportedWikiHosts = map[string]*wikiInfo{
	"en.wikipedia.org":         newWikiInfo("https://en.wikipedia.org/w/api.php", "/wiki/", "mediawiki"),
	"oldschool.runescape.wiki": newWikiInfo("https://oldschool.runescape.wiki/api.php", "/w/", "mediawiki"),
}

type wikiNotSupportedError struct {
	Code string
	Info string
}

func (e *wikiNotSupportedError) Error() string {
	return fmt.Sprintf("WikiNotSupportedError: [code] %s [info] %s", e.Code, e.Info)
}

func GetQueryDataFromURL(rawURL string) (*queryData, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	if info, ok := supportedWikiHosts[parsedURL.Host]; ok {
		pageName, err := getPageNameFromPath(parsedURL.Path, info.PagePathPrefix)
		if err != nil {
			return nil, err
		}
		return &queryData{
			Page: pageName,
			Info: info,
		}, nil
	}
	return nil, &wikiNotSupportedError{
		Code: "hostnotfound",
		Info: "The provided host is not yet supported (unknown api endpoint or page prefix)",
	}
}

func GetQueryDataFromName(pageName string, wikiName string) (*queryData, error) {
	if info, ok := supportedWikiNames[wikiName]; ok {
		return &queryData{
			Page: pageName,
			Info: info,
		}, nil
	}
	return nil, &wikiNotSupportedError{
		Code: "namenotfound",
		Info: "The provided name is not yet supported (unknown api endpoint or page prefix)",
	}
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
