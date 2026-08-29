package fingerprinter

import "math"

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
