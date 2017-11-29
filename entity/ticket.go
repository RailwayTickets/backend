package entity

import "time"

type (
	Ticket struct {
		From      string    `json:"from"`
		To        string    `json:"to"`
		Departure time.Time `json:"departure"`
		Carriage  int       `json:"carriage"`
		Seat      int       `json:"seat"`
	}

	TicketSearchParams struct {
		From string    `json:"from"`
		To   string    `json:"to"`
		Date time.Time `json:"date"`
	}

	TicketSearchResult struct {
		Tickets []Ticket `json:"tickets"`
	}

	AvailableLocations struct {
		Locations []string `json:"locations"`
	}
)
