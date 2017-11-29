package entity

import (
	"strings"
	"time"
)

type (
	DepartureDay time.Time

	Ticket struct {
		From      string    `json:"from"`
		To        string    `json:"to"`
		Departure time.Time `json:"departure"`
		Carriage  int       `json:"carriage"`
		Seat      int       `json:"seat"`
	}

	TicketSearchParams struct {
		From string       `json:"from"`
		To   string       `json:"to"`
		Date DepartureDay `json:"date"`
	}

	TicketSearchResult struct {
		Tickets []Ticket `json:"tickets"`
	}

	AvailableLocations struct {
		Locations []string `json:"locations"`
	}
)

func (d *DepartureDay) UnmarshalJSON(data []byte) error {
	if string(data) == `"null"` {
		return nil
	}
	stringDate := strings.Trim(string(data), `"`)
	parsed, err := time.Parse("2006-01-02", stringDate)
	*d = DepartureDay(parsed)
	return err
}
