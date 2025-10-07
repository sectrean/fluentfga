package gen

import (
	"strings"
	"text/template"
)

var templateFunctions = template.FuncMap{
	"abbr":      Abbr,
	"titleCase": TitleCase,
	"camelCase": CamelCase,
}

func Abbr(name string) string {
	if len(name) == 0 {
		return ""
	}

	return strings.ToLower(string(name[0]))
}

func TitleCase(name string) string {
	words := strings.Split(name, "_")
	titleCased := make([]string, len(words))

	for i, word := range words {
		if len(word) > 0 {
			titleCased[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(titleCased, "")
}

func CamelCase(name string) string {
	name = TitleCase(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

func ID(name string) string {
	return "ID"
}

func NameID(name string) string {
	return TitleCase(name) + "ID"
}
