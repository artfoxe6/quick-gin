package main

import (
	"context"
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 加载配置
	config.Load("./config/config.ini")

	//gin 运行模式
	gin.SetMode(config.App.AppMode)

	// 启动http server
	srv := &http.Server{
		Addr:              config.App.Listen,
		Handler:           app.Handler(),
		ReadTimeout:       config.App.ReadTimeout * time.Second,
		ReadHeaderTimeout: config.App.ReadTimeout * time.Second,
		WriteTimeout:      config.App.WriteTimeout * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
