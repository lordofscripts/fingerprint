package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
)

func TestStreamFingerprintMatchesTextFingerprint(t *testing.T) {
	options := fingerprinter.Options{
		GuaranteeThreshold: 4,
		NoiseThreshold:     4,
		Normalize:          true, // @todo update test to try both
	}

	text := "The quick brown ß Fox jumped over the lazy dog."
	chunks := []string{
		"The quick ",
		"brown ß Fox ",
		"jumped over ",
		"the lazy dog.",
	}

	// Test with both normalize conditions
	for _, normalize := range []bool{true, false} {
		options.Normalize = normalize
		expected := fingerprinter.TextFingerprint(text, options)

		stream := fingerprinter.NewStreamFingerprinter(options)
		for _, chunk := range chunks {
			if _, err := stream.Write([]byte(chunk)); err != nil {
				t.Fatal(err)
			}
		}
		actual := stream.Finish()

		// 1. For slices: !fingerprintsEqual(expected, actual)
		// 2. For sets  : !reflect.DeepEqual(expected, actual)
		if !fingerprintsEqual(expected, actual) {
			t.Fatalf("fingerprints differ. Normalize:%t\nexpected: %v\nactual:   %v",
				normalize, expected, actual)
		}
	}
}

func fingerprintsEqual(a, b fingerprinter.Fingerprint) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
