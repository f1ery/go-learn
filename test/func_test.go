package test

import (
	"github.com/pkg/profile"
	"os"
	"runtime/pprof"
	"testing"
)
import _ "runtime/pprof"

func array() [1024]int {
	var arr [1024]int
	for i := 0; i < 1024; i++ {
		arr[i] = i
	}
	return arr
}

func slice() []int {
	sli := make([]int, 1024)
	for i := 0; i < 1024; i++ {
		sli[i] = i
	}
	return sli
}

func BenchmarkArray(b *testing.B) {
	f, _ := os.OpenFile("array_cpu.profile", os.O_CREATE|os.O_RDWR, 0755)
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < b.N; i++ {
		array()
	}
}

func BenchmarkSlice(b *testing.B) {
	//f, _ := os.OpenFile("slice_cpu.profile", os.O_CREATE|os.O_RDWR, 0755)
	//pprof.StartCPUProfile(f)
	defer profile.Start(profile.MemProfile).Stop()
	defer pprof.StopCPUProfile()
	for i := 0; i < b.N; i++ {
		slice()
	}
}
