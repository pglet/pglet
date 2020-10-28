package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pglet/pglet/internal/page"
)

const (
	DefaultServerPort   int    = 5000
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./tests" //"./client/build"
	siteDefaultDocument string = "index.html"
)

func Start(serverPort int) {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(contentRootFolder, true)))

	// WebSockets
	router.GET("/ws", func(c *gin.Context) {
		page.WebsocketHandler(c.Writer, c.Request)
	})

	// Setup route group for the API
	api := router.Group(apiRoutePrefix)
	{
		api.GET("/", func(c *gin.Context) {
			time.Sleep(4 * time.Second)
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	api.GET("/users/:userID", userHandler)
	api.GET("/pages/:accountName/:pageName", pageHandler)

	// unknown API routes - 404, all the rest - index.html
	router.NoRoute(func(c *gin.Context) {
		log.Println(c.Request.RequestURI)
		if !strings.HasPrefix(c.Request.RequestURI, apiRoutePrefix+"/") {
			c.File(contentRootFolder + "/" + siteDefaultDocument)
		}
	})

	// Start and run the server
	router.Run(fmt.Sprintf(":%d", serverPort))
}

func userHandler(c *gin.Context) {

	time.Sleep(2 * time.Second)

	if userID, err := strconv.Atoi(c.Param("userID")); err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"id":       userID,
			"username": "admin",
		})
	} else {
		// Joke ID is invalid
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func pageHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	accountName := c.Param("accountName")
	pageName := c.Param("pageName")
	sessionID := c.Query("sessionID")
	log.Println("sessionID:", sessionID)

	fullPageName := fmt.Sprintf("%s/%s", accountName, pageName)
	page := page.Pages().Get(fullPageName)
	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
		return
	}
	session := page.GetSession(sessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func removeElementAt(source []int, pos int) []int {
	copy(source[pos:], source[pos+1:]) // Shift a[i+1:] left one index.
	source[len(source)-1] = 0          // Erase last element (write zero value).
	return source[:len(source)-1]      // Truncate slice.
}
