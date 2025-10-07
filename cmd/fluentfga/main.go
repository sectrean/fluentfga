package main

import (
	"os"

	"github.com/spf13/cobra"
)

const ConfigFlag = "config"

func main() {
	cmd := &cobra.Command{
		Use:   "fluentfga",
		Short: "fluentfga is a tool for generating strongly typed code from an OpenFGA authorization model.",
	}
	cmd.AddCommand(NewGenerateCommand())

	cmd.PersistentFlags().String(ConfigFlag, "", "path to the configuration file")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
