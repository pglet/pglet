package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./client/build"
	siteDefaultDocument string = "index.html"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(contentRootFolder, true)))

	// Setup route group for the API
	api := router.Group(apiRoutePrefix)
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// unknown API routes - 404, all the rest - index.html
	router.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.RequestURI)
		if !strings.HasPrefix(c.Request.RequestURI, apiRoutePrefix+"/") {
			c.File(contentRootFolder + "/" + siteDefaultDocument)
		}
	})

	// Start and run the server
	router.Run(":5000")
}
