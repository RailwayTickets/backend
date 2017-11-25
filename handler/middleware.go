package handler

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/RailwayTickets/backend-go/controller"
)

// Middleware is a wrapper type for http handlers
type Middleware func(http.Handler) http.Handler

const (
	protoHeaderName        = "X-Forwarded-Proto"
	eTagHeaderName         = "ETag"
	ifNoneMatchHeaderName  = "If-None-Match"
	cacheControlHeaderName = "Cache-Control"
	tokenHeaderName        = "auth-token"
)

// NoCache is a wrapper that adds 'Cache-Control: no-cache' header in response
func NoCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheControlHeaderName, "no-cache")
		handler.ServeHTTP(w, r)
	})
}

// CachePublic is a wrapper that adds 'Cache-Control: public, max-age=31536000' header in response
func CachePublic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheControlHeaderName, "public, max-age=31536000")
		handler.ServeHTTP(w, r)
	})
}

// AllowCORS is a wrapper that adds 'Access-Control-Allow' headers in response
func AllowCORS(subdomain, domain string) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var allowedOrigin string
			if subdomain == "dev." {
				allowedOrigin = "*"
			} else {
				allowedOrigin = r.Header.Get("Origin")
				if !strings.Contains(allowedOrigin, domain) {
					return
				}
			}
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, df-token")
			handler.ServeHTTP(w, r)
		})
	}
}

// FilterOptions is a wrapper that filters any OPTIONS request
func FilterOptions(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			handler.ServeHTTP(w, r)
		}
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
		setTokenInfoHeaders(r.Header, ti)
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

// CheckAndAddETagHeader check Etag header and sets it if necessary
func CheckAndAddETagHeader(baseDir string) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fi, _ := os.Stat(path.Join(baseDir, r.URL.Path))
			if fi.IsDir() {
				handler.ServeHTTP(w, r)
				return
			}
			lastModified := []byte(fi.ModTime().String())
			eTag := fmt.Sprintf("%x", md5.Sum(lastModified))
			incomingETag := r.Header.Get(ifNoneMatchHeaderName)
			if eTag == incomingETag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set(eTagHeaderName, eTag)
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

// ReturnGzipIfPossible is a wrapper that optimises file serving: if server has a gzipped
// version of requested file it will serve it
func ReturnGzipIfPossible(baseDir string) Middleware {
	contentTypeFromExtension := map[string]string{
		"js":   "application/javascript",
		"html": "text/html",
		"css":  "text/css",
	}
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gzipPath := path.Join(baseDir, r.URL.Path+".gz")
				ext := r.URL.Path[strings.LastIndex(r.URL.Path, ".")+1:]
				if _, err := os.Stat(gzipPath); err == nil {
					r.URL.Path += ".gz"
					w.Header().Set("Content-Encoding", "gzip")
					if contentType, ok := contentTypeFromExtension[ext]; ok {
						w.Header().Set("Content-Type", contentType)
					}
				}
			}
			handler.ServeHTTP(w, r)
		})
	}
}

// RedirectToHTTPS redirects any non secure traffic to https
func RedirectToHTTPS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protocol := r.Header.Get(protoHeaderName)
		if protocol == "https" {
			handler.ServeHTTP(w, r)
		} else {
			redirectURL := formHTTPSRedirectURL(r)
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		}
	})
}

// RedirectWWW redirects any requests started with www to same url without www
func RedirectWWW(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Host, "www.") {
			handler.ServeHTTP(w, r)
		} else {
			redirectURL := formWWWRedirectURL(r)
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		}
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
