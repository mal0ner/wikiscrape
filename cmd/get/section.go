package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sectionCmd = &cobra.Command{
	Use:   "section",
	Short: "Get a single page section",
	Long:  "Get a single page section",
	RunE:  section,
}

func section(cmd *cobra.Command, args []string) error {
	fmt.Println("Executing section command")
	return nil
}
