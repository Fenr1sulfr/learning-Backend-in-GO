package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"sulfur.test.net/internal/data"
)

func configureAppForTest() *application {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:pa55word@localhost:5433/greenlight?sslmode=disable", "PostgreSQL dsn")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()
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
	// Create a mock application instance
	testapp := configureAppForTest()

	// Define an expected response envelope
	expectedEnv := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": testapp.config.env,
			"version":     version, // Assuming version is defined elsewhere
		},
	}

	// Create a recorder to capture the response
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	// Call the handler with the mock app and recorder
	testapp.healthcheckHandler(w, r)

	// Assert on the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Assert on the response body (consider using a JSON marshalling library)
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
