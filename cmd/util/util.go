package util

import (
	"unicode"
)

// ToCamelCase converts snake_case to CamelCase
func ToCamelCase(str string) string {
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

// ToSnakeCase converts CamelCase to snake_case
func ToSnakeCase(str string) string {
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
