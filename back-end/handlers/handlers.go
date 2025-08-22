package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runway/config"
	"runway/logger"
	"runway/services"
	"strconv"
	"time"
)

// Handlers struct holds the dependencies for all HTTP handlers.
type Handlers struct {
	AppService services.AppServiceInterface
	Config     *config.Config
	Logger     *logger.SimpleLogger
}

// NewHandlers creates a new Handlers instance with the provided dependencies.
func NewHandlers(appService services.AppServiceInterface, cfg *config.Config, log *logger.SimpleLogger) *Handlers {
	return &Handlers{
		AppService: appService,
		Config:     cfg,
		Logger:     log,
	}
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().String(),
		"version":   os.Getenv("APP_VERSION"),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// AppListHandler is the handler for the /app/list endpoint.
// It fetches a list of apps and returns them as a JSON response.
func (h *Handlers) AppListHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Processing app list request")
	apps, err := h.AppService.GetApps()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching apps: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		h.Logger.Error("Failed to encode JSON response", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
	h.Logger.Info("Successfully returned app list", "count", len(apps))
}

// AppReviewsHandler is the handler for the /app/reviews endpoint.
// It retrieves app reviews based on the provided app ID and filters them by a time window.
// The 'hours' parameter is now optional.
func (h *Handlers) AppReviewsHandler(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("id")
	if appID == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	hoursStr := r.URL.Query().Get("hours")
	h.Logger.Info("Processing app reviews request", "appID", appID, "hours", hoursStr)
	hours := 0
	if hoursStr != "" {
		var err error
		hours, err = strconv.Atoi(hoursStr)
		if err != nil || hours < 0 {
			h.Logger.Error("Invalid hours parameter", err, "hours", hoursStr)
			http.Error(w, "Invalid 'hours' parameter", http.StatusBadRequest)
			return
		}
	}

	reviews, err := h.AppService.GetReviews(appID, hours)
	if err != nil {
		h.Logger.Error("Failed to fetch reviews", err, "appID", appID)
		http.Error(w, fmt.Sprintf("Error fetching reviews: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		h.Logger.Error("Failed to encode JSON response", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
	h.Logger.Info("Successfully returned reviews", "count", len(reviews), "appID", appID)
}
