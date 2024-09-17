package model

import "time"

type GlobalOptions struct {
	Config string `json:"config" yaml:"config"` // Where to read gencoder.yaml
}

type BuildInfo struct {
	Version   string
	BuildTime time.Time
}
