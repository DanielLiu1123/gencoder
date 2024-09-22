package util

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestToJson(t *testing.T) {
	v := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	expectedJson := `{
  "age": 30,
  "name": "John"
}`

	jsonStr := ToJson(v)

	assert.Equal(t, expectedJson, jsonStr)
}

func TestToMap(t *testing.T) {
	v := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	mappedValue := ToMap(v)
	assert.True(t, reflect.DeepEqual(mappedValue, v))

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	person := Person{
		Name: "John",
		Age:  30,
	}
	expectedMap := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
	}
	mappedValueFromStruct := ToMap(person)
	assert.True(t, reflect.DeepEqual(mappedValueFromStruct, expectedMap))
}

func TestToYaml(t *testing.T) {
	v := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	expectedYaml := "age: 30\nname: John\n"

	yamlStr := ToYaml(v)

	assert.Equal(t, expectedYaml, yamlStr)
}
