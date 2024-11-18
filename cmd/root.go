package cmd

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"nerd-fonts-cli/pkg/cache"
)

var rootCmd = &cobra.Command{
	Use: "nf-cli",
	Short: "An unofficial CLI tool to manage Nerd Fonts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Use one of the available subcommands")
	},
}

var Cache cache.Cache

func init() {
	// NOTE We ignore the error to fall back to the empty string
	Cache, _ = cache.Default()
	rootCmd.PersistentFlags().StringVarP(&Cache.Path, "cache", "c", Cache.Path, "path to the cache to save Nerd Font data")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
