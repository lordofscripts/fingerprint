/* -----------------------------------------------------------------
 *
 *
 *                       GoFingerprint
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *                      U n i t   T e s t
 *-----------------------------------------------------------------*/
package tests

import (
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/stretchr/testify/assert"
)

// Was: TestKgramGenerations
func TestFingerprinter_KgramGenerations(t *testing.T) {
	t.Run("should return empty list for empty string", func(t *testing.T) {
		assert.Empty(t, fingerprinter.KGram(0, ""))
	})
	t.Run("should return empty list for strings less than k", func(t *testing.T) {
		assert.Empty(t, fingerprinter.KGram(4, "abs"))
	})

	t.Run("should return a single kgram if k matches the string length", func(t *testing.T) {
		assert.Equal(t, []string{"abs"}, fingerprinter.KGram(3, "abs"))
	})
	t.Run("should return a single kgram if k matches the string length", func(t *testing.T) {
		expected := []string{
			"adoru", "dorun", "orunr", "runru", "unrun", "nrunr", "runru", "unrun", "nruna", "runad", "unado", "nador", "adoru", "dorun", "orunr", "runru", "unrun",
		}
		assert.Equal(t, expected, fingerprinter.KGram(5, "adorunrunrunadorunrun"))
	})
}
