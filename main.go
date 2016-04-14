package main

import (
	"net/http"
	_ "net/http/pprof"
)

var globalStatSink chan GlobalStatRecord
var defaultConfig Config

func main() {
	globalStatSink = make(chan GlobalStatRecord)
	proxy := &ReverseProxy{Director: director}
	defaultConfig = NewConfig("localhost:9091", []BackendServer{NewBackendServer("localhost:9091"), NewBackendServer("localhost:9092")})

	go statProcessor()
	go reqsPrinter()
	go UIServer()
	http.ListenAndServe(":9090", proxy)
}

func director(req *http.Request) {
	next := <-defaultConfig.NextBackendServer
	req.URL.Scheme = "http"
	req.URL.Host = next.Host
}
