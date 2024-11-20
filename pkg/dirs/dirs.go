// Package dirs provides utilities for accessing font directories.
package dirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// UserFontDir returns the user font directory.
func UserFontDir() (string, error) {
	switch os := runtime.GOOS; os {
	case "linux":
		return linuxUserFontDir()
	case "windows":
		return windowsUserFontDir()
	default:
		return "", UnsupportedOS{os: os}
	}

}

// LinuxUserFontDir returns the user font directory on Linux.
func linuxUserFontDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".local", "share", "fonts"), nil
}

// WindowsUserFontDir returns the user font directory on Windows.
func windowsUserFontDir() (string, error) {
	localAppData := os.Getenv("LocalAppData")
	if localAppData == "" {
		return "", ErrNoLocalAppData
	}
	return filepath.Join(localAppData, "Microsoft", "Windows", "Fonts"), nil
}

// UnsupportedOS is returned when the OS is not supported.
type UnsupportedOS struct {
	os string
}

func (err UnsupportedOS) Error() string {
	return fmt.Sprintf("unsupported OS: %s", err.os)
}

// ErrNoLocalAppData is returned when the LocalAppData environment variable is not set
// on Windows.
var ErrNoLocalAppData = errors.New("LocalAppData environment variable not set")
