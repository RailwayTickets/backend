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

func TestValidation(t *testing.T) {
	tt := []testCaseEntity{
		{"empty token", TokenInfo{}, true},
		{"all correct", TokenInfo{Token: "test"}, false},
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
