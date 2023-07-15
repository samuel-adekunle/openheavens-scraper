package main

import (
	"fmt"
	"strings"
)

type Post struct {
	Title               string   `json:"title"`
	MemoryVerse         string   `json:"memoryVerse"`
	BibleReadingHeading string   `json:"bibleReadingHeading"`
	BibleReadingBody    []string `json:"bibleReadingBody"`
	MessageBody         []string `json:"messageBody"`
	ActionPoint         string   `json:"actionPoint"`
	HymnTitle           string   `json:"hymnTitle"`
	HymnBody            []string `json:"hymnBody"`
	BibleInOneYear      string   `json:"bibleInOneYear"`
	Authors             []string `json:"authors"`
}

// string representation of a post
func (p Post) String() string {
	return fmt.Sprintf("*%s*\n\n*MEMORIZE:*\n%s\n\n*%s*\n%s\n\n*MESSAGE:*\n%s\n\n*ACTION POINT:*\n%s\n\n*HYMN:*\n%s\n%s\n\n*%s*\n\n*AUTHORS:*\n%s", p.Title, p.MemoryVerse, p.BibleReadingHeading, strings.Join(p.BibleReadingBody, "\n\n"), strings.Join(p.MessageBody, "\n\n"), p.ActionPoint, p.HymnTitle, strings.Join(p.HymnBody, "\n"), p.BibleInOneYear, strings.Join(p.Authors, "\n"))
}
