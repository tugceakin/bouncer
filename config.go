package main

type Config struct {
	Port  int
	Rules []ConfigRule
}

type ConfigRule struct {
	Host           string
	BackendServers []string
}
