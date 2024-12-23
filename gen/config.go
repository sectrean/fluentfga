package gen

type Config struct {
	Package string                `yaml:"package"`
	Types   map[string]TypeConfig `yaml:"types"`
}

type TypeConfig struct {
	Type   string `yaml:"type"`
	IDName string `yaml:"id_name"`
	IDType string `yaml:"id_type"`
}

func (c Config) TypeName(name string) (bool, string) {
	tc, ok := c.Types[name]
	if !ok {
		return false, ""
	}

	return true, tc.Type
}
