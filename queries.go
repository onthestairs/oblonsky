package main

import (
	"time"
)

const oxenholme = "OXN"
const manchester = "115"
const wakefield = "164"
const londonStations = "182"

func generateTimes(numberOfDays int) []time.Time {
	var times []time.Time
	now := time.Now()
	for i := 1; i <= (2 * 24 * numberOfDays); i++ {
		delta := 30 * time.Duration(i) * time.Minute
		t := now.Add(delta)
		if t.Hour() > 6 && t.Hour() < 22 {
			times = append(times, t)
		}
	}
	return times
}

func generateQueries() []JourneyQuery {
	northStations := []string{wakefield, manchester, oxenholme}
	times := generateTimes(60)
	var queries []JourneyQuery
	for _, t := range times {
		for _, northStation := range northStations {
			southQuery := makeQuery(northStation, londonStations, t)
			northQuery := makeQuery(londonStations, northStation, t)
			queries = append(queries, southQuery, northQuery)
		}
	}
	return queries
}
