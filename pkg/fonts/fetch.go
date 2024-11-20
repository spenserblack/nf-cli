package fonts

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ReleaseEndpoint is the endpoint for downloading a release asset.
const releaseEndpoint = "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/%s"

// FilenamePattern is the pattern for a release asset. The zip file's name should be the
// folder name.
const filenamePattern = "%s.zip"

// FetchTmp downloads the fonts' zip files to a temporary directory.
//
// It returns the temporary directory and a list of filenames for the downloaded
// fonts.
func FetchTmp(fonts []Font) (string, []string, error) {
	dir, err := tmpdir()
	if err != nil {
		return "", nil, err
	}
	filenames := make([]string, 0, len(fonts))
	fetchErrs := BulkFetchErr{}

	for _, font := range fonts {
		path, err := font.fetch(dir)
		if err != nil {
			fetchErrs.errors = append(fetchErrs.errors, err)
		} else {
			filenames = append(filenames, path)
		}
	}

	finalErr := error(nil)
	if len(fetchErrs.errors) > 0 {
		finalErr = fetchErrs
	}

	return dir, filenames, finalErr
}

// Fetch downloads the font from the latest Nerd Font release to the target directory.
//
// It returns the file's path after download.
func (font Font) fetch(target string) (string, error) {
	filename := fmt.Sprintf(filenamePattern, font.Folder)
	endpoint := fmt.Sprintf(releaseEndpoint, filename)
	dest := filepath.Join(target, filename)

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return dest, nil
}

// Tmpdir creates a temporary directory for downloading fonts to.
func tmpdir() (string, error) {
	return os.MkdirTemp(os.TempDir(), "nf-cli-*")
}

// FontFetchFailureErr represents an individual fetch error.
type fontFetchFailureErr struct {
	name   string
	reason error
}

func (err fontFetchFailureErr) Error() string {
	return fmt.Sprintf("Failed to fetch %s: %v", err.name, err.reason)
}

// BulkFetchErr is a collection of errors when fetching fonts.
type BulkFetchErr struct {
	errors []error
}

func (err BulkFetchErr) Error() string {
	messages := make([]string, 0, len(err.errors))
	for _, e := range err.errors {
		messages = append(messages, e.Error())
	}
	return fmt.Sprintf("Failed to fetch fonts:\n%s", strings.Join(messages, "\n"))
}
