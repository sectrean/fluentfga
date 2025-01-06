package model

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	proto "github.com/openfga/api/proto/openfga/v1"
	language "github.com/openfga/language/pkg/go/transformer"
)

func ReadModelFromFile(file string) (*proto.AuthorizationModel, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", file, err)
	}

	ext := filepath.Ext(file)
	data := string(bytes)

	switch {
	case ext == ".fga":
		return language.TransformDSLToProto(data)
	case ext == ".json":
		return language.LoadJSONStringToProto(data)
	case file == "fga.mod":
		return transformModularModelToProto(file, data)

	default:
		return nil, fmt.Errorf("unsupported file format")
	}
}

func transformModularModelToProto(modFile string, data string) (*proto.AuthorizationModel, error) {
	parsedModFile, err := language.TransformModFile(data)
	if err != nil {
		return nil, fmt.Errorf("failed to transform fga.mod file: %w", err)
	}

	var moduleFiles []language.ModuleFile
	var errs []error
	directory := path.Dir(modFile)

	for _, fileName := range parsedModFile.Contents.Value {
		filePath := path.Join(directory, fileName.Value)

		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			err = fmt.Errorf("failed to read module file %s: %w", fileName.Value, err)
			errs = append(errs, err)

			continue
		}

		moduleFiles = append(moduleFiles, language.ModuleFile{
			Name:     fileName.Value,
			Contents: string(fileContents),
		})
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return language.TransformModuleFilesToModel(moduleFiles, parsedModFile.Schema.Value)
}
