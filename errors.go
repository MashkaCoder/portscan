package portscan

import "errors"

var (
	errNoHost = errors.New("no hosts given")
	errNoPort = errors.New("no ports given")
)
