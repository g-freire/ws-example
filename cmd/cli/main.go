package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"ws-example/internal/profile"
	msg_http "ws-example/pkg/msg/transport/http"
	msg_ws "ws-example/pkg/msg/transport/websocket"

	"ws-example/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)



func main() {
	defer profile.Start("main")()

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}


func run() error{
	// CONFIG
	conf := config.GetConfig()

	// WS POOL
	var wg sync.WaitGroup
	p := msg_ws.NewPool(wg)
	h := msg_ws.NewHandler(*p, conf)
	go p.Start("", "")

	// WEB SERVER SETUP
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// ROUTING
	// REST
	msg_http.ApplyRoutes(r)
	// WEBSOCKET
	msg_ws.ApplyRoutes(r, *h)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg.Wait()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}
