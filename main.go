package main

import (
	"context"
	"fmt"
	"github.com/bells307/qwatro/port_scanner"
	"time"
)

func main() {
	scanner := port_scanner.
		NewScannerBuilder().
		IP("192.168.65.82").
		PortRange(port_scanner.OrderedRange(8000, 9000)).
		Tcp(300 * time.Millisecond).
		NumWorkers(500).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ch := scanner.Run(ctx)

	for {
		r, more := <-ch
		if more {
			fmt.Printf("%s:%d\n", r.IP, r.Port)
		} else {
			break
		}
	}

}
