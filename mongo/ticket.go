package mongo

import (
	"time"

	"github.com/RailwayTickets/backend-go/entity"
	"gopkg.in/mgo.v2/bson"
)

type ticket struct{}

func (ticket) Search(params *entity.TicketSearchParams) ([]entity.Ticket, error) {
	var found []entity.Ticket
	query := bson.M{}
	if params.From != "" {
		query["from"] = params.From
	}
	if params.To != "" {
		query["to"] = params.To
	}
	departure := time.Time(params.Date)
	if !departure.IsZero() {
		query["departure"] = bson.M{
			"$gte": departure,
			"$lte": departure.Add(time.Hour * 24),
		}
	}
	err := tickets.Find(query).All(&found)
	return found, err
}

func (ticket) AllDirections() ([]string, error) {
	var directions []string
	err := tickets.Find(nil).Distinct("to", &directions)
	return directions, err
}

func (ticket) AllDepartures() ([]string, error) {
	var departures []string
	err := tickets.Find(nil).Distinct("from", &departures)
	return departures, err
}
