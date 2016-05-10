package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type CurrentStats struct {
	sync.RWMutex
	Hits           int64
	TotalHits      int64
	LastCalculated time.Time
	ResponseTimes  []time.Duration
	ResponseCodes  []int
}

type GlobalStatRecord struct {
	Config               *Config
	StartTime            time.Time
	EndTime              time.Time
	AverageResponseTime  time.Duration
	ResponseStatusCounts map[int]int
	TotalRequests        int64
}

var globalCurrentStats CurrentStats
var configCurrentStats map[*Config]*CurrentStats
var globalStatSink chan GlobalStatRecord
var globalStatSinkSubscribers []chan GlobalStatRecord

func recordStat(config *Config, res *http.Response, elapsed time.Duration) {
	_, exists := configCurrentStats[config]
	if !exists {
		configCurrentStats[config] = &CurrentStats{}
	}
	addToStat(&globalCurrentStats, res, elapsed)
	addToStat(configCurrentStats[config], res, elapsed)
}

func addToStat(stat *CurrentStats, res *http.Response, elapsed time.Duration) {
	stat.Lock()
	stat.Hits++
	stat.TotalHits++
	stat.ResponseTimes = append(stat.ResponseTimes, elapsed)
	stat.ResponseCodes = append(stat.ResponseCodes, res.StatusCode)
	stat.Unlock()
}

func processStats() {
	globalStatSink <- processStat(nil, &globalCurrentStats)
	for config, stat := range configCurrentStats {
		globalStatSink <- processStat(config, stat)
	}
}

func processStat(config *Config, currentStat *CurrentStats) GlobalStatRecord {
	var statRecord GlobalStatRecord
	statRecord.Config = config

	currentStat.Lock()
	statRecord.TotalRequests = currentStat.Hits
	statRecord.StartTime = currentStat.LastCalculated
	statRecord.EndTime = time.Now()

	currentStat.Hits = 0
	currentStat.LastCalculated = statRecord.EndTime
	responseTimes := currentStat.ResponseTimes
	responseCodes := currentStat.ResponseCodes
	currentStat.ResponseTimes = []time.Duration{}
	currentStat.ResponseCodes = []int{}
	currentStat.Unlock()

	statRecord.AverageResponseTime = GetAverageDuration(responseTimes)
	statRecord.ResponseStatusCounts = MakeFrequencyMap(responseCodes)
	return statRecord
}

func statProcessor() {
	for {
		processStats()
		time.Sleep(5 * time.Second)
	}
}

func SubscribeGlobalStats(c chan GlobalStatRecord) {
	globalStatSinkSubscribers = append(globalStatSinkSubscribers, c)
}

var configStatSubscribers map[*Config][]chan GlobalStatRecord

func SubscribeConfigStats(config *Config, c chan GlobalStatRecord) {
	subscribers, exists := configStatSubscribers[config]
	if !exists {
		subscribers = make([]chan GlobalStatRecord, 1)
	}
	configStatSubscribers[config] = append(subscribers, c)
}

func GlobalStatBroadcaster() {
	for {
		astat := <-globalStatSink
		if astat.Config == nil {
			for _, c := range globalStatSinkSubscribers {
				go func() {
					c <- astat
				}()
			}
		} else {
			for _, c := range configStatSubscribers[astat.Config] {
				go func() {
					c <- astat
				}()
			}
		}
	}
}

func reqsPrinter() {
	c := make(chan GlobalStatRecord)
	nc := make(chan GlobalStatRecord)
	SubscribeGlobalStats(c)
	SubscribeConfigStats(defaultConfig, nc)
	for {
		select {
		case astat := <-c:
			log.Println(astat)
		case astat := <-nc:
			log.Println(astat)
		}
	}
}
