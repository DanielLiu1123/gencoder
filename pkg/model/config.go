package model

type Config struct {
	TemplatesDir string            `yaml:"templatesDir"`
	OutputMarker string            `yaml:"outputMarker"`
	BlockMarker  BlockMarker       `yaml:"blockMarker"`
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
	IgnoreColumns []string          `yaml:"ignoreColumns"`
}

type BlockMarker struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

func (c Config) GetTemplatesDir() string {
	if c.TemplatesDir == "" {
		return "templates"
	}
	return c.TemplatesDir
}

func (c Config) GetOutputMarker() string {
	if c.OutputMarker == "" {
		return "@gencoder.generated:"
	}
	return c.OutputMarker
}

func (e BlockMarker) GetStart() string {
	if e.Start == "" {
		return "@gencoder.block.start:"
	}
	return e.Start
}

func (e BlockMarker) GetEnd() string {
	if e.End == "" {
		return "@gencoder.block.end:"
	}
	return e.End
}
