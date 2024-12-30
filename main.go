package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputDir := "./txt"
	outputDir := "./output"
	createOutputDir(outputDir)
	copyCSSFile(outputDir)
	htmlFiles := processTxtFiles(inputDir, outputDir)
	generateIndex(outputDir, htmlFiles)

	fmt.Println("Site generation complete!")
}

// createOutputDir ensures the output directory exists
func createOutputDir(outputDir string) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
}

// copyCSSFile copies the CSS file to the output directory
func copyCSSFile(outputDir string) {
	cssContent := `
/* Center the entire page */
body {
    font-family: monospace;
    background-color: #000;
    color: #fff;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    margin: 0;
}

/* Main content container */
.container {
    max-width: 79ch;
    width: calc(100% - 40px);
    padding: 20px;
    border: 1px solid #444;
    background-color: #111;
    text-align: left;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
}

/* Title styling */
h1 {
    font-size: 1.5em;
    text-align: center;
    margin: 0 0 20px;
}

/* Content block styling */
pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    margin: 0;
    line-height: 1.6;
}

/* Footer styling */
footer {
    text-align: center;
    margin-top: 20px;
    font-size: 0.9em;
    color: #666;
}

a {
    color: #ccc;
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}
`
	cssFilePath := filepath.Join(outputDir, "style.css")
	if err := os.WriteFile(cssFilePath, []byte(cssContent), 0644); err != nil {
		log.Fatalf("Failed to write CSS file: %v", err)
	}
}

// processTxtFiles processes all .txt files in the input directory
func processTxtFiles(inputDir, outputDir string) []string {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	var htmlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && isTxtFile(entry.Name()) {
			fmt.Printf("Processing: %s\n", entry.Name())
			outputFile := processFile(filepath.Join(inputDir, entry.Name()), outputDir)
			htmlFiles = append(htmlFiles, filepath.Base(outputFile))
		}
	}
	return htmlFiles
}

// isTxtFile checks if a file has a .txt extension
func isTxtFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt")
}

// processFile reads a .txt file, generates HTML, and writes it to the output directory
func processFile(inputFilePath, outputDir string) string {
	// Read file content
	content, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", inputFilePath, err)
		return ""
	}

	// Generate HTML
	file := filepath.Base(inputFilePath)
	title := strings.TrimSuffix(file, ".txt")
	outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(file, ".txt")+".html")
	htmlContent := generateHTML(string(content), title)

	// Write HTML to the output directory
	writeOutputFile(outputFilePath, htmlContent)
	return outputFilePath
}

// writeOutputFile writes content to a specified file
func writeOutputFile(filePath, content string) {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Printf("Failed to write file %s: %v\n", filePath, err)
	}
}

// generateHTML generates an HTML string from the provided content and title
func generateHTML(content, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <pre>%s</pre>
        <footer>
            <p>&copy; 2024 eax. All rights reserved.</p>
            <p><a href="index.html">Back to Index</a></p>
        </footer>
    </div>
</body>
</html>`, title, title, content)
}

// generateIndex generates the index.html file
func generateIndex(outputDir string, files []string) {
	var links string
	for _, file := range files {
		title := strings.TrimSuffix(filepath.Base(file), ".html")
		links += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, file, title)
	}

	indexContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Index</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div class="container">
        <h1>Plaintext Blog</h1>
        <ul>
            %s
        </ul>
        <footer>
            <p>&copy; 2024 eax. All rights reserved.</p>
        </footer>
    </div>
</body>
</html>`, links)

	indexFile := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(indexFile, []byte(indexContent), 0644); err != nil {
		log.Fatalf("Failed to write index.html: %v", err)
	}
}

