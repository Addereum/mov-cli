package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// GetPluginSearchDirs returns all directories where plugins might be stored
func GetPluginSearchDirs() []string {
	var dirs []string

	// 1. User config directory (~/.config/mov-cli/plugins)
	dirs = append(dirs, GetPluginsDir())

	// 2. System-wide directories (installed via AUR/pacman)
	dirs = append(dirs, "/usr/share/mov-cli/plugins", "/usr/local/share/mov-cli/plugins")

	// 3. Current working directory / relative fallback for development
	dirs = append(dirs, "plugins")

	return dirs
}

// FindPlugin finds the full path of a plugin by name (with or without .js)
func FindPlugin(name string) (string, error) {
	if !strings.HasSuffix(name, ".js") {
		name += ".js"
	}

	// Check if name is already an explicit file path
	if info, err := os.Stat(name); err == nil && !info.IsDir() {
		return name, nil
	}

	for _, dir := range GetPluginSearchDirs() {
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("plugin '%s' not found (searched in: %s)", name, strings.Join(GetPluginSearchDirs(), ", "))
}

type PluginInfo struct {
	Name string
	Path string
}

// ListInstalledPlugins returns all unique plugins found in the search directories
func ListInstalledPlugins() []PluginInfo {
	seen := make(map[string]bool)
	var list []PluginInfo

	for _, dir := range GetPluginSearchDirs() {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".js") {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), ".js")
			if !seen[name] {
				seen[name] = true
				list = append(list, PluginInfo{
					Name: name,
					Path: filepath.Join(dir, entry.Name()),
				})
			}
		}
	}
	return list
}
