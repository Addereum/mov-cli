package tui

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Addereum/mov-cli/internal/engine"
)

func SelectResult(results []engine.SearchResult) (engine.SearchResult, error) {
	if len(results) == 0 {
		return engine.SearchResult{}, fmt.Errorf("no results found")
	}

	var options []string
	for _, res := range results {
		options = append(options, res.Title)
	}

	prompt := &survey.Select{
		Message: "Choose Result:",
		Options: options,
		PageSize: 15,
	}

	var selectedIndex int
	err := survey.AskOne(prompt, &selectedIndex)
	if err != nil {
		return engine.SearchResult{}, err
	}

	return results[selectedIndex], nil
}
