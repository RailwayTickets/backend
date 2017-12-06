package entity

import "errors"

type (
	User struct {
		LoginInfo `bson:",inline" json:"-"`
		FirstName string `bson:"firstName" json:"first_name"`
		LastName  string `bson:"lastName" json:"last_name"`
		Phone     string `bson:"phone" json:"phone"`
		Email     string `bson:"email" json:"email"`
		DocType   string `bson:"doc_type" json:"doc_type"`
		DocNumber string `bson:"doc_number" json:"doc_number"`
	}

	RegistrationInfo struct {
		LoginInfo
	}

	LoginInfo struct {
		Login    string `bson:"login" json:"login"`
		Password string `bson:"password" json:"password"`
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
