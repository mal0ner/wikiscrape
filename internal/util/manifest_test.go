package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mal0ner/wikiscrape/internal/util"
)

func TestReadManifestFrom(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "manifest.json")
	content := []byte(`["page1", "page2", "page3"]`)
	err := os.WriteFile(jsonFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary JSON file: %v", err)
	}

	// Test 1: Valid json
	manifest, err := util.ReadManifestFrom(jsonFile)
	if err != nil {
		t.Errorf("Invalid JSON format: %v", err)
	}
	expectedManifest := util.Manifest{"page1", "page2", "page3"}
	if len(manifest) != len(expectedManifest) {
		t.Errorf("Manifest length mismatch. Got: %d, Want: %d", len(manifest), len(expectedManifest))
	}
	for i, page := range manifest {
		if page != expectedManifest[i] {
			t.Errorf("Manifest content mismatch at index %d. Got: %s, Want: %s", i, page, expectedManifest[i])
		}
	}

	// Test 2: Non-existent file
	_, err = util.ReadManifestFrom("fake.json")
	if err == nil {
		t.Error("Expected an error for a non-existent json file, but got nil")
	}

	// Test 3: Invalid JSON
	invalidJSONFile := filepath.Join(tmpDir, "invalid.json")
	invalidContent := []byte(`["page1", "page2", "page3"`)
	err = os.WriteFile(invalidJSONFile, invalidContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary invalid JSON file: %v", err)
	}
	_, err = util.ReadManifestFrom(invalidJSONFile)
	if err == nil {
		t.Error("Expected an error for invalid JSON file, but got nil")
	}
}
