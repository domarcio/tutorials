package reverse

func WordReverse(s string) string {
	var (
		runes  = []rune(s)
		strLen = len(runes)
	)

	for i, j := 0, strLen-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
