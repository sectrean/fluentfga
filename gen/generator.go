package gen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"text/template"

	"golang.org/x/tools/imports"
)

//go:embed templates/*.tmpl
var templates embed.FS

var fileTemplates = []struct {
	Template string
	FileName string
}{
	{"relations", "_relations_gen.go"},
	{"types", "_types_gen.go"},
}

type Generator struct {
	config *Config
	tmpl   *template.Template
}

func NewGenerator(config *Config) (*Generator, error) {
	tmpl := template.New("templates")
	tmpl.Funcs(templateFunctions)

	tmpl, err := tmpl.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("new generator: %w", err)
	}

	return &Generator{
		config: config,
		tmpl:   tmpl,
	}, nil
}

func (g *Generator) Generate(model *Model, output *os.Root) error {
	for _, ft := range fileTemplates {
		fileName := g.config.FilePrefix + ft.FileName
		err := g.renderTemplate(ft.Template, model, output, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderTemplate(name string, data any, output *os.Root, fileName string) error {
	const fileMode = 0644

	buf := new(bytes.Buffer)
	err := g.tmpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %q: %w", name, err)
	}

	generated := buf.Bytes()
	formatted, err := imports.Process(fileName, generated, nil)
	if err != nil {
		_ = output.WriteFile(fileName+".dump", generated, fileMode)
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	if err := output.WriteFile(fileName, formatted, fileMode); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (g *Generator) CleanOutput(output *os.Root) error {
	files, err := fs.Glob(output.FS(), g.config.FilePrefix+"_*_gen.go")
	if err != nil {
		return err
	}

	dumpFiles, err := fs.Glob(output.FS(), g.config.FilePrefix+"_*_gen.go.dump")
	if err != nil {
		return err
	}
	files = append(files, dumpFiles...)

	for _, f := range files {
		err := output.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}
