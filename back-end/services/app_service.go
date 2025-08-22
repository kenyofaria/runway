package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runway/config"
	"runway/logger"
	"runway/models"
	"time"
)

type AppServiceInterface interface {
	GetApps() ([]*models.AppResponse, error)
	GetAppReviewsFromApi(appID string) ([]models.Review, error)
	GetReviews(appID string, hours int) ([]models.ReviewResponse, error)
}

// AppService handles fetching app data.
type AppService struct {
	Client *http.Client
	Config *config.Config
	Logger *logger.SimpleLogger
}

func NewAppService(client *http.Client, cfg *config.Config, log *logger.SimpleLogger) *AppService {
	return &AppService{
		Client: client,
		Config: cfg,
		Logger: log,
	}
}

// GetApps fetches a list of apps from a given URL and deserializes
// the JSON response into an array of App structs.
func (s *AppService) GetApps() ([]*models.AppResponse, error) {
	s.Logger.Info("Fetching apps from API", "url", s.Config.AppsApiUrl)
	existingApps, err := s.loadAppsFromFile(s.Config.AppsStorageFile)

	if err != nil {
		s.Logger.Debug("Failed to load apps from file, will fetch from API", "error", err)
		_ = fmt.Errorf("failed to load apps from apps.json: %w", err)
	} else if len(existingApps) != 0 {
		s.Logger.Info("Loaded apps from cache file", "count", len(existingApps))
		appResponses := make([]*models.AppResponse, len(existingApps))
		for i, app := range existingApps {
			response, _ := app.ToAppResponse()
			appResponses[i] = response
		}
		return appResponses, nil
	}

	resp, err := s.Client.Get(s.Config.AppsApiUrl)
	//resp, err := s.Client.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.Logger.Error("API returned non-200 status", nil, "status", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error("Failed to read response body", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var root models.Root
	err = json.Unmarshal(body, &root)
	if err != nil {
		s.Logger.Error("Failed to unmarshal JSON response", err)
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	err = s.saveAppsToFile(root.Feed.Entries, s.Config.AppsStorageFile)
	if err != nil {
		s.Logger.Error("Failed to save apps to file", err)
		fmt.Errorf("failed to write apps.json file: %w", err)
	} else {
		s.Logger.Info("Successfully saved apps to cache file", "count", len(root.Feed.Entries))
	}
	appResponses := s.convertRootToAppResponse(root)
	s.Logger.Info("Successfully fetched apps from API", "count", len(root.Feed.Entries))
	return appResponses, nil
}

func (s *AppService) convertRootToAppResponse(root models.Root) []*models.AppResponse {
	var appResponses []*models.AppResponse
	for _, app := range root.Feed.Entries {
		response, _ := app.ToAppResponse()
		appResponses = append(appResponses, response)
	}
	return appResponses
}

// saveAppsToFile marshals the provided slice of App structs and saves it to a JSON file.
func (s *AppService) saveAppsToFile(apps []models.App, filename string) error {
	return saveDataToFile(apps, filename)
}

// loadAppsFromFile reads a JSON file, unmarshal the data, and returns a slice of App structs.
func (s *AppService) loadAppsFromFile(filename string) ([]models.App, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var apps []models.App
	err = json.Unmarshal(jsonData, &apps)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from file: %w", err)
	}

	return apps, nil
}

// GetAppReviewsFromApi fetches a list of reviews for a specific app ID.
func (s *AppService) GetAppReviewsFromApi(appID string) ([]models.Review, error) {
	url := fmt.Sprintf("%s/id=%s/sortBy=mostRecent/page=1/json", s.Config.ReviewsBaseUrl, appID)
	s.Logger.Info("Fetching reviews from API", "appID", appID)
	resp, err := s.Client.Get(url)
	if err != nil {
		s.Logger.Error("HTTP request failed", err)
		return nil, fmt.Errorf("failed to make HTTP request for reviews: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.Logger.Error("API returned non-200 status", nil, "status", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 status code for reviews: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error("Failed to read response body", err)
		return nil, fmt.Errorf("failed to read reviews response body: %w", err)
	}

	var reviewResponse models.ReviewFeed
	err = json.Unmarshal(body, &reviewResponse)
	if err != nil {
		s.Logger.Error("Failed to unmarshal JSON response", err)
		return nil, fmt.Errorf("failed to unmarshal reviews JSON: %w", err)
	}
	err = s.saveReviewsToFile(reviewResponse.Feed.Entries, s.Config.ReviewsStorageFile)
	if err != nil {
		s.Logger.Error("failed to write reviews.json file: %w", err)
	}
	s.Logger.Info("Successfully fetched reviews from API", "count", len(reviewResponse.Feed.Entries))
	return reviewResponse.Feed.Entries, nil
}

// saveReviewsToFile marshals the provided slice of Review structs and saves it to a JSON file.
func (s *AppService) saveReviewsToFile(reviews []models.Review, filename string) error {
	return saveDataToFile(reviews, filename)
}

// loadReviewsFromFile reads a JSON file, unmarshal the data, and returns a slice of Review structs.
func (s *AppService) loadReviewsFromFile(filename string) ([]models.Review, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var reviews []models.Review
	err = json.Unmarshal(jsonData, &reviews)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from file: %w", err)
	}

	return reviews, nil
}

func convertReviews(reviews []models.Review) ([]models.ReviewResponse, error) {
	var reviewResponses []models.ReviewResponse
	for _, review := range reviews {
		response, err := review.ToReviewResponse()
		if err != nil {
			return nil, err
		}
		reviewResponses = append(reviewResponses, *response)
	}
	return reviewResponses, nil
}

func (s *AppService) GetReviews(appID string, hours int) ([]models.ReviewResponse, error) {
	s.Logger.Info("Starting GetReviews operation", "appID", appID, "hours", hours)
	allReviews, err := s.GetAppReviewsFromApi(appID)
	if err != nil {
		s.Logger.Error("Failed to get reviews from API", err, "appID", appID)
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	if hours == 0 {
		reviews, err := convertReviews(allReviews)
		if err != nil {
			s.Logger.Error("Failed to convert reviews", err)
			return nil, err
		}
		s.Logger.Info("Returning all reviews", "total", len(reviews))
		return reviews, nil
	}
	var recentReviews []models.Review
	cutoff := time.Now().Add(time.Duration(-hours) * time.Hour)
	for _, review := range allReviews {
		reviewTime, err := time.Parse(time.RFC3339, review.Timestamp.Label)
		if err != nil {
			s.Logger.Debug("Failed to parse review timestamp, skipping", "error", err, "timestamp", review.Timestamp.Label)
			continue // Skip this review if its timestamp is invalid
		}
		if reviewTime.After(cutoff) {
			recentReviews = append(recentReviews, review)
		}
	}
	reviews, err := convertReviews(recentReviews)
	if err != nil {
		s.Logger.Error("Failed to convert filtered reviews", err)
		return nil, err
	}

	s.Logger.Info("Successfully filtered reviews by time", "total", len(allReviews), "filtered", len(reviews))
	return reviews, nil
}

// saveDataToFile is a generic function that marshals a slice of any type T to a pretty-printed JSON file.
// It creates the directory if it doesn't exist and writes the data to the specified filename.
func saveDataToFile[T any](data []T, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the JSON data to the specified file.
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}
	return nil
}
