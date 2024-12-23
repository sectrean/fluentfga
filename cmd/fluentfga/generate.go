package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	proto "github.com/openfga/api/proto/openfga/v1"
	language "github.com/openfga/language/pkg/go/transformer"
	"github.com/spf13/cobra"

	"github.com/johnrutherford/fluentfga/gen"
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

	protoModel, err := readModelFromFile(file)
	if err != nil {
		return err
	}

	if clean {
		// TODO: Also clean .go.dump files
		file, err := filepath.Glob(filepath.Join(outDir, filePrefix+"_*_gen.go"))
		if err != nil {
			return err
		}
		for _, f := range file {
			err := os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	generator, err := gen.NewGenerator()
	if err != nil {
		return err
	}

	model := gen.NewModel(protoModel, config)
	return generator.Generate(model, gen.NewWriteFS(outDir))
}

func readModelFromFile(file string) (*proto.AuthorizationModel, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q due to %w", file, err)
	}

	if strings.HasSuffix(file, ".fga") {
		return language.TransformDSLToProto(string(fileBytes))
	} else {
		return nil, fmt.Errorf("unsupported file format")
	}
}
