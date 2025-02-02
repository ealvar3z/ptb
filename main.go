package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BlogPost struct {
	Filename  string
	Title     string
	Timestamp time.Time
}

var postTemplate = template.Must(template.New("post").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        <pre>{{.Content}}</pre>
        <footer>
            <p>&copy; 2024 eax. All rights reserved.</p>
            <p><a href="index.html">Back to Index</a></p>
        </footer>
    </div>
</body>
</html>`))

var indexTemplate = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Index</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div class="container">
        <h1><b>[P]</b>lain <b>[T]</b>ext <b>[B]</b>log</h1>
        <ul>
            {{range .}}
            <li><a href="{{.Filename}}">{{.Title}}</a> - {{.Timestamp.Format "2006-01-02"}}</li>
            {{end}}
        </ul>
        <footer>
            <p>&copy; 2024 eax. All rights reserved.</p>
        </footer>
    </div>
</body>
</html>`))

func main() {
	inputDir := "./txt"
	outputDir := "./output"
	createOutputDir(outputDir)
	posts := processTxtFiles(inputDir, outputDir)
	sortPostsByDate(posts)
	generateIndex(outputDir, posts)
	fmt.Println("Site generation complete!")
}

func createOutputDir(outputDir string) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
}

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
			timestamp, err := determineTimestamp(filePath)
			if err != nil {
				log.Fatalf("Failed to get timestamp: %v", err)
			}
			outputFile := processFile(filePath, outputDir)
			post := BlogPost{
				Filename:  filepath.Base(outputFile),
				Title:     strings.TrimSuffix(entry.Name(), ".txt"),
				Timestamp: timestamp,
			}
			posts = append(posts, post)
		}
	}
	return posts
}

func isTxtFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt")
}

func determineTimestamp(filePath string) (time.Time, error) {
	args := []string{"log", "--follow", "--diff-filter=A", "--format=%aI", filePath}
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		return time.Parse(time.RFC3339, lines[len(lines)-1])
	}
	return time.Time{}, nil
}

func processFile(inputFilePath, outputDir string) string {
	content, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v", inputFilePath, err)
		return ""
	}
	file := filepath.Base(inputFilePath)
	title := strings.TrimSuffix(file, ".txt")
	outputFilePath := filepath.Join(outputDir, title+".html")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create output file %s: %v", outputFilePath, err)
		return ""
	}
	defer outputFile.Close()
	postTemplate.Execute(outputFile, map[string]string{"Title": title, "Content": string(content)})
	return outputFilePath
}

func sortPostsByDate(posts []BlogPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(posts[j].Timestamp)
	})
}

func generateIndex(outputDir string, posts []BlogPost) {
	indexFile := filepath.Join(outputDir, "index.html")
	outputFile, err := os.Create(indexFile)
	if err != nil {
		log.Fatalf("Failed to write index.html: %v", err)
	}
	defer outputFile.Close()
	indexTemplate.Execute(outputFile, posts)
}
