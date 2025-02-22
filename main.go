package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BlogPost struct {
	Filename  string
	Title     string
	Timestamp time.Time
	Content   string
}

var postTemplate = template.Must(template.New("post").Funcs(template.FuncMap{
	"urlquery": func(s string) string {
		return template.URLQueryEscaper(s)
	},
}).Parse(`<!DOCTYPE html>
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
            <p>
		<a href="index.html">Back to Index</a> |
		<a href="https://github.com/ealvar3z/ptb/issues/new?title={{.Title | urlquery}}&body={{printf \"Comments for the post: %s (published on %s)\" .Title. Timestamp.Format \"2006-01-02\" | urlquery}}" target="_blank">Comments</a>
	    </p>
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
			post := parseFilename(filePath)
			post.Content = readFileContent(filePath)
			post.Filename = post.Title + ".html"
			processFile(&post, outputDir)
			posts = append(posts, post)
		}
	}
	return posts
}

func isTxtFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt")
}

func parseFilename(filePath string) BlogPost {
	filename := filepath.Base(filePath)
	filename = strings.TrimSuffix(filename, ".txt")
	parts := strings.SplitN(filename, "_", 2)

	if len(parts) < 2 {
		log.Fatalf("Invalid filename format: %s", filename)
	}

	timestamp, err := time.Parse("20060102", parts[0])
	if err != nil {
		log.Fatalf("Failed to parse date from filename: %s", filename)
	}

	title := parts[1]
	return BlogPost{
		Title:     title,
		Timestamp: timestamp,
	}
}

func readFileContent(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	return string(content)
}

func processFile(post *BlogPost, outputDir string) {
	outputFilePath := filepath.Join(outputDir, post.Filename)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", post.Filename, err)
	}
	defer outputFile.Close()

	postTemplate.Execute(outputFile, post)
	post.Filename = post.Filename // TODO: hacky way to ensure proper filename (gotta fix this)
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
