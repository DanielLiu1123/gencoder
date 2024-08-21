package info

type Config struct {
	TemplatesDir string       `yaml:"templates-dir"`
	BlockMarker  *BlockMarker `yaml:"block-marker"`
	Databases    []*struct {
		Name       string      `yaml:"name"`
		Dsn        string      `yaml:"dsn"`
		Schema     string      `yaml:"schema"`
		Properties []*Property `yaml:"properties"`
		Tables     []*struct {
			Schema          string      `yaml:"schema"`
			Name            string      `yaml:"name"`
			ExtraProperties []*Property `yaml:"properties"`
			IgnoreColumns   []string    `yaml:"ignore-columns"`
		}
	}
}

type Property struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type BlockMarker struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}
