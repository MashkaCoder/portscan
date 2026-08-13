package portscan

import (
	"time"
)

type config struct {
	concurrency int
	timeout     time.Duration
	allIps      bool
	bufferSize  int64
}

type Option func(*config)

func defaultConfig() *config {
	return &config{
		concurrency: 10,
		timeout:     500 * time.Millisecond,
		allIps:      false,
		bufferSize:  100,
	}
}

func WithConcurrency(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.concurrency = n
		}
	}

}

func WithTimeout(t time.Duration) Option {
	return func(c *config) {
		if t > 0 {
			c.timeout = t
		}
	}
}

func WithResolvedAllIPs(all bool) Option {
	return func(c *config) {
		if all {
			c.allIps = true
		}
	}
}

func WithBufferSize(size int64) Option {
	return func(c *config) {
		if size > 0 {
			c.bufferSize = size
		}
	}
}
