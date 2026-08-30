package similarity

import (
	"os"

	"github.com/lordofscripts/fingerprint/fingerprinter"
)

// Computes the similarity between two strings
func StringSimilarity(s1, s2 string, options fingerprinter.Options) (float64, *SimilarityScores) {
	fp1 := fingerprinter.TextFingerprint(s1, options)
	fp2 := fingerprinter.TextFingerprint(s2, options)

	var scores *SimilarityScores = nil
	var coverageScore float64
	if len(s1) < len(s2) {
		coverageScore = Coverage(fp1, fp2)
	} else {
		coverageScore = Coverage(fp2, fp1)
	}

	scores = &SimilarityScores{
		RawScore:     Compare(fp1, fp2),
		SizeAdjusted: SizeAdjustedSimilarity(fp1, fp2),
		MultiSet:     MultisetSimilarity(fp1, fp2),
		Coverage:     coverageScore,

		UsrFingerprintA: fp1.HashStr(),
		UsrFingerprintB: fp2.HashStr(),
	}

	return scores.RawScore, scores
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

		UsrFingerprintA: fp1.HashStr(),
		UsrFingerprintB: fp2.HashStr(),
	}

	return scores.RawScore, scores, nil
}

// Raw similarity score between two fingerprints
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
