package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	// TODO: delete this comment, testing auto deploy
	port := flag.Int("port", 8080, "port to serve")
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("../")))
	fmt.Printf("Server is running at port %d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
