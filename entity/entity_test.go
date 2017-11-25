package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseEntity struct {
	name      string
	ti        TokenInfo
	waitError bool
}

func TestTokenValidation(t *testing.T) {
	tt := []testCaseEntity{
		{"empty token", TokenInfo{UserID: 666}, true},
		{"non empty token", TokenInfo{Token: "test", UserID: 666}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ti.ValidateTokenOnly()
			if tc.waitError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	tt := []testCaseEntity{
		{"empty token", TokenInfo{}, true},
		{"empty companyURL", TokenInfo{Token: "test"}, true},
		{"empty companyID", TokenInfo{Token: "test", CompanyURL: "unet"}, true},
		{"empty role", TokenInfo{Token: "test", CompanyURL: "unet", CompanyID: 1}, true},
		{"empty userID", TokenInfo{Token: "test", CompanyURL: "unet", CompanyID: 1, Role: "owner"}, true},
		{"all correct", TokenInfo{Token: "test", CompanyURL: "unet", CompanyID: 1, Role: "owner", UserID: 1}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ti.Validate()
			if tc.waitError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
