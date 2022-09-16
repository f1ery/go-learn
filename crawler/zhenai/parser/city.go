package parser

import (
	"go-learn/crawler/engine"
	"regexp"
)

var cityRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/\d+)"[^>]*>([^<]+)</a>`)

func ParseCity(contents []byte) engine.ParseResult {

	match := cityRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}

	for _, v := range match {
		name := string(v[2]) //将name拷贝出来送给闭包，不这样做会导致一个城市的人名称相同
		result.Items = append(result.Items, "User "+string(v[2]))
		result.Requests = append(result.Requests,
			engine.Request{
				Url: string(v[1]),
				ParserFunc: func(contents []byte) engine.ParseResult {
					return ParseProfile(contents, name)
				},
			},
		)
	}
	return result
}