package model

type Config struct {
	Templates    string            `json:"templates,omitempty" yaml:"templates,omitempty" jsonschema:"description=The dir or URL to store templates,example=templates"`
	OutputMarker string            `json:"outputMarker,omitempty" yaml:"outputMarker,omitempty" jsonschema:"description=The magic comment to identify the generated file,example=@gencoder.generated:"`
	BlockMarker  BlockMarker       `json:"blockMarker,omitempty" yaml:"blockMarker,omitempty" jsonschema:"description=The block marker to identify the generated block"`
	Databases    []*DatabaseConfig `json:"databases,omitempty" yaml:"databases,omitempty" jsonschema:"description=The list of databases"`
	Properties   map[string]string `json:"properties,omitempty" yaml:"properties,omitempty" jsonschema:"description=The global properties,will be overridden by properties in databases and tables"`
	Output       string            `json:"output,omitempty" yaml:"output,omitempty" jsonschema:"description=The output directory for generated files,example=./output"`
}

type DatabaseConfig struct {
	Name       string            `json:"name,omitempty" yaml:"name,omitempty" jsonschema:"description=The name of the database,example=mydb"`
	Dsn        string            `json:"dsn,omitempty" yaml:"dsn,omitempty" jsonschema:"description=The database connection string\\, gencoder uses [xo/dburl](https://github.com/xo/dburl) to provides a uniform way to parse database connections,example=mysql://user:password@localhost:3306/dbname,required"`
	Schema     string            `json:"schema,omitempty" yaml:"schema,omitempty" jsonschema:"description=The schema of the database,example=public"`
	Properties map[string]string `json:"properties,omitempty" yaml:"properties,omitempty" jsonschema:"description=Properties specific to the database"`
	Tables     []*TableConfig    `json:"tables,omitempty" yaml:"tables,omitempty" jsonschema:"description=The list of tables in the database"`
}

type TableConfig struct {
	Schema        string            `json:"schema,omitempty" yaml:"schema,omitempty" jsonschema:"description=The schema of the table,example=public"`
	Name          string            `json:"name,omitempty" yaml:"name,omitempty" jsonschema:"description=The name of the table,example=user,required"`
	Properties    map[string]string `json:"properties,omitempty" yaml:"properties,omitempty" jsonschema:"description=Properties specific to the table"`
	IgnoreColumns []string          `json:"ignoreColumns,omitempty" yaml:"ignoreColumns,omitempty" jsonschema:"description=The list of columns to ignore,example=password"`
}

type BlockMarker struct {
	Start string `json:"start,omitempty" yaml:"start,omitempty" jsonschema:"description=The start marker for code block,example=@gencoder.block.start:"`
	End   string `json:"end,omitempty" yaml:"end,omitempty" jsonschema:"description=The end marker for code block,example=@gencoder.block.end:"`
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
