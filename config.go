package main

import (
	"net/url"
)

type Config struct {
	Port  int
	Rules []ConfigRule
}

type ConfigRuleId int64

type ConfigRule struct {
	Id             ConfigRuleId
	Host           string
	BackendServers []BackendServer
}

type BackendServer struct {
	Url url.URL
}

func MakeBackendChannelMap(rules []ConfigRule) map[ConfigRuleId]chan string {
	ret := make(map[ConfigRuleId]chan string)
	for _, rule := range rules {
		ret[rule.Id] = make(chan string)
	}
	return ret
}
