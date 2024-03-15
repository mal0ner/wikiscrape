package util

import (
	"os"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Manifest []string

// ReadManifestFrom accepts a path to a JSON file containing
// an array of page names and reads its content.
func ReadManifestFrom(filepath string) (Manifest, error) {
	manFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var man Manifest
	err = json.Unmarshal(manFile, &man)
	if err != nil {
		return nil, err
	}
	return man, nil
}
