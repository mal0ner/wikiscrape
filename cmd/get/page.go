package get

import (
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

func init() {
	pageCmd.Flags().StringVarP(&wikiName, "wiki", "w", "", "Name of the wiki you wish to scrape")
	pageCmd.MarkFlagRequired("wiki")
}
