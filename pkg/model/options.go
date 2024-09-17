package model

type GlobalOptions struct {
	Config string `json:"config" yaml:"config"` // Where to read gencoder.yaml
}

type BuildInfo struct {
	Version string
}
