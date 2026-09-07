package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	
	// Legacy / Helper: Run commands
	coreObj.Set("runCmd", func(name string, args ...string) string {
		out, _ := exec.Command(name, args...).CombinedOutput()
		return string(out)
	})

	// NEW: Native HTTP Request
	coreObj.Set("fetch", func(url string) string {
		resp, err := http.Get(url)
		if err != nil {
			return ""
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return string(b)
	})

	// NEW: Native HTML Parser (GoQuery)
	coreObj.Set("parseHTML", func(html string, selector string, attr string) []string {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil
		}
		
		var results []string
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if attr == "text" {
				results = append(results, s.Text())
			} else {
				if val, exists := s.Attr(attr); exists {
					results = append(results, val)
				}
			}
		})
		return results
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
