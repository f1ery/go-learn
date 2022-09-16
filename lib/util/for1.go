package util

func For1() {
	s := "abcdef"
	for i, l := 0, len(s); i < l; i++ {
		println(i, s[i])
	}
}

func For2() {
	s := "abcdef"
	for i:= 0; i < len(s); i++ {
		println(i, s[i])
	}
}