package get

import (
	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Get a single page",
	Long:  "Get a single page",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		pageName := args[0]
		err := getPageFromName(pageName, wikiName)
		if err != nil {
			return err
		}
		return nil
	},
}

var wikiName string

// TODO: Maybe make the default wiki come from config or something?
func init() {
	pageCmd.Flags().StringVarP(&wikiName, "wiki", "w", "", "Name of the wiki you wish to scrape")
	pageCmd.MarkFlagRequired("wiki")
}

func getPageFromName(pageName string, wikiName string) error {
	queryData, err := util.GetQueryDataFromName(pageName, wikiName)
	if err != nil {
		return err
	}
	w, err := getScraperFromQueryData(queryData)
	if err != nil {
		return err
	}
	err = w.Page(queryData.Page)
	if err != nil {
		return err
	}
	return nil
}
