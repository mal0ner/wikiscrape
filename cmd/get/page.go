package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Get a single page",
	Long:  "Get a single page",
	RunE:  page,
}

func init() {}

func page(cmd *cobra.Command, args []string) error {
	fmt.Println("Executing page command")
	return nil
}
