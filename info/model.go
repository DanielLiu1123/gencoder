package info

import "database/sql"

type TableInfo struct {
	Table        *Table         `json:"table"`
	Columns      []*Column      `json:"columns"`
	Indexes      []*Index       `json:"indexes"`
	IndexColumns []*IndexColumn `json:"index_columns"`
}

type Table struct {
	Schema    string `json:"schema"`
	TableName string `json:"table_name"`
	Comment   string `json:"comment"`
}

type Column struct {
	Schema       string         `json:"schema"`
	TableName    string         `json:"table_name"`
	Ordinal      int            `json:"ordinal"`
	ColumnName   string         `json:"column_name"`
	ColumnType   string         `json:"column_type"`
	IsNullable   bool           `json:"is_nullable"`
	DefaultValue sql.NullString `json:"default_value"`
	IsPrimaryKey bool           `json:"is_primary_key"`
	Comment      sql.NullString `json:"comment"`
}

type Index struct {
	Schema    string `json:"schema"`
	TableName string `json:"table_name"`
	IndexName string `json:"index_name"`
	IsUnique  bool   `json:"is_unique"`
	IsPrimary bool   `json:"is_primary"`
}

type IndexColumn struct {
	Schema     string `json:"schema"`
	TableName  string `json:"table_name"`
	Ordinal    int    `json:"ordinal"`
	IndexName  string `json:"index_name"`
	ColumnName string `json:"column_name"`
}
