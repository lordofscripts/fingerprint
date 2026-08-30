/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package similarity

import (
	"github.com/lordofscripts/fingerprint/fingerprinter"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Scores the similarity of lengths.
func LengthSimilarity(a, b int) float64 {
	if a == 0 && b == 0 {
		return 1
	}
	if a == 0 || b == 0 {
		return 0
	}

	if a > b {
		a, b = b, a
	}

	return float64(a) / float64(b)
}

// Computes the overall similarity by adjusting the
// content similarity metric with the length similarity.
func AdjustedLengthSimilarity(
	contentSimilarity float64,
	lenA int,
	lenB int,
) float64 {
	sizeSimilarity := LengthSimilarity(lenA, lenB)

	return contentSimilarity * sizeSimilarity
}

// Compares two fingerprints using Set Similarity and adjusts it
// for the length of both sets.
func SizeAdjustedSimilarity(a, b fingerprinter.Fingerprint) float64 {
	jaccard := SetSimilarity(a, b)

	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	sizeRatio := float64(minimumOf(len(a), len(b))) /
		float64(maximumOf(len(a), len(b)))

	return jaccard * sizeRatio
}

func minimumOf(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maximumOf(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// The variation of score accounts for multiplicity, like
// how many times it occurs. For example 20 identical paragraphs
// vs. 15 identical paragraphs would account for approx. 0.750.
// Uses:
// - Are these files the same overall size and content?
// - Repetitions count
func MultisetSimilarity(a, b fingerprinter.Fingerprint) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}

	countA := make(map[fingerprinter.Mark]int)
	countB := make(map[fingerprinter.Mark]int)

	for _, mark := range a {
		countA[mark]++
	}

	for _, mark := range b {
		countB[mark]++
	}

	var intersection int
	var union int

	allMarks := make(map[fingerprinter.Mark]struct{})

	for mark := range countA {
		allMarks[mark] = struct{}{}
	}

	for mark := range countB {
		allMarks[mark] = struct{}{}
	}

	for mark := range allMarks {
		aCount := countA[mark]
		bCount := countB[mark]

		if aCount < bCount {
			intersection += aCount
			union += bCount
		} else {
			intersection += bCount
			union += aCount
		}
	}

	if union == 0 {
		return 1
	}

	return float64(intersection) / float64(union)
}

// Asymmetric coverage scoring. The shorter vs longer should be determined
// by the caller, not by comparing Fingerprint lengths but by the document sizes.
// Uses: How much text was copied?
func Coverage(shorter, longer fingerprinter.Fingerprint) float64 {
	shortSet := make(map[fingerprinter.Mark]struct{})
	longSet := make(map[fingerprinter.Mark]struct{})

	for _, mark := range shorter {
		shortSet[mark] = struct{}{}
	}

	for _, mark := range longer {
		longSet[mark] = struct{}{}
	}

	if len(shortSet) == 0 {
		return 0
	}

	var matches int
	for mark := range shortSet {
		if _, ok := longSet[mark]; ok {
			matches++
		}
	}

	return float64(matches) / float64(len(shortSet))
}

// Compares the similarity of two fingerprint sets by removing duplicates
// inside each set.
func SetSimilarity(a, b fingerprinter.Fingerprint) float64 {
	// @note the conversion to map[Mark]struct{} removes duplicate fingerprints
	setA := make(map[fingerprinter.Mark]struct{}, len(a))
	setB := make(map[fingerprinter.Mark]struct{}, len(b))

	for _, mark := range a {
		setA[mark] = struct{}{}
	}

	for _, mark := range b {
		setB[mark] = struct{}{}
	}

	if len(setA) == 0 && len(setB) == 0 {
		return 1
	}

	if len(setA) == 0 || len(setB) == 0 {
		return 0
	}

	intersection := 0
	union := make(map[fingerprinter.Mark]struct{},
		len(setA)+len(setB),
	)

	for mark := range setA {
		union[mark] = struct{}{}

		if _, exists := setB[mark]; exists {
			intersection++
		}
	}

	for mark := range setB {
		union[mark] = struct{}{}
	}

	return float64(intersection) / float64(len(union))
}
