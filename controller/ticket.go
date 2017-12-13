package controller

import (
	"errors"

	"time"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
	"gopkg.in/mgo.v2"
)

var (
	AlreadyTaken  = errors.New("ticket is already taken")
	NotYourTicket = errors.New("ticket is not yours")
	AlreadyIssued = errors.New("return of issued ticket is forbidden")
)

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

func Return(login, ticketID string) error {
	ticket, err := mongo.Tickets.ByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.Owner != login {
		return NotYourTicket
	}
	if ticket.Departure.Before(time.Now()) {
		return AlreadyIssued
	}
	return mongo.Tickets.Return(login, ticketID)
}
func ValidReturn(login string) (entity.TicketSearchResult, error) {
	tickets, err := mongo.Tickets.ValidReturnForUser(login)
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
