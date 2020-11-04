package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pglet/pglet/internal/page"
)

const (
	DefaultServerPort   int    = 5000
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "client/build"
	siteDefaultDocument string = "index.html"
)

func Start(ctx context.Context, serverPort int) {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", BinaryFileSystem(contentRootFolder, "")))

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
		if !strings.HasPrefix(c.Request.RequestURI, apiRoutePrefix+"/") {
			// SPA index.html
			indexData, _ := Asset(contentRootFolder + "/" + siteDefaultDocument)
			c.Data(http.StatusOK, "text/html", indexData)
		} else {
			// API not found
			c.JSON(http.StatusNotFound, gin.H{
				"message": "API endpoint not found",
			})
		}
	})

	// Start and run the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
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
		// User ID is invalid
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
