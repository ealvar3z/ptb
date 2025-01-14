package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BlogPost struct {
	Filename string		// Name of the HTML file
	Title string		// Title of the blog post
	Timestamp time.Time	// Creation or mod time of the file
}

func main() {
	inputDir := "./txt"
	outputDir := "./output"
	createOutputDir(outputDir)
	copyCSSFile(outputDir)

	posts := processTxtFiles(inputDir, outputDir)
	sortPostsByDate(posts)

	generateIndex(outputDir, posts)

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

h1 b {
	color: 	#4caf50; /* green */
	font-weight: bold;
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
// and returns a list of blog posts
func processTxtFiles(inputDir, outputDir string) []BlogPost {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	var posts []BlogPost
	for _, entry := range entries {
		if !entry.IsDir() && isTxtFile(entry.Name()) {
			fmt.Printf("Processing: %s\n", entry.Name())
			filePath := filepath.Join(inputDir, entry.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				log.Printf("Failed to get file info for %s: %v\n",
				entry.Name(), err)
				continue
			}
			outputFile := processFile(filePath, outputDir)
			post := BlogPost{
				Filename: filepath.Base(outputFile),
				Title: strings.TrimSuffix(entry.Name(), ".txt"),
				Timestamp: info.ModTime(),
			}
			posts = append(posts, post)
		}
	}
	return posts
}

// isTxtFile checks if a file has a .txt extension
func isTxtFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt")
}

// processFile reads a .txt file, generates HTML, and writes it to the output directory
func processFile(inputFilePath, outputDir string) string {
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

//sortPostsByDate sorts the blog posts by their creation
// or modified time
func sortPostsByDate(posts []BlogPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(posts[j].Timestamp)
	})
}

// generateIndex generates the index.html file
// w/ blog posts sorted by date
func generateIndex(outputDir string, posts []BlogPost) {
	var links string
	for _, post := range posts {
		links += fmt.Sprintf(
			`<li><a href="%s">%s</a> - %s</li>`, 
			post.Filename,
			post.Title,
			post.Timestamp.Format("2006-01-02 15:04:05"),
		)
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
        <h1><b>[P]</b>lain <b>[T]</b>ext <b>[B]</b>log || [PTB]</h1>
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

