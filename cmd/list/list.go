package list

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

// Long message
var listMsg = "List name and other relevant information about wikis that are supported by wikiscrape"

// Flag vars
var backendFilter string

// Command
var ListCmd = &cobra.Command{
	Use:   "list [-b backend]",
	Short: "List supported wiki information",
	Args:  cobra.NoArgs,
	Long:  listMsg,
	Run: func(_ *cobra.Command, _ []string) {
		wikiStrings := util.GetWikiInfoStrings(backendFilter)
		for _, w := range wikiStrings {
			fmt.Println(w)
		}
	},
}

func init() {
	ListCmd.Flags().StringVarP(&backendFilter, "filter-backend", "b", "", "Name of backend you wish to filter by")
}
