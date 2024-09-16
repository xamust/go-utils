package slice

import "strings"

var replacerXml = strings.NewReplacer(
	"<", "&lt;",
	">", "&gt;")

func StringContains(str string, arr []string) bool {
	for index := range arr {
		if arr[index] == str {
			return true
		}
	}
	return false
}

func EscapeXmlString(s string) string {
	return replacerXml.Replace(s)
}

func EscapeXml(s []byte) string {
	return EscapeXmlString(string(s))
}
