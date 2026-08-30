/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A model containing the aggregated similarity scores and
 * fingerprints.
 *-----------------------------------------------------------------*/
package similarity

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ fmt.Stringer = (*SimilarityScores)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// An object that reports a summary of Similarity scores.
type SimilarityScores struct {
	RawScore     float64
	SizeAdjusted float64
	MultiSet     float64
	Coverage     float64

	UsrFingerprintA string // human-friendly Document A's fingerprint
	UsrFingerprintB string // human-friendly Document B's fingerprint
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer to show a formatted string on the console
// with the current scores.
func (ss SimilarityScores) String() string {
	const LEADER rune = '\t'
	var sb strings.Builder
	fmt.Fprintf(&sb, "%cRaw          : %.1f%%\n", LEADER, ss.RawScore*100)
	fmt.Fprintf(&sb, "%cSize-adjusted: %.1f%%\n", LEADER, ss.SizeAdjusted*100)
	fmt.Fprintf(&sb, "%cMulti-set    : %.1f%%\n", LEADER, ss.MultiSet*100)
	fmt.Fprintf(&sb, "%cCoverage     : %.1f%%\n", LEADER, ss.Coverage*100)
	fmt.Fprintf(&sb, "%cFingerprint-1: %16s\n", LEADER, ss.UsrFingerprintA)
	fmt.Fprintf(&sb, "%cFingerprint-2: %16s\n", LEADER, ss.UsrFingerprintB)
	return sb.String()
}
