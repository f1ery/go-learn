package util

import (
	"fmt"
	"testing"
)

func BenchmarkFor1(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := "abcdef"
		for i, l := 0, len(s); i < l; i++ {
			//println(i, s[i])
		}
	}
	b.StopTimer()
}

func BenchmarkFor2(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s := "abcdef"
		for i:= 0; i < len(s); i++ {
			//println(i, s[i])
		}
	}
	b.StopTimer()
}

func BenchmarkArgs1(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	var sli1 []int
	for i := 0; i < 1000000; i++ {
		sli1 = append(sli1, i)
	}

	sli2 := []int{11,22,33}
	sli2 = append(sli2, sli1...)
	b.StopTimer()
}

func BenchmarkArgs2(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	var sli1 []int
	for i := 0; i < 1000000; i++ {
		sli1 = append(sli1, i)
	}

	sli2 := []int{11,22,33}
	for _, v := range sli1 {
		sli2 = append(sli2, v)
	}
	b.StopTimer()
}

func BenchmarkArgs3(b *testing.B) {
	b.ResetTimer()
	b.StartTimer()
	var sli1 []int
	for i := 0; i < 5; i++ {
		sli1 = append(sli1, i)
	}
fmt.Println(sli1)
	testArgs(sli1...)
	b.StopTimer()
}

func testArgs(sli ...int)  {
	for v := range sli {
		fmt.Println(v)
	}
}

