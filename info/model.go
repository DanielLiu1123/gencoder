package info

type Table struct {
	Schema  string
	Name    string
	Comment string
	Columns []*Column
	Indexes []*Index
}

type Column struct {
	Ordinal      int
	Name         string
	Type         string
	IsNullable   bool
	DefaultValue *string
	IsPrimaryKey bool
	Comment      *string
}

type Index struct {
	Name         string
	IsUnique     bool
	IsPrimary    bool
	IndexColumns []*IndexColumn
}

type IndexColumn struct {
	Ordinal    int
	ColumnName string
}
