package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spenserblack/nf-cli/pkg/fonts"
	"github.com/spenserblack/nf-cli/internal/cache"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "Display a list of the available Nerd Fonts in the terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cache.RefreshIfOld(Cache, MaxCacheAge)
		if err != nil {
			return err
		}
		fonts, err := fonts.UnmarshalCache(Cache)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Original Name", "Version", "License", "Description"})

		for _, font := range fonts {
			t.AppendRow(table.Row{font.PatchedName, font.UnpatchedName, font.Version, font.License, font.Description})
		}
		t.Render()

		return nil
	},
}
