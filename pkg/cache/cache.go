// Package cache provides utilities for interfacing with the cached Nerd Fonts data.
package cache

import (
	"io"
	"math"
	"os"
	"path/filepath"
	"time"
)

// Cache handles interacting with the cache.
type Cache struct {
	// Path is the path to the cache.
	Path string
}

// Default returns the default cache.
func Default() (Cache, error) {
	path, err := defaultPath()
	// HACK As long as the struct is simple, it should be safe to use the zero
	//		value of the path.
	cache := Cache{
		Path: path,
	}
	return cache, err
}

// DefaultPath gets the path to the cache.
func defaultPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nerd-fonts-cli", "nerdfonts.json"), nil
}

// Write saves a reader to the cache.
func (cache Cache) Write(r io.Reader) error {
	path := cache.Path
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
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
func (cache Cache) Open() (*os.File, error) {
	path := cache.Path
	return os.Open(path)
}

// Age gets the age of the cache. Returns max int value if the cache could not be read.
func (cache Cache) Age() (time.Duration, error) {
	const max = time.Duration(math.MaxInt64)
	path := cache.Path
	fileinfo, err := os.Stat(path)
	if err != nil {
		return max, err
	}

	return time.Since(fileinfo.ModTime()), nil
}
