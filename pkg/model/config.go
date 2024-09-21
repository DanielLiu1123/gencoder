package model

type Config struct {
	Templates    string            `json:"templates" yaml:"templates"` // Can be a directory or URL
	OutputMarker string            `json:"outputMarker" yaml:"outputMarker"`
	BlockMarker  BlockMarker       `json:"blockMarker" yaml:"blockMarker"`
	Databases    []*DatabaseConfig `json:"databases" yaml:"databases"`
	Properties   map[string]string `json:"properties" yaml:"properties"`
}

type DatabaseConfig struct {
	Name       string            `json:"name" yaml:"name"`
	Dsn        string            `json:"dsn" yaml:"dsn"`
	Schema     string            `json:"schema" yaml:"schema"`
	Properties map[string]string `json:"properties" yaml:"properties"`
	Tables     []*TableConfig    `json:"tables" yaml:"tables"`
}

type TableConfig struct {
	Schema        string            `json:"schema" yaml:"schema"`
	Name          string            `json:"name" yaml:"name"`
	Properties    map[string]string `json:"properties" yaml:"properties"`
	IgnoreColumns []string          `json:"ignoreColumns" yaml:"ignoreColumns"`
}

type BlockMarker struct {
	Start string `json:"start" yaml:"start"`
	End   string `json:"end" yaml:"end"`
}

func (c Config) GetTemplates() string {
	if c.Templates == "" {
		return "templates"
	}
	return c.Templates
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
