package util

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"log"
)

// ToJson converts a value to a JSON string.
func ToJson(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// ToMap converts a value to a map[string]interface{}.
// If the value is already a map[string]interface{}, it is returned as is.
func ToMap(v any) map[string]interface{} {
	m, ok := v.(map[string]interface{})
	if ok {
		return m
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(ToJson(v)), &result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// ToYaml converts a value to a YAML string.
func ToYaml(v any) string {
	out, err := yaml.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
