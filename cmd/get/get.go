package get

import (
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get [page|section]",
	Short: "Get a page or section from a wiki",
	Long:  "Get a page or section from a wiki",
}

func init() {
	GetCmd.AddCommand(pageCmd)
	GetCmd.AddCommand(sectionCmd)
}
