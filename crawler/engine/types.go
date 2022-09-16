package engine

import "fmt"

type Request struct {
	Url string
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items []interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
func init()  {
	fmt.Println("types")
}