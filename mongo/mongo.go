package mongo

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

const (
	mongoURLEnv = "MONGODB_URI"
	dbNameEnv   = "DB_NAME"
)

var (
	User  = user{}
	users *mgo.Collection

	Tickets = ticket{}
	tickets *mgo.Collection

	Token  = token{}
	tokens *mgo.Collection
)

func connectToMongo(connectionURL, dbName string) (*mgo.Database, error) {
	session, err := mgo.Dial(connectionURL)
	if err != nil {
		return nil, err
	}
	return session.DB(dbName), session.Ping()
}

func init() {
	mongoURL, ok := os.LookupEnv(mongoURLEnv)
	if !ok {
		log.Fatalf("please, specify %s variable", mongoURLEnv)
	}
	dbName, ok := os.LookupEnv(dbNameEnv)
	if !ok {
		log.Fatalf("please, specify %s variable", dbNameEnv)
	}
	db, err := connectToMongo(mongoURL, dbName)
	if err != nil {
		log.Fatal(err)
	}
	users = db.C("users")
	tokens = db.C("tokens")
	tickets = db.C("tickets")
	indexes := []struct {
		index      mgo.Index
		collection *mgo.Collection
	}{
		{
			index: mgo.Index{
				Name:   "token_index",
				Key:    []string{"token"},
				Unique: true,
			},
			collection: tokens,
		},
		{
			index: mgo.Index{
				Name:   "user_login_index",
				Key:    []string{"login"},
				Unique: true,
			},
			collection: users,
		},
	}
	for i := range indexes {
		col := indexes[i].collection
		err := col.EnsureIndex(indexes[i].index)
		if err != nil {
			log.Fatal(err)
		}
	}

}
