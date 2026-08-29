package tests

import (
	"math"
	"testing"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/stretchr/testify/assert"
)

func TestWinnow(t *testing.T) {
	tests := []struct {
		name     string
		g        int
		expected [][]uint32
		values   []uint32
	}{
		{name: "k-gram less that g", expected: nil, g: 0, values: []uint32{}},
		{name: "returns a single k-gram if length equal to g", g: 3, expected: [][]uint32{{1, 2, 3}}, values: []uint32{1, 2, 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, fingerprinter.Winnow(test.g, test.values))
		})
	}
}

func TestWinnowFingerprinting(t *testing.T) {
	tests := []struct {
		name        string
		fingerprint fingerprinter.Fingerprint
		g           int
		kgrams      []uint32
	}{
		{name: "empty k-grams", fingerprint: fingerprinter.Fingerprint{fingerprinter.NewMark(0x2, 2), fingerprinter.NewMark(0x1, 2)}, g: 3, kgrams: []uint32{4, 3, 2, 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.fingerprint, fingerprinter.WinnowFingerprint(test.g, test.kgrams))
		})
	}
}

func TestRightmostLowestValue(t *testing.T) {
	tests := []struct {
		name     string
		expected fingerprinter.Mark
		values   []uint32
	}{
		{"empty", fingerprinter.NewMark(math.MaxUint32, 0), []uint32{}},
		{"single entry", fingerprinter.NewMark(1, 0), []uint32{1}},
		{"min in the right most position", fingerprinter.NewMark(1, 2), []uint32{100, 10, 1}},
		{"min in the right most position repeated", fingerprinter.NewMark(1, 4), []uint32{1, 100, 10, 1, 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, fingerprinter.RightmostLowestValue(test.values))
		})
	}
}
