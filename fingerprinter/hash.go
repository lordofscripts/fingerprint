package fingerprinter

import "hash/fnv"

func KGramHash(kgrams []string) []uint32 {
	hashes := make([]uint32, len(kgrams))
	for i, kgram := range kgrams {
		hashes[i] = hash(kgram)
	}
	return hashes
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
