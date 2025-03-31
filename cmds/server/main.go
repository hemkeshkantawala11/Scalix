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

    nodes := []string{"Node1", "Node2", "Node3"} // Initial nodes
    cacheSize := 1000 

    keyValueCache := cache.New(nodes, cacheSize)
    cacheHandler := handlers.NewCacheHandler(keyValueCache)
    router := gin.Default()

    router.POST("/put", cacheHandler.SetHandler)
    router.GET("/get", cacheHandler.GetHandler)

    srv := &http.Server{
        Addr:           ":7171",
        Handler:        router,
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   5 * time.Second,
        IdleTimeout:    30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    log.Println("Starting server on :7171")
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
    }
}
