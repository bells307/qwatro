package port_scanner

import (
	"fmt"
	"net"
	"time"
)

// TCP-сканирование
type tcpStrategy struct {
	// Максимальное время ожидания установления подключения
	respTimeout time.Duration
}

func newTcpStrategy(respTimeout time.Duration) *tcpStrategy {
	return &tcpStrategy{respTimeout}
}

func (s *tcpStrategy) scan(task *scanTask) *scanTaskSuccess {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", task.ip, task.port), s.respTimeout)
	if err == nil {
		conn.Close()
		return &scanTaskSuccess{
			IP:   task.ip,
			Port: task.port,
		}
	} else {
		return nil
	}
}
