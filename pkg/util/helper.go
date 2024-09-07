package util

import (
	"strings"
	"unicode"
)

// ToSnakeCase converts "userName" or "UserName" to "user_name"
//
// Example:
// ToSnakeCase("userName") => "user_name"
// ToSnakeCase("UserName") => "user_name"
func ToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) {
			// Add an underscore before uppercase letters, except at the start
			if i > 0 && !unicode.IsUpper(rune(input[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ToCamelCase converts "user_name" to "userName"
func ToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			// Convert the first character of each word to uppercase
			words[i] = strings.ToUpper(string(words[i][0])) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

// ToPascalCase converts "user_name" or "userName" to "UserName"
func ToPascalCase(input string) string {
	// Convert to CamelCase first if input is snake_case
	if strings.Contains(input, "_") {
		input = ToCamelCase(input)
	}
	// Capitalize the first character
	if len(input) > 0 {
		input = strings.ToUpper(string(input[0])) + input[1:]
	}
	return input
}
