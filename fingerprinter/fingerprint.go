package fingerprinter

import (
	"io"
	"os"
)

type Options struct {
	GuaranteeThreshold int  // t
	NoiseThreshold     int  // k
	Normalize          bool // convert text to lowercase
	LettersOnly        bool // strip non-letters from the
}

func (o Options) IsValid() bool {
	return o.NoiseThreshold <= o.GuaranteeThreshold && o.NoiseThreshold > 0
}

// If the current instance isn't valid, it replaces the invalid
// values for defaults.
func (o Options) VerifyOrDefault() Options {
	if !o.IsValid() {
		return Options{GuaranteeThreshold: 4, NoiseThreshold: 4, Normalize: false, LettersOnly: true}
	}
	return o
}

// computes the hashes of the KGrams
func Record(k int, sourceText string) []uint32 {
	return KGramHash(KGram(k, sourceText))
}

// Fingerprints a text string
func TextFingerprint(input string, options Options) Fingerprint {
	//normalized := text.Clean(input)
	//options = options.VerifyOrDefault()
	//return WinnowFingerprint(options.GuaranteeThreshold, record(options.NoiseThreshold, normalized))
	//normalized := text.Clean(input)
	options = options.VerifyOrDefault()
	stream := NewStreamFingerprinter(options)
	_, _ = stream.Write([]byte(input))
	return stream.Finish()
}

// Fingerprints a text file
// Usage: FileFingerprint("document.txt", options)
func FileFingerprint(path string, options Options) (Fingerprint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stream := NewStreamFingerprinter(options)

	if _, err := io.Copy(stream, f); err != nil {
		return nil, err
	}

	return stream.Finish(), nil
}
