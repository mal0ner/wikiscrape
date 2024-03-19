package util_test

import (
	"fmt"
	"testing"

	"github.com/mal0ner/wikiscrape/internal/util"
)

type utilTestCase struct {
	Input string
	Want  string
}

func TestTrimLower(t *testing.T) {
	cases := []utilTestCase{
		{"PIZZA", "pizza"},
		{"  PIZZA", "pizza"},
		{"  ", ""},
		{" Ronald Raegan ", "ronald raegan"},
		{"ronald i love you", "ronald i love you"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("lower+trimspace(%s)=%s", tc.Input, tc.Want), func(t *testing.T) {
			got := util.TrimLower(tc.Input)
			if tc.Want != got {
				t.Errorf("Expected %s, got %s", tc.Want, got)
			}
		})
	}
}
