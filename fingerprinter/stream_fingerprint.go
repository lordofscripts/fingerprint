package fingerprinter

import (
	"errors"
	"hash/fnv"
	"io"
	"unicode"
	"unicode/utf8"
)

type candidate struct {
	hash  Mark
	index int
}

type StreamFingerprinter struct {
	options Options

	// Bytes that form an incomplete UTF-8 sequence at the end of a Write.
	bytePending []byte

	// The final k-1 runes from the previous chunk.
	runePending []rune

	// Monotonic deque containing candidates in the current Winnowing window.
	window []candidate

	fingerprint Fingerprint

	// Number of k-grams processed so far.
	gramIndex int

	// Prevent writes after Finish.
	finished bool

	// Prevent duplicate selection when the minimum remains unchanged.
	lastSelected Mark
	haveLast     bool
}

func NewStreamFingerprinter(options Options) *StreamFingerprinter {
	return &StreamFingerprinter{
		options: options,
	}
}

var _ io.Writer = (*StreamFingerprinter)(nil)

func (s *StreamFingerprinter) Write(p []byte) (int, error) {
	if s.finished {
		return 0, errors.New("fingerprinter: Write called after Finish")
	}

	if len(p) == 0 {
		return 0, nil
	}

	// Combine any incomplete UTF-8 sequence from the previous Write with
	// the new bytes.
	data := make([]byte, 0, len(s.bytePending)+len(p))
	data = append(data, s.bytePending...)
	data = append(data, p...)
	s.bytePending = nil

	runes := make([]rune, 0, len(data))

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)

		// An incomplete sequence may occur at the end of this chunk.
		if r == utf8.RuneError && size == 1 &&
			!utf8.FullRune(data) {
			s.bytePending = append(s.bytePending[:0], data...)
			break
		}

		if s.options.Normalize && unicode.IsLetter(r) {
			r = unicode.ToLower(r)
		}
		runes = append(runes, r)
		data = data[size:]
	}

	s.processRunes(runes)

	return len(p), nil
}

func (s *StreamFingerprinter) Finish() Fingerprint {
	if s.finished {
		return s.fingerprint
	}

	s.finished = true

	// Decode any final incomplete UTF-8 sequence as U+FFFD, matching the
	// behavior of strings.ToValidUTF8-style processing.
	if len(s.bytePending) > 0 {
		s.processRunes([]rune{utf8.RuneError})
		s.bytePending = nil
	}

	// No additional data is needed here: runePending contains the suffix
	// needed to form k-grams across Write calls, and processRunes has already
	// processed every complete k-gram.
	s.runePending = nil

	return s.fingerprint
}

func (s *StreamFingerprinter) processRunes(input []rune) {
	k := s.options.NoiseThreshold
	if k <= 0 {
		return
	}

	// Join the suffix from the previous chunk to the current input.
	runes := make([]rune, 0, len(s.runePending)+len(input))
	runes = append(runes, s.runePending...)
	runes = append(runes, input...)

	if len(runes) < k {
		s.runePending = append(s.runePending[:0], runes...)
		return
	}

	// Every k-gram starting before this position is complete. Retain the
	// final k-1 runes for the next Write.
	complete := len(runes) - (k - 1)

	for i := 0; i < complete; i++ {
		gram := runes[i : i+k]
		hash := hashRunes(gram)

		s.addCandidate(candidate{
			hash:  hash,
			index: s.gramIndex,
		})

		s.gramIndex++
	}

	s.runePending = append(
		s.runePending[:0],
		runes[complete:]...,
	)
}

func (s *StreamFingerprinter) addCandidate(c candidate) {
	windowSize :=
		s.options.GuaranteeThreshold -
			s.options.NoiseThreshold +
			1

	if windowSize < 1 {
		windowSize = 1
	}

	// Maintain the deque in ascending hash order. The newest candidate
	// wins ties, which is the usual rightmost-minimum Winnowing rule.
	for len(s.window) > 0 &&
		s.window[len(s.window)-1].hash >= c.hash {
		s.window = s.window[:len(s.window)-1]
	}

	s.window = append(s.window, c)

	// Remove candidates that have fallen outside the current window.
	firstValid := c.index - windowSize + 1

	for len(s.window) > 0 &&
		s.window[0].index < firstValid {
		s.window = s.window[1:]
	}

	// A fingerprint exists only once a complete window has been seen.
	if c.index+1 < windowSize {
		return
	}

	minimum := s.window[0].hash

	// The same minimum can remain selected across multiple windows.
	if !s.haveLast || minimum != s.lastSelected {
		s.fingerprint = append(s.fingerprint, minimum)
		s.lastSelected = minimum
		s.haveLast = true
	}
}

/*
func hashRunes(runes []rune) Mark {
	h := fnv.New64a()

	// Hash UTF-8 bytes rather than Go's implementation-dependent rune
	// representation.
	for _, r := range runes {
		var buf [utf8.UTFMax]byte
		n := utf8.EncodeRune(buf[:], r) // handled different in earlier GO versions!
		_, _ = h.Write(buf[:n])
	}

	return Mark(h.Sum64())
}
*/

// more portable implementation than the above
func hashRunes(runes []rune) Mark {
	h := fnv.New64a()

	for _, r := range runes {
		var buf [utf8.UTFMax]byte
		n := utf8.RuneLen(r)

		if n < 0 {
			r = utf8.RuneError
			n = utf8.RuneLen(r)
		}

		utf8Bytes := buf[:0]
		utf8Bytes = utf8.AppendRune(utf8Bytes, r)

		_, _ = h.Write(utf8Bytes)
	}

	return Mark(h.Sum64())
}
