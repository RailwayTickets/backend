package controller

import (
	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
)

func Search(query *entity.TicketSearchParams) (entity.TicketSearchResult, error) {
	tickets, err := mongo.Tickets.Search(query)
	return entity.TicketSearchResult{tickets}, err
}
