package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/heptiolabs/healthcheck"
)

var port *string

func init() {
	port = flag.String("port", "8080", "port to serve")
}

func getPostsHtml() []string {
	files, err := os.ReadDir("../posts/")
	if err != nil {
		log.Fatalln(err)
	}
	var posts []string
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".html") {
			posts = append(posts, file.Name())
		}
	}
	return posts
}

func createHealthHandler() http.Handler {
	health := healthcheck.NewHandler()
	return http.HandlerFunc(health.LiveEndpoint)
}

func createPostsFileServer(prefix string) http.Handler {
	fs := http.FileServer(http.Dir("../posts"))
	return http.StripPrefix(prefix, fs)
}

func createIndexHandler(posts []string) http.Handler {
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		err = tmpl.ExecuteTemplate(w, "index", posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func main() {
	flag.Parse()
	posts := getPostsHtml()
	mux := http.NewServeMux()
	mux.Handle("/", createIndexHandler(posts))
	mux.Handle("/healthz", createHealthHandler())
	mux.Handle("/posts/", createPostsFileServer("/posts/"))
	addr := fmt.Sprintf("0.0.0.0:%s", *port)
	log.Printf("Server is running at %s\n", addr)
	http.ListenAndServe(addr, mux)
}
