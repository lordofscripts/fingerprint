package tests

import (
	"fmt"
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/lordofscripts/fingerprint/similarity"
	"github.com/stretchr/testify/assert"
)

// Was (original repo): TestFingerprintingAcceptance
func TestSimilarity_Acceptance(t *testing.T) {
	tests := []struct {
		name     string
		input1   string
		input2   string
		options  fingerprinter.Options
		expected float64
	}{
		{
			name:     "identical text",
			input1:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
			input2:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
			options:  fingerprinter.Options{GuaranteeThreshold: 4, NoiseThreshold: 4, Normalize: true},
			expected: 0,
		},
		{
			name:     "small changes are similar",
			input1:   `Lorem ipsum dolor sit amet, consectetur adipiscing foo, bar do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
			input2:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
			options:  fingerprinter.Options{GuaranteeThreshold: 4, NoiseThreshold: 4, Normalize: true},
			expected: 0.1,
		},
		{
			name:     "very different text block",
			input1:   `Lorem ipsum dolor sit amet, consectetur adipiscing foo, bar do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
			input2:   `The quick brown fox jumped over the lazy dog`,
			options:  fingerprinter.Options{GuaranteeThreshold: 4, NoiseThreshold: 4, Normalize: true},
			expected: 1.00,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fingerprint1 := fingerprinter.TextFingerprint(test.input1, test.options)
			fingerprint2 := fingerprinter.TextFingerprint(test.input2, test.options)
			diff := 1.0 - similarity.Compare(fingerprint1, fingerprint2)
			assert.True(t, diff <= test.expected, fmt.Sprintf("%f not less than threshold %f", diff, test.expected))
		})
	}
}
