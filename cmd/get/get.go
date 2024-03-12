package get

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/mal0ner/wikiscrape/internal/wiki"
	"github.com/spf13/cobra"
)

// Long message
var getMsg = "Get and export page/s or sections of pages from a wiki. Subcommands are available both for retrieval of a single page given a page name and wiki provider, or a list of pages given a path to a manifest file.\n\nYou can also provide a URL directly to the get command to scrape a whole page directly. No need to name the provider or page name; if the wiki is supported, it will just work!\nUsage: wikiscrape get <URL>.\n\nFor a list of supported wikis and export formats, please see \"wikiscrape list -h\"."

// Flag vars
var section string

// Command
var GetCmd = &cobra.Command{
	Use: "get [url|page|section]",
	// SilenceUsage: true,
	Short: "Get a page or section from a wiki",
	Long:  getMsg,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return getPageFromURL(args[0], section)
	},
}

func init() {
	GetCmd.AddCommand(pageCmd)
	GetCmd.PersistentFlags().StringVarP(&section, "section", "s", "", "section heading you wish to scrape")
}

// getWikiFromQueryData identifies and returns the appropriate wiki scraper for the wiki provider listed
// in the generated queryData from a 'get' command request. Returns a WikiNotSupportedError in the case
// that the provider is not explicitly supported.
func getWikiFromQueryData(queryData *util.QueryData) (wiki.Wiki, *util.WikiNotSupportedError) {
	backend := util.TrimLower(queryData.Info.Backend)
	switch backend {
	case "mediawiki":
		return wiki.NewMediaWiki(backend, queryData.Info.APIPath), nil
	}
	return nil, &util.WikiNotSupportedError{
		Code: "backendnotsupported",
		Info: fmt.Sprintf("The detected backend %s is not yet a supported wiki provider", backend),
	}
}

// getPageFromURL parses the provided URL to check for explicit support for the wiki's backend provider before
// initializing the appropriate scraper and retrieving the requested page. Returns an error indicating the success
// of page retrieval.
func getPageFromURL(rawURL string, section string) error {
	queryData, err := util.GetQueryDataFromURL(rawURL)
	if err != nil {
		return err
	}
	w, err := getWikiFromQueryData(queryData)
	if err != nil {
		return err
	}

	if section != "" {
		return w.Section(queryData.Page, section)
	}
	return w.Page(queryData.Page)
}

// getPageFromName checks wikiscrape support for the provided wikiName. If supported, the appropriate
// scraper is initialized and the page with the provided name is scraped. Returns an error indicating
// the success of page retrieval.
func getPageFromName(pageName string, wikiName string, section string) error {
	// TODO: Maybe make the section into an array and support it in the command line apps later,
	// could give users opportunity to ask for multiple sections. would just need to loop over them
	// in the wiki.GetSection() function.
	queryData, err := util.GetQueryDataFromName(pageName, wikiName)
	if err != nil {
		return err
	}
	w, err := getWikiFromQueryData(queryData)
	if err != nil {
		return err
	}
	if section != "" {
		return w.Section(queryData.Page, section)
	}
	return w.Page(queryData.Page)
}
