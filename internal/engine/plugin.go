package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/dop251/goja"
)

type SearchResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func RunJSPlugin(filepath string, query string) ([]SearchResult, error) {
	scriptBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin: %v", err)
	}

	vm := goja.New()

	coreObj := vm.NewObject()
	coreObj.Set("runCmd", func(name string, args ...string) string {
		out, err := exec.Command(name, args...).CombinedOutput()
		if err != nil {
			// Print standard error to help debug JS plugin issues if yt-dlp fails
			fmt.Fprintf(os.Stderr, "[Plugin Debug] runCmd failed: %v\nOutput: %s\n", err, string(out))
			return ""
		}
		return string(out)
	})
	vm.Set("__core", coreObj)

	_, err = vm.RunString(string(scriptBytes))
	if err != nil {
		return nil, fmt.Errorf("javascript compilation error: %v", err)
	}

	var searchFn func(string) string
	err = vm.ExportTo(vm.Get("search"), &searchFn)
	if err != nil {
		return nil, fmt.Errorf("plugin must export a 'search' function returning JSON string: %v", err)
	}

	// Catch any panics from the JavaScript execution (e.g. JSON.parse on empty string)
	var jsonOutput string
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("javascript runtime panic: %v", r)
			}
		}()
		jsonOutput = searchFn(query)
	}()

	if err != nil {
		return nil, err
	}

	var results []SearchResult
	if err := json.Unmarshal([]byte(jsonOutput), &results); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from plugin: %v\nRaw output: %s", err, jsonOutput)
	}

	return results, nil
}
