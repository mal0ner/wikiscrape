package list

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

// Long message
var backendsMsg = "List all backends currently supported by wikiscrape. A \"Backend\" here refers to the framework used to create the wiki site, as different providers will each have their own unique patterns for accessing and querying their APIs (if one exists at all) and so must be supported explicitly and on an individual basis."

// Command
var backendsCmd = &cobra.Command{
	Use:   "backends",
	Short: "list supported backends",
	Long:  backendsMsg,
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		backends := util.GetSupportedBackends()
		for _, b := range backends {
			fmt.Println(b)
		}
	},
}
