package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pglet/pglet/page"
)

const (
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./tests" //"./client/build"
	siteDefaultDocument string = "index.html"
)

func runServer() {
	createTestPages()

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
		fmt.Println(c.Request.RequestURI)
		if !strings.HasPrefix(c.Request.RequestURI, apiRoutePrefix+"/") {
			c.File(contentRootFolder + "/" + siteDefaultDocument)
		}
	})

	// Start and run the server
	router.Run(fmt.Sprintf(":%d", serverPort))
}

func userHandler(c *gin.Context) {

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
	fmt.Println("sessionID:", sessionID)

	fullPageName := fmt.Sprintf("%s/%s", accountName, pageName)
	c.JSON(http.StatusOK, page.Pages().Get(fullPageName))
}

func removeElementAt(source []int, pos int) []int {
	copy(source[pos:], source[pos+1:]) // Shift a[i+1:] left one index.
	source[len(source)-1] = 0          // Erase last element (write zero value).
	return source[:len(source)-1]      // Truncate slice.
}

func createTestPage() *page.Session {
	p := page.NewPage("test-1", false)

	s := page.NewSession(p, page.ZeroSession)

	s.AddControl(page.NewControl("Row", "0", "1"))
	s.AddControl(page.NewControl("Column", "1", "2"))
	s.AddControl(page.NewControl("Column", "1", "3"))

	ctl3 := page.NewControl("Text", "2", "4")
	s.AddControl(ctl3)

	ctl4 := page.NewControl("Button", "3", "5")
	ctl4["text"] = "Click me!"
	s.AddControl(ctl4)

	ctl5, err := page.NewControlFromJSON(`{
		"i": "myBtn",
		"p": "2",
		"t": "Button",
		"text": "Cancel"
	  }`)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(ctl5)

	s.AddControl(ctl5)

	return s
}

func createTestPages() {
	//fmt.Printf("string: %s", "sss")

	p := createTestPage()

	//fmt.Println(ctl3)

	//ctl1 := page.controls["ctl_1"]

	var jsonPage string
	j, err := json.MarshalIndent(&p, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonPage = string(j)

	fmt.Printf("----------------\n%+v\n--------------\n", jsonPage)

	// p := page.NewPage("test page 2")

	// fmt.Println(page.Pages())

	// p2 := &page.Page{}

	// err = json.Unmarshal([]byte(jsonPage), p2)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Printf("%+v\n", p2)

	// arr := []int{1, 2, 3, 4, 5, 6}

	// arr = removeElementAt(arr, 1)
	// fmt.Println(arr)
}
