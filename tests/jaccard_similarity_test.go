package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/lordofscripts/fingerprint/similarity"
	"github.com/stretchr/testify/assert"
)

func TestStringSimilarity(t *testing.T) {
	t.Run("should return 1 for identical strings", func(t *testing.T) {
		text := "the quick brown fox jumped over the lazy dog"
		assert.Equal(t, 1.0, similarity.StringSimilarity(text, text, fingerprinter.Options{}))
	})
	t.Run("should return 0 for completely different strings", func(t *testing.T) {
		s1 := "aaaaaaaaaaaaaaaa"
		s2 := "bbbbbbbbbbbbbbbb"
		sim := similarity.StringSimilarity(s1, s2, fingerprinter.Options{})
		assert.Less(t, sim, 0.1)
	})
}

func TestCompare(t *testing.T) {
	t.Run("two empty fingerprints returns 0", func(t *testing.T) {
		f1 := fingerprinter.Fingerprint{}
		f2 := fingerprinter.Fingerprint{}
		assert.Equal(t, 0.0, similarity.Compare(f1, f2))
	})
	t.Run("identical fingerprints returns 1", func(t *testing.T) {
		f1 := fingerprinter.Fingerprint{fingerprinter.NewMark(1, 2), fingerprinter.NewMark(3, 4)}
		assert.Equal(t, 1.0, similarity.Compare(f1, f1))
	})
	t.Run("disjoint fingerprints returns 0", func(t *testing.T) {
		f1 := fingerprinter.Fingerprint{fingerprinter.NewMark(1, 2)}
		f2 := fingerprinter.Fingerprint{fingerprinter.NewMark(5, 6)}
		assert.Equal(t, 0.0, similarity.Compare(f1, f2))
	})
	t.Run("partial overlap", func(t *testing.T) {
		f1 := fingerprinter.Fingerprint{fingerprinter.NewMark(1, 2), fingerprinter.NewMark(3, 4)}
		f2 := fingerprinter.Fingerprint{fingerprinter.NewMark(1, 2), fingerprinter.NewMark(5, 6)}
		assert.InDelta(t, 1.0/3.0, similarity.Compare(f1, f2), 0.001)
	})
}
