package similarity

import (
	"github.com/lordofscripts/fingerprint/fingerprinter"
)

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
