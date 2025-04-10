package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.example.ini", "config file path")
	flag.Parse()
	config.Setup(configPath)

	gin.SetMode(config.App.AppMode)

	srv := &http.Server{
		Addr:              config.App.Listen,
		Handler:           router.Handler(),
		ReadTimeout:       config.App.ReadTimeout * time.Second,
		ReadHeaderTimeout: config.App.ReadTimeout * time.Second,
		WriteTimeout:      config.App.WriteTimeout * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server exited gracefully")
	}
}
