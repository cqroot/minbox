// Copyright (c) 2026 Keith Chu
package cmd

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var subnetColors = struct {
	header lipgloss.Style
	label  lipgloss.Style
	value  lipgloss.Style
}{
	header: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
	label:  lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true),
	value:  lipgloss.NewStyle().Foreground(lipgloss.Color("228")),
}

var subnetStyles = struct {
	result lipgloss.Style
}{
	result: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2),
}

var cidrRegex = regexp.MustCompile(`(?i)^(\d+\.\d+\.\d+\.\d+)/(\d+)$`)
var prefixRegex = regexp.MustCompile(`(?i)^(\d+)$`)

func maskFromPrefix(prefix int) net.IPMask {
	mask := net.IPMask{}
	for i := 0; i < 4; i++ {
		if prefix >= 8 {
			mask = append(mask, 255)
			prefix -= 8
		} else if prefix > 0 {
			mask = append(mask, ^uint8(0)<<(8-prefix))
			prefix = 0
		} else {
			mask = append(mask, 0)
		}
	}
	return net.IPMask{mask[0], mask[1], mask[2], mask[3]}
}

func parseSubnetInput(s string) (net.IP, net.IPMask, error) {
	s = strings.TrimSpace(s)

	if m := cidrRegex.FindStringSubmatch(s); len(m) == 3 {
		ipStr := m[1]
		prefix, _ := strconv.Atoi(m[2])
		if prefix < 0 || prefix > 32 {
			return nil, nil, fmt.Errorf("invalid input")
		}
		ip := net.ParseIP(ipStr).To4()
		if ip == nil {
			return nil, nil, fmt.Errorf("invalid input")
		}
		return ip, maskFromPrefix(prefix), nil
	}

	if m := prefixRegex.FindStringSubmatch(s); len(m) == 2 {
		prefix, _ := strconv.Atoi(m[1])
		if prefix < 0 || prefix > 32 {
			return nil, nil, fmt.Errorf("invalid input")
		}
		ip := net.IPv4(0, 0, 0, 0)
		return ip, maskFromPrefix(prefix), nil
	}

	return nil, nil, fmt.Errorf("invalid input")
}

func calcSubnetInfo(ip net.IP, mask net.IPMask) (network, broadcast net.IP, usableIPs []net.IP) {
	network = make(net.IP, 4)
	broadcast = make(net.IP, 4)
	for i := 0; i < 4; i++ {
		network[i] = ip[i] & mask[i]
		broadcast[i] = ip[i]&mask[i] | ^mask[i]
	}

	first := make(net.IP, 4)
	last := make(net.IP, 4)
	copy(first, network)
	copy(last, broadcast)
	for i := 3; i >= 0; i-- {
		if last[i] > 0 {
			last[i]--
			break
		}
	}
	for i := 3; i >= 0; i-- {
		if first[i] < 255 {
			first[i]++
			break
		}
	}
	usableIPs = []net.IP{first, last}
	return network, broadcast, usableIPs
}

func runSubnetCmd(cmd *cobra.Command, args []string) {
	ip, mask, err := parseSubnetInput(args[0])
	if err != nil {
		return
	}

	network, broadcast, usable := calcSubnetInfo(ip, mask)

	ones, bits := mask.Size()
	wildcard := net.IPMask{}
	for i := 0; i < 4; i++ {
		wildcard = append(wildcard, ^mask[i])
	}

	totalHosts := int64(1) << uint(bits-ones)
	var usableHosts string
	if totalHosts <= 2 {
		usableHosts = "0"
	} else {
		usableHosts = fmt.Sprintf("%d", totalHosts-2)
	}

	rows := [][]string{
		{"Network", network.String()},
		{"Subnet Mask", fmt.Sprintf("%s/%d", mask.String(), ones)},
		{"Wildcard Mask", wildcard.String()},
		{"Broadcast", broadcast.String()},
		{"First Usable IP", usable[0].String()},
		{"Last Usable IP", usable[1].String()},
		{"Total Hosts", fmt.Sprintf("%d", totalHosts)},
		{"Usable Hosts", usableHosts},
	}

	header := "Subnet Calculation Results"
	lines := []string{subnetColors.header.Render(header) + "\n"}

	labelWidth := 0
	valueWidth := 0
	for _, row := range rows {
		if len(row[0]) > labelWidth {
			labelWidth = len(row[0])
		}
		if len(row[1]) > valueWidth {
			valueWidth = len(row[1])
		}
	}

	for _, row := range rows {
		label := subnetColors.label.Render(fmt.Sprintf("%-*s", labelWidth, row[0]))
		value := subnetColors.value.Render(fmt.Sprintf("%*s", valueWidth, row[1]))
		lines = append(lines, fmt.Sprintf("%s  %s", label, value))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	contentWidth := lipgloss.Width(content)

	headerWidth := lipgloss.Width(subnetColors.header.Render(header))
	if contentWidth < headerWidth {
		contentWidth = headerWidth
	}

	fmt.Println(subnetStyles.result.Width(contentWidth + 4).Render(content))
}

func newSubnetCmd() *cobra.Command {
	subnetCmd := cobra.Command{
		Use:   "subnet <input>",
		Short: "Calculate subnet information",
		Long:  `Calculate subnet information from CIDR notation (e.g., 192.168.1.0/24) or prefix length (e.g., 24).`,
		Args:  cobra.ExactArgs(1),
		Run:   runSubnetCmd,
	}
	return &subnetCmd
}
