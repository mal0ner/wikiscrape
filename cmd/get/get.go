package get

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/mal0ner/wikiscrape/internal/wiki"
	"github.com/spf13/cobra"
)

var longMsg = "Get a page or section from a wiki. Subcommands are available both for retrieval of an entire page given a name and wiki provider, or a single section given a specified heading on top of the previous information.\n\nYou can also provide a URL directly to the get command to scrape a whole page directly. No need to name the provider or page name; if the wiki is supported, it will just work!\nUsage: wikiscrape get <URL>.\n\nFor a list of supported wikis please see \"wikiscrape list\"."

var GetCmd = &cobra.Command{
	Use: "get [url|page|section]",
	// SilenceUsage: true,
	Short: "Get a page or section from a wiki",
	Long:  longMsg,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		err := getPageFromURL(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	GetCmd.AddCommand(pageCmd)
	GetCmd.AddCommand(sectionCmd)
}

func getScraperFromQueryData(queryData *util.QueryData) (wiki.Wiki, error) {
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

func getPageFromURL(rawURL string) error {
	queryData, err := util.GetQueryDataFromURL(rawURL)
	if err != nil {
		return err
	}
	w, err := getScraperFromQueryData(queryData)
	if err != nil {
		return err
	}
	err = w.Page(queryData.Page)
	return nil
}
