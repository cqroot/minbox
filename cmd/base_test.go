// Copyright (c) 2026 Keith Chu
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		base     int
	}{
		{"decimal plain", "100", 100, 10},
		{"decimal with spaces", "  255  ", 255, 10},
		{"hex 0x prefix", "0xFF", 255, 16},
		{"hex 0X prefix", "0XFF", 255, 16},
		{"binary 0b prefix", "0b1010", 10, 2},
		{"binary 0B prefix", "0B1010", 10, 2},
		{"octal 0o prefix", "0o77", 63, 8},
		{"octal 0O prefix", "0O77", 63, 8},
		{"hex lowercase", "0xff", 255, 16},
		{"large decimal", "9223372036854775807", 9223372036854775807, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, base := parseInput(tt.input)
			assert.Equal(t, tt.expected, val)
			assert.Equal(t, tt.base, base)
		})
	}
}
