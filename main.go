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
	"github.com/jamolh/notice-board/handlers"

	"github.com/julienschmidt/httprouter"
	"github.com/swaggo/http-swagger"

	_ "github.com/joho/godotenv/autoload"
)

var (
	router = httprouter.New()
)

// @title Notice-Board API
// @version v0.0.1
// @description This is simple API to interacting with the Notice-Board server
// @host http://localhost:3000
// @licence.name Custom
// @schemes http
func main() {
	port, found := os.LookupEnv("SERVER_ADDR")
	if !found {
		log.Fatal("apiservice:main SERVER_ADDR not found")
	}

	dbConnection, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("db:Connect DATABASE_URL not found")
	}

	err := db.Connect(dbConnection)
	if err != nil {
		log.Fatal("Connecting database failed", err)
	}
	defer db.Close()

	initRoutes()
	srv := &http.Server{
		Addr:         port, //port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}
	go func() {
		stopSignal := make(chan os.Signal)
		signal.Notify(stopSignal, os.Interrupt, os.Kill)
		s := <-stopSignal
		log.Println("server received signal", s)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("server: couldn't shutdown because of ", err)
		}
	}()

	log.Println("http server is running on port", port)
	err = srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("closing server error:", err)
		}
	}
	log.Println("server closed")
}

// declaring our routes
func initRoutes() {
	router.HandlerFunc("GET", "/swagger/*any", httpSwagger.WrapHandler)

	router.POST("/v1/notices", handlers.CreateNoticeHandler)
	router.GET("/v1/notices/:id", handlers.GetNoticeHandler)
	router.GET("/v1/notices", handlers.GetNoticesHandler)
}
