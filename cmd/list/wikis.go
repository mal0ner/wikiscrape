package list

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

var longMsg = "List by name all supported wikis in the current version of wikiscrape. The items outputted by this command are designed as valid input for the \"--wiki\" flag.\n\nTo say a wiki is \"supported\" means that we have a record of the website's api query path, backend framework, and page prefix stored in internal/util/wikisupport.go AND that the backend is explicitly supported in internal/wiki"

var wikisCmd = &cobra.Command{
	Use:   "wikis",
	Short: "List supported wikis by name",
	Long:  longMsg,
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		wikis := util.GetSupportedWikis()
		for _, w := range wikis {
			fmt.Println(w)
		}
	},
}
