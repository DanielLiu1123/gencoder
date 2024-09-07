package model

type RenderContext struct {
	Table      *Table            `json:"table" yaml:"table"`
	Properties map[string]string `json:"properties" yaml:"properties"`
}
