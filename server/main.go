package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/heptiolabs/healthcheck"
)

func addHealthChecks(health healthcheck.Handler) {
	// TODO(samuel-adekunle): add liveness and readiness checks
	// TODO(samuel-adekunle): add prometheus metrics
}

func addHealthEndpoints(health healthcheck.Handler) {
	http.Handle("/healthz", http.HandlerFunc(health.LiveEndpoint))
	http.Handle("/statusz", http.HandlerFunc(health.ReadyEndpoint))
}

func addEndpoints(health healthcheck.Handler) {
	addHealthEndpoints(health)
	http.Handle("/", http.FileServer(http.Dir("../")))
}

func main() {
	port := flag.Int("port", 8080, "port to serve")
	flag.Parse()
	health := healthcheck.NewHandler()
	addHealthChecks(health)
	addEndpoints(health)
	fmt.Printf("Server is running at port %d\n", *port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil)
}
