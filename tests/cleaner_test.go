package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/text"
	"github.com/stretchr/testify/assert"
)

// NOTE: Used in the original repository. The new (my enhancements) no
// longer invoke this as it is done based on "Options.Normalize"
func TestText_Clean(t *testing.T) {
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
