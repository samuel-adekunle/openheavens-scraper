package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var (
	date     *string
	textPath *string
	htmlPath *string
)

func init() {
	date = flag.String("date", "", "date to scrape")
	textPath = flag.String("text", "post.txt", "path to save post text")
	htmlPath = flag.String("html", "post.html", "path to save post html")
}

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
			post.MemoryVerse = sanitizeString(el.Text)
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
				post.IsActionPoint = strings.Contains(el.Text, "ACTION")
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
			html, _ := el.DOM.Html()
			tmp := sanitizeString(strings.Replace(html, "<br/>", "\n", -1))
			if (len(tmp) > 0 && (tmp[0] >= '0' && tmp[0] <= '9')) || strings.Contains(tmp, "Refrain") || strings.Contains(tmp, "Chorus") {
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

func scrapePost() *Post {
	var post *Post
	postURL := fmt.Sprintf("https://flatimes.com/open-heaven-%s/", *date)
	c := colly.NewCollector()
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
		if !strings.Contains(e.Request.URL.String(), *date) {
			return
		}
		post = parsePostHTML(e)
	})
	c.Visit(postURL)
	c.Wait()
	return post
}

func savePostText(post *Post, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	_, err = file.WriteString(post.String())
	if err != nil {
		log.Fatalln(err)
	}
}

func savePostHTML(post *Post, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	err = post.HTMl(file)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	if *date == "" {
		log.Fatalln("Date not provided!")
	}
	post := scrapePost()
	if post == nil {
		log.Println("Post not scraped successfully!")
		return
	}
	log.Println("Post scraped successfully!")
	if len(*textPath) > 0 {
		savePostText(post, *textPath)
		log.Println("Post text saved successfully!")
	}
	if len(*htmlPath) > 0 {
		savePostHTML(post, *htmlPath)
		log.Println("Post html saved successfully!")
	}
	err := sendNotification(post)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Notification sent successfully!")
	}
}
