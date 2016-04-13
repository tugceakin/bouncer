package main

import (
	"net/http"
	_ "net/http/pprof"
	"net/url"
)

var defaultConfigRule ConfigRule
var ruleBackendChannelMap map[ConfigRuleId]chan string

var globalStatSink chan GlobalStatRecord

func main() {
	globalStatSink = make(chan GlobalStatRecord)

	proxy := NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:9091",
		Path:   "/abcd/",
	})

	defaultConfigRule = ConfigRule{
		Host:           "localhost:9091",
		BackendServers: []string{"localhost:9090", "localhost:9089"},
	}
	ruleBackendChannelMap = MakeBackendChannelMap([]ConfigRule{defaultConfigRule})

	go statProcessor()
	go reqsPrinter()
	http.ListenAndServe(":9090", proxy)
}
