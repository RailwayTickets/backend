package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseRedirect struct {
	url            string
	expectedResult string
}

func TestFormHTTPSRedirectURL(t *testing.T) {
	tt := []testCaseRedirect{
		{"http://www.railways.io", "https://www.railways.io"},
		{"http://railways.io", "https://railways.io"},
		{"https://railways.io", "https://railways.io"},
		{"http://dev.railways.io", "https://dev.railways.io"},
		{"http://dev.railways.io/v1/login", "https://dev.railways.io/v1/login"},
	}

	for _, tc := range tt {
		t.Run(tc.url, func(t *testing.T) {
			r := httptest.NewRequest("", tc.url, nil)
			result := formHTTPSRedirectURL(r)
			assert.Equal(t, result, tc.expectedResult)
		})
	}
}

func TestFormWWWRedirectURL(t *testing.T) {
	tt := []testCaseRedirect{
		{"http://www.railways.io", "http://railways.io"},
		{"http://railways.io", "http://railways.io"},
		{"https://www.railways.io", "https://railways.io"},
		{"http://www.api.railways.io", "http://api.railways.io"},
		{"http://www.dev.railways.io/v1/login", "http://dev.railways.io/v1/login"},
	}

	for _, tc := range tt {
		t.Run(tc.url, func(t *testing.T) {
			r := httptest.NewRequest("", tc.url, nil)
			result := formWWWRedirectURL(r)
			assert.Equal(t, result, tc.expectedResult)
		})
	}
}
