package info

import "database/sql"

type Table struct {
	Schema  string    `json:"schema"`
	Name    string    `json:"name"`
	Comment string    `json:"comment"`
	Columns []*Column `json:"columns"`
	Indexes []*Index  `json:"indexes"`
}

type Column struct {
	Ordinal      int            `json:"ordinal"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	IsNullable   bool           `json:"is_nullable"`
	DefaultValue sql.NullString `json:"default_value"`
	IsPrimaryKey bool           `json:"is_primary_key"`
	Comment      sql.NullString `json:"comment"`
}

type Index struct {
	Schema   string   `json:"schema"`
	Table    string   `json:"table"`
	Name     string   `json:"name"`
	IsUnique bool     `json:"is_unique"`
	Columns  []string `json:"columns"`
}
