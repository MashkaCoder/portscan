package portscan

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"
)

type Scanner struct {
	config *config
}

func New(opts ...Option) *Scanner {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return &Scanner{
		config: cfg,
	}
}

func (s *Scanner) Scan(ctx context.Context, hosts []string, ports Ports) (<-chan Result, error) {
	if len(hosts) == 0 {
		return nil, errNoHost
	}

	if ports.empty() {
		return nil, errNoPort
	}

	uniqHost := s.uniqueHosts(hosts)

	tasks := make(chan task, s.config.bufferSize)
	results := make(chan Result, s.config.bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < s.config.concurrency; i++ {
		wg.Add(1)
		go s.worker(ctx, tasks, results, &wg)
	}

	go func() {
		defer close(tasks)

		for _, host := range uniqHost {
			select {
			case <-ctx.Done():
				return
			default:
			}

			ips, err := s.resolveHost(host)
			if err != nil {
				for _, port := range ports.ports {
					select {
					case <-ctx.Done():
						return
					case tasks <- task{
						host: host,
						port: port,
						ip:   nil,
						err:  err,
					}:
					}
				}
				continue
			}

			for _, ip := range ips {
				for _, port := range ports.ports {
					select {
					case <-ctx.Done():
						return
					case tasks <- task{
						host: host,
						port: port,
						ip:   ip,
					}:
					}

				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil

}

func (s *Scanner) worker(ctx context.Context, tasks <-chan task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		select {
		case <-ctx.Done():
			results <- Result{
				Host:     task.host,
				Port:     task.port,
				State:    Canceled,
				Duration: 0,
				Err:      ctx.Err(),
			}
			return
		default:
		}

		if task.err != nil {
			results <- Result{
				Host:     task.host,
				Port:     task.port,
				State:    Unreachable,
				Duration: 0,
				Err:      task.err,
			}
			continue
		}

		result := s.checkPort(ctx, task.host, task.ip, task.port)
		results <- result
	}
}

func (s *Scanner) checkPort(ctx context.Context, host string, ip net.IP, port uint16) Result {
	addr := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
	dialer := &net.Dialer{
		Timeout: s.config.timeout,
	}

	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	duration := time.Since(start)

	var state State
	if err == nil {
		conn.Close()
		state = Open
	} else {
		state = processErrTcp(err)
	}

	return Result{
		Host:     host,
		IP:       ip,
		Port:     port,
		State:    state,
		Duration: duration,
		Err:      err,
	}
}

func processErrTcp(err error) State {
	if errors.Is(err, context.Canceled) {
		return Canceled
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return Timeout
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return Timeout
	}

	if errors.Is(err, syscall.ECONNREFUSED) {
		return Closed
	}

	if errors.Is(err, syscall.ENETUNREACH) || errors.Is(err, syscall.EHOSTUNREACH) {
		return Unreachable
	}

	if dnsErr, ok := err.(*net.DNSError); ok {
		if dnsErr.IsNotFound || dnsErr.IsTemporary {
			return Unreachable
		}
	}
	return Error
}
