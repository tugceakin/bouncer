package main

import (
	"net/http"
	_ "net/http/pprof"
	"net/url"
)

func main() {
	proxy := NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:9091",
		Path:   "/abcd/",
	})
	go reqsPrinter()
	http.ListenAndServe(":9090", proxy)
}
