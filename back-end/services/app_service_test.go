package services

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runway/config"
	"runway/logger"
	"testing"
)

// mockRoundTripper is a mock implementation of http.RoundTripper for testing.
type mockRoundTripper func(req *http.Request) (*http.Response, error)

func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

// setupTestService creates a new AppService with a mock HTTP client and a mock config.
func setupTestService(responseBody string, statusCode int, t *testing.T) (*AppService, *config.Config) {
	mockClient := &http.Client{
		Transport: mockRoundTripper(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}, nil
		}),
	}

	tempDir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	testConfig := &config.Config{
		AppsApiUrl:         "http://mock-api.com/apps",
		AppsStorageFile:    filepath.Join(tempDir, "apps.json"),
		ReviewsBaseUrl:     "http://mock-api.com/reviews",
		ReviewsStorageFile: filepath.Join(tempDir, "reviews"),
	}
	log, err := logger.NewSimpleLogger(testConfig.Logger)

	return NewAppService(mockClient, testConfig, log), testConfig
}

// TestGetApps tests the GetApps method of the AppService.
func TestGetApps(t *testing.T) {
	// Sample valid API response body
	validAppsJSON := `{"feed": {"results": [{"id":"1", "artistName":"Artist1", "releaseDate":"2023-01-01", "name":"App1"}, {"id":"2", "artistName":"Artist2", "releaseDate":"2023-02-02", "name":"App2"}]}}`

	// Sample invalid API response body
	invalidJSON := `{"feed": {"results": "invalid"}}`

	t.Run("successful fetch from API", func(t *testing.T) {
		s, cfg := setupTestService(validAppsJSON, http.StatusOK, t)
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		apps, err := s.GetApps()
		if err != nil {
			t.Fatalf("GetApps() failed unexpectedly: %v", err)
		}
		if len(apps) != 2 {
			t.Fatalf("Expected 2 apps, but got %d", len(apps))
		}
		if apps[0].ID != "1" {
			t.Errorf("Expected app ID '1', got '%s'", apps[0].ID)
		}
	})

	t.Run("fetch from file when it exists", func(t *testing.T) {
		s, cfg := setupTestService("", http.StatusInternalServerError, t) // Mock client returns an error
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		// Create a dummy file to be read
		mockFileContent := `[{"id":"3", "artistName":"Artist3", "releaseDate":"2023-03-03", "name":"App3"}]`
		if err := os.WriteFile(cfg.AppsStorageFile, []byte(mockFileContent), 0644); err != nil {
			t.Fatalf("Failed to write mock app file: %v", err)
		}

		apps, err := s.GetApps()
		if err != nil {
			t.Fatalf("GetApps() failed unexpectedly: %v", err)
		}
		if len(apps) != 1 {
			t.Fatalf("Expected 1 app, but got %d", len(apps))
		}
		if apps[0].ID != "3" {
			t.Errorf("Expected app ID '3', got '%s'", apps[0].ID)
		}
	})

	t.Run("API returns a non-200 status code", func(t *testing.T) {
		s, cfg := setupTestService("", http.StatusNotFound, t)
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		_, err := s.GetApps()
		if err == nil {
			t.Fatal("GetApps() was expected to return an error, but it did not.")
		}
		expectedErr := "received non-200 status code: 404"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("API returns invalid JSON", func(t *testing.T) {
		s, cfg := setupTestService(invalidJSON, http.StatusOK, t)
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		_, err := s.GetApps()
		if err == nil {
			t.Fatal("GetApps() was expected to return an error, but it did not.")
		}
	})
}

// TestGetReviews tests the GetReviews method of the AppService.
func TestGetReviews(t *testing.T) {
	// Sample valid API review response body
	validReviewsJSON := `{
		"feed": {
			"entry": [
				{"id":{"label":"1"},"author":{"name":{"label":"User1"}},"content":{"label":"Great app!"},"im:rating":{"label":"5"},"updated":{"label":"2023-08-21T09:00:00Z"}},
				{"id":{"label":"2"},"author":{"name":{"label":"User2"}},"content":{"label":"It's ok."},"im:rating":{"label":"3"},"updated":{"label":"2023-08-21T08:00:00Z"}},
				{"id":{"label":"3"},"author":{"name":{"label":"User3"}},"content":{"label":"Terrible."},"im:rating":{"label":"1"},"updated":{"label":"2023-08-20T00:00:00Z"}}
			]
		}
	}`

	t.Run("return all reviews with hours=0", func(t *testing.T) {
		s, cfg := setupTestService(validReviewsJSON, http.StatusOK, t)
		defer os.RemoveAll(filepath.Dir(cfg.ReviewsStorageFile))

		reviews, err := s.GetReviews("123", 0)
		if err != nil {
			t.Fatalf("GetReviews() failed unexpectedly: %v", err)
		}
		if len(reviews) != 3 {
			t.Fatalf("Expected 3 reviews, but got %d", len(reviews))
		}
		if reviews[0].ID != "1" {
			t.Errorf("Expected review ID '1', got '%s'", reviews[0].ID)
		}
	})
}
