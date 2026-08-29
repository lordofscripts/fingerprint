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

	text := "The quick brown fox jumped over the lazy dog."

	expected := fingerprinter.TextFingerprint(text, options)

	stream := fingerprinter.NewStreamFingerprinter(options)

	chunks := []string{
		"The quick ",
		"brown fox ",
		"jumped over ",
		"the lazy dog.",
	}

	for _, chunk := range chunks {
		if _, err := stream.Write([]byte(chunk)); err != nil {
			t.Fatal(err)
		}
	}

	actual := stream.Finish()

	// 1. For slices: !fingerprintsEqual(expected, actual)
	// 2. For sets  : !reflect.DeepEqual(expected, actual)
	if !fingerprintsEqual(expected, actual) {
		t.Fatalf("fingerprints differ:\nexpected: %v\nactual:   %v",
			expected, actual)
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
