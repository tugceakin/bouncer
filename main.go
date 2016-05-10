package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

var defaultConfig *Config

func main() {
	globalStatSink = make(chan GlobalStatRecord)
	globalStatSinkSubscribers = make([]chan GlobalStatRecord, 0)
	proxy := &ReverseProxy{Director: director}
	config := NewConfig(
		"localhost:9090",
		[]BackendServer{NewBackendServer("localhost:9091"), NewBackendServer("localhost:9092")},
		"",
		"",
		10,
		200)
	defaultConfig = &config
	go statProcessor()
	go GlobalStatBroadcaster()
	go reqsPrinter()
	go UIServer()
	http.ListenAndServe(":9090", proxy)
}

func director(req *http.Request) (*Config, *BackendServer) {
	config := defaultConfig
	select {
	case next := <-config.NextBackendServer:
		req.URL.Scheme = "http"
		req.URL.Host = next.Host
		return config, &next
	case <-time.After(120 * time.Second):
		return nil, nil

	}
}
