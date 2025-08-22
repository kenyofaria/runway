package main

import (
	"fmt"
	"net/http"
	"os"
	"runway/config"
	"runway/handlers"
	"runway/logger"
	"runway/middleware" // Import the new middleware package
	"runway/services"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}
	log, err := logger.NewSimpleLogger(cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.TimeoutSecs * 10000000000000000),
	}
	appService := services.NewAppService(httpClient, cfg, log)
	apiHandlers := handlers.NewHandlers(appService, cfg, log)
	http.Handle("/app/list", middleware.CORS(http.HandlerFunc(apiHandlers.AppListHandler)))
	http.Handle("/app/reviews", middleware.CORS(http.HandlerFunc(apiHandlers.AppReviewsHandler)))
	fmt.Printf("Server starting on port %d...\n", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		os.Exit(1)
	}
}
