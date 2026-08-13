package portscan

import (
	"sort"
)

type Ports struct {
	ports []uint16
}

func Range(start, end uint16) Ports {
	if start > end {
		start, end = end, start
	}

	if start < 1 {
		start = 1
	}

	ports := make([]uint16, 0, end-start+1)
	for p := start; p <= end; p++ {
		ports = append(ports, p)
	}

	return Ports{ports: ports}
}

func List(ports ...uint16) Ports {
	return newPorts(ports)
}

func Combine(sets ...Ports) Ports {
	all := make([]uint16, 0)
	for _, set := range sets {
		all = append(all, set.ports...)
	}

	return newPorts(all)
}

func newPorts(ports []uint16) Ports {
	if len(ports) == 0 {
		return Ports{ports: []uint16{}}
	}

	seen := make(map[uint16]bool)
	valid := make([]uint16, 0, len(ports))

	for _, p := range ports {
		if p < 1 {
			continue
		}

		if !seen[p] {
			seen[p] = true
			valid = append(valid, p)
		}
	}

	sort.Slice(valid, func(i, j int) bool { return valid[i] < valid[j] })
	return Ports{ports: valid}
}

func (p Ports) empty() bool {
	return len(p.ports) == 0
}
