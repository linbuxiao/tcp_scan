package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// worker 由用户 自行调用处理， 暴露一个 worker 方法 Run

type Scanner struct {
	portChan chan int
	wg       sync.WaitGroup
	ports    []int
}

func (s *Scanner) StartWorker(ctx context.Context, t time.Duration) {
	// go routine 退出
	for {
		select {
		case <-ctx.Done():
			fmt.Println("goroutine quit")
			return
		case port := <-s.portChan:
			if _, err := net.DialTimeout("tcp", fmt.Sprintf("www.baidu.com:%d", port), t); err == nil {
				go fmt.Println("success:", port)
			}
			s.wg.Done()
		}
	}
}

func (s *Scanner) Run() {
	fmt.Println("scanner is beginning")
	for i := 0; i < len(s.ports); i++ {
		port := s.ports[i]
		s.wg.Add(1)
		s.portChan <- port
	}
	s.wg.Wait()
	fmt.Println("scanner is closed")
}

func NewRangeScanner(start, end int) *Scanner {
	var ports []int
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return NewScanner(ports)
}

// NewScanner 根据提供的 ports 范围返回扫描器
func NewScanner(ports []int) *Scanner {
	return &Scanner{
		portChan: make(chan int),
		ports:    ports,
	}
}
