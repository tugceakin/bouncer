package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// func MakeBackendChannelMap(rules []ConfigRule) map[ConfigRuleId]chan string {
// 	ret := make(map[ConfigRuleId]chan string)
// 	for _, rule := range rules {
// 		ret[rule.Id] = make(chan string)
// 	}
// 	ret
// }

var db *sql.DB

type Config struct {
	Id                int
	Host              string
	BackendServers    []BackendServer
	NextBackendServer chan BackendServer
	done              chan struct{}
}

type BackendServer struct {
	Id       int
	Host     string
	ConfigId int
}

func (config *Config) NextBackendServerRoutine() {
	for {
		for _, next := range config.BackendServers {
			select {
			case config.NextBackendServer <- next:

			case <-config.done:
				return
			}
		}

	}
}

func (config *Config) Destroy() {
	config.done <- struct{}{}
}

func NewConfig(hostname string, backendServers []BackendServer) Config {
	var config Config
	config.BackendServers = backendServers
	config.NextBackendServer = make(chan BackendServer)
	config.done = make(chan struct{})
	go config.NextBackendServerRoutine()
	return config
}

func NewBackendServer(url string) BackendServer {
	var backendServer BackendServer
	backendServer.Host = url
	return backendServer
}
