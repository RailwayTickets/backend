package handler

import (
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
