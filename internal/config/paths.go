package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config", "mov-cli")
	os.MkdirAll(path, 0755)
	return path
}

func GetPluginsDir() string {
	path := filepath.Join(GetConfigDir(), "plugins")
	os.MkdirAll(path, 0755)
	return path
}
