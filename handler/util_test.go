package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RailwayTickets/backend-go/entity"
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

func TestSetTokenInfoHeaders(t *testing.T) {
	ti := &entity.TokenInfo{UserID: 666, CompanyURL: "unet", Role: "user", CompanyID: 0}
	expectedHeaders := http.Header(map[string][]string{
		http.CanonicalHeaderKey("df-userID"):     {"666"},
		http.CanonicalHeaderKey("df-companyURL"): {"unet"},
		http.CanonicalHeaderKey("df-role"):       {"user"},
		http.CanonicalHeaderKey("df-companyID"):  {"0"},
	})
	r := httptest.NewRequest("", "localhost:8080", nil)
	setTokenInfoHeaders(r.Header, ti)
	assert.Equal(t, expectedHeaders, r.Header)
}
