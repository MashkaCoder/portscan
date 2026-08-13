package portscan

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint16
		expected []uint16
	}{
		{
			name:     "single port",
			input:    []uint16{80},
			expected: []uint16{80},
		},
		{
			name:     "multiple ports",
			input:    []uint16{80, 443, 8080},
			expected: []uint16{80, 443, 8080},
		},
		{
			name:     "duplicate ports",
			input:    []uint16{80, 443, 80, 443, 8080},
			expected: []uint16{80, 443, 8080},
		},
		{
			name:     "contains invalid ports",
			input:    []uint16{0, 80, 443},
			expected: []uint16{80, 443},
		},
		{
			name:     "all invalid ports",
			input:    []uint16{0},
			expected: []uint16{},
		},
		{
			name:     "empty list",
			input:    []uint16{},
			expected: []uint16{},
		},
		{
			name:     "max port",
			input:    []uint16{65535},
			expected: []uint16{65535},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ports := List(tt.input...)
			if !reflect.DeepEqual(ports.ports, tt.expected) {
				t.Errorf("got %v, want %v", ports.ports, tt.expected)
			}
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		start    uint16
		end      uint16
		expected []uint16
	}{
		{
			name:     "normal range",
			start:    1,
			end:      5,
			expected: []uint16{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse range",
			start:    5,
			end:      1,
			expected: []uint16{1, 2, 3, 4, 5},
		},
		{
			name:     "range with 0",
			start:    0,
			end:      5,
			expected: []uint16{1, 2, 3, 4, 5},
		},
		{
			name:     "single port range",
			start:    80,
			end:      80,
			expected: []uint16{80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ports := Range(tt.start, tt.end)
			if !reflect.DeepEqual(ports.ports, tt.expected) {
				t.Errorf("got %v, want %v", ports.ports, tt.expected)
			}

		})
	}
}

func TestCombine(t *testing.T) {
	tests := []struct {
		name     string
		sets     []Ports
		expected []uint16
	}{
		{
			name: "combine two lists",
			sets: []Ports{
				List(80, 443),
				List(8080, 8443),
			},
			expected: []uint16{80, 443, 8080, 8443},
		},
		{
			name: "combine list and range",
			sets: []Ports{
				List(80, 443),
				Range(3000, 3002),
			},
			expected: []uint16{80, 443, 3000, 3001, 3002},
		},
		{
			name: "combine lists with duplicates",
			sets: []Ports{
				List(80, 443),
				List(80, 8080),
			},
			expected: []uint16{80, 443, 8080},
		},
		{
			name: "combine multiple sets",
			sets: []Ports{
				List(80, 443),
				Range(3000, 3002),
				List(8080, 8443),
			},
			expected: []uint16{80, 443, 3000, 3001, 3002, 8080, 8443},
		},
		{
			name: "combine empty sets",
			sets: []Ports{
				List(),
				List(80),
				List(),
			},
			expected: []uint16{80},
		},
		{
			name:     "combine no sets",
			sets:     []Ports{},
			expected: []uint16{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ports := Combine(tt.sets...)
			if !reflect.DeepEqual(ports.ports, tt.expected) {
				t.Errorf("got %v, want %v", ports.ports, tt.expected)
			}
		})
	}
}

func TestPortsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		ports    Ports
		expected bool
	}{
		{
			name:     "empty list",
			ports:    List(),
			expected: true,
		},
		{
			name:     "list with valid port",
			ports:    List(80),
			expected: false,
		},
		{
			name:     "list with invalid only",
			ports:    List(0),
			expected: true,
		},
		{
			name:     "range with valid ports",
			ports:    Range(1, 10),
			expected: false,
		},
		{
			name:     "empty range",
			ports:    Range(0, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ports.empty()
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
