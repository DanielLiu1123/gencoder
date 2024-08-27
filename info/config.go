package info

type Config struct {
	TemplatesDir string            `yaml:"templates-dir"`
	OutputMarker string            `yaml:"output-marker"`
	BlockMarker  *BlockMarker      `yaml:"block-marker"`
	Databases    []*DatabaseConfig `yaml:"databases"`
}

type DatabaseConfig struct {
	Name       string            `yaml:"name"`
	Dsn        string            `yaml:"dsn"`
	Schema     string            `yaml:"schema"`
	Properties map[string]string `yaml:"properties"`
	Tables     []*TableConfig    `yaml:"tables"`
}

type TableConfig struct {
	Schema        string            `yaml:"schema"`
	Name          string            `yaml:"name"`
	Properties    map[string]string `yaml:"properties"`
	IgnoreColumns []string          `yaml:"ignore-columns"`
}

type BlockMarker struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}
