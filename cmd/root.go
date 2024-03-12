package cmd

import (
	"os"

	"github.com/mal0ner/wikiscrape/cmd/get"
	"github.com/mal0ner/wikiscrape/cmd/list"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wikiscrape",
	Short: "Scrape and export wiki pages!",
	Long:  "A tool for scraping wikis running on a number of different backends and then exporting the data to different formats",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
}

func init() {
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(list.ListCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
