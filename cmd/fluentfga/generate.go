package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sectrean/fluentfga/gen"
	"github.com/sectrean/fluentfga/internal/model"
)

const CleanFlag = "clean"

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate Code",
		Long:    "Generate code from an OpenFGA authorization model",
		Example: `fluentfga generate model.fga ./output`,
		Args:    cobra.ExactArgs(2),
		RunE:    run,
	}

	cmd.PersistentFlags().Bool(CleanFlag, false, "clean the output directory before generating code")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	file := args[0]
	outDir := args[1]

	const packageName = "model"
	const filePrefix = "fga"

	config := &gen.Config{
		Package: packageName,
	}

	clean, err := cmd.PersistentFlags().GetBool(CleanFlag)
	if err != nil {
		return err
	}

	protoModel, err := model.ReadModelFromFile(file)
	if err != nil {
		return err
	}

	if clean {
		err := cleanOutDir(outDir, filePrefix)
		if err != nil {
			return err
		}
	}

	generator, err := gen.NewGenerator()
	if err != nil {
		return err
	}

	model := gen.NewModel(protoModel, config)
	return generator.Generate(model, gen.NewWriteFS(outDir))
}

func cleanOutDir(outDir string, filePrefix string) error {
	files, err := filepath.Glob(filepath.Join(outDir, filePrefix+"_*_gen.go"))
	if err != nil {
		return err
	}

	dumpFiles, err := filepath.Glob(filepath.Join(outDir, filePrefix+"_*_gen.go.dump"))
	if err != nil {
		return err
	}
	files = append(files, dumpFiles...)

	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}
