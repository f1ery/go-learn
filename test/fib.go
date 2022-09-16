package test

import (
	"fmt"
	"sync"
)

func add11() {
	for i := 0; i < 5000; i++ {
		//lock.Lock()
		x = x + 1
		//lock.Unlock()
	}
	wg.Done()
}

var x int64
var wg sync.WaitGroup
var lock sync.Mutex
func RaceLock()  {
	wg.Add(2)
	go add11()
	go add11()
	wg.Wait()
	fmt.Println(x)
	fmt.Println("end")
}
