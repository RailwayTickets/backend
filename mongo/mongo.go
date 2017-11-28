package mongo

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

const (
	mongoURLEnv = "MONGO_URL"
	dbNameEnv   = "DB_NAME"
)

var (
	// User is used for interaction with users collection and method separation
	User  = user{}
	users *mgo.Collection

	// Token is used for interaction with token collection and method separation
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
	userLoginIndex := mgo.Index{
		Name:   "user_login_index",
		Key:    []string{"login"},
		Unique: true,
	}
	err = users.EnsureIndex(userLoginIndex)
	if err != nil {
		log.Fatal(err)
	}

	tokens = db.C("tokens")
	tokenIndex := mgo.Index{
		Name:   "token_index",
		Key:    []string{"token"},
		Unique: true,
	}
	err = tokens.EnsureIndex(tokenIndex)
	if err != nil {
		log.Fatal(err)
	}
}
