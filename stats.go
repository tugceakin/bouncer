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
}

type StatRecord struct {
	StartTime            time.Time
	EndTime              time.Time
	AverageResponseTime  time.Duration
	ResponseStatusCounts map[int]int
	TotalRequests        int64
}

var currentStats CurrentStats

func recordStat(res *http.Response, elapsed time.Duration) {
	// TODO Mutex vs Channels?
	currentStats.Lock()
	currentStats.Hits++
	currentStats.TotalHits++
	currentStats.ResponseTimes = append(currentStats.ResponseTimes, elapsed)
	currentStats.Unlock()
}

func printReqsec() {
	var statRecord StatRecord

	currentStats.Lock()
	statRecord.TotalRequests = currentStats.Hits
	statRecord.StartTime = currentStats.LastCalculated
	statRecord.EndTime = time.Now()

	currentStats.Hits = 0
	currentStats.LastCalculated = statRecord.EndTime
	responseTimes := currentStats.ResponseTimes
	currentStats.ResponseTimes = []time.Duration{}

	currentStats.Unlock()

	statRecord.AverageResponseTime = AverageDuration(responseTimes)

	log.Println(statRecord)
}

func reqsPrinter() {
	for {
		printReqsec()
		<-time.After(10 * time.Second)
	}
}
