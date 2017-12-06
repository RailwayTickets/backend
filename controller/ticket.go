package controller

import (
	"errors"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
	"gopkg.in/mgo.v2"
)

var AlreadyTaken = errors.New("ticket is already taken")

func Search(query *entity.TicketSearchParams) (entity.TicketSearchResult, error) {
	tickets, err := mongo.Tickets.Search(query)
	return entity.TicketSearchResult{tickets}, err
}

func Buy(login, ticketID string) error {
	err := mongo.Tickets.Buy(login, ticketID)
	if err == mgo.ErrNotFound {
		return AlreadyTaken
	}
	return err
}

func GetDirections() (entity.AvailableLocations, error) {
	locations, err := mongo.Tickets.AllDirections()
	return entity.AvailableLocations{locations}, err
}

func GetDepartures() (entity.AvailableLocations, error) {
	locations, err := mongo.Tickets.AllDepartures()
	return entity.AvailableLocations{locations}, err
}
