package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Addereum/mov-cli/internal/config"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage mov-cli JavaScript plugins",
}

var addCmd = &cobra.Command{
	Use:   "add [url]",
	Short: "Install a plugin from a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rawURL := args[0]

		// Auto-convert standard GitHub blob URLs to raw URLs
		downloadURL := rawURL
		if strings.Contains(rawURL, "github.com") && strings.Contains(rawURL, "/blob/") {
			downloadURL = strings.Replace(rawURL, "github.com", "raw.githubusercontent.com", 1)
			downloadURL = strings.Replace(downloadURL, "/blob/", "/", 1)
		}

		fmt.Printf("📥 Downloading plugin from %s...\n", downloadURL)
		resp, err := http.Get(downloadURL)
		if err != nil {
			fmt.Printf("❌ Failed to download: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("❌ Failed to download: HTTP %d\n", resp.StatusCode)
			return
		}

		// Extract clean filename from URL (stripping query parameters)
		parsedURL, err := url.Parse(downloadURL)
		var filename string
		if err == nil {
			filename = filepath.Base(parsedURL.Path)
		} else {
			parts := strings.Split(downloadURL, "/")
			filename = parts[len(parts)-1]
		}

		if filename == "" || filename == "." || filename == "/" {
			filename = "plugin.js"
		}
		if !strings.HasSuffix(filename, ".js") {
			filename += ".js"
		}

		destPath := filepath.Join(config.GetPluginsDir(), filename)
		out, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("❌ Failed to create file: %v\n", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("❌ Failed to save file: %v\n", err)
			return
		}

		fmt.Printf("✅ Plugin '%s' installed successfully to %s!\n", filename, destPath)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		plugins := config.ListInstalledPlugins()
		if len(plugins) == 0 {
			fmt.Println("No plugins installed.")
			fmt.Println("Install one using: mov-cli plugin add <url>")
			return
		}

		fmt.Println("🔌 Installed plugins:")
		for _, p := range plugins {
			fmt.Printf("  • %s (%s)\n", p.Name, p.Path)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove an installed plugin",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, ".js") {
			name += ".js"
		}

		userPluginPath := filepath.Join(config.GetPluginsDir(), name)
		if _, err := os.Stat(userPluginPath); os.IsNotExist(err) {
			fmt.Printf("❌ Plugin '%s' not found in user directory (%s)\n", name, config.GetPluginsDir())
			return
		}

		if err := os.Remove(userPluginPath); err != nil {
			fmt.Printf("❌ Failed to remove plugin: %v\n", err)
			return
		}

		fmt.Printf("✅ Plugin '%s' removed successfully.\n", name)
	},
}

func init() {
	pluginCmd.AddCommand(addCmd)
	pluginCmd.AddCommand(listCmd)
	pluginCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(pluginCmd)
}
