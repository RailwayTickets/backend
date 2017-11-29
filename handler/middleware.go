package handler

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/RailwayTickets/backend-go/controller"
)

// Middleware is a wrapper type for http handlers
type Middleware func(http.Handler) http.Handler

const (
	tokenHeaderName = "token"
)

// AllowCORS is a wrapper that adds 'Access-Control-Allow' headers in response
func AllowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, df-token")
		handler.ServeHTTP(w, r)
	})
}

// FilterOptions is a wrapper that filters any OPTIONS request
func FilterOptions(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			handler.ServeHTTP(w, r)
		}
	})
}

// RequiredPost is a wrapper that filters out any non post requests
func RequiredPost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// SetContentTypeJSON is a wrapper that set response content type to json
func SetContentTypeJSON(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

// CheckAndUpdateToken is a wrapper that checks and updates token
func CheckAndUpdateToken(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(tokenHeaderName)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ti, err := controller.Token.GetInfo(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		controller.Token.UpdateTTL(ti)
		handler.ServeHTTP(w, r)
	})
}

// CheckPathOrSetDefault is a wrapper that helps to redirect non existing paths of file server
func CheckPathOrSetDefault(baseDir string) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(path.Join(baseDir, r.URL.Path)); os.IsNotExist(err) {
				r.URL.Path = "/"
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// NoListFiles is a wrapper that forbids directory listing on file server
func NoListFiles(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// NoListFilesWithException is a wrapper that forbids directory listing on file server except root directory
func NoListFilesWithException(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// Chain is a helper function that chains middleWares and wraps handler
func Chain(handler http.Handler, middleware ...Middleware) http.Handler {
	h := handler
	for _, m := range middleware {
		h = m(h)
	}
	return h
}
