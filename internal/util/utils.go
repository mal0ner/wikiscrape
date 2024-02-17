package util

import "strings"

func TrimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
