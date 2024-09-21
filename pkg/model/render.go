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

type Tpl struct {
	TemplateName      string     // template file name
	GeneratedFileName string     // generated file name, if empty, it's a partial template
	Source            string     // template source code
	Template          goja.Value // compiled template
}
