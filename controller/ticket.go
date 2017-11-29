package controller

import (
	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
)

func Search(query *entity.TicketSearchParams) (entity.TicketSearchResult, error) {
	tickets, err := mongo.Tickets.Search(query)
	return entity.TicketSearchResult{tickets}, err
}

func GetDirections() (entity.AvailableLocations, error) {
	locations, err := mongo.Tickets.AllDirections()
	return entity.AvailableLocations{locations}, err
}

func GetDepartures() (entity.AvailableLocations, error) {
	locations, err := mongo.Tickets.AllDepartures()
	return entity.AvailableLocations{locations}, err
}
