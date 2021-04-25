package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func init() {
}

// @title Notice-Board API
// @version v0.0.1
// @description This is simple API to interacting with the Notice-Board server
// @host http://localhost:3000
// @schemes http
func main() {
	port, found := os.LookupEnv("SERVER_ADDR")
	if !found {
		log.Fatal("apiservice:main SERVER_ADDR not found")
	}

	db.Connect()
	defer db.Close()

	initRoutes()
	srv := &http.Server{
		Addr:         port, //port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	// TODO: refactor!
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
		s := <-sigint
		log.Println("server received signal", s)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("server: couldn't shutdown because of " + err.Error())
		}
	}()

	log.Println("http server is running on port", port)
	log.Fatal("Closing Server error:", srv.ListenAndServe())

}

// declaring our routes
func initRoutes() {
	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Write([]byte("i'm ok"))
	})
	router.POST("/v1/notices", handlers.CreateNoticeHandler)
	router.GET("/v1/notices/:id", handlers.GetNoticesByIDHandler)
	router.GET("/v1/notices", handlers.GetNoticesHandler)
	router.HandlerFunc("GET", "/swagger/*any", httpSwagger.WrapHandler)
	// router.HandlerFunc("GET", "/swagger/*any", httpSwagger.Handler(
	// 	httpSwagger.URL("docs/swagger.json"),
	// ))
	// router.GET("/swagger/", httpSwagger.Handler(
	// 	httpSwagger.URL("/docs/docs.json"),
	// ))
	// router.HandlerFunc("/swagger/", httpSwagger.Handler(
	// 	httpSwagger.URL("/docs/docs.json"),
	// ))
}
