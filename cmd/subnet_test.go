// Copyright (c) 2026 Keith Chu
package cmd

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSubnetInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedIP  string
		expectedErr bool
	}{
		{"CIDR /24", "192.168.1.0/24", "192.168.1.0", false},
		{"CIDR /16", "10.0.0.0/16", "10.0.0.0", false},
		{"CIDR /8", "10.0.0.0/8", "10.0.0.0", false},
		{"CIDR /32", "192.168.1.1/32", "192.168.1.1", false},
		{"CIDR /0", "0.0.0.0/0", "0.0.0.0", false},
		{"prefix only 24", "24", "0.0.0.0", false},
		{"prefix only 16", "16", "0.0.0.0", false},
		{"prefix only 8", "8", "0.0.0.0", false},
		{"prefix 0", "0", "0.0.0.0", false},
		{"invalid IP", "999.999.999.999/24", "", true},
		{"invalid prefix 33", "192.168.1.0/33", "", true},
		{"invalid prefix -1", "192.168.1.0/-1", "", true},
		{"invalid format", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, mask, err := parseSubnetInput(tt.input)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedIP, ip.String())
			assert.NotNil(t, mask)
		})
	}
}

func TestMaskFromPrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  int
		maskStr string
	}{
		{"prefix 0", 0, "00000000"},
		{"prefix 8", 8, "ff000000"},
		{"prefix 16", 16, "ffff0000"},
		{"prefix 24", 24, "ffffff00"},
		{"prefix 32", 32, "ffffffff"},
		{"prefix 1", 1, "80000000"},
		{"prefix 17", 17, "ffff8000"},
		{"prefix 25", 25, "ffffff80"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask := maskFromPrefix(tt.prefix)
			assert.Equal(t, tt.maskStr, mask.String())
		})
	}
}

func TestCalcSubnetInfo(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		prefix         int
		expectedNet    string
		expectedBroad  string
		expectedFirst  string
		expectedLast   string
		expectedHosts  int64
		expectedUsable int64
	}{
		{
			name:           "/24 network",
			ip:             "192.168.1.0",
			prefix:         24,
			expectedNet:    "192.168.1.0",
			expectedBroad:  "192.168.1.255",
			expectedFirst:  "192.168.1.1",
			expectedLast:   "192.168.1.254",
			expectedHosts:  256,
			expectedUsable: 254,
		},
		{
			name:           "/16 network",
			ip:             "10.0.0.0",
			prefix:         16,
			expectedNet:    "10.0.0.0",
			expectedBroad:  "10.0.255.255",
			expectedFirst:  "10.0.0.1",
			expectedLast:   "10.0.255.254",
			expectedHosts:  65536,
			expectedUsable: 65534,
		},
		{
			name:           "/8 network",
			ip:             "10.0.0.0",
			prefix:         8,
			expectedNet:    "10.0.0.0",
			expectedBroad:  "10.255.255.255",
			expectedFirst:  "10.0.0.1",
			expectedLast:   "10.255.255.254",
			expectedHosts:  16777216,
			expectedUsable: 16777214,
		},
		{
			name:           "/32 single host",
			ip:             "192.168.1.1",
			prefix:         32,
			expectedNet:    "192.168.1.1",
			expectedBroad:  "192.168.1.1",
			expectedFirst:  "192.168.1.2",
			expectedLast:   "192.168.1.0",
			expectedHosts:  1,
			expectedUsable: 0,
		},
		{
			name:           "/31 point to point",
			ip:             "192.168.1.0",
			prefix:         31,
			expectedNet:    "192.168.1.0",
			expectedBroad:  "192.168.1.1",
			expectedFirst:  "192.168.1.1",
			expectedLast:   "192.168.1.0",
			expectedHosts:  2,
			expectedUsable: 0,
		},
		{
			name:           "/30 subnet",
			ip:             "192.168.1.0",
			prefix:         30,
			expectedNet:    "192.168.1.0",
			expectedBroad:  "192.168.1.3",
			expectedFirst:  "192.168.1.1",
			expectedLast:   "192.168.1.2",
			expectedHosts:  4,
			expectedUsable: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip).To4()
			require.NotNil(t, ip, "failed to parse IP: %s", tt.ip)
			mask := maskFromPrefix(tt.prefix)
			network, broadcast, usable := calcSubnetInfo(ip, mask)

			assert.Equal(t, tt.expectedNet, network.String())
			assert.Equal(t, tt.expectedBroad, broadcast.String())
			assert.Equal(t, tt.expectedFirst, usable[0].String())
			assert.Equal(t, tt.expectedLast, usable[1].String())

			_, bits := mask.Size()
			totalHosts := int64(1) << uint(bits-tt.prefix)
			var usableHosts int64
			if totalHosts > 2 {
				usableHosts = totalHosts - 2
			}
			assert.Equal(t, tt.expectedHosts, totalHosts)
			assert.Equal(t, tt.expectedUsable, usableHosts)
		})
	}
}
