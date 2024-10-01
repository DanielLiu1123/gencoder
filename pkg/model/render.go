package model

import (
	"github.com/dop251/goja"
)

type RenderContext struct {
	Table          *Table            `json:"table" yaml:"table"`
	Properties     map[string]string `json:"properties" yaml:"properties"` // Merged properties
	Config         *Config           `json:"config" yaml:"config"`
	DatabaseConfig *DatabaseConfig   `json:"databaseConfig" yaml:"databaseConfig"`
	TableConfig    *TableConfig      `json:"tableConfig" yaml:"tableConfig"`
}

type FileType int

const (
	FileTypeNormal FileType = iota
	FileTypePartial
	FileTypeTemplate
)

type File struct {
	Name         string
	RelativePath string
	Content      []byte
	Type         FileType
	Output       string     // for Template FileType
	Template     goja.Value // for Template/Partial FileType
}
