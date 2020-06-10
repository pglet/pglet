package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pglet/pglet"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./client/build"
	siteDefaultDocument string = "index.html"
)

func removeElementAt(source []int, pos int) []int {
	copy(source[pos:], source[pos+1:]) // Shift a[i+1:] left one index.
	source[len(source)-1] = 0          // Erase last element (write zero value).
	return source[:len(source)-1]      // Truncate slice.
}

func createPage() pglet.Page {
	page := pglet.Page{}
	page.Name = "test page 1"
	page.Controls = make(map[string]pglet.Control)

	page.AddControl(pglet.NewControl("Row", "", "0"))
	page.AddControl(pglet.NewControl("Column", "0", "1"))
	page.AddControl(pglet.NewControl("Column", "0", "2"))

	ctl3 := pglet.NewControl("Text", "1", "3")
	page.AddControl(ctl3)

	ctl4 := pglet.NewControl("Button", "2", "4")
	ctl4["text"] = "Click me!"
	page.AddControl(ctl4)
	return page
}

func main() {

	page := createPage()

	//fmt.Println(ctl3)

	//ctl1 := page.controls["ctl_1"]

	var jsonPage string
	j, err := json.MarshalIndent(&page, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonPage = string(j)

	fmt.Printf("%+v\n\n\n", jsonPage)

	p2 := pglet.Page{}

	err = json.Unmarshal([]byte(jsonPage), &p2)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", p2)

	arr := []int{1, 2, 3, 4, 5, 6}

	arr = removeElementAt(arr, 1)
	fmt.Println(arr)

	return

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile(contentRootFolder, true)))

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
