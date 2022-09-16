package regex

import (
	"fmt"
	"regexp"
)

const text = `My email is edgedroper@163.com@gmail.com
email is abc@def.org
im golang@11.com.cn bd
`

func main1()  {
	//re, err := regexp.Compile("edgedroper@163.com")
	re := regexp.MustCompile(`(\w+)@([\w.]+)(\.[\w]+)`)
	//match := re.FindString(text)
	//match := re.FindAllString(text, -1)
	match := re.FindAllStringSubmatch(text, -1)
	for _, m := range match {
		fmt.Println(m)
	}
	//fmt.Println(match)
}
