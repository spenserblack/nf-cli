package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
	"nerd-fonts-cli/internal/cache"
	"nerd-fonts-cli/internal/prompts"
	"nerd-fonts-cli/pkg/fonts"
	"github.com/MakeNowJust/heredoc"
)

var PreviewFontName string

func init() {
	previewCmd.PersistentFlags().StringVar(&PreviewFontName, "font", "", "Name of the font to preview")
	rootCmd.AddCommand(previewCmd)
}

var previewCmd = &cobra.Command{
	Use: "preview",
	Short: "Preview the font in the browser",
	Long: heredoc.Doc(`
		Preview the font in the browser. Note that not all fonts have available
		previews.

		The $BROWSER environment variable must be set.
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		var selected fonts.Font
		browser := os.Getenv("BROWSER")
		if browser == "" {
			return errors.New("The environment variable $BROWSER must be set")
		}
		if err := cache.RefreshIfOld(Cache, MaxCacheAge); err != nil {
			return err
		}
		all, err := fonts.UnmarshalCache(Cache)
		if err != nil {
			return err
		}
		previewable := make([]fonts.Font, 0, len(all))
		for _, font := range all {
			if _, ok := font.LinkPreviewFont.(string); ok {
				previewable = append(previewable, font)
			}
		}

		if PreviewFontName != "" {
			names := make([]string, 0, len(previewable) * 2)
			found := false
			for _, font := range previewable {
				if  PreviewFontName == font.PatchedName || PreviewFontName == font.UnpatchedName {
					selected = font
					found = true
					break
				}
				names = append(names, font.PatchedName)
				if font.UnpatchedName != font.PatchedName {
					names = append(names, font.UnpatchedName)
				}
			}

			if !found {
				return fmt.Errorf(
					"--font must be one of: %s\n",
					strings.Join(names, ", "),
				)
			}
		} else {
			font, err := prompts.PromptForFont(previewable)
			if err != nil {
				return err
			}
			selected = font
		}

		exepath, err := exec.LookPath(browser)
		if err != nil {
			return err
		}

		fmt.Printf("%#v\n", selected)
		// NOTE Already safe and don't need URL encoding?
		url := fmt.Sprintf("https://www.programmingfonts.org/#%s", selected.LinkPreviewFont)

		fmt.Printf("Opening %s in the browser...\n", url)

		browserCmd := exec.Command(exepath, url)
		browserCmd.Stdout = os.Stdout
		browserCmd.Stderr = os.Stderr

		if err := browserCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}
