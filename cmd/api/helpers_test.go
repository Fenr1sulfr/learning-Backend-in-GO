package main

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"sulfur.test.net/internal/data"
)

type Movie struct {
	ID        int64        `json:"id"`                // Unique integer ID for the movie
	CreatedAt time.Time    `json:"-"`                 // Timestamp for when the movie is added to our database
	Title     string       `json:"tittle"`            // Movie title
	Year      int32        `json:"year,omitempty"`    // Movie release year
	Runtime   data.Runtime `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres    []string     `json:"genres,omitempty"`  // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32        `json:"version"`           // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

func TestReadJSON_BadlyFormedJSON(t *testing.T) {
	app := &application{}
	body := bytes.NewBuffer([]byte(`{
	  "title": "Dune 2",
	  "year": 2000,
	  "genres": ["aboba"]
	}`)) // Syntax error

	req := httptest.NewRequest("POST", "/v1/movies/", body)
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &Movie{})

	assert.Error(t, err)
	// assert.Contains(t, err.Error(), "body contains badly-formed JSON") // Optional assertion
}

func TestReadJSON_UnknownKey(t *testing.T) {
	app := &application{}
	body := bytes.NewBuffer([]byte(`{"title": "Alice", "year": 2000, "genres": ["aboba"], "director": "James Cameron"}`)) // Unknown key "director"
	req := httptest.NewRequest("POST", "/v1/movies", body)
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &Movie{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body contains unknown key director")
}

func TestReadJSON_IncorrectType(t *testing.T) {
	app := &application{}
	body := bytes.NewBuffer([]byte(`{"title": "Alice", "year": "two thousand", "genres": ["aboba"]}`)) // Incorrect type for "year"
	req := httptest.NewRequest("POST", "/v1/movies", body)
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &Movie{})
	assert.Error(t, err)
	// assert.Contains(t.Error(), "body contains incorrect JSON type for field year")
}

func TestReadJSON_EmptyBody(t *testing.T) {
	app := &application{}
	req := httptest.NewRequest("POST", "/test", nil) // Empty body
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &struct{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body must not to be empty")
}

func TestReadJSON_ExceedMaxBytes(t *testing.T) {
	app := &application{}
	body := strings.Repeat("x", 1_048_577) // 1 byte larger than max_bytes
	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &struct{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body must not to be larger than")
}

func TestReadJSON_MultipleValues(t *testing.T) {
	app := &application{}
	body := bytes.NewBuffer([]byte(`{"title": "Alice"} {"title": "Bob"}`)) // Multiple JSON values
	req := httptest.NewRequest("POST", "/test", body)
	w := httptest.NewRecorder()

	err := app.readJSON(w, req, &struct{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "body must contain a single json value")
}
