package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"
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
var MaxCacheAge time.Duration

func init() {
	// NOTE We ignore the error to fall back to the empty string
	Cache, _ = cache.Default()
	// NOTE 7 days
	defaultCacheAge, err := time.ParseDuration("168h")
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVarP(&Cache.Path, "cache", "c", Cache.Path, "path to the cache to save Nerd Font data")
	rootCmd.PersistentFlags().DurationVarP(&MaxCacheAge, "max-cache-age", "A", defaultCacheAge, "how long the cache should exist before being automatically refreshed")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
