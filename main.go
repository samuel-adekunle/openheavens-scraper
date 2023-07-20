package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var DATE_FORMATS = []string{
	"2-January-2006-Monday",
	"2-January-2006",
}

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
const BASE_URL = "https://rccgonline.org/open-heaven-"

func parsePostHTML(e *colly.HTMLElement) (post *Post) {
	post = new(Post)
	post.Title = e.DOM.Find("div.entry-content > h1").Clone().Children().Remove().End().Text()
	removeMediaRegexp := regexp.MustCompile(`@media.*;`)
	// mode 0 - memory verse
	// mode 1 - bible reading heading
	// mode 2 - bible reading body
	// mode 3 - wait for empty line
	// mode 4 - message body & skip action point header
	// mode 5 - action point body
	// mode 6 - wait for empty line
	// mode 7 - skip hymn header
	// mode 8 - hymn title & (maybe) bible in one year
	// mode 9 - (if hymn) hymn body & bible in one year
	// mode 10 - post.Authors & finish
	mode := 0
	e.ForEach("div.entry-content > p", func(_ int, el *colly.HTMLElement) {
		switch mode {
		case 0:
			if strings.Contains(el.Text, "MEMORIZE") {
				post.MemoryVerse = strings.TrimSpace(el.DOM.Find("em, span").Clone().Children().Remove().End().Text())
				post.MemoryVerse = post.MemoryVerse + " " + strings.TrimSpace(el.DOM.Find("strong:nth-child(3)").Clone().Children().Remove().End().Text())
				mode = 1
				return
			}
		case 1:
			if strings.Contains(el.Text, "BIBLE READING") {
				post.BibleReadingHeading = strings.TrimSpace(el.DOM.Find("strong").Clone().Children().Remove().End().Text())
				el.ForEach("span", func(_ int, ell *colly.HTMLElement) {
					post.BibleReadingBody = append(post.BibleReadingBody, strings.TrimSpace(removeMediaRegexp.ReplaceAllString(ell.Text, "")))
				})
				mode = 2
			}
		case 2:
			if el.Text == string([]byte{194, 160}) {
				mode = 3
			}
		case 3:
			if el.Text == string([]byte{194, 160}) {
				mode = 4
			}
		case 4:
			if strings.Contains(el.Text, "POINT") {
				mode = 5
			} else {
				tmp := strings.TrimSpace(el.DOM.Find("span").Clone().Children().Remove().End().Text())
				if len(tmp) > 0 {
					post.MessageBody = append(post.MessageBody, strings.TrimSpace(
						removeMediaRegexp.ReplaceAllString(el.ChildText("span"), "")),
					)
				}
			}
		case 5:
			post.ActionPoint = strings.TrimSpace(el.DOM.Find("span").Clone().Children().Remove().End().Text())
			mode = 6
		case 6:
			if el.Text == string([]byte{194, 160}) {
				mode = 7
			} else if strings.Contains(el.Text, "HYMN") {
				mode = 8
			}
		case 7:
			if strings.Contains(el.Text, "HYMN") {
				mode = 8
			}
		case 8:
			if strings.Contains(el.Text, "ONE YEAR") {
				post.BibleInOneYear = strings.TrimSpace(el.DOM.Find("strong").Clone().Children().Remove().End().Text())
				mode = 10
			} else {
				post.HymnTitle = strings.TrimSpace(el.DOM.Find("strong").Clone().Children().Remove().End().Text())
				mode = 9
			}
		case 9:
			if strings.Contains(el.Text, "ONE YEAR") {
				post.BibleInOneYear = strings.TrimSpace(el.DOM.Find("strong").Clone().Children().Remove().End().Text())
				post.HymnBody = post.HymnBody[:len(post.HymnBody)-1]
				mode = 10
			} else {
				el.ForEach("strong", func(_ int, ell *colly.HTMLElement) {
					post.HymnBody = append(post.HymnBody, strings.TrimSpace(ell.Text))
				})
				post.HymnBody = append(post.HymnBody, " ")
			}
		case 10:
			el.ForEach("strong", func(_ int, ell *colly.HTMLElement) {
				post.Authors = append(post.Authors, strings.TrimSpace(strings.Trim(ell.Text, "AUTHORS: ")))
			})
			mode = 11
		default:
			log.Println("Something went wrong: mode = ", mode)
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
	c.OnHTML("article", func(e *colly.HTMLElement) {
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
	for _, format := range DATE_FORMATS {
		date = strings.ToLower(time.Now().Format(format))
		post = scrapePost(date)
		if post != nil {
			return post, date
		}
	}
	return nil, ""
}

func savePost(post *Post, date string) {
	f, err := os.Create(date + ".txt")
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
	log.Fatalln("Post not found!")
}
