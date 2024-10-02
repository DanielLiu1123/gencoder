package model

type Config struct {
	Templates    string            `json:"templates" yaml:"templates" jsonschema:"description=The dir or URL to store templates,example=templates"`
	OutputMarker string            `json:"outputMarker" yaml:"outputMarker" jsonschema:"description=The magic comment to identify the generated file,example=@gencoder.generated:"`
	BlockMarker  BlockMarker       `json:"blockMarker" yaml:"blockMarker" jsonschema:"description=The block marker to identify the generated block"`
	Databases    []*DatabaseConfig `json:"databases" yaml:"databases" jsonschema:"description=The list of databases"`
	Properties   map[string]string `json:"properties" yaml:"properties" jsonschema:"description=The global properties,will be overridden by properties in databases and tables"`
	Output       string            `json:"output" yaml:"output" jsonschema:"description=The output directory for generated files,example=./output"`
}

type DatabaseConfig struct {
	Name       string            `json:"name" yaml:"name" jsonschema:"description=The name of the database,example=mydb"`
	Dsn        string            `json:"dsn" yaml:"dsn" jsonschema:"description=The database connection string\\, gencoder uses [xo/dburl](https://github.com/xo/dburl) to provides a uniform way to parse database connections,example=mysql://user:password@localhost:3306/dbname,required"`
	Schema     string            `json:"schema" yaml:"schema" jsonschema:"description=The schema of the database,example=public"`
	Properties map[string]string `json:"properties" yaml:"properties" jsonschema:"description=Properties specific to the database"`
	Tables     []*TableConfig    `json:"tables" yaml:"tables" jsonschema:"description=The list of tables in the database"`
}

type TableConfig struct {
	Schema        string            `json:"schema" yaml:"schema" jsonschema:"description=The schema of the table,example=public"`
	Name          string            `json:"name" yaml:"name" jsonschema:"description=The name of the table,example=user,required"`
	Properties    map[string]string `json:"properties" yaml:"properties" jsonschema:"description=Properties specific to the table"`
	IgnoreColumns []string          `json:"ignoreColumns" yaml:"ignoreColumns" jsonschema:"description=The list of columns to ignore,example=[password,secret]"`
}

type BlockMarker struct {
	Start string `json:"start" yaml:"start" jsonschema:"description=The start marker for code block,example=@gencoder.block.start:"`
	End   string `json:"end" yaml:"end" jsonschema:"description=The end marker for code block,example=@gencoder.block.end:"`
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
