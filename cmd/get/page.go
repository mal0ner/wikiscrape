package get

import (
	"github.com/spf13/cobra"
)

// Long message
var pageMsg = "Get and export a single page from a wiki, given the page name and wiki provider.\n\nFor a list of supported wikis and export formats, please see \"wikiscrape list -h\".\n"

// Flag vars
var wikiName string

// Command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Get a single page",
	Long:  "Get a single page",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		pageName := args[0]
		return getPageFromName(pageName, wikiName, section)
	},
}

func init() {
	pageCmd.Flags().StringVarP(&wikiName, "wiki", "w", "", "name of the wiki you wish to scrape")
	pageCmd.MarkFlagRequired("wiki")
}
