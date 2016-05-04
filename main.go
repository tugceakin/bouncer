package main

import (
	"net/http"
	_ "net/http/pprof"
)

var defaultConfig Config

func main() {
	globalStatSink = make(chan GlobalStatRecord)
	globalStatSinkSubscribers = make([]chan GlobalStatRecord, 0)
	proxy := &ReverseProxy{Director: director}
	defaultConfig = NewConfig("localhost:9090", []BackendServer{NewBackendServer("localhost:9091"), NewBackendServer("localhost:9092")})

	go statProcessor()
	go GlobalStatBroadcaster()
	go reqsPrinter()
	go UIServer()
	http.ListenAndServe(":9090", proxy)
}

func director(req *http.Request) {
	next := <-defaultConfig.NextBackendServer
	req.URL.Scheme = "http"
	req.URL.Host = next.Host
}
