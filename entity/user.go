package entity

import "errors"

type User struct {
	ID       string `bson:"id"`
	Login    string `bson:"login"`
	Password string `bson:"password"`
}

type RegistrationInfo struct {
	Login                string `json:"login"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"confirmation"`
}

type LoginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginCredentials struct {
	TokenInfo
}

func (r *RegistrationInfo) Validate() error {
	if r.Login == "" {
		return errors.New("login cannot be empty")
	}
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	if r.PasswordConfirmation == "" {
		return errors.New("password confirmation cannot be empty")
	}

	if r.Password != r.PasswordConfirmation {
		return errors.New("password confirmation doesn't match password")
	}
	return nil
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
