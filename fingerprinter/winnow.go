package fingerprinter

import (
	"fmt"
	"hash/fnv"
	"math"
)

func Winnow(g int, kgrams []uint32) [][]uint32 {
	length := len(kgrams)
	n := length - g + 1
	if g <= 0 || n <= 0 {
		return nil
	}
	winnow := make([][]uint32, n)
	for i := range n {
		winnow[i] = kgrams[i : i+g]
	}
	return winnow
}

func WinnowFingerprint(g int, kgrams []uint32) Fingerprint {
	windows := Winnow(g, kgrams)
	if len(windows) == 0 {
		return nil
	}
	finger := make(Fingerprint, len(windows))
	for i, window := range windows {
		finger[i] = RightmostLowestValue(window)
	}
	return finger
}

type Mark uint64

func NewMark(minValue uint32, index uint32) Mark {
	return Mark(uint64(minValue) | uint64(index)>>32)
}

type Fingerprint []Mark

func RightmostLowestValue(values []uint32) (w Mark) {
	var MinValue uint32 = math.MaxUint32
	var Index uint32 = 0
	for i, v := range values {
		if v <= MinValue {
			MinValue = v
			Index = uint32(i)
		}
	}
	return NewMark(MinValue, Index)
}

type MarkSet map[Mark]bool

func (f Fingerprint) AsSet() MarkSet {
	set := make(MarkSet, len(f))
	for _, mark := range f {
		set[mark] = true
	}
	return set
}

// HashStr condenses the similarity signature into a single
// unique string ID. This reduces the "massive" vector down to a clean
// unique 16-character hex string that uniquely identifies that specific
// similarity structural profile.
func (f Fingerprint) HashStr() string {
	h := fnv.New64a()
	for _, val := range f {
		// convert uint64 to bytes
		bytes := []byte{
			byte(val), byte(val >> 8),
			byte(val >> 16), byte(val >> 24),
			byte(val >> 32), byte(val >> 40),
			byte(val >> 48), byte(val >> 56),
		}
		h.Write(bytes)
	}
	return fmt.Sprintf("%x", h.Sum64())
}
