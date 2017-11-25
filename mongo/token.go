package mongo

import (
	"github.com/RailwayTickets/backend-go/entity"
	"gopkg.in/mgo.v2"
)

type token struct{}

// Add inserts token into database
func (token) Add(token *entity.TokenInfo) error {
	return tokens.Insert(token)
}

// Update updates token expiration time in db
func (token) Update(token *entity.TokenInfo) error {
	selector := map[string]interface{}{
		"token": token.Token,
	}
	document := map[string]interface{}{
		"$set": map[string]interface{}{
			"expires": token.Expires,
		},
	}
	return tokens.Update(selector, document)
}

// Remove deletes given token from database
func (token) Remove(token string) error {
	selector := map[string]interface{}{
		"token": token,
	}
	return tokens.Remove(selector)
}

// GetInfo returns all information related to given token
func (token) Get(token string) (*entity.TokenInfo, error) {
	query := map[string]interface{}{
		"token": token,
	}
	ti := new(entity.TokenInfo)
	err := tokens.Find(query).One(ti)
	return ti, err
}

// Index ensures that index is present in database
func (token) Index(index mgo.Index) error {
	return tokens.EnsureIndex(index)
}
