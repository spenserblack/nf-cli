// Package data provides utilities to get Nerd Fonts data.
package data

import "net/http"

const endpoint = "https://raw.githubusercontent.com/ryanoasis/nerd-fonts/refs/heads/master/bin/scripts/lib/fonts.json"

// FetchData fetches the Nerd Fonts JSON file content from the repository.
func FetchData() (*http.Response, error) {
	return http.Get(endpoint)
}
