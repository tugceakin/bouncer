package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var defaultConfig *Config
var configStore ConfigStore

func main() {
	configStore = make(ConfigStore)
	globalStatSink = make(chan GlobalStatRecord, 10)
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
	configCurrentStats = make(map[*Config]*CurrentStats)
	configStatSubscribers = make(map[*Config][]chan GlobalStatRecord)
	go statProcessor()
	go GlobalStatBroadcaster()
	go reqsPrinter()
	go UIServer()
	defaultConfig.MaxConcurrentPerBackendServer = 20
	defaultConfig.ReqPerSecond = 10
	defaultConfig.Reload()
	http.ListenAndServe(":9090", proxy)
}

func director(req *http.Request) (*Config, *BackendServer) {
	config := defaultConfig
	<-config.Throttle
	log.Println("XX", config)
	select {
	case next := <-config.NextBackendServer:
		req.URL.Scheme = "http"
		req.URL.Host = next.Host
		return config, &next
	case <-time.After(120 * time.Second):
		return nil, nil

	}
}
