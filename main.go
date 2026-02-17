package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

var tmpl = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"urlquery": template.URLQueryEscaper,
			"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		}).
		ParseFS(templateFS, "templates/*.html"),
)

func main() {
	inputDir := "./txt"
	outputDir := "./output"
	ensureDir(outputDir)

	posts := collectPosts(inputDir, outputDir)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(posts[j].Timestamp)
	})
	writeIndex(outputDir, posts)
	writeRSS(outputDir, posts, defaultRSSConfig())
	fmt.Println("Site generation complete!")
}

func ensureDir(dir string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
}

func collectPosts(inputDir, outputDir string) []BlogPost {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}
	var posts []BlogPost
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		src := filepath.Join(inputDir, e.Name())
		post := parseFilename(src)
		post.Content = mustReadFile(src)
		post.Filename = post.Title + ".html"

		writePost(outputDir, &post)
		posts = append(posts, post)
	}
	return posts
}

func parseFilename(path string) BlogPost {
	base := strings.TrimSuffix(filepath.Base(path), ".txt")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) < 2 {
		log.Fatalf("Invalid filename format: %s", base)
	}

	ts, err := time.Parse("20060102", parts[0])
	if err != nil {
		log.Fatalf("Failed to parse date from filename: %s", base)
	}
	return BlogPost{Title: parts[1], Timestamp: ts}
}

func mustReadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Reading %s: %v", path, err)
	}
	return string(b)
}

func writePost(outDir string, post *BlogPost) {
	inputFile := filepath.Join(outDir, post.Filename)
	outputFile, err := os.Create(inputFile)
	if err != nil {
		log.Fatalf("Creating %s: %v", inputFile, err)
	}
	defer outputFile.Close()

	if err := tmpl.ExecuteTemplate(outputFile, "post.html", post); err != nil {
		log.Fatalf("Rendering post %s: %v", post.Title, err)
	}
}

func writeIndex(outputDir string, posts []BlogPost) {
	indexFile := filepath.Join(outputDir, "index.html")
	outputFile, err := os.Create(indexFile)
	if err != nil {
		log.Fatalf("Failed to write index.html: %v", err)
	}
	defer outputFile.Close()

	if err := tmpl.ExecuteTemplate(outputFile, "index.html", posts); err != nil {
		log.Fatalf("Rendering index: %v", err)
	}
}
