package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"sulfur.test.net/internal/data"
)

func configureAppForTest() *application {

	cfg := config{
		port: 4000,
		env:  "dev",
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}{
			dsn:          `postgres://greenlight:pa55word@localhost:5433/greenlight?sslmode=disable`,
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  "15m",
		},
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection established")
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	return app

}

var testapp = configureAppForTest()

func TestMovieHandler(t *testing.T) {
	testapp := configureAppForTest()

	t.Run("IsListWorking", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/v1/movies", nil)
		w := httptest.NewRecorder()
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		testapp.listMovieHandler(w, req)

	})
	t.Run("isCreateMovieWorks", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{
			"title":"Dune 2",
			"runtime":"102 mins",
			"year":2000,
			"genres":["aboba"]
		}`))
		req := httptest.NewRequest(http.MethodPost, "/v1/movies", body)
		w := httptest.NewRecorder()
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		testapp.createMovieHandler(w, req)
	})
}

func TestHealthcheckHandler(t *testing.T) {

	expectedEnv := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": testapp.config.env,
			"version":     version,
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	testapp.healthcheckHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var actualEnv envelope
	err := json.NewDecoder(w.Body).Decode(&actualEnv)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
		return
	}
	if diff := cmp.Diff(expectedEnv, actualEnv); diff != "" {
		t.Errorf("Response body mismatch (-expected +actual):\n%s", diff)
	}
}
