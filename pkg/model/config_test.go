package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetHelpers(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected []string
	}{
		{
			name: "when new Helpers field is set, should return Helpers",
			config: Config{
				Helpers:       []string{"helper1.js", "helper2.js"},
				ImportHelpers: []string{"old1.js", "old2.js"},
			},
			expected: []string{"helper1.js", "helper2.js"},
		},
		{
			name: "when only ImportHelpers field is set, should return ImportHelpers",
			config: Config{
				ImportHelpers: []string{"old1.js", "old2.js"},
			},
			expected: []string{"old1.js", "old2.js"},
		},
		{
			name: "when both fields are empty, should return empty slice",
			config: Config{
				Helpers:       []string{},
				ImportHelpers: []string{},
			},
			expected: []string{},
		},
		{
			name: "when both fields are nil, should return nil",
			config: Config{
				Helpers:       nil,
				ImportHelpers: nil,
			},
			expected: nil,
		},
		{
			name: "when Helpers is empty but ImportHelpers has values, should return ImportHelpers",
			config: Config{
				Helpers:       []string{},
				ImportHelpers: []string{"old1.js"},
			},
			expected: []string{"old1.js"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetHelpers()
			assert.Equal(t, tt.expected, result)
		})
	}
}
