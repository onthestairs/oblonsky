package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func getDb() (*sql.DB, error) {
	return sql.Open("sqlite3", "./train-prices.db")
}

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS journeyPrices (
		id INTEGER PRIMARY KEY,
    observedTime DATETIME,
		departureStation TEXT,
		departureTime DATETIME,
		arrivalStation TEXT,
		arrivalTime DATETIME,
		price FLOAT,
    numberOfLegs TEXT,
    ticketType TEXT,
    ticketClass TEXT
	)
`

func initDatabase(db *sql.DB) error {
	statement, err := db.Prepare(createTableQuery)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

const insertSQL = `
  INSERT INTO journeyPrices (
    observedTime,
    departureStation,
    departureTime,
    arrivalStation,
    arrivalTime,
    price,
    numberOfLegs,
    ticketType,
    ticketClass
  ) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
  )
`

func insertJourney(db *sql.DB, journey NiceJourney) error {
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		time.Now(),
		journey.departureStation,
		journey.departureTime,
		journey.arrivalStation,
		journey.arrivalTime,
		journey.price,
		journey.numberOfLegs,
		journey.ticketType,
		journey.ticketClass,
	)
	if err != nil {
		return err
	}
	return nil
}
