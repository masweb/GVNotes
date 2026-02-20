package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const appName = "gvnotes"

// DBPath returns the OS-appropriate path for the SQLite database file.
// macOS:   ~/Library/Application Support/gvnotes/gvnotes.db
// Linux:   ~/.local/share/gvnotes/gvnotes.db  (XDG_DATA_HOME)
// Windows: %APPDATA%\gvnotes\gvnotes.db
func DBPath() (string, error) {
	dir := filepath.Join(xdg.DataHome, appName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, appName+".db"), nil
}

// ImagesDir returns the directory where note images are stored and ensures it exists.
// macOS:   ~/Library/Application Support/gvnotes/images
// Linux:   ~/.local/share/gvnotes/images
// Windows: %APPDATA%\gvnotes\images
func ImagesDir() (string, error) {
	dir := filepath.Join(xdg.DataHome, appName, "images")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
