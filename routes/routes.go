package routes

import (
	"github.com/Shashankm886/url-shortener/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/shorten", controllers.ShortenURL)
	router.GET("/:shortUrl", controllers.RedirectURL)
	router.GET("/usage/:shortUrl", controllers.GetUsage)
}
