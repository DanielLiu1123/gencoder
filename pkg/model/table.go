package model

type Table struct {
	Name    string    `json:"name" yaml:"name"`
	Schema  string    `json:"schema" yaml:"schema"`
	Comment *string   `json:"comment" yaml:"comment"`
	Columns []*Column `json:"columns" yaml:"columns"`
	Indexes []*Index  `json:"indexes" yaml:"indexes"`
}

type Column struct {
	Name         string  `json:"name" yaml:"name"`
	Ordinal      int     `json:"ordinal" yaml:"ordinal"`
	Type         string  `json:"type" yaml:"type"`
	IsNullable   bool    `json:"isNullable" yaml:"isNullable"`
	DefaultValue *string `json:"defaultValue" yaml:"defaultValue"`
	IsPrimaryKey bool    `json:"isPrimaryKey" yaml:"isPrimaryKey"`
	Comment      *string `json:"comment" yaml:"comment"`
}

type Index struct {
	Name      string         `json:"name" yaml:"name"`
	IsUnique  bool           `json:"isUnique" yaml:"isUnique"`
	IsPrimary bool           `json:"isPrimary" yaml:"isPrimary"`
	Columns   []*IndexColumn `json:"columns" yaml:"columns"`
}

type IndexColumn struct {
	Ordinal int    `json:"ordinal" yaml:"ordinal"`
	Name    string `json:"name" yaml:"name"`
}
