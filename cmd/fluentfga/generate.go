package main

import (
	"fmt"
	"os"
	"strings"

	proto "github.com/openfga/api/proto/openfga/v1"
	language "github.com/openfga/language/pkg/go/transformer"
	"github.com/spf13/cobra"

	"github.com/johnrutherford/fluentfga"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate Code",
		Long:    "Generate code from an OpenFGA authorization model",
		Example: `fluentfga generate model.fga ./output`,
		Args:    cobra.ExactArgs(2),
		RunE:    run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	file := args[0]
	outDir := args[1]

	protoModel, err := readModelFromFile(file)
	if err != nil {
		return err
	}

	generator, err := fluentfga.NewGenerator()
	if err != nil {
		return err
	}

	config := &fluentfga.Config{
		Package: "fga",
	}

	model := fluentfga.NewModel(protoModel, config)
	return generator.Generate(model, fluentfga.NewWriteFS(outDir))
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
