package cmd

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "nf-cli",
	Short: "An unofficial CLI tool to manage Nerd Fonts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Use one of the available subcommands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
