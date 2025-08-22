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
	t.Run("successful fetch from API", func(t *testing.T) {
		s, cfg := setupTestService(getValidAppsJSON(), http.StatusOK, t)
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		apps, err := s.GetApps()
		if err != nil {
			t.Fatalf("GetApps() failed unexpectedly: %v", err)
		}
		if len(apps) != 2 {
			t.Fatalf("Expected 2 apps, but got %d", len(apps))
		}
		if apps[0].AppID != "123456789" {
			t.Errorf("Expected app ID '123456789', got '%s'", apps[0].AppID)
		}
		if apps[0].Name != "Test App 1" {
			t.Errorf("Expected app name 'Test App 1', got '%s'", apps[0].Name)
		}
		if apps[0].Author != "Test Artist 1" {
			t.Errorf("Expected app author 'Test Artist 1', got '%s'", apps[0].Author)
		}
	})

	t.Run("fetch from file when it exists", func(t *testing.T) {
		s, cfg := setupTestService("", http.StatusInternalServerError, t) // Mock client returns an error
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		// Create a dummy file to be read
		if err := os.WriteFile(cfg.AppsStorageFile, []byte(getMockFileContentJSON()), 0644); err != nil {
			t.Fatalf("Failed to write mock app file: %v", err)
		}

		apps, err := s.GetApps()
		if err != nil {
			t.Fatalf("GetApps() failed unexpectedly: %v", err)
		}
		if len(apps) != 1 {
			t.Fatalf("Expected 1 app, but got %d", len(apps))
		}
		if apps[0].AppID != "123456791" {
			t.Errorf("Expected app ID '123456791', got '%s'", apps[0].AppID)
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
		s, cfg := setupTestService(getInvalidJSON(), http.StatusOK, t)
		defer os.RemoveAll(filepath.Dir(cfg.AppsStorageFile))

		_, err := s.GetApps()
		if err == nil {
			t.Fatal("GetApps() was expected to return an error, but it did not.")
		}
	})
}

