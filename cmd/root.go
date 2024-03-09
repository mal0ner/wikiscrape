package cmd

import (
	"fmt"
	"os"

	"github.com/mal0ner/wikiscrape/cmd/get"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wikiscrape",
	Short: "Scrape and export wiki pages!",
	Long:  "A tool for scraping wikis running on a number of different backends and then exporting the data to different formats",
	RunE:  getPageFromURL,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
}

func init() {
	rootCmd.AddCommand(get.GetCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getPageFromURL(cmd *cobra.Command, args []string) error {
	// parse url for hostname and path
	// use internal scraper to get data
	// output to format
	fmt.Println("Executing get page from url")
	return nil
}
