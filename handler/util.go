package handler

import (
	"net/http"
	"strings"
)

func formHTTPSRedirectURL(r *http.Request) string {
	r.URL.Scheme = "https"
	r.URL.Host = r.Host
	return r.URL.String()
}

func formWWWRedirectURL(r *http.Request) string {
	r.Host = strings.TrimPrefix(r.Host, "www.")
	r.URL.Host = r.Host
	return r.URL.String()
}
