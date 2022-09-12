package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	msg_http "ws-example/pkg/msg/transport/http"
	msg_ws "ws-example/pkg/msg/transport/websocket"

	"ws-example/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	// Config loads info from env files
	conf := config.GetConfig()

	// WEB SERVER SETUP
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// ROUTING
	// REST
	msg_http.ApplyRoutes(r)

	// WEBSOCKET
	go msg_ws.Pool.Start("", "")
	msg_ws.ApplyRoutes(r)

	srv := &http.Server{
		Addr:    ":" + conf.PortHTTP,
		Handler: r,
	}

	// GRACEFULL SHUTDOWN
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
