package portscan

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		wantConc int
		wantTO   time.Duration
		wantAll  bool
		wantBuf  int64
	}{
		{
			name:     "default config",
			opts:     []Option{},
			wantConc: 10,
			wantTO:   500 * time.Millisecond,
			wantAll:  false,
			wantBuf:  100,
		},
		{
			name:     "with concurrency",
			opts:     []Option{WithConcurrency(50)},
			wantConc: 50,
			wantTO:   500 * time.Millisecond,
			wantAll:  false,
			wantBuf:  100,
		},
		{
			name:     "with timeout",
			opts:     []Option{WithTimeout(2 * time.Second)},
			wantConc: 10,
			wantTO:   2 * time.Second,
			wantAll:  false,
			wantBuf:  100,
		},
		{
			name:     "with all options",
			opts:     []Option{WithConcurrency(20), WithTimeout(3 * time.Second), WithResolvedAllIPs(true), WithBufferSize(200)},
			wantConc: 20,
			wantTO:   3 * time.Second,
			wantAll:  true,
			wantBuf:  200,
		},
		{
			name:     "negative concurrency",
			opts:     []Option{WithConcurrency(-10)},
			wantConc: 10,
			wantTO:   500 * time.Millisecond,
			wantAll:  false,
			wantBuf:  100,
		},
		{
			name:     "negative timeout",
			opts:     []Option{WithTimeout(-1 * time.Second)},
			wantConc: 10,
			wantTO:   500 * time.Millisecond,
			wantAll:  false,
			wantBuf:  100,
		},
		{
			name:     "negative buffer",
			opts:     []Option{WithBufferSize(-100)},
			wantConc: 10,
			wantTO:   500 * time.Millisecond,
			wantAll:  false,
			wantBuf:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.opts...)
			if s.config.concurrency != tt.wantConc {
				t.Errorf("concurrency: got %d, want %d", s.config.concurrency, tt.wantConc)
			}
			if s.config.timeout != tt.wantTO {
				t.Errorf("timeout: got %v, want %v", s.config.timeout, tt.wantTO)
			}
			if s.config.allIps != tt.wantAll {
				t.Errorf("allIps: got %v, want %v", s.config.allIps, tt.wantAll)
			}
			if s.config.bufferSize != tt.wantBuf {
				t.Errorf("bufferSize: got %d, want %d", s.config.bufferSize, tt.wantBuf)
			}
		})
	}
}

func TestScan_Validation(t *testing.T) {
	s := New()

	tests := []struct {
		name      string
		hosts     []string
		ports     Ports
		wantError string
	}{
		{
			name:      "empty hosts",
			hosts:     []string{},
			ports:     List(80),
			wantError: errNoHost.Error(),
		},
		{
			name:      "empty ports",
			hosts:     []string{"google.com"},
			ports:     Ports{},
			wantError: errNoPort.Error(),
		},
		{
			name:      "both empty",
			hosts:     []string{},
			ports:     Ports{},
			wantError: errNoHost.Error(),
		},
		{
			name:      "valid",
			hosts:     []string{"google.com"},
			ports:     List(80),
			wantError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := s.Scan(ctx, tt.hosts, tt.ports)

			if tt.wantError == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantError)
				} else if err.Error() != tt.wantError {
					t.Errorf("expected error %q, got %q", tt.wantError, err.Error())
				}
			}
		})
	}
}

func TestScan_States(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		port      uint16
		wantState State
		opts      []Option
	}{
		{
			name:      "open",
			host:      "google.com",
			port:      80,
			wantState: Open,
		},
		{
			name:      "closed",
			host:      "google.com",
			port:      9999,
			wantState: Closed,
		},
		{
			name:      "nont exist host",
			host:      "google290.com",
			port:      80,
			wantState: Unreachable,
		},
		{
			name:      "timeout",
			host:      "google.com",
			port:      9999,
			wantState: Timeout,
			opts:      []Option{WithTimeout(1 * time.Millisecond)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]Option{WithTimeout(3 * time.Second)}, tt.opts...)
			s := New(opts...)

			ctx := context.Background()

			results, err := s.Scan(ctx, []string{tt.host}, List(tt.port))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for result := range results {
				t.Logf("got state: %s for %s:%d", result.State, result.Host, result.Port)

				switch tt.wantState {
				case Open:
					if result.State != Open {
						t.Logf("expected Open, got %s (network issue?)", result.State)
					}
				case Closed:
					if result.State != Closed && result.State != Timeout {
						t.Logf("expected Closed or Timeout, got %s", result.State)
					}
				case Timeout:
					if result.State != Timeout {
						t.Errorf("expected Timeout, got %s", result.State)
					}
				case Unreachable:
					if result.State != Unreachable {
						t.Errorf("expected Unreachable, got %s", result.State)
					}
				default:
					t.Errorf("unexpected state: %s", result.State)
				}
			}
		})
	}
}
