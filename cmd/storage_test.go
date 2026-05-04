// Copyright (c) 2026 Keith Chu
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStorageInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"plain bytes", "100", 0},                  // regex requires unit suffix
		{"plain bytes with spaces", "  1024  ", 0}, // regex requires unit
		{"kilobyte KB", "1KB", 1024},
		{"kilobyte Kb", "1Kb", 128}, // 1Kb = 1024 bits = 128 bytes
		{"kilobyte kB", "1kB", 1024},
		{"kilobyte kb lowercase", "1kb", 128},
		{"kilobyte k lowercase", "1k", 0}, // regex requires unit suffix (B/b/bit)
		{"megabyte MB", "1MB", 1 << 20},
		{"megabyte Mb", "1Mb", 1 << 20 / 8},
		{"gigabyte GB", "1GB", 1 << 30},
		{"gigabyte Gb", "1Gb", 1 << 30 / 8},
		{"terabyte TB", "1TB", 1 << 40},
		{"terabyte Tb", "1Tb", 1 << 40 / 8},
		{"petabyte PB", "1PB", 1 << 50},
		{"petabyte Pb", "1Pb", 1 << 50 / 8},
		{"exabyte EB", "1EB", 1 << 60},
		{"exabyte Eb", "1Eb", 1 << 60 / 8},
		{"bit with bit suffix", "1bit", 0}, // 1bit = 1/8 bytes = 0 bytes (truncated)
		{"kilobit with bit suffix", "1Kbit", 128},
		{"megabit with bit suffix", "1Mbit", 1 << 20 / 8},
		{"gigabit with bit suffix", "1Gbit", 1 << 30 / 8},
		{"terabit with bit suffix", "1Tbit", 1 << 40 / 8},
		{"petabit with bit suffix", "1Pbit", 1 << 50 / 8},
		{"exabit with bit suffix", "1Ebit", 1 << 60 / 8},
		{"decimal value", "1.5MB", int64(1.5 * float64(1<<20))},
		{"decimal kilobyte", "1.5KB", int64(1.5 * float64(1<<10))},
		{"invalid input", "invalid", 0},
		{"empty input", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStorageInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrimTrailingZeros(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no trailing zeros", "123.456", "123.456"},
		{"trailing zeros", "123.4500", "123.45"},
		{"trailing zeros single", "123.40", "123.4"},
		{"ends with dot", "123.", "123"},
		{"whole number", "123", "123"},
		{"complex case", "0.09765625", "0.09765625"},
		{"many zeros", "1.00000", "1"},
		{"single zero", "0", "0"},
		{"decimal zero", "0.0000", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimTrailingZeros(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
