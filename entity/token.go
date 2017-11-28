package entity

import (
	"errors"
	"time"
)

// TokenInfo holds information about user and his token
type TokenInfo struct {
	Token   string    `bson:"token" json:"token"`
	Expires time.Time `bson:"expires" json:"expires"`
	Login   string    `bson:"login" json:"login"`
}

func (ti *TokenInfo) checkToken() error {
	if ti.Token == "" {
		return errors.New("token cannot be empty")
	}
	return nil
}

// Validate validates all fields
func (ti *TokenInfo) Validate() error {
	if err := ti.checkToken(); err != nil {
		return err
	}
	if ti.Login == "" {
		return errors.New("login cannot be empty")
	}
	return nil
}

// ValidateTokenOnly validates only token
func (ti *TokenInfo) ValidateTokenOnly() error {
	return ti.checkToken()
}
