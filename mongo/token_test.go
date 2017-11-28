package mongo

import (
	"testing"
	"time"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
)

func TestTokenCRUD(t *testing.T) {
	ti := entity.TokenInfo{
		Token:   "lala",
		Expires: time.Now(),
		Login:   "asas",
	}
	err := Token.Add(&ti)
	assert.NoError(t, err, "unexpected error on add")

	actual, err := Token.Get(ti.Token)
	assert.NoError(t, err, "unexpected error on get")
	assert.Equal(t, ti.Token, actual.Token, "read token does not match inserted")

	ti.Expires = time.Now()
	err = Token.Update(&ti)
	assert.NoError(t, err, "unexpected error on update")

	actual, err = Token.Get(ti.Token)
	assert.NoError(t, err, "unexpected error on get")
	assert.NotEqual(t, ti.Expires, actual.Expires, "read expires match inserted")

	err = Token.Remove(ti.Token)
	assert.NoError(t, err, "unexpected error on remove")
}

func TestTokenIndex(t *testing.T) {
	index := mgo.Index{
		Name: "_test_index",
		Key:  []string{"some_key"},
	}
	assert.NoError(t, Token.Index(index))
	indexes, err := tokens.Indexes()
	assert.NoError(t, err)
	for i := range indexes {
		if indexes[i].Name == index.Name {
			assert.Equal(t, index, indexes[i])
			return
		}
	}
	assert.Fail(t, "index not found")
}
