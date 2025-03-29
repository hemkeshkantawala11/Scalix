package main

import (
	"log"
	"net/http"
	"time"

	"HLD-REDIS-ASSIGNMENT/internal/cache"
	"HLD-REDIS-ASSIGNMENT/internal/cache_handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	keyValueCache := cache.New(cache.WithCapacity(1000))
	cacheHandler := handlers.NewCacheHandler(keyValueCache)
	router := gin.Default()

	// Middleware to enable Keep-Alive
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/cache", cacheHandler.SetHandler)
		v1.GET("/cache", cacheHandler.GetHandler)
	}

	// Background task to clean expired cache items
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			keyValueCache.CleanExpired()
			log.Println("Expired cache items cleaned")
		}
	}()

	// Custom HTTP server with Keep-Alive settings
	srv := &http.Server{
		Addr:              ":7171",
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second, // Allows Keep-Alive for 60 sec
		MaxHeaderBytes:    1 << 20,          // 1MB header limit
	}

	log.Println("Starting server on :7171")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
