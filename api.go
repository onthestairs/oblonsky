package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Journey struct {
	ID                int       `json:"id"`
	Origin            string    `json:"origin"`
	Destination       string    `json:"destination"`
	DepartureDateTime time.Time `json:"departureDateTime"`
	ArrivalDateTime   time.Time `json:"arrivalDateTime"`
	Direction         string    `json:"direction"`
	Legs              []struct {
		ID     int `json:"id"`
		Origin struct {
			StationCode    string    `json:"stationCode"`
			ScheduledTime  time.Time `json:"scheduledTime"`
			RealTimeStatus string    `json:"realTimeStatus"`
			Platform       string    `json:"platform"`
			PlatformStatus string    `json:"platformStatus"`
		} `json:"origin"`
		Destination struct {
			StationCode    string    `json:"stationCode"`
			ScheduledTime  time.Time `json:"scheduledTime"`
			RealTimeStatus string    `json:"realTimeStatus"`
			Platform       string    `json:"platform"`
			PlatformStatus string    `json:"platformStatus"`
		} `json:"destination"`
		TransportMode         string   `json:"transportMode"`
		ReservationFlag       string   `json:"reservationFlag"`
		RetailTrainIdentifier string   `json:"retailTrainIdentifier,omitempty"`
		TrainID               string   `json:"trainId,omitempty"`
		ServiceProviderCode   string   `json:"serviceProviderCode"`
		ServiceProviderName   string   `json:"serviceProviderName"`
		SeatingClass          string   `json:"seatingClass"`
		IsCancelled           bool     `json:"isCancelled"`
		FinalDestinations     []string `json:"finalDestinations,omitempty"`
		BusyData              struct {
			Coaches []interface{} `json:"coaches"`
		} `json:"busyData,omitempty"`
	} `json:"legs"`
	WalkUpFareCategory string `json:"walkUpFareCategory"`
	JourneyStatus      string `json:"journeyStatus"`
}

type Ticket struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	TicketType             string `json:"ticketType"`
	TicketClass            string `json:"ticketClass"`
	ThroughLondon          bool   `json:"throughLondon"`
	TotalFare              int    `json:"totalFare"`
	RouteRestriction       string `json:"routeRestriction"`
	FareRestrictionCode    string `json:"fareRestrictionCode"`
	FareOriginNlc          string `json:"fareOriginNlc"`
	FareOriginStation      string `json:"fareOriginStation"`
	FareDestinationNlc     string `json:"fareDestinationNlc"`
	FareDestinationStation string `json:"fareDestinationStation"`
	FareSource             string `json:"fareSource"`
	Fares                  []struct {
		Code               string `json:"code"`
		Price              int    `json:"price"`
		NumberOfPassengers int    `json:"numberOfPassengers"`
		FareType           string `json:"fareType"`
		PassengerType      string `json:"passengerType"`
	} `json:"fares"`
	OutboundValidity struct {
		Days   int `json:"days"`
		Months int `json:"months"`
		Years  int `json:"years"`
	} `json:"outboundValidity"`
	InboundValidity struct {
		Days   int `json:"days"`
		Months int `json:"months"`
		Years  int `json:"years"`
	} `json:"inboundValidity"`
	ReservationRequired bool   `json:"reservationRequired"`
	FareCategory        string `json:"fareCategory"`
	IsPromotional       bool   `json:"isPromotional"`
	AdvanceTiers        []int  `json:"advanceTiers,omitempty"`
}

type JourneySearchResult struct {
	JourneySearchID int64     `json:"journeySearchId"`
	Journeys        []Journey `json:"journeys"`
	Tickets         []Ticket  `json:"tickets"`
	JourneyTickets  []struct {
		JourneyID            int      `json:"journeyId"`
		TicketID             int      `json:"ticketId"`
		SeatAvailabilityCode string   `json:"seatAvailabilityCode"`
		DeliveryOptions      []string `json:"deliveryOptions"`
		SeatsRemaining       string   `json:"seatsRemaining,omitempty"`
	} `json:"journeyTickets"`
	RouteRestrictions []struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"routeRestrictions"`
	BookingHorizon struct {
		Mobile  int `json:"Mobile"`
		Kiosk   int `json:"Kiosk"`
		ETicket int `json:"ETicket"`
	} `json:"bookingHorizon"`
}

type OutboundJourney struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type JourneyQuery struct {
	Origin                   string          `json:"origin"`
	Destination              string          `json:"destination"`
	JourneyType              string          `json:"journeyType"`
	MaxJourneys              int             `json:"maxJourneys"`
	Adults                   int             `json:"adults"`
	Children                 int             `json:"children"`
	OutboundJourney          OutboundJourney `json:"outboundJourney"`
	Railcards                []interface{}   `json:"railcards"`
	TimeTableOnly            bool            `json:"timeTableOnly"`
	EticketSupported         int             `json:"eticketSupported"`
	PricePredictionSupported int             `json:"pricePredictionSupported"`
	ShowCancelledTrains      int             `json:"ShowCancelledTrains"`
}

func makeQuery(origin string, destination string, departAtTime time.Time) JourneyQuery {
	outboundJourney := OutboundJourney{"LeaveAfter", departAtTime.Format(time.RFC3339)}
	var railcards []interface{}
	return JourneyQuery{origin, destination, "Single", 5, 1, 0, outboundJourney, railcards, false, 1, 0, 1}
}

func doQuery(query JourneyQuery) JourneySearchResult {
	url := "https://api.thetrainline.com/mobile/journeys"

	body, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("X-Feature", "Plan and buy")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-gb")
	req.Header.Set("X-Api-Version", "2.0")
	req.Header.Set("X-Platform-Type", "iOS")
	req.Header.Set("Content-Length", "406")
	req.Header.Set("User-Agent", "Trainline/16772 CFNetwork/887 Darwin/17.0.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Consumer-Version", "923")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result := &JourneySearchResult{}
	json.NewDecoder(resp.Body).Decode(result)
	return *result
}

type NiceJourney struct {
	departureTime    time.Time
	departureStation string
	arrivalTime      time.Time
	arrivalStation   string
	price            float32
	numberOfLegs     int
	ticketType       string
	ticketClass      string
}

func makeJourneyMap(journeys []Journey) map[int]Journey {
	var m map[int]Journey
	m = make(map[int]Journey)
	for _, journey := range journeys {
		m[journey.ID] = journey
	}
	return m
}

func makeTicketMap(ticket []Ticket) map[int]Ticket {
	var m map[int]Ticket
	m = make(map[int]Ticket)
	for _, ticket := range ticket {
		m[ticket.ID] = ticket
	}
	return m
}

func resultToNiceJourneys(searchResult JourneySearchResult) []NiceJourney {

	journeyMap := makeJourneyMap(searchResult.Journeys)
	ticketMap := makeTicketMap(searchResult.Tickets)

	niceJourneys := []NiceJourney{}

	for _, journeyTicket := range searchResult.JourneyTickets {
		journey := journeyMap[journeyTicket.JourneyID]
		ticket := ticketMap[journeyTicket.TicketID]
		price := float32(ticket.TotalFare) / 100
		niceJourney := NiceJourney{
			journey.DepartureDateTime,
			journey.Origin,
			journey.ArrivalDateTime,
			journey.Destination,
			price,
			len(journey.Legs),
			ticket.TicketType,
			ticket.TicketClass,
		}
		niceJourneys = append(niceJourneys, niceJourney)
	}

	return niceJourneys

}
