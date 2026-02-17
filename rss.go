package main

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language,omitempty"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func defaultRSSConfig() RSSConfig {
	return RSSConfig{
		Title:       "Plain Text Blog",
		Link:        "https://example.com",
		Description: "Posts from Plain Text Blog",
		Language:    "en-us",
		FeedPath:    "rss.xml",
	}
}

func writeRSS(outputDir string, posts []BlogPost, cfg RSSConfig) {
	feed := rss{
		Version: "2.0",
		Channel: rssChannel{
			Title:       cfg.Title,
			Link:        cfg.Link,
			Description: cfg.Description,
			Language:    cfg.Language,
		},
	}
	if len(posts) > 0 {
		feed.Channel.LastBuildDate = posts[0].Timestamp.UTC().Format(time.RFC1123Z)
	}

	for _, post := range posts {
		link := joinURL(cfg.Link, post.Filename)
		feed.Channel.Items = append(feed.Channel.Items, rssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        link,
			PubDate:     post.Timestamp.UTC().Format(time.RFC1123Z),
			Description: summarize(post.Content),
		})
	}

	b, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatalf("Failed to render RSS: %v", err)
	}

	out := filepath.Join(outputDir, cfg.FeedPath)
	content := []byte(xml.Header + string(b) + "\n")
	if err := os.WriteFile(out, content, 0o644); err != nil {
		log.Fatalf("Failed to write RSS file: %v", err)
	}
}

func joinURL(base, name string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(name, "/")
}

func summarize(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	const maxLen = 280
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
