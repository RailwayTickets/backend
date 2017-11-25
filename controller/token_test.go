package controller

import (
	"os"
	"testing"

	"github.com/RailwayTickets/backend-go/entity"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	realTTL := os.Getenv(tokenTTLEnv)
	os.Setenv(tokenTTLEnv, "")
	_, err := newTokenController()
	assert.NoError(t, err)
	os.Setenv(tokenTTLEnv, realTTL)
}

func TestInsertToken(t *testing.T) {
	mc, _ := newTokenController()
	err := mc.Insert(&entity.TokenInfo{
		Token:  "test",
		UserID: 0,
	})
	assert.NoError(t, err)
}

func TestRemoveToken(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token:  "test",
		UserID: 0,
	}
	mc, _ := newTokenController()
	assert.NoError(t, mc.Insert(tokenInfo))
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}

func TestGetInfoByToken(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token:      "testSuperToken",
		UserID:     1,
		CompanyURL: "qwe",
	}
	mc, _ := newTokenController()
	mc.Insert(tokenInfo)
	ti, err := mc.GetInfo(tokenInfo.Token)
	assert.NoError(t, err)
	assert.Equal(t, ti.Token, tokenInfo.Token)
	assert.Equal(t, ti.UserID, tokenInfo.UserID)
	assert.Equal(t, ti.CompanyURL, tokenInfo.CompanyURL)
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}

func TestUpdateTokenTTL(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token:  "test",
		UserID: 0,
	}
	mc, _ := newTokenController()
	assert.NoError(t, mc.Insert(tokenInfo))
	assert.NoError(t, mc.UpdateTTL(tokenInfo))
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}
