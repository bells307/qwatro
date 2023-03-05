package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "qwatro",
	Short: "qwatro is a simple network tool which can scan tcp ports",
	Long:  `qwatro is a simple network tool which can scan tcp ports`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Errorf("%s\n", err)
		os.Exit(1)
	}
}