// TestGetReviews tests the GetReviews method of the AppService.
func TestGetReviews(t *testing.T) {
	t.Run("return all reviews with hours=0", func(t *testing.T) {
		s, cfg := setupTestService(getValidReviewsJSON(), http.StatusOK, t)
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

// getValidAppsJSON returns a mock JSON response that matches the iTunes API structure
func getValidAppsJSON() string {
	return `{
		"feed": {
			"entry": [
				{
					"id": {
						"label": "https://apps.apple.com/us/app/test-app-1/id123456789?uo=2",
						"attributes": {
							"im:id": "123456789",
							"im:bundleId": "com.test.app1"
						}
					},
					"im:name": {
						"label": "Test App 1"
					},
					"im:artist": {
						"label": "Test Artist 1",
						"attributes": {
							"href": "https://apps.apple.com/us/developer/test-artist-1/id987654321?uo=2"
						}
					},
					"im:releaseDate": {
						"label": "2023-01-01T00:00:00-07:00",
						"attributes": {
							"label": "January 1, 2023"
						}
					},
					"category": {
						"attributes": {
							"im:id": "6007",
							"term": "Productivity",
							"scheme": "https://apps.apple.com/us/genre/ios-productivity/id6007?uo=2",
							"label": "Productivity"
						}
					},
					"im:image": [
						{
							"label": "https://example.com/icon.png",
							"attributes": {
								"height": "100"
							}
						}
					],
					"summary": {
						"label": "A test app for testing purposes"
					},
					"im:price": {
						"label": "Get",
						"attributes": {
							"amount": "0.00",
							"currency": "USD"
						}
					},
					"im:contentType": {
						"attributes": {
							"term": "Application",
							"label": "Application"
						}
					},
					"rights": {
						"label": "© 2023 Test Company"
					},
					"title": {
						"label": "Test App 1 - Test Artist 1"
					},
					"link": {
						"attributes": {
							"rel": "alternate",
							"type": "text/html",
							"href": "https://apps.apple.com/us/app/test-app-1/id123456789?uo=2"
						}
					}
				},
				{
					"id": {
						"label": "https://apps.apple.com/us/app/test-app-2/id123456790?uo=2",
						"attributes": {
							"im:id": "123456790",
							"im:bundleId": "com.test.app2"
						}
					},
					"im:name": {
						"label": "Test App 2"
					},
					"im:artist": {
						"label": "Test Artist 2",
						"attributes": {
							"href": "https://apps.apple.com/us/developer/test-artist-2/id987654322?uo=2"
						}
					},
					"im:releaseDate": {
						"label": "2023-02-02T00:00:00-07:00",
						"attributes": {
							"label": "February 2, 2023"
						}
					},
					"category": {
						"attributes": {
							"im:id": "6008",
							"term": "Photo & Video",
							"scheme": "https://apps.apple.com/us/genre/ios-photo-video/id6008?uo=2",
							"label": "Photo & Video"
						}
					},
					"im:image": [
						{
							"label": "https://example.com/icon2.png",
							"attributes": {
								"height": "100"
							}
						}
					],
					"summary": {
						"label": "Another test app for testing purposes"
					},
					"im:price": {
						"label": "Get",
						"attributes": {
							"amount": "0.00",
							"currency": "USD"
						}
					},
					"im:contentType": {
						"attributes": {
							"term": "Application",
							"label": "Application"
						}
					},
					"rights": {
						"label": "© 2023 Test Company 2"
					},
					"title": {
						"label": "Test App 2 - Test Artist 2"
					},
					"link": {
						"attributes": {
							"rel": "alternate",
							"type": "text/html",
							"href": "https://apps.apple.com/us/app/test-app-2/id123456790?uo=2"
						}
					}
				}
			]
		}
	}`
}

// getMockFileContentJSON returns a mock JSON for file cache testing
func getMockFileContentJSON() string {
	return `[
		{
			"id": {
				"label": "https://apps.apple.com/us/app/test-app-3/id123456791?uo=2",
				"attributes": {
					"im:id": "123456791",
					"im:bundleId": "com.test.app3"
				}
			},
			"im:name": {
				"label": "Test App 3"
			},
			"im:artist": {
				"label": "Test Artist 3",
				"attributes": {
					"href": "https://apps.apple.com/us/developer/test-artist-3/id987654323?uo=2"
				}
			},
			"im:releaseDate": {
				"label": "2023-03-03T00:00:00-07:00",
				"attributes": {
					"label": "March 3, 2023"
				}
			},
			"category": {
				"attributes": {
					"im:id": "6007",
					"term": "Productivity",
					"scheme": "https://apps.apple.com/us/genre/ios-productivity/id6007?uo=2",
					"label": "Productivity"
				}
			},
			"im:image": [
				{
					"label": "https://example.com/icon3.png",
					"attributes": {
						"height": "100"
					}
				}
			],
			"summary": {
				"label": "A third test app for testing purposes"
			},
			"im:price": {
				"label": "Get",
				"attributes": {
					"amount": "0.00",
					"currency": "USD"
				}
			},
			"im:contentType": {
				"attributes": {
					"term": "Application",
					"label": "Application"
				}
			},
			"rights": {
				"label": "© 2023 Test Company 3"
			},
			"title": {
				"label": "Test App 3 - Test Artist 3"
			},
			"link": {
				"attributes": {
					"rel": "alternate",
					"type": "text/html",
					"href": "https://apps.apple.com/us/app/test-app-3/id123456791?uo=2"
				}
			}
		}
	]`
}

// getInvalidJSON returns invalid JSON for error testing
func getInvalidJSON() string {
	return `{"feed": {"entry": "invalid"}}`
}

// getValidReviewsJSON returns a mock reviews JSON response
func getValidReviewsJSON() string {
	return `{
		"feed": {
			"entry": [
				{
					"id": {"label": "1"},
					"author": {"name": {"label": "User1"}},
					"content": {"label": "Great app!"},
					"im:rating": {"label": "5"},
					"updated": {"label": "2023-08-21T09:00:00Z"}
				},
				{
					"id": {"label": "2"},
					"author": {"name": {"label": "User2"}},
					"content": {"label": "It's ok."},
					"im:rating": {"label": "3"},
					"updated": {"label": "2023-08-21T08:00:00Z"}
				},
				{
					"id": {"label": "3"},
					"author": {"name": {"label": "User3"}},
					"content": {"label": "Terrible."},
					"im:rating": {"label": "1"},
					"updated": {"label": "2023-08-20T00:00:00Z"}
				}
			]
		}
	}`
}
