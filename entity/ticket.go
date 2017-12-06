package entity

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	DepartureDay time.Time

	Ticket struct {
		ID        bson.ObjectId `bson:"_id" json:"id"`
		From      string        `bson:"from" json:"from"`
		To        string        `bson:"to" json:"to"`
		Departure time.Time     `bson:"departure" json:"departure"`
		Carriage  int           `bson:"carriage" json:"carriage"`
		Seat      int           `bson:"seat" json:"seat"`
		Type      string        `bson:"type" json:"type"`
		Owner     string        `bson:"owner" json:"-"`
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
