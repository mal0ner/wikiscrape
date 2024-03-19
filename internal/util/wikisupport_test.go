package util_test

import (
	"testing"

	"github.com/mal0ner/wikiscrape/internal/util"
)

func TestGetQueryDataFromURL(t *testing.T) {
	// Test 1: Invalid URL
	invalidURL := "cheesebiscuit./restaurant"
	_, err := util.GetQueryDataFromURL(invalidURL)
	if err == nil {
		t.Error("Expected an error for invalid URL, but got nil")
	}

	// Test 2: Invalid URL Page Prefix
	invalidPrefixURL := "https://en.wikipedia.org/weeee/Page"
	_, err = util.GetQueryDataFromURL(invalidPrefixURL)
	if err == nil {
		t.Error("Expected an error for invalid URL, but got nil")
	}

	// Test 3: Valid URL
	validURL := "https://en.wikipedia.org/wiki/Jesus"
	got, err := util.GetQueryDataFromURL(validURL)
	if err != nil {
		t.Errorf("Failed to generate query data from valid url: %v", err)
	}
	_, err = util.GetWikiInfoFromHost("en.wikipedia.org")
	if err != nil {
		t.Errorf("Failed to get wiki info for known supported host: %s", err)
	}
	want := "Jesus"
	if got.Page != want {
		t.Errorf("Failed to parse page name from URL, Got: %s, Want: %s", got.Page, want)
	}
}

func TestGetQueryDataFromName(t *testing.T) {
	// Test 1: Invalid wiki name
	invalidWikiName := "cheesebiscuit"
	pageName := "Jesus"
	_, err := util.GetQueryDataFromName(pageName, invalidWikiName)
	if err == nil {
		t.Errorf("Expected an error for invalid wiki name, got nil")
	}
}
