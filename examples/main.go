package main

import (
	"context"
	"fmt"
	"github.com/MashkaCoder/portscan"
	"time"
)

func main() {
	scanner := portscan.New(
		portscan.WithConcurrency(10),        // кол-во одновременных проверок, опционально, по умолчанию 10
		portscan.WithTimeout(2*time.Second), // таймаут подключения, опционально, по умолчанию 500 * time.Millisecond
		portscan.WithResolvedAllIPs(true),   // проверять все  IP для DNS-имени, опционально, по умолчанию false
		portscan.WithBufferSize(10),         // размер буфера каналов, оционально, по умолчанию 100
	)

	hosts := []string{
		"google.com",
		"github.com",
	}

	ports := portscan.Combine( // объедение нескольких наборов портов
		portscan.List(80, 443),     // создание набора портов из списка
		portscan.Range(8000, 8002), // создание набора портов из диапазона
	)

	ctx := context.Background()
	results, err := scanner.Scan(ctx, hosts, ports) // сканирование достпуности портов
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	count := 0
	for result := range results {
		count++
		fmt.Printf("[%d] %s:%d -> %s (IP: %s, за %v)\n",
			count,
			result.Host,
			result.Port,
			result.State,
			result.IP,
			result.Duration,
		)
		if result.Err != nil {
			fmt.Printf("Ошибка: %v\n", result.Err)
		}
	}
}
