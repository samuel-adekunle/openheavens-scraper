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
	PrayerPoint         string   `json:"prayerPoint"`
	IsActionPoint       bool     `json:"isActionPoint"`
	HymnTitle           string   `json:"hymnTitle"`
	HymnBody            []string `json:"hymnBody"`
	BibleInOneYear      string   `json:"bibleInOneYear"`
}

// string representation of a post
func (p Post) String() string {
	pointString := "PRAYER POINT"
	if p.IsActionPoint {
		pointString = "ACTION POINT"
	}
	return fmt.Sprintf("*%s*\n\n*MEMORY VERSE:*\n%s\n\n*%s*\n%s\n\n*%s*\n\n*MESSAGE:*\n%s\n\n*%s:*\n%s\n\n*%s*\n%s\n", p.Title, p.MemoryVerse, p.BibleReadingHeading, strings.Join(p.BibleReadingBody, "\n\n"), p.BibleInOneYear, strings.Join(p.MessageBody, "\n\n"), pointString, p.PrayerPoint, p.HymnTitle, strings.Join(p.HymnBody, "\n\n"))
}
