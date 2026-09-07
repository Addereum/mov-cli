package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"mov-cli-go/internal/config"
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
		url := args[0]
		
		fmt.Printf("📥 Downloading plugin from %s...\n", url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ Failed to download: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("❌ Failed to download: HTTP %d\n", resp.StatusCode)
			return
		}

		// Extract filename from URL
		parts := strings.Split(url, "/")
		filename := parts[len(parts)-1]
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

func init() {
	pluginCmd.AddCommand(addCmd)
	rootCmd.AddCommand(pluginCmd)
}
