package handler

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/RailwayTickets/backend-go/controller"
	"github.com/RailwayTickets/backend-go/entity"
	"github.com/stretchr/testify/assert"
)

const (
	testDirPath      = "./_df_test_dir"
	testInnerDirName = "_df_inner_test_dir/"
	testFileName     = "_df_test"
)

type testCaseCORS struct {
	name   string
	method string
}

type testCaseListFiles struct {
	name           string
	path           string
	expectedStatus int
}

func createTestDir() error {
	err := os.Mkdir(testDirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(testDirPath, testInnerDirName), os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	for _, ext := range []string{"", ".txt", ".txt.gz", ".js", ".js.gz",
		".html", ".html.gz", ".css", ".css.gz"} {
		_, err = os.Create(path.Join(testDirPath, testFileName+ext))
		if err != nil {
			return err
		}
	}
	return nil
}

func removeTestDir() error {
	return os.RemoveAll(testDirPath)
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func TestFilterOptions(t *testing.T) {
	tt := []struct {
		name         string
		method       string
		expectedBody string
	}{
		{"method GET", http.MethodGet, "ok"},
		{"method OPTIONS", http.MethodOptions, ""},
		{"method POST", http.MethodPost, "ok"},
	}

	handler := Chain(http.HandlerFunc(okHandler), FilterOptions)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "http://localhost:8080/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestAllowCORS(t *testing.T) {
	tt := []struct {
		name           string
		subdomain      string
		origin         string
		expectedOrigin string
	}{
		{
			name:           "development",
			subdomain:      "dev.",
			origin:         "http://localhost:8080",
			expectedOrigin: "*",
		},
		{
			name:           "staging",
			subdomain:      "staging.",
			origin:         "http://staging.railways.io",
			expectedOrigin: "http://staging.railways.io",
		},
		{
			name:           "not df staging",
			subdomain:      "staging.",
			origin:         "http://staging.sashalala.io",
			expectedOrigin: "",
		},
		{
			name:           "production",
			subdomain:      "",
			origin:         "http://railways.io",
			expectedOrigin: "http://railways.io",
		},
		{
			name:           "not df production",
			subdomain:      "",
			origin:         "http://sashalala.io",
			expectedOrigin: "",
		},
	}
	domain := "railways.io"
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := Chain(http.HandlerFunc(okHandler), AllowCORS(tc.subdomain, domain))
			req := httptest.NewRequest("", tc.origin, nil)
			req.Header.Set("Origin", tc.origin)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.expectedOrigin, res.Header.Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestNoCache(t *testing.T) {
	noCacheHeaders := map[string]string{
		"Cache-Control": "no-cache",
	}

	tt := []testCaseCORS{
		{"method GET", http.MethodGet},
		{"method OPTIONS", http.MethodOptions},
		{"method POST", http.MethodPost},
	}

	handler := Chain(http.HandlerFunc(okHandler), NoCache)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "http://localhost:8080/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			for k, v := range noCacheHeaders {
				assert.Equal(t, v, res.Header.Get(k))
			}
		})
	}
}

func TestCachePublic(t *testing.T) {
	noCacheHeaders := map[string]string{
		"Cache-Control": "public, max-age=31536000",
	}

	tt := []testCaseCORS{
		{"method GET", http.MethodGet},
		{"method OPTIONS", http.MethodOptions},
		{"method POST", http.MethodPost},
	}

	handler := Chain(http.HandlerFunc(okHandler), CachePublic)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "http://localhost:8080/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			for k, v := range noCacheHeaders {
				assert.Equal(t, v, res.Header.Get(k))
			}
		})
	}
}

func TestCheckPathOrSetDefault(t *testing.T) {
	assert.NoError(t, createTestDir(), "failed to create test directory")
	tt := []struct {
		name         string
		path         string
		expectedPath string
	}{
		{"file exists", testFileName, "/" + testFileName},
		{"file does not exist", "aba1aba", "/"},
	}

	handler := Chain(http.HandlerFunc(okHandler), CheckPathOrSetDefault(testDirPath))
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/"+tc.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.expectedPath, req.URL.Path)
		})
	}
	assert.NoError(t, removeTestDir(), "failed to remove test directory")
}

