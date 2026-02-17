package main

import "time"

type BlogPost struct {
	Filename  string
	Title     string
	Timestamp time.Time
	Content   string
}

type RSSConfig struct {
	Title       string
	Link        string
	Description string
	Language    string
	FeedPath    string
}
