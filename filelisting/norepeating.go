package filelisting

func LenghtOfNoRepeatingSubStr(s string) int {
	start, maxLength := 0, 0
	lastOccurred := make(map[rune]int)
	for i, v := range []rune(s) {
		if lastI, ok := lastOccurred[v]; ok && lastI >= start {
			start = lastI + 1
		}
		if i - start + 1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[v] = i
	}
	return maxLength
}


