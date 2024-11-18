// Package cache provides utilities for interfacing with the cached Nerd Fonts data.
package cache

import (
	"io"
	"math"
	"os"
	"path"
	"time"
)

// GetPath gets the path to the cache.
func GetPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "nerd-fonts-cli", "nerdfonts.json"), nil
}

// Write saves a reader to the cache.
func Write(r io.Reader) error {
	path, err := GetPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return err
	}
	return nil
}

// Open opens the cache file.
func Open() (*os.File, error) {
	path, err := GetPath()
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

// Age gets the age of the cache. Returns max int value if the cache could not be read.
func Age() (time.Duration, error) {
	const max = time.Duration(math.MaxInt64)
	path, err := GetPath()
	if err != nil {
		return max, err
	}
	fileinfo, err := os.Stat(path)
	if err != nil {
		return max, err
	}

	return time.Since(fileinfo.ModTime()), nil
}
