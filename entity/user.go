package entity

import "errors"

type User struct {
	Login    string `bson:"login"`
	Password string `bson:"password"`
}

type (
	RegistrationInfo struct {
		LoginInfo
	}

	LoginInfo struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	LoginCredentials struct {
		TokenInfo
	}
)

func (r *RegistrationInfo) Validate() error {
	return r.LoginInfo.Validate()
}

func (r *LoginInfo) Validate() error {
	if r.Login == "" {
		return errors.New("login cannot be empty")
	}
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}
