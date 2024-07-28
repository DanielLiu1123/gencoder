package info

type Config struct {
	TemplatesDir string `yaml:"templates-dir"`
	BlockMarker  *struct {
		Start string `yaml:"start"`
		End   string `yaml:"end"`
	} `yaml:"block-marker"`
	Databases []*struct {
		Name            string             `yaml:"name"`
		Dsn             string             `yaml:"dsn"`
		ExtraProperties []*ExtraProperties `yaml:"extra-properties"`
		Tables          []*struct {
			Name            string             `yaml:"name"`
			ExtraProperties []*ExtraProperties `yaml:"extra-properties"`
			IgnoreColumns   []string           `yaml:"ignore-columns"`
		}
	}
}

type ExtraProperties struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
