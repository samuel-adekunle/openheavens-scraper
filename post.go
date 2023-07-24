package main

import (
	"fmt"
	"html/template"
	"io"
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

func (p Post) String() string {
	pointString := "PRAYER POINT"
	if p.IsActionPoint {
		pointString = "ACTION POINT"
	}
	return fmt.Sprintf("*%s*\n\n*MEMORY VERSE*\n%s\n\n*BIBLE READING*\n%s\n\n%s\n\n*MESSAGE*\n%s\n\n*%s*\n%s\n\n*HYMN*\n%s\n\n%s\n\n*BIBLE IN ONE YEAR*\n%s\n", p.Title, p.MemoryVerse, p.BibleReadingHeading, strings.Join(p.BibleReadingBody, "\n\n"), strings.Join(p.MessageBody, "\n\n"), pointString, p.PrayerPoint, p.HymnTitle, strings.Join(p.HymnBody, "\n\n"), p.BibleInOneYear)
}

func (p *Post) HTMl(w io.Writer) (err error) {
	t := template.Must(template.ParseFiles("templates/post.gohtml"))
	hymnBodyCopy := make([]string, len(p.HymnBody))
	copy(hymnBodyCopy, p.HymnBody)
	newHymnBody := []string{}
	for i, line := range p.HymnBody {
		newHymnBody = append(newHymnBody, strings.Split(line, "\n")...)
		if i != len(p.HymnBody)-1 {
			newHymnBody = append(newHymnBody, "")
		}
	}
	p.HymnBody = newHymnBody
	err = t.ExecuteTemplate(w, "post", p)
	p.HymnBody = hymnBodyCopy
	return
}