func TestReturnGzipIfPossible(t *testing.T) {
	assert.NoError(t, createTestDir(), "failed to create test directory")
	handler := Chain(http.HandlerFunc(okHandler), ReturnGzipIfPossible(testDirPath))
	tt := []struct {
		ext          string
		gzipExt      string
		expectedType string
	}{
		{".txt", ".txt.gz", "text/plain; charset=utf-8"},
		{".js", ".js.gz", "application/javascript"},
		{".html", ".html.gz", "text/html"},
		{".css", ".css.gz", "text/css"},
	}

	for _, tc := range tt {
		t.Run(tc.ext, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/"+testFileName+tc.ext, nil)
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "/"+testFileName+tc.gzipExt, req.URL.Path)
			assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
			assert.Equal(t, tc.expectedType, res.Header.Get("Content-Type"))
		})
	}
	assert.NoError(t, removeTestDir(), "failed to remove test directory")
}

func TestNoListFiles(t *testing.T) {
	assert.NoError(t, createTestDir(), "failed to create test directory")
	tt := []testCaseListFiles{
		{"file request", testFileName, http.StatusOK},
		{"non root directory request", testInnerDirName, http.StatusForbidden},
		{"root directory request", "/", http.StatusForbidden},
	}

	handler := Chain(http.HandlerFunc(okHandler), NoListFiles)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/"+tc.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
	assert.NoError(t, removeTestDir(), "failed to remove test directory")
}

func TestNoListFilesWithException(t *testing.T) {
	assert.NoError(t, createTestDir(), "failed to create test directory")
	tt := []testCaseListFiles{
		{"file request", testFileName, http.StatusOK},
		{"non root directory request", testInnerDirName, http.StatusForbidden},
		{"root directory request", "", http.StatusOK},
	}

	handler := Chain(http.HandlerFunc(okHandler), NoListFilesWithException)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/"+tc.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
	assert.NoError(t, removeTestDir(), "failed to remove test directory")
}

func TestCheckAndUpdateToken(t *testing.T) {
	tt := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"empty token", "", http.StatusUnauthorized},
		{"non existing token", "aba1aba", http.StatusUnauthorized},
		{"existing token", "test", http.StatusOK},
	}

	controller.Token.Insert(&entity.TokenInfo{Token: "test"})
	handler := Chain(http.HandlerFunc(okHandler), CheckAndUpdateToken)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
			req.Header.Set("token", tc.token)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestRedirectToHTTPS(t *testing.T) {
	tt := []struct {
		name           string
		url            string
		protoHeader    string
		expectedStatus int
	}{
		{"no proto header + http", "http://www.railways.io", "", http.StatusMovedPermanently},
		{"no proto header + https", "https://www.railways.io", "", http.StatusMovedPermanently},
		{"proto header + http", "http://www.railways.io", "https", http.StatusOK},
		{"proto header + https", "https://www.railways.io", "https", http.StatusOK},
	}

	handler := Chain(http.HandlerFunc(okHandler), RedirectToHTTPS)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("", tc.url, nil)
			req.Header.Set("X-Forwarded-Proto", tc.protoHeader)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestRedirectWWW(t *testing.T) {
	tt := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"www + http", "http://www.railways.io", http.StatusMovedPermanently},
		{"www + https", "https://www.railways.io", http.StatusMovedPermanently},
		{"http", "http://railways.io", http.StatusOK},
		{"https", "https://railways.io", http.StatusOK},
	}

	handler := Chain(http.HandlerFunc(okHandler), RedirectWWW)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("", tc.url, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestCheckAndAddETagHeader(t *testing.T) {
	assert.NoError(t, createTestDir(), "failed to create test directory")
	fi, _ := os.Stat(fmt.Sprintf("%s/%s", testDirPath, testFileName))
	lastModified := []byte(fi.ModTime().String())
	correctETag := fmt.Sprintf("%x", md5.Sum(lastModified))
	tt := []struct {
		name           string
		file           string
		sentETag       string
		expectedETag   string
		expectedStatus int
	}{
		{
			name:           "first request",
			file:           testFileName,
			sentETag:       "",
			expectedETag:   correctETag,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty filename",
			file:           "",
			sentETag:       "",
			expectedETag:   "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "correct etag",
			file:           testFileName,
			sentETag:       correctETag,
			expectedETag:   "",
			expectedStatus: http.StatusNotModified,
		},
		{
			name:           "wrong etag",
			file:           testFileName,
			sentETag:       correctETag + "sdas",
			expectedETag:   correctETag,
			expectedStatus: http.StatusOK,
		},
	}

	handler := Chain(http.HandlerFunc(okHandler), CheckAndAddETagHeader(testDirPath))
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost:8080/"+tc.file, nil)
			req.Header.Set("If-None-Match", tc.sentETag)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			actualETag := res.Header.Get("ETag")
			assert.Equal(t, tc.expectedETag, actualETag)
		})
	}
	assert.NoError(t, removeTestDir(), "failed to remove test directory")
}
