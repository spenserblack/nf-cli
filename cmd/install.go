package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	"github.com/spenserblack/nf-cli/internal/cache"
	"github.com/spenserblack/nf-cli/internal/prompts"
	"github.com/spenserblack/nf-cli/pkg/dirs"
	"github.com/spenserblack/nf-cli/pkg/fonts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install fonts",
	Long: heredoc.Doc(`
		Install fonts. This command will install the fonts in the user's font
		directory.
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		// NOTE Ensure that the font directory exists
		destDir, err := dirs.UserFontDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(destDir, os.ModeDir|os.ModePerm); err != nil {
			return err
		}

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

		srcDir, paths, err := fonts.FetchTmp(selected)
		if err != nil {
			var bulkFetchErr *fonts.BulkFetchErr
			if errors.As(err, bulkFetchErr) {
				// NOTE Non-fatal error, some fonts might still be successfully
				//		fetched.
				fmt.Fprintf(os.Stderr, err.Error())
			} else {
				return err
			}
		}
		defer os.RemoveAll(srcDir)
		fmt.Fprintf(os.Stdout, "Downloaded fonts to %s\n", srcDir)

		for _, path := range paths {
			installZip(path, destDir)
		}

		switch runtime.GOOS {
		case "linux":
			fmt.Fprintln(os.Stdout, "You may need to run `fc-cache -f` to complete installation")
		}

		return nil
	},
}

func installZip(src string, dest string) {
	zf, err := zip.OpenReader(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer zf.Close()

	for _, file := range zf.File {
		tryInstallFile(file, dest)
	}
}

func tryInstallFile(file *zip.File, dest string) {
	if ext := filepath.Ext(file.Name); !(ext == ".otf" || ext == ".ttf") {
		return
	}
	fmt.Fprintf(os.Stdout, "Installing %s to %s ... ", file.Name, dest)
	r, err := file.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer r.Close()

	destFile, err := os.Create(filepath.Join(dest, file.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer destFile.Close()
	if _, err := io.Copy(destFile, r); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Println("done!")
}
