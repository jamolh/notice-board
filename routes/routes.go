package routes

import (
	_ "github.com/jamolh/notice-board/docs"
	"github.com/jamolh/notice-board/handlers"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes() *httprouter.Router {
	// declaring our routes
	var router = httprouter.New()

	// swagger doc api
	router.HandlerFunc("GET", "/swagger/*any", httpSwagger.WrapHandler)

	router.POST("/v1/notices", handlers.CreateNoticeHandler)
	router.GET("/v1/notices/:id", handlers.GetNoticeHandler)
	router.GET("/v1/notices", handlers.GetNoticesHandler)
	return router
}
