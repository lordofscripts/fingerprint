package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/text"
	"github.com/stretchr/testify/assert"
)

func TestTestCleaning(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		cleaned string
	}{
		{name: "should return an empty string for an empty string", raw: "", cleaned: ""},
		{name: "should return an empty string for all punctuation strings", raw: `,/[] `, cleaned: ""},
		{name: "should return significant characters from string", cleaned: "adorunrunrunadorunrun", raw: `A do run run run, a do run run`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.cleaned, text.Clean(test.raw))
		})
	}
}
