package utils

func Reverse(s string) string {
	runes := []rune(s)
	size := len(runes)
	for i, j := 0, size-1; i < size>>1; i, j = i+1, j-1 { // i < size >> 1 same as i < j
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
