package handler

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/RailwayTickets/backend-go/entity"
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

func setTokenInfoHeaders(headers http.Header, ti *entity.TokenInfo) {
	v := reflect.ValueOf(ti).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		headerName := fieldType.Tag.Get("df")
		if headerName != "" {
			headerValue := field.String()
			if fieldType.Type.Kind() != reflect.String {
				switch fieldType.Type.Kind() {
				case reflect.Int64:
					headerValue = strconv.FormatInt(field.Int(), 10)
				}
			}
			headers.Set(headerName, headerValue)
		}
	}
}
