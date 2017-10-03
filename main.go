package main

import (
	"fmt"
	"time"
)

// const timeFormat = "Jan 2, 2006 at 3:04pm (MST)"
// const trainlineTimeFormat = "2006-01-02T15:04:05"

func main() {

	db := getDb()
	defer db.Close()
	initDatabase(db)

	t := time.Date(2017, time.October, 13, 12, 0, 0, 0, time.UTC)
	query := makeQuery("OXN", "182", t)
	result := doQuery(query)
	journeys := resultToNiceJourneys(result)

	for i, journey := range journeys {
		insertJourney(db, journey)
		fmt.Printf("Journey %d: %s (%s) -> %s (%s) £%f [%d legs %s, %s]\n", i, journey.departureTime.Format("2 Jan 15:04"), journey.departureStation, journey.arrivalTime.Format("2 Jan 15:04"), journey.arrivalStation, journey.price, journey.numberOfLegs, journey.ticketType, journey.ticketClass)
	}
}
