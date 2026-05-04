// Copyright (c) 2026 Keith Chu
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var baseColors = struct {
	prefix lipgloss.Style
	label  lipgloss.Style
	value  lipgloss.Style
}{
	prefix: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
	label:  lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true),
	value:  lipgloss.NewStyle().Foreground(lipgloss.Color("228")),
}

var baseStyles = struct {
	result lipgloss.Style
}{
	result: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2),
}

func parseInput(s string) (int64, int) {
	s = strings.TrimSpace(s)
	base := 10
	if len(s) > 2 {
		switch strings.ToLower(s[:2]) {
		case "0x":
			base, s = 16, s[2:]
		case "0b":
			base, s = 2, s[2:]
		case "0o":
			base, s = 8, s[2:]
		}
	}
	val, _ := strconv.ParseInt(s, base, 64)
	return val, base
}

func runBaseCmd(cmd *cobra.Command, args []string) {
	val, _ := parseInput(args[0])

	rows := [][]string{
		{"DEC", fmt.Sprintf("%d", val)},
		{"HEX", fmt.Sprintf("0x%X", val)},
		{"BIN", fmt.Sprintf("0b%b", val)},
		{"OCT", fmt.Sprintf("0o%o", val)},
	}

	header := "Base Conversion Results"
	lines := []string{baseColors.prefix.Render(header) + "\n"}

	valueWidth := 0
	for _, row := range rows {
		if len(row[1]) > valueWidth {
			valueWidth = len(row[1])
		}

		label := baseColors.label.Render(row[0])
		value := baseColors.value.Render(row[1])
		lines = append(lines, fmt.Sprintf("%s  %s", label, value))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	contentWidth := 3 + valueWidth + 6
	if contentWidth < len(header)+4 {
		contentWidth = len(header) + 4
	}
	fmt.Println(baseStyles.result.Width(contentWidth).Render(content))
}

func newBaseCmd() *cobra.Command {
	baseCmd := cobra.Command{
		Use:   "base <input>",
		Short: "Convert numbers between bases",
		Long:  `Detect the base from prefix (0x=hex, 0b=binary, 0o=octal, plain=decimal) and show all conversions.`,
		Args:  cobra.ExactArgs(1),
		Run:   runBaseCmd,
	}
	return &baseCmd
}
