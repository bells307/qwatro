package main

import (
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

type Task struct {
	ip   string
	port uint16
}

const workersNum = 500

func main() {
	tasks := make(chan Task, workersNum)
	successes := make(chan Task)

	go portScanning("localhost", tasks)

	done := make(chan struct{})
	go processSuccesses(successes, done)

	wg := sync.WaitGroup{}
	for i := 0; i < workersNum; i++ {
		go func() {
			wg.Add(1)
			for {
				task, more := <-tasks
				if more {
					scan(task, successes)
				} else {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(successes)
	<-done
}

func portScanning(host string, tasksChan chan Task) {
	for i := 1; i <= math.MaxUint16; i++ {
		port := uint16(i)
		tasksChan <- Task{
			ip:   host,
			port: port,
		}
	}
	close(tasksChan)
}

func processSuccesses(successes chan Task, done chan struct{}) {
	count := 0
	for {
		s, more := <-successes
		if more {
			fmt.Printf("%s:%d/tcp\n", s.ip, s.port)
			count++
		} else {
			fmt.Printf("available ports: %d", count)
			done <- struct{}{}
			return
		}
	}
}

func scan(task Task, successes chan Task) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", task.ip, task.port), 300*time.Millisecond)
	if err == nil {
		conn.Close()
		successes <- task
	}
}
