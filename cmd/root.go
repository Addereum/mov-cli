package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/briandowns/spinner"

	"mov-cli-go/internal/config"
	"mov-cli-go/internal/engine"
	"mov-cli-go/internal/player"
	"mov-cli-go/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "mov-cli [query]",
	Short: "Watch everything from your terminal",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		for i := 1; i < len(args); i++ {
			query += " " + args[i]
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" 🔍 Scraping for '%s'...", query)
		s.Color("cyan", "bold")
		s.Start()

		// For now, load youtube.js from local or config folder
		pluginPath := filepath.Join(config.GetPluginsDir(), "youtube.js")
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			pluginPath = "plugins/youtube.js" // fallback for development
		}

		results, err := engine.RunJSPlugin(pluginPath, query)
		s.Stop()

		if err != nil {
			fmt.Printf("\n[mov-cli] ❌ Plugin crashed gracefully!\nDetails: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Println("[mov-cli] ❌ No results found.")
			os.Exit(1)
		}

		selected, err := tui.SelectResult(results)
		if err != nil {
			fmt.Println("\n[mov-cli] ❌ Cancelled.")
			os.Exit(0)
		}

		fmt.Printf("[mov-cli] 🎬 Playing '%s' via MPV...\n", selected.Title)

		err = player.PlayVideo(selected.URL)
		if err != nil {
			fmt.Printf("[mov-cli] ❌ Playback Error: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
