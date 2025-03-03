package main

import (
	"github.com/Shashankm886/url-shortener/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.RegisterRoutes(router)

	router.Run(":8080")
}
