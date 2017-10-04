package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"
)

// const timeFormat = "Jan 2, 2006 at 3:04pm (MST)"
// const trainlineTimeFormat = "2006-01-02T15:04:05"

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func main() {

	f, err := os.OpenFile("train-api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	db, err := getDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = initDatabase(db)
	if err != nil {
		panic(err)
	}

	for {
		queries := generateQueries()
		for _, query := range queries {
			log.Printf("Doing %s -> %s (%s)\n", query.Origin, query.Destination, query.OutboundJourney.Time)
			err = runQuery(db, query)
			if err != nil {
				log.Printf("ERROR!: %s", err.Error())
			}
			secondsToWait := random(3, 15)
			time.Sleep(time.Duration(secondsToWait) * time.Second)
		}
		time.Sleep(1 * time.Minute)
	}

}

func runQuery(db *sql.DB, query JourneyQuery) error {
	result, err := doQuery(query)
	if err != nil {
		return err
	}

	journeys := resultToNiceJourneys(result)

	for _, journey := range journeys {
		err = insertJourney(db, journey)
		if err != nil {
			return err
		}
	}
	// fmt.Printf("Journey %d: %s (%s) -> %s (%s) £%f [%d legs %s, %s]\n", i, journey.departureTime.Format("2 Jan 15:04"), journey.departureStation, journey.arrivalTime.Format("2 Jan 15:04"), journey.arrivalStation, journey.price, journey.numberOfLegs, journey.ticketType, journey.ticketClass)
	log.Printf("Saved %d journeys\n", len(journeys))
	return nil
}
