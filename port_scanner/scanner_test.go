package port_scanner

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTcpScan(t *testing.T) {
	ip := "127.0.0.1"
	port := uint16(9999)
	go runServer(t, fmt.Sprintf("%s:%d", ip, port))

	scanner := NewScannerBuilder().
		IP(ip).
		PortRange(SpecificRange(port)).
		Tcp(300 * time.Millisecond).
		Build()

	ch := scanner.Run(context.Background())
	res := <-ch

	expectedRes := scanTaskSuccess{
		IP:   ip,
		Port: port,
	}

	if *res != expectedRes {
		t.Errorf("expected result: %v, got: %v", expectedRes, *res)
	}
}

func runServer(t *testing.T, addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
		}

		go func() {
			reader := bufio.NewReader(conn)
			for {
				message, err := reader.ReadString('\n')
				if err != nil {
					conn.Close()
					return
				}
				fmt.Printf("Message incoming: %s", string(message))
				conn.Write([]byte("Message received.\n"))
			}
		}()
	}
}
