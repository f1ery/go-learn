package lib

import (
	"go-learn/lib/internal"
	"os"
)

var Str1 string

func Hello(name string) {
	//fmt.Printf("param is %s!\n", name)
	internal.Hello(os.Stdout, name)
}
