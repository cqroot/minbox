// Copyright (c) 2026 Keith Chu
package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var storageColors = struct {
	header lipgloss.Style
	label  lipgloss.Style
	value  lipgloss.Style
}{
	header: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
	label:  lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true),
	value:  lipgloss.NewStyle().Foreground(lipgloss.Color("228")),
}

var storageStyles = struct {
	result lipgloss.Style
}{
	result: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2),
}

type unitInfo struct {
	name  string
	label string
	exp   int64
}

var byteUnits = []unitInfo{
	{"B", "Byte", 1},
	{"KB", "Kilobyte", 1 << 10},
	{"MB", "Megabyte", 1 << 20},
	{"GB", "Gigabyte", 1 << 30},
	{"TB", "Terabyte", 1 << 40},
	{"PB", "Petabyte", 1 << 50},
	{"EB", "Exabyte", 1 << 60},
}

var bitUnits = []unitInfo{
	{"b", "Bit", 1},
	{"Kb", "Kilobit", 1 << 10},
	{"Mb", "Megabit", 1 << 20},
	{"Gb", "Gigabit", 1 << 30},
	{"Tb", "Terabit", 1 << 40},
	{"Pb", "Petabit", 1 << 50},
	{"Eb", "Exabit", 1 << 60},
}

var unitRegex = regexp.MustCompile(`(?i)^(\d+(?:\.\d+)?)\s*([KMGTPEkmgtpe]?(?:B|b|bit))$`)

func parseStorageInput(s string) int64 {
	s = strings.TrimSpace(s)
	matches := unitRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return 0
	}

	value, _ := strconv.ParseFloat(matches[1], 64)

	unit := matches[2]
	isBit := strings.HasSuffix(strings.ToLower(unit), "bit") || unit[len(unit)-1] == 'b'
	if isBit {
		unit = strings.TrimSuffix(strings.TrimSuffix(unit, "bit"), "b")
	} else {
		unit = strings.TrimSuffix(unit, "B")
	}
	unit = strings.ToUpper(unit)

	var exp int64
	switch unit {
	case "":
		exp = 1
	case "K":
		exp = 1 << 10
	case "M":
		exp = 1 << 20
	case "G":
		exp = 1 << 30
	case "T":
		exp = 1 << 40
	case "P":
		exp = 1 << 50
	case "E":
		exp = 1 << 60
	default:
		return 0
	}

	val := value * float64(exp)
	if isBit {
		val = val / 8
	}
	return int64(val)
}

func trimTrailingZeros(s string) string {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

func formatFull(val float64, unit string) string {
	if val >= 1 {
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d %s", int64(val), unit)
		}
		s := fmt.Sprintf("%.6f", val)
		return trimTrailingZeros(s) + " " + unit
	}
	if val == 0 {
		return "0 " + unit
	}
	s := fmt.Sprintf("%.12f", val)
	s = trimTrailingZeros(s)
	return s + " " + unit
}

func runStorageCmd(cmd *cobra.Command, args []string) {
	inputBytes := parseStorageInput(args[0])

	bytes := inputBytes

	byteCols := make([][2]string, len(byteUnits))
	for i, u := range byteUnits {
		val := float64(bytes) / float64(u.exp)
		byteCols[i] = [2]string{u.label, formatFull(val, u.name)}
	}

	bitCols := make([][2]string, len(bitUnits))
	for i, u := range bitUnits {
		val := float64(bytes*8) / float64(u.exp)
		bitCols[i] = [2]string{u.label, formatFull(val, u.name)}
	}

	labelColWidth := 0
	valColWidth := 0
	for _, col := range byteCols {
		if len(col[0]) > labelColWidth {
			labelColWidth = len(col[0])
		}
		if len(col[1]) > valColWidth {
			valColWidth = len(col[1])
		}
	}
	for _, col := range bitCols {
		if len(col[0]) > labelColWidth {
			labelColWidth = len(col[0])
		}
		if len(col[1]) > valColWidth {
			valColWidth = len(col[1])
		}
	}

	colGap := 4
	sep := strings.Repeat(" ", colGap)

	header := "Storage Conversion Results"
	lines := []string{storageColors.header.Render(header) + "\n"}
	for i := range byteCols {
		bLabel := storageColors.label.Render(fmt.Sprintf("%-*s", labelColWidth, byteCols[i][0]))
		bVal := storageColors.value.Render(fmt.Sprintf("%*s", valColWidth, byteCols[i][1]))
		bitLabel := storageColors.label.Render(fmt.Sprintf("%-*s", labelColWidth, bitCols[i][0]))
		bitVal := storageColors.value.Render(fmt.Sprintf("%*s", valColWidth, bitCols[i][1]))
		lines = append(lines, fmt.Sprintf("%s  %s%s%s  %s", bLabel, bVal, sep, bitLabel, bitVal))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	contentWidth := lipgloss.Width(content)

	headerWidth := lipgloss.Width(storageColors.header.Render(header))
	if contentWidth < headerWidth {
		contentWidth = headerWidth
	}

	fmt.Println(storageStyles.result.Width(contentWidth + 4).Render(content))
}

func newStorageCmd() *cobra.Command {
	storageCmd := cobra.Command{
		Use:   "storage <input>",
		Short: "Convert storage units",
		Long:  `Convert storage units (B, KB, MB, GB, TB, PB, EB, b, Kb, Mb, Gb, Tb, Pb, Eb, bit, Kbit, Mbit, Gbit, Tbit, Pbit, Ebit) to all other units.`,
		Args:  cobra.ExactArgs(1),
		Run:   runStorageCmd,
	}
	return &storageCmd
}
