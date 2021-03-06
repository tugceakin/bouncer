package main

import "time"

func GetAverageDuration(durations []time.Duration) time.Duration {
	length := len(durations)
	var sum time.Duration
	for _, i := range durations {
		sum = sum + i
	}
	if sum == 0 {
		return time.Duration(0)
	}
	return time.Duration(int64(sum) / int64(length))
}

func MakeFrequencyMap(responseCodes []int) map[int]int {
	ret := make(map[int]int)
	for _, code := range responseCodes {
		ret[code]++
	}
	return ret
}
