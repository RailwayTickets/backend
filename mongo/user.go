package mongo

import (
	"github.com/RailwayTickets/backend-go/entity"
	"gopkg.in/mgo.v2/bson"
)

type user struct{}

func (user) Add(user *entity.User) error {
	return users.Insert(user)
}

func (user) ByLogin(login string) (*entity.User, error) {
	user := new(entity.User)
	err := users.Find(bson.M{
		"login": login,
	}).One(user)
	return user, err
}
