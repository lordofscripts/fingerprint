package similarity

import (
	"fmt"
	"os"
	"strings"

	"github.com/lordofscripts/fingerprint/fingerprinter"
)

// An object that reports a summary of Similarity scores.
type SimilarityScores struct {
	RawScore     float64
	SizeAdjusted float64
	MultiSet     float64
	Coverage     float64
}

// implements fmt.Stringer to show a formatted string on the console
// with the current scores.
func (ss SimilarityScores) String() string {
	const LEADER rune = '\t'
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%cRaw          : %.3f\n", LEADER, ss.RawScore))
	sb.WriteString(fmt.Sprintf("%cSize-adjusted: %.3f\n", LEADER, ss.SizeAdjusted))
	sb.WriteString(fmt.Sprintf("%cMulti-set    : %.3f\n", LEADER, ss.MultiSet))
	sb.WriteString(fmt.Sprintf("%cCoverage     : %.3f\n", LEADER, ss.Coverage))
	return sb.String()
}

// Computes the similarity between two strings
func StringSimilarity(s1, s2 string, options fingerprinter.Options) float64 {
	f1 := fingerprinter.TextFingerprint(s1, options)
	f2 := fingerprinter.TextFingerprint(s2, options)

	return Compare(f1, f2)
}

// Computes the similarity between two text files.
func FileSimilarity(
	path1, path2 string,
	options fingerprinter.Options,
) (float64, *SimilarityScores, error) {
	// just get the file size or -1 on error
	getFileSize := func(filename string) int64 {
		if fi, err := os.Stat(filename); err != nil {
			return -1
		} else {
			return fi.Size()
		}
	}

	var scores *SimilarityScores = nil
	fp1, err := fingerprinter.FileFingerprint(path1, options)
	if err != nil {
		return 0, nil, err
	}
	fs1 := getFileSize(path1)

	fp2, err := fingerprinter.FileFingerprint(path2, options)
	if err != nil {
		return 0, nil, err
	}
	fs2 := getFileSize(path2)

	var coverageScore float64
	if fs1 < fs2 {
		coverageScore = Coverage(fp1, fp2)
	} else {
		coverageScore = Coverage(fp2, fp1)
	}

	scores = &SimilarityScores{
		RawScore:     Compare(fp1, fp2),
		SizeAdjusted: SizeAdjustedSimilarity(fp1, fp2),
		MultiSet:     MultisetSimilarity(fp1, fp2),
		Coverage:     coverageScore,
	}

	return scores.RawScore, scores, nil
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

func Compare(f1, f2 fingerprinter.Fingerprint) float64 {
	set1 := f1.AsSet()
	set2 := f2.AsSet()

	unionCount := len(set1)
	intersectCount := 0
	for mark := range set2 {
		if set1[mark] {
			intersectCount++
		} else {
			unionCount++
		}
	}
	if unionCount == 0 {
		return 0
	}
	return float64(intersectCount) / float64(unionCount)
}
