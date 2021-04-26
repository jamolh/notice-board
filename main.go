package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jamolh/notice-board/db"
	_ "github.com/jamolh/notice-board/docs"
	"github.com/jamolh/notice-board/helpers"
	"github.com/jamolh/notice-board/routers"

	"github.com/swaggo/http-swagger"

	_ "github.com/joho/godotenv/autoload"
)

const (
	serverAddr  = "SERVER_ADDR"
	databaseURL = "DATABASE_URL"
)

// @title Notice-Board API
// @version v0.0.1
// @description This is simple API to interacting with the Notice-Board server
// @host http://localhost:3000
// @licence.name Custom
// @schemes http
func main() {
	port := helpers.GetEnv(serverAddr, true)
	dbConnection := helpers.GetEnv(databaseURL, true)

	err := db.Connect(dbConnection)
	if err != nil {
		log.Fatal("Connecting database failed", err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:         port, //port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      routers.InitRoutes(),
	}

	go gracefulShutdown(srv)

	log.Println("http server is running on port", port)
	err = srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("closing server error:", err)
		}
	}
	log.Println("server closed")
}

func gracefulShutdown(srv *http.Server) {
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, os.Interrupt, os.Kill)
	s := <-stopSignal
	log.Println("server received signal", s)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server: couldn't shutdown because of ", err)
	}
}
