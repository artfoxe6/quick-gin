package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quick_gin/config/env"
	"quick_gin/route"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// 初始化路由
	route.Init()
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(env.Server().Port),
		Handler:      route.Route,
		ReadTimeout:  time.Duration(env.Server().ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(env.Server().WriteTimeout) * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		}
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownTime := env.Server().ShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTime)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		os.Exit(-1)
	}
}
