package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"mov-cli-go/internal/engine"
	"mov-cli-go/internal/player"
	"mov-cli-go/internal/tui"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mov-cli-go <search query>")
		os.Exit(1)
	}
	query := strings.Join(os.Args[1:], " ")
	
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" 🔍 Scraping YouTube for '%s'...", query)
	s.Color("cyan", "bold")
	s.Start()

	results, err := engine.RunJSPlugin("plugins/youtube.js", query)
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
}
