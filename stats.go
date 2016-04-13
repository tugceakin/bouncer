package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Stat struct {
	ResponseTime time.Duration
	StatusCode   int
}

type Stats struct {
	records []Stat
}

var stats Stats

type CurrentStats struct {
	sync.RWMutex
	Hits           int64
	TotalHits      int64
	LastCalculated time.Time
	ResponseTimes  []time.Duration
	ResponseCodes  []int
}

type GlobalStatRecord struct {
	StartTime            time.Time
	EndTime              time.Time
	AverageResponseTime  time.Duration
	ResponseStatusCounts map[int]int
	TotalRequests        int64
}

var globalCurrentStats CurrentStats

func recordStat(res *http.Response, elapsed time.Duration) {
	// TODO Mutex vs Channels?
	globalCurrentStats.Lock()
	globalCurrentStats.Hits++
	globalCurrentStats.TotalHits++
	globalCurrentStats.ResponseTimes = append(globalCurrentStats.ResponseTimes, elapsed)
	globalCurrentStats.ResponseCodes = append(globalCurrentStats.ResponseCodes, res.StatusCode)
	globalCurrentStats.Unlock()
}

func processGlobalStats() {
	var statRecord GlobalStatRecord

	globalCurrentStats.Lock()
	statRecord.TotalRequests = globalCurrentStats.Hits
	statRecord.StartTime = globalCurrentStats.LastCalculated
	statRecord.EndTime = time.Now()

	globalCurrentStats.Hits = 0
	globalCurrentStats.LastCalculated = statRecord.EndTime
	responseTimes := globalCurrentStats.ResponseTimes
	responseCodes := globalCurrentStats.ResponseCodes
	globalCurrentStats.ResponseTimes = []time.Duration{}
	globalCurrentStats.ResponseCodes = []int{}

	globalCurrentStats.Unlock()

	statRecord.AverageResponseTime = GetAverageDuration(responseTimes)
	statRecord.ResponseStatusCounts = MakeFrequencyMap(responseCodes)

	globalStatSink <- statRecord
}

func statProcessor() {
	for {
		processGlobalStats()
		time.Sleep(10 * time.Second)
	}
}

func reqsPrinter() {
	for {
		select {
		case astat := <-globalStatSink:
			log.Println(astat)
		}
	}
}
