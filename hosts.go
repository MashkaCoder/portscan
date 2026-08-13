package portscan

import (
	"fmt"
	"net"
	"strings"
)

func (s *Scanner) resolveHost(host string) ([]net.IP, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		return []net.IP{ip}, nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("DNS resolution failed for %s: %w", host, err)
	}

	if s.config.allIps {
		return ips, nil
	}

	return ips[0:1], nil
}

func (s *Scanner) uniqueHosts(hosts []string) []string {
	if len(hosts) == 0 {
		return nil
	}

	unique := make(map[string]bool)
	uniqueHosts := make([]string, 0, len(hosts))
	for _, h := range hosts {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		if !unique[h] {
			unique[h] = true
			uniqueHosts = append(uniqueHosts, h)
		}
	}

	return uniqueHosts
}
