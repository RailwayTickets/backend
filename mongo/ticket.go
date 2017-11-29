package mongo

import (
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
	if !params.Date.IsZero() {
		query["departure"] = bson.M{
			"$gte": params.Date,
		}
	}
	err := tickets.Find(query).All(&found)
	return found, err
}
