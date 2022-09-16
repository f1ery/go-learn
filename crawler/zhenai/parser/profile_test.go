package parser

import (
	"fmt"
	"go-learn/crawler/fetcher"
	"testing"
)

func TestParseProfile(t *testing.T) {
	contents, err := fetcher.Fetch("http://album.zhenai.com/u/1846280757")
	if err != nil {
		panic(err)
	}
	fmt.Println(ParseProfile(contents, "summer"))
}
