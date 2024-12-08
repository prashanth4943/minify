package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	// Define the directory containing the JavaScript files
	jsDir := "public"
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	// Process the JS files in the specified directory
	if err := processJSFiles(jsDir, m); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("All JavaScript files minified successfully.")
}

// Function to walk through the directory and minify all .js files
func processJSFiles(jsDir string, m *minify.M) error {
	// Track errors during the minification process
	errorOccurred := false

	// Walk through the directory to find .js files
	err := filepath.Walk(jsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}
		// Process only .js files, excluding already minified ones
		if !info.IsDir() && filepath.Ext(path) == ".js" && !strings.Contains(path, ".min") {
			fmt.Printf("Minifying: %s\n", path)

			// Call the minify function for each JS file
			if minifyErr := minifyJSFile(path, m); minifyErr != nil {
				errorOccurred = true
				fmt.Println(minifyErr)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the directory: %v", err)
	}

	// Return an error if minification failed for any file
	if errorOccurred {
		return fmt.Errorf("minification process failed")
	}

	return nil
}

// Function to handle the minification of a single JavaScript file
func minifyJSFile(path string, m *minify.M) error {
	// Read the source file
	input, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return fmt.Errorf("failed to read file %s: %v", path, readErr)
	}

	// Minify the JavaScript content
	minified, minifyErr := m.Bytes("application/javascript", input)
	if minifyErr != nil {
		return fmt.Errorf("failed to minify file %s: %v", path, minifyErr)
	}

	// Write the minified content to a new file with .min.js extension
	minifiedFile := path[:len(path)-3] + ".min.js"
	writeErr := ioutil.WriteFile(minifiedFile, minified, 0644)
	if writeErr != nil {
		return fmt.Errorf("failed to write file %s: %v", minifiedFile, writeErr)
	}

	fmt.Printf("Minified file created: %s\n", minifiedFile)
	return nil
}

func endsWithMinified(path string) bool {
	return filepath.Ext(path) == ".min.js"
}




