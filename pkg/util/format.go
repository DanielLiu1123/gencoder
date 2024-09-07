package util

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

func ToJson(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToYaml(v any) (string, error) {
	out, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
