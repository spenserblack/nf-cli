package cmd

import (
	"fmt"

	"github.com/spenserblack/nf-cli/internal/cache"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refreshCacheCmd)
}

var refreshCacheCmd = &cobra.Command{
	Use:   "refresh-cache",
	Short: "Force the cached Nerd Font data to be refreshed",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cache.Refresh(Cache)
		if err != nil {
			return err
		}

		fmt.Println("Cache refreshed")
		return nil
	},
}
