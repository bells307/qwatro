package cmd

import (
	"context"
	"fmt"
	"github.com/bells307/qwatro/port_scanner"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"time"
)

var portRangeStr string
var workers int
var tcpTimeout time.Duration
var timeout time.Duration

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Port scanning",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Определяем диапазон портов из строки
		portRange, err := port_scanner.RangeFromString(portRangeStr)
		if err != nil {
			log.Fatal(err)
		}

		// Создаем сканер
		scanner := port_scanner.
			NewScannerBuilder().
			IP(args[0]).
			PortRange(portRange).
			NumWorkers(workers).
			Tcp(tcpTimeout).
			Build()

		var rootCtx context.Context
		if timeout != 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			rootCtx = ctx
		} else {
			rootCtx = context.Background()
		}

		ctx, cancel := context.WithCancel(rootCtx)
		defer cancel()

		// Запускаем сканер
		resChan := scanner.Run(ctx)

		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt)

		for {
			select {
			case res, more := <-resChan:
				if more {
					fmt.Printf("%s:%d\n", res.IP, res.Port)
				} else {
					return
				}
			case _ = <-sigChan:
				log.Println("got interrupt signal")
				cancel()
			}
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&portRangeStr, "port-range", "p", "1-65535", "Port range")
	scanCmd.Flags().IntVarP(&workers, "workers", "w", port_scanner.DefaultNumWorkers, "Workers for scanning")
	scanCmd.Flags().DurationVar(&tcpTimeout, "tcp-timeout", 300*time.Millisecond, "Maximum response time for tcp scanning")
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "t", 0, "General scan timeout")

	rootCmd.AddCommand(scanCmd)
}
