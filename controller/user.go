package controller

import (
	"time"

	"crypto/rand"

	"fmt"

	"errors"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(info *entity.RegistrationInfo) (*entity.LoginCredentials, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	user := &entity.User{
		Login:    info.Login,
		Password: string(passwordHash),
	}
	id := make([]byte, 25)
	rand.Read(id)
	user.ID = fmt.Sprintf("%x", id)
	if err := mongo.User.Add(user); err != nil {
		return nil, err
	}
	token, _ := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	creds := &entity.LoginCredentials{
		TokenInfo: entity.TokenInfo{
			Token:  string(token),
			UserID: user.ID,
		},
	}
	return creds, Token.Insert(&creds.TokenInfo)
}

func Login(info *entity.LoginInfo) (*entity.LoginCredentials, error) {
	user, err := mongo.User.ByLogin(info.Login)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user: %s", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("login or password invalid")
	}
	if err != nil {
		return nil, fmt.Errorf("could not check password: %s", err)
	}
	token, _ := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	creds := &entity.LoginCredentials{
		TokenInfo: entity.TokenInfo{
			Token:  string(token),
			UserID: user.ID,
		},
	}
	return creds, Token.Insert(&creds.TokenInfo)
}
