package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	// Define the directory containing the JavaScript files
	jsDir := "public/js"
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	// Track errors during the minification process
	errorOccurred := false

	// Walk through the directory to find .js files
	err := filepath.Walk(jsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		// Process only .js files, excluding already minified ones
		if !info.IsDir() && filepath.Ext(path) == ".js" && filepath.Ext(path) != ".min.js" {
			fmt.Printf("Minifying: %s\n", path)

			// Read the source file
			input, readErr := ioutil.ReadFile(path)
			if readErr != nil {
				fmt.Printf("Failed to read file %s: %v\n", path, readErr)
				errorOccurred = true
				return nil
			}

			// Minify the JavaScript content
			minified, minifyErr := m.Bytes("application/javascript", input)
			if minifyErr != nil {
				fmt.Printf("Failed to minify file %s: %v\n", path, minifyErr)
				errorOccurred = true
				return nil
			}

			// Write the minified content to a new file with .min.js extension
			minifiedFile := path[:len(path)-3] + ".min.js"
			writeErr := ioutil.WriteFile(minifiedFile, minified, 0644)
			if writeErr != nil {
				fmt.Printf("Failed to write file %s: %v\n", minifiedFile, writeErr)
				errorOccurred = true
				return nil
			}

			fmt.Printf("Minified file created: %s\n", minifiedFile)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		os.Exit(1)
	}

	// Exit with a non-zero status code if any errors occurred
	if errorOccurred {
		fmt.Println("Minification process failed.")
		os.Exit(1)
	}

	fmt.Println("All JavaScript files minified successfully.")
}
