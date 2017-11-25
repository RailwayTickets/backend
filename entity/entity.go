package entity

import (
	"errors"
	"time"
)

// TokenInfo holds information about user and his token
type TokenInfo struct {
	Token      string    `bson:"token" json:"token"`
	Expires    time.Time `bson:"expires" json:"expires"`
	UserID     int64     `bson:"userID" json:"userID" df:"df-userID"`
	CompanyURL string    `bson:"companyURL" json:"companyURL" df:"df-companyURL"`
	CompanyID  int64     `bson:"companyID" json:"companyID" df:"df-companyID"`
	Role       string    `bson:"role" json:"role" df:"df-role"`
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
	if ti.CompanyURL == "" {
		return errors.New("companyURL cannot be empty")
	}
	if ti.CompanyID == 0 {
		return errors.New("companyID cannot be empty")
	}
	if ti.Role == "" {
		return errors.New("role cannot be empty")
	}
	if ti.UserID == 0 {
		return errors.New("userID cannot be empty")
	}
	return nil
}

// ValidateTokenOnly validates only token
func (ti *TokenInfo) ValidateTokenOnly() error {
	return ti.checkToken()
}
