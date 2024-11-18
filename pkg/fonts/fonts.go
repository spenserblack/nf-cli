// Package fonts provides utilities for parsing Nerd Font data.
package fonts

import (
	"bytes"
	"encoding/json"

	"nerd-fonts-cli/pkg/cache"
)

// Font represents a Nerd Font's data.
type Font struct {
	// The Nerd Font name.
	PatchedName string `json:"patchedName"`
	// The original font's name.
	UnpatchedName string `json:"unpatchedName"`
	Version string `json:"version"`
	Description string `json:"description"`
	License string `json:"licenseId"`
	// The name of the folder in release assets.
	Folder string `json:"folderName"`
	LinkPreviewFont string `json:"linkPreviewFont"`
}

// Unmarshal loads Nerd Font data from JSON bytes.
func Unmarshal(data []byte) ([]Font, error) {
	fonts := struct{
		fonts []Font `json:"fonts"`
	}{}

	err := json.Unmarshal(data, &fonts)
	return fonts.fonts, err
}

// UnmarshalCache loads Nerd Font data from the cached data.
func UnmarshalCache(cache cache.Cache) ([]Font, error) {
	f, err := cache.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(f); err != nil {
		return nil, err
	}
	return Unmarshal(b.Bytes())
}
