package mongo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToMongo(t *testing.T) {
	tt := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "wrong connection URL",
			url:         "_something_wrong_",
			expectError: true,
		},
		{
			name:        "correct connection URL",
			url:         os.Getenv(mongoURLEnv),
			expectError: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := connectToMongo(tc.url, os.Getenv(dbNameEnv))
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
