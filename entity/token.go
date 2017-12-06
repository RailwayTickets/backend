package entity

import (
	"errors"
	"time"
)

// TokenInfo holds information about user and his token
type TokenInfo struct {
	Token   string    `bson:"token" json:"token"`
	Login   string    `bson:"login" json:"-"`
	Expires time.Time `bson:"expires" json:"expires"`
}

// Validate validates all fields
func (ti *TokenInfo) Validate() error {
	if ti.Token == "" {
		return errors.New("token cannot be empty")
	}
	return nil
}
