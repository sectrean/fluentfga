package main

import (
	"os"

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

	config := gen.NewConfig()

	configFile, err := cmd.Flags().GetString(ConfigFlag)
	if err != nil {
		return err
	}
	if configFile != "" {
		config, err = gen.LoadConfig(configFile)
		if err != nil {
			return err
		}
	}

	clean, err := cmd.PersistentFlags().GetBool(CleanFlag)
	if err != nil {
		return err
	}

	protoModel, err := model.ReadModelFromFile(file)
	if err != nil {
		return err
	}

	generator, err := gen.NewGenerator(config)
	if err != nil {
		return err
	}

	output, err := os.OpenRoot(outDir)
	if err != nil {
		return err
	}

	if clean {
		err := generator.CleanOutput(output)
		if err != nil {
			return err
		}
	}

	model := gen.NewModel(protoModel, config)
	return generator.Generate(model, output)
}
