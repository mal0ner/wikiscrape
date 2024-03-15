package get

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

// Long message
var pagesMsg = "Get and export a list of pages whose names are defined in a 'manifest' file. The persistent 'section' flag available to the get command and all its subcommands allows for retrieving specific sections from all pages in the manifest."

// Flag vars
var manFile string

// Command
var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Get all pages listed in a file",
	Long:  pagesMsg,
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		pageNames, err := util.ReadManifestFrom(manFile)
		if err != nil {
			return err
		}
		queryData, err := util.GetQueryDataFromName(pageNames[0], wikiName)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		w, err := getWikiFromQueryData(queryData)
		if err != nil {
			return err
		}
		err = w.ScrapeAndExport(pageNames)
		return nil
	},
}

func init() {
	pagesCmd.Flags().StringVarP(&manFile, "from-manifest", "f", "", "path to the manifest file")
	pagesCmd.Flags().StringVarP(&wikiName, "wiki", "w", "", "name of the wiki you wish to scrape")
	pagesCmd.MarkPersistentFlagRequired("wiki")
	pagesCmd.MarkFlagRequired("from-manifest")
	pagesCmd.MarkFlagRequired("wiki")
	pagesCmd.MarkFlagFilename("from-manifest", "json")
}
