package manifest

import (
	"encoding/json"
	"os"
)

type Manifest []string

func ReadFrom(filepath string) (Manifest, error) {
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
