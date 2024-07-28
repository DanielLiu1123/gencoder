package util

import (
	"strings"
	"unicode"
)

// PascalCase converts string to PascalCase
// e.g. "hello_world" -> "HelloWorld"
func PascalCase(str string) string {
	var result []rune
	capitalizeNext := true
	for _, r := range str {
		if r == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext {
				result = append(result, unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

// CamelCase converts string to camelCase
// e.g. "hello_world" -> "helloWorld"
func CamelCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func SnakeCase(str string) string {
	var snakeCase []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase = append(snakeCase, '_')
		}
		snakeCase = append(snakeCase, unicode.ToLower(r))
	}
	return strings.ReplaceAll(string(snakeCase), " ", "_")
}

func UpperSnakeCase(str string) string {
	var snakeCase []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase = append(snakeCase, '_')
		}
		snakeCase = append(snakeCase, unicode.ToUpper(r))
	}
	return strings.ReplaceAll(string(snakeCase), " ", "_")
}
