package gen

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"golang.org/x/tools/imports"
)

//go:embed templates/*.tmpl
var templates embed.FS

var fileTemplates = []struct {
	Template string
	FileName string
}{
	{"relations", "fga_relations_gen.go"},
	{"types", "fga_types_gen.go"},
}

type Generator struct {
	tmpl *template.Template
}

func NewGenerator() (*Generator, error) {
	tmpl := template.New("templates")
	tmpl.Funcs(TemplateFunctions)

	tmpl, err := tmpl.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("new generator: %w", err)
	}

	return &Generator{tmpl: tmpl}, nil
}

func (g *Generator) Generate(model *Model, output WriteFS) error {
	for _, ft := range fileTemplates {
		err := g.renderTemplate(ft.Template, model, output, ft.FileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderTemplate(name string, data any, fs WriteFS, fileName string) error {
	buf := new(bytes.Buffer)
	err := g.tmpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %q: %w", name, err)
	}

	generated := buf.Bytes()
	formatted, err := imports.Process(fileName, generated, nil)
	if err != nil {
		fs.WriteFile(fileName+".dump", generated)
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	if err := fs.WriteFile(fileName, formatted); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
