package cmd

import (
	"fmt"
	"os"

	"github.com/mal0ner/wikiscrape/cmd/get"
	"github.com/mal0ner/wikiscrape/cmd/list"
	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
)

// Flag vars
var printVersion bool

// Command
var rootCmd = &cobra.Command{
	Use:   "wikiscrape",
	Short: "Scrape and export wiki pages!",
	Long:  "A tool for scraping wikis running on a number of different backends and then exporting the data to different formats",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			fmt.Println(version)
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(list.ListCmd)

	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "print version")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
