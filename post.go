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
	HymnTitle           string   `json:"hymnTitle"`
	HymnBody            []string `json:"hymnBody"`
	BibleInOneYear      string   `json:"bibleInOneYear"`
}

// string representation of a post
func (p Post) String() string {
	return fmt.Sprintf("*%s*\n\n*MEMORISE:*\n%s\n\n*%s*\n%s\n\n*%s*\n\n*MESSAGE:*\n%s\n\n*PRAYER POINT:*\n%s\n\n*%s*\n%s\n", p.Title, p.MemoryVerse, p.BibleReadingHeading, strings.Join(p.BibleReadingBody, "\n\n"), p.BibleInOneYear, strings.Join(p.MessageBody, "\n\n"), p.PrayerPoint, p.HymnTitle, strings.Join(p.HymnBody, "\n\n"))
}
