package portscan

import (
	"net"
	"time"
)

type State string

const (
	Open        State = "open"
	Closed      State = "closed"
	Timeout     State = "timeout"
	Unreachable State = "unreachable"
	Canceled    State = "canceled"
	Error       State = "error"
)

type Result struct {
	Host     string
	IP       net.IP
	Port     uint16
	State    State
	Duration time.Duration
	Err      error
}

type task struct {
	host string
	ip   net.IP
	port uint16
	err  error
}
