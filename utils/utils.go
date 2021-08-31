package utils

import (
	"fmt"
)

func HandlePanicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func UnsliceStrings(strings []string, separator string) string {
	result := ""
	for _, s := range strings {
		result = result + s + separator
	}
	return result[len(result)-len(separator):]
}

/*
String separator : $
Escape character : !
Examples : Date$Name$Content -> {"Date", "Name", "Content"}
		   Date$ -> {"Date", ""}
		   Date!$ -> {"Date$"}
		   Date!! -> {"Date!"}

*/
func ParseDatabaseStringList(dbText string) []string {
	result := make([]string, 1)
	stringIndex := 0
	for i := 0; i < len(dbText); i++ {
		c := dbText[i]
		switch c {
		case '$':
			result = append(result, "")
			stringIndex++
		case '!':
			i++
			result[stringIndex] += string(dbText[i])
		default:
			result[stringIndex] += string(c)
		}
	}
	return result
}

func Esc(s string) string {
	return "\"" + s + "\""
}
