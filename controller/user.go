package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/RailwayTickets/backend-go/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(info *entity.RegistrationInfo) (*entity.TokenInfo, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	user := &entity.User{
		LoginInfo: entity.LoginInfo{
			Login:    info.Login,
			Password: string(passwordHash),
		},
	}
	if err := mongo.User.Add(user); err != nil {
		return nil, err
	}
	token, _ := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	creds := &entity.TokenInfo{
		Token: string(token),
		Login: info.Login,
	}
	return creds, Token.Insert(creds)
}

func Login(info *entity.LoginInfo) (*entity.TokenInfo, error) {
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
	creds := &entity.TokenInfo{
		Token: string(token),
		Login: info.Login,
	}
	return creds, Token.Insert(creds)
}

func UpdateProfile(login string, profile *entity.User) error {
	return mongo.User.Update(login, profile)
}

func GetProfile(login string) (*entity.User, error) {
	return mongo.User.ByLogin(login)
}

func GetMyTickets(login string) (entity.TicketSearchResult, error) {
	tickets, err := mongo.Tickets.ForUser(login)
	return entity.TicketSearchResult{tickets}, err
}
