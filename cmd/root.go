package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/briandowns/spinner"

	"github.com/Addereum/mov-cli/internal/config"
	"github.com/Addereum/mov-cli/internal/engine"
	"github.com/Addereum/mov-cli/internal/player"
	"github.com/Addereum/mov-cli/internal/tui"
)

var scraperFlag string

var rootCmd = &cobra.Command{
	Use:   "mov-cli [query]",
	Short: "Watch everything from your terminal",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		for i := 1; i < len(args); i++ {
			query += " " + args[i]
		}

		pluginPath, err := config.FindPlugin(scraperFlag)
		if err != nil {
			fmt.Printf("[mov-cli] ❌ %v\n", err)
			fmt.Println("Tip: Check available plugins with 'mov-cli plugin list' or install one with 'mov-cli plugin add <url>'")
			os.Exit(1)
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" 🔍 Scraping [%s] for '%s'...", scraperFlag, query)
		s.Color("cyan", "bold")
		s.Start()

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

func init() {
	rootCmd.Flags().StringVarP(&scraperFlag, "scraper", "s", "youtube", "Scraper / plugin to use (default: youtube)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
