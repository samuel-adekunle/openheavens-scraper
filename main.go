package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	DATE_FORMAT = "02-January-2006"
	USER_AGENT  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	BASE_URL    = "https://flatimes.com/open-heaven-"
)

func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func parsePostHTML(e *colly.HTMLElement) (post *Post) {
	post = new(Post)
	state := 0

	e.ForEach(".et_pb_text_inner h2, .et_pb_text_inner p:not(.has-text-align-center)", func(_ int, el *colly.HTMLElement) {
		switch state {
		case 0:
			if el.Name == "h2" {
				state += 1
			} else {
				state = -1
			}
		case 1:
			post.Title = sanitizeString(el.Text)
			state += 1
		case 2:
			post.MemoryVerse = sanitizeString(el.DOM.Clone().Children().Remove().End().Text())
			state += 1
		case 3:
			post.BibleReadingHeading = sanitizeString(el.Text)
			state += 1
		case 4:
			if strings.Contains(el.Text, "YEAR") {
				post.BibleInOneYear = sanitizeString(el.Text)
				state += 1
			} else {
				post.BibleReadingBody = append(post.BibleReadingBody, sanitizeString(el.Text))
			}
		case 5:
			if el.Name == "h2" {
				state += 1
			} else {
				state = -1
			}
		case 6:
			if el.Name == "h2" {
				state += 1
			} else {
				post.MessageBody = append(post.MessageBody, sanitizeString(el.Text))
			}
		case 7:
			post.PrayerPoint = sanitizeString(el.Text)
			state += 1
		case 8:
			if el.Name == "h2" {
				post.HymnTitle = sanitizeString(el.Text)
				state += 1
			} else {
				state = -1
			}
		case 9:
			tmp := sanitizeString(el.Text)
			if len(tmp) > 0 && tmp[0] >= '0' && tmp[0] <= '9' {
				post.HymnBody = append(post.HymnBody, tmp)
			} else {
				return
			}
		default:
			log.Println("Invalid post!")
			return
		}
	})
	return post
}

func scrapePost(date string) *Post {
	var post *Post
	postURL := BASE_URL + date + "/"
	c := colly.NewCollector()
	c.UserAgent = USER_AGENT
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Page visited: ", r.Request.URL)
	})
	c.OnScraped(func(r *colly.Response) {
		log.Println(r.Request.URL, " scraped!")
	})
	c.OnHTML(".et_pb_text_inner", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.String(), date) {
			return
		}
		post = parsePostHTML(e)
	})
	c.Visit(postURL)
	c.Wait()
	return post
}

func scrapeToday() (post *Post, date string) {
	date = strings.ToLower(time.Now().Format(DATE_FORMAT))
	post = scrapePost(date)
	return post, date
}

func savePost(post *Post, date string) {
	f, err := os.Create(fmt.Sprintf("./posts/%s.txt", date))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(post.String())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	post, date := scrapeToday()
	if post != nil {
		savePost(post, date)
		log.Println("Post saved!")
		return
	}
	log.Println("Post not found!")
}
