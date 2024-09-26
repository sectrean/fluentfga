package fluentfga

import (
	"strings"
	"text/template"
)

var TemplateFunctions = template.FuncMap{
	"abbr":      abbr,
	"titleCase": titleCase,
	"camelCase": camelCase,
}

func abbr(name string) string {
	if len(name) == 0 {
		return ""
	}

	return strings.ToLower(string(name[0]))
}

func titleCase(name string) string {
	words := strings.Split(name, "_")
	titleCased := make([]string, len(words))

	for i, word := range words {
		if len(word) > 0 {
			titleCased[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(titleCased, "")
}

func camelCase(name string) string {
	return strings.ToLower(string(name[0])) + name[1:]
}
