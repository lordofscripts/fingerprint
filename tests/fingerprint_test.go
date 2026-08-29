package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/stretchr/testify/assert"
)

func TestFingerprinting(t *testing.T) {
	t.Run("returns empty fingerprinter for empty text", func(t *testing.T) {
		assert.Empty(t, fingerprinter.Record(4, ""))
	})

	t.Run("returns fingerprinter vector for given text", func(t *testing.T) {
		expected := []uint32{0xf765d270, 0xce399891, 0x9d4bf9df, 0x88294ce7, 0x2f4243db, 0x7116aba4, 0x88294ce7, 0x2f4243db, 0x8416c98d, 0x8908facf, 0x43870890, 0x1a820a81, 0xf765d270, 0xce399891, 0x9d4bf9df, 0x88294ce7, 0x2f4243db}
		assert.Equal(t, expected, fingerprinter.Record(5, "adorunrunrunadorunrun"))
	})
}

func TestFingerprintOptions(t *testing.T) {
	tests := []struct {
		name   string
		option fingerprinter.Options
		result bool
		error  error
	}{
		{name: "k > t", option: fingerprinter.Options{GuaranteeThreshold: 0, NoiseThreshold: 1}, result: false, error: nil},
		{name: "k == t", option: fingerprinter.Options{GuaranteeThreshold: 1, NoiseThreshold: 1}, result: true, error: nil},
		{name: "k < t", option: fingerprinter.Options{GuaranteeThreshold: 2, NoiseThreshold: 1}, result: true, error: nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid := test.option.IsValid()
			assert.Equal(t, test.result, valid)
		})
	}
}

func TestFingerprint_AsSet(t *testing.T) {
	t.Run("Empty fingerprinter empty set", func(t *testing.T) {
		fingerprint := make(fingerprinter.Fingerprint, 0)
		set := fingerprint.AsSet()
		assert.Len(t, set, 0)
	})

	t.Run("set contains all marks from fingerprinter", func(t *testing.T) {
		fingerprint := make(fingerprinter.Fingerprint, 0)
		fingerprint = append(fingerprint, fingerprinter.NewMark(1, 2))
		fingerprint = append(fingerprint, fingerprinter.NewMark(3, 4))
		set := fingerprint.AsSet()
		assert.Len(t, set, 2)
		assert.Contains(t, fingerprint, fingerprinter.NewMark(1, 2))
		assert.Contains(t, fingerprint, fingerprinter.NewMark(3, 4))
	})
}
