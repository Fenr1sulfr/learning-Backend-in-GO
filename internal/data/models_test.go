package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sulfur.test.net/internal/data"
)

type testConfig struct {
	dsn         string
	maxOpenCons int
	maxIdleCons int
	maxIdleTime string
}

var ConfigObject = testConfig{
	dsn:         `postgres://greenlight:pa55word@localhost:5433/greenlight?sslmode=disable`,
	maxOpenCons: 25,
	maxIdleCons: 25,
	maxIdleTime: "15m",
}

func openDB(cfg testConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.maxOpenCons)
	db.SetMaxIdleConns(cfg.maxIdleCons)
	duration, err := time.ParseDuration(cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestInsert(t *testing.T) {
	dataBase, err := openDB(ConfigObject)

	if err != nil {
		t.Errorf("Error connecting database %s", err)
	}
	movieModel := data.MovieModel{DB: dataBase}

	testMovie := &data.Movie{
		Title:   "Test Movie",
		Year:    2023,
		Runtime: 120,
		Genres:  []string{"Comedy"},
	}

	err = movieModel.Insert(testMovie)

	require.NoError(t, err)
	testMovie = &data.Movie{
		Title:   "Test Movie",
		Year:    1200,
		Runtime: 120,
		Genres:  []string{"Comedy"},
	}
	err = movieModel.Insert(testMovie)
	require.Error(t, err)

}

func TestGetAll(t *testing.T) {
	dataBase, err := openDB(ConfigObject)

	if err != nil {
		t.Errorf("Error connecting database %s", err)
	}
	movieModel := data.MovieModel{DB: dataBase}

	_, err = movieModel.GetAll("", nil, data.Filters{})

	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	// Create a MovieModel instance with a real database connection (replace with your actual database connection)
	dataBase, err := openDB(ConfigObject)
	movieModel := data.MovieModel{DB: dataBase}
	// Define a valid movie ID (assuming some movies exist in the database)
	testID := int64(1)

	// Call the Get function
	movie, err := movieModel.Get(testID)

	// Here, we cannot verify the exact SQL query or returned data without mocking.
	// However, we can check for errors and handle different scenarios.

	if err != nil {
		// Handle the error based on the expected behavior (e.g., ErrNoRecordFound for non-existent ID)
		require.Equal(t, sql.ErrNoRows, err) // Replace with your expected error
		return
	}

	// If no error, assert that a movie is returned
	require.NotNil(t, movie)

	// Optionally, you can perform additional assertions on the returned movie properties.
}
