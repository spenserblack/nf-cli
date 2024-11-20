package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spenserblack/nf-cli/internal/cache"
	"github.com/spenserblack/nf-cli/internal/prompts"
	"github.com/spenserblack/nf-cli/pkg/fonts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use: "install",
	Short: "Install fonts",
	Long: heredoc.Doc(`
		Install fonts. This command will install the fonts in the user's font
		directory.
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cache.RefreshIfOld(Cache, MaxCacheAge); err != nil {
			return err
		}
		all, err := fonts.UnmarshalCache(Cache)
		if err != nil {
			return err
		}

		selected, err := prompts.MultiPromptForFonts(all)
		if err != nil {
			return err
		}

		fmt.Println(selected)

		return nil
	},
}
