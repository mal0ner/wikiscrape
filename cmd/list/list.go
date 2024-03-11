package list

import (
	"github.com/spf13/cobra"
)

var listMsg = "The List command is a wrapper for subcommands \"list wikis\" and \"list backends\"."

var ListCmd = &cobra.Command{
	Use:   "list [wikis|backends]",
	Short: "List things",
	Long:  listMsg,
}

func init() {
	ListCmd.AddCommand(wikisCmd)
	ListCmd.AddCommand(backendsCmd)
}
