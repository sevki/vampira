package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"runtime"

	"sevki.org/namespace"
)

type sitemap map[string]string

func (s sitemap) String() string {
	return fmt.Sprint(map[string]string(s))
}

func (s sitemap) Set(value string) error {
	u, err := url.Parse("https://" + value)
	if err != nil {
		panic(err)
	}
	s[u.Host] = u.Path
	return nil
}

func (s sitemap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir := "/tmp/vampira"
	if d, ok := s[r.Host]; ok {
		dir = d
	}
	w.Header().Add("Server-Agent", fmt.Sprintf("vampira/master (%s; %s);", runtime.GOOS, runtime.GOARCH))
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	return
}

func main() {
	var sitemap = sitemap{}
	nsfile := flag.String("n", "/lib/namespace.httpd", "namespace file")
	addr := flag.String("http", ":3000", "http address to listen to")
	flag.Var(&sitemap, "map", "urls")
	flag.Parse()

	if err := namespace.AddNS(*nsfile, "web"); err != nil {
		fmt.Printf("couldn't build namespace: %v\n", err)
		return
	}

	fmt.Printf("Listening on %s...\n", *addr)

	if err := http.ListenAndServe(*addr, sitemap); err != nil {
		fmt.Printf("couldn't start httpserver: %v\n", err)
		return
	}
}
