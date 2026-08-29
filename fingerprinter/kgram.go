package fingerprinter

func KGram(k int, input string) []string {
	if k < 1 {
		return nil
	}
	l := len(input)
	n := l - k + 1
	if n <= 0 {
		return nil
	}
	result := make([]string, n)
	for i := range n {
		result[i] = input[i : i+k]
	}
	return result
}
