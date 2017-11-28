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

func TestInsertAndRemoveToken(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token: "TestInsertAndRemoveToken",
		Login: "asas",
	}
	mc, _ := newTokenController()
	assert.NoError(t, mc.Insert(tokenInfo))
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}

func TestGetInfoByToken(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token: "TestGetInfoByToken",
		Login: "asas",
	}
	mc, _ := newTokenController()
	mc.Insert(tokenInfo)
	ti, err := mc.GetInfo(tokenInfo.Token)
	assert.NoError(t, err)
	assert.Equal(t, ti.Token, tokenInfo.Token)
	assert.Equal(t, ti.Login, tokenInfo.Login)
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}

func TestUpdateTokenTTL(t *testing.T) {
	tokenInfo := &entity.TokenInfo{
		Token: "TestUpdateTokenTTL",
		Login: "asas",
	}
	mc, _ := newTokenController()
	assert.NoError(t, mc.Insert(tokenInfo))
	assert.NoError(t, mc.UpdateTTL(tokenInfo))
	assert.NoError(t, mc.Remove(tokenInfo.Token))
}
