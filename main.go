package main

import (
	"flag"
	"fmt"
	"net/http"

	"sevki.org/namespace"
)

func main() {
	nsfile := flag.String("n", "/lib/namespace.httpd", "namespace file")
	root := flag.String("w", "/", "webroot")
	addr := flag.String("http", ":3000", "http address to listen to")
	flag.Parse()

	if err := namespace.AddNS(*nsfile); err != nil {
		fmt.Printf("couldn't build namespace: %v\n", err)
		return
	}

	http.Handle("/", http.FileServer(http.Dir(*root)))

	fmt.Printf("Listening on %s...\n", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Printf("couldn't start httpserver: %v\n", err)
		return
	}
}
