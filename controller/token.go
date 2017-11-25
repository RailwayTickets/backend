package controller

import (
	"log"
	"time"

	"os"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
	"gopkg.in/mgo.v2"
)

const (
	tokenTTLEnv = "TOKEN_TTL"
	defaultTTL  = "30m"
)

// Token controls all operations that involve user token
var Token *token

type token struct {
	ttl time.Duration
}

func newTokenController() (*token, error) {
	ttlString := os.Getenv(tokenTTLEnv)
	ttl, err := time.ParseDuration(ttlString)
	if err != nil {
		log.Printf("unable to parse %s: %s. using default of %s", tokenTTLEnv, err.Error(), defaultTTL)
		ttl, _ = time.ParseDuration(defaultTTL)
	}

	ttlExpireIndex := mgo.Index{
		Name:        "ttl_index",
		Key:         []string{"expires"},
		ExpireAfter: time.Second,
	}
	err = mongo.Token.Index(ttlExpireIndex)
	return &token{ttl}, err
}

// Insert sets token expiration time and inserts it into database
func (tc *token) Insert(token *entity.TokenInfo) error {
	token.Expires = time.Now().Add(tc.ttl)
	return mongo.Token.Add(token)
}

// UpdateTTL refreshes token expiration time
func (tc *token) UpdateTTL(token *entity.TokenInfo) error {
	token.Expires = time.Now().Add(tc.ttl)
	return mongo.Token.Update(token)
}

// Remove deletes given token from database
func (tc *token) Remove(token string) error {
	return mongo.Token.Remove(token)
}

// GetInfo returns all information related to given token
func (tc *token) GetInfo(token string) (*entity.TokenInfo, error) {
	return mongo.Token.Get(token)
}

func init() {
	var err error
	Token, err = newTokenController()
	if err != nil {
		log.Fatal(err)
	}
}
