package portscan

import (
	"testing"
)

func TestResolveHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		allIps  bool
		wantErr bool
		minIPs  int
		wantIP  string
	}{
		{
			name:    "IPv4 address",
			host:    "8.8.8.8",
			allIps:  false,
			wantErr: false,
			minIPs:  1,
			wantIP:  "8.8.8.8",
		},
		{
			name:    "IPv6 address",
			host:    "::1",
			allIps:  false,
			wantErr: false,
			minIPs:  1,
			wantIP:  "::1",
		},
		{
			name:    "DNS name with allIps=false",
			host:    "google.com",
			allIps:  false,
			wantErr: false,
			minIPs:  1,
			wantIP:  "",
		},
		{
			name:    "wron DNS name",
			host:    "google390.com",
			allIps:  false,
			wantErr: true,
			minIPs:  1,
			wantIP:  "",
		},
		{
			name:    "DNS name with allIps=true",
			host:    "google.com",
			allIps:  true,
			wantErr: false,
			minIPs:  2,
			wantIP:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(WithResolvedAllIPs(tt.allIps))
			ips, err := s.resolveHost(tt.host)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(ips) < tt.minIPs {
				t.Errorf("expected at least %d IPs, got %d", tt.minIPs, len(ips))
			}

			if tt.wantIP != "" {
				found := false
				for _, ip := range ips {
					if ip.String() == tt.wantIP {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected IP %s, got %v", tt.wantIP, ips)
				}
			}

			t.Logf("resolved %s to %v", tt.host, ips)
		})
	}
}

func TestUniqueHosts(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"google.com", "github.com", "yahoo.com"},
			expected: []string{"google.com", "github.com", "yahoo.com"},
		},
		{
			name:     "with duplicates",
			input:    []string{"google.com", "github.com", "google.com", "yahoo.com"},
			expected: []string{"google.com", "github.com", "yahoo.com"},
		},
		{
			name:     "with empty strings",
			input:    []string{"google.com", "", "github.com", " ", "yahoo.com"},
			expected: []string{"google.com", "github.com", "yahoo.com"},
		},
		{
			name:     "all empty",
			input:    []string{"", " ", "\t"},
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "with spaces",
			input:    []string{" google.com ", "github.com", " google.com "},
			expected: []string{"google.com", "github.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			got := s.uniqueHosts(tt.input)

			if tt.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}

			if len(got) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(got))
				return
			}

			expectedMap := make(map[string]bool)
			for _, h := range tt.expected {
				expectedMap[h] = true
			}
			for _, h := range got {
				if !expectedMap[h] {
					t.Errorf("unexpected host %s", h)
				}
			}
		})
	}
}
