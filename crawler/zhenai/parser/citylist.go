package parser

import (
	"go-learn/crawler/engine"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/\w+)"[^>]*>([^<]+)</a>`



func ParseCityList(contents []byte) engine.ParseResult {
	//cityMap := make(map[string]string)
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range matches {
		//fmt.Println(string(m))
		//fmt.Printf("%s\n", m)
		//cityMap[string(m[1])] = string(m[2])
		//for _, subMatch := range m {
		//	fmt.Printf("%s ", subMatch)
		//}
		//fmt.Println()
		result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			ParserFunc: ParseCity,
		},
		)
	}
	//fmt.Println(cityMap)
	//fmt.Println(len(matches))
	return result
}