package list

import (
	"fmt"

	"github.com/mal0ner/wikiscrape/internal/util"
	"github.com/spf13/cobra"
)

var backendsMsg = "List all backends currently supported by wikiscrape. A \"Backend\" here refers to the framework used to create the wiki site, as different providers will each have their own unique patterns for accessing and querying their APIs (if one exists at all) and so must be supported explicitly and on an individual basis."

var backendsCmd = &cobra.Command{
	Use:   "backends",
	Short: "List supported backends by name",
	Long:  backendsMsg,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		backends := util.GetSupportedBackends()
		for _, b := range backends {
			fmt.Println(b)
		}
	},
}
