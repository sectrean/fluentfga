package gen

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// NewConfig returns a default configuration.
func NewConfig() *Config {
	c := &Config{}
	c.applyDefaults()

	return c
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *Config
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&config)
	case ".json":
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
	default:
		return nil, errors.New("unsupported config file format")
	}

	if config != nil {
		config.applyDefaults()
	}

	return config, err
}

type Config struct {
	Package    string                 `yaml:"package" json:"package"`
	FilePrefix string                 `yaml:"file_prefix" json:"filePrefix"`
	Types      map[string]*TypeConfig `yaml:"types" json:"types"`

	FormatTypeName func(string) string `yaml:"-" json:"-"`
	FormatIDName   func(string) string `yaml:"-" json:"-"`
}

type TypeConfig struct {
	Type   string `yaml:"type" json:"type"`
	IDName string `yaml:"id_name" json:"idName"`

	// TODO: Support types with a package name, e.g., uuid.UUID
	IDType string `yaml:"id_type" json:"idType"`
}

func (c *Config) applyDefaults() {
	if c.Package == "" {
		c.Package = "model"
	}
	if c.FilePrefix == "" {
		c.FilePrefix = "fga"
	}
	if c.FormatTypeName == nil {
		c.FormatTypeName = TitleCase
	}
	if c.FormatIDName == nil {
		c.FormatIDName = ID
	}
}

func (c *Config) TypeName(name string) string {
	if tc, ok := c.Types[name]; ok {
		if tc.Type != "" {
			return tc.Type
		}
	}
	return c.FormatTypeName(name)
}
