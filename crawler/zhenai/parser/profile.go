package parser

import (
	"encoding/json"
	"fmt"
	"go-learn/crawler/engine"
	"go-learn/crawler/model"
	"regexp"
)

func init()  {
	fmt.Println("profile")
}

//var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)</td>`)
//var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>(\w+)</td>`)

var profileRe = regexp.MustCompile(`window.__INITIAL_STATE=(\w+);\(func`)

func ParseProfile(contents []byte, name string) engine.ParseResult {
	match := profileRe.FindSubmatch(contents)
	if len(match) <= 1 {
		return engine.ParseResult{
			Items: nil,
		}
	}

	userProfile := model.UserProfile{}
	err := json.Unmarshal(match[1], &userProfile)
	if err != nil {
		panic(err)
	}

	profile := model.Profile{
		Age: userProfile.ObjectInfo.Age,
		Gender: userProfile.ObjectInfo.GenderString,
		Height:userProfile.ObjectInfo.HeightString,
		Income:userProfile.ObjectInfo.SalaryString,
		Marriage:userProfile.ObjectInfo.MarriageString,
		Education:userProfile.ObjectInfo.EducationString,
	}
	profile.Name = name

	//age, err := strconv.Atoi(extractString(contents, ageRe))
	//if err != nil {
	//	profile.Age = age
	//}
	//profile.Marriage = extractString(contents, marriageRe)
	//
	result := engine.ParseResult{
		Items: []interface{}{profile},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}