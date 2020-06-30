package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pglet/pglet/page"
	"github.com/pglet/pglet/utils"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	apiRoutePrefix      string = "/api"
	contentRootFolder   string = "./client/build"
	siteDefaultDocument string = "index.html"
)

var (
	isServer   bool
	serverPort int
	serverAddr string
	pageName   string
	sessionID  string
)

func main() {

	fmt.Println(utils.GenerateRandomString(16))
	//sha1 := utils.SHA1(strings.ToLower("Hello, world!"))
	//pipeName := fmt.Sprintf("pglet_pipe_%s", sha1)

	// fi, err := os.Stat(pipeName)
	// if os.IsNotExist(err) {
	// 	// create pipe
	// 	// check for IsExist
	// 	// start sub-process with WS client
	// }

	// send command to a named pipe

	if sessionID != "" {
		runProxy()
	} else if isServer {
		runServer()
	} else {
		runClient()
	}
}

func removeElementAt(source []int, pos int) []int {
	copy(source[pos:], source[pos+1:]) // Shift a[i+1:] left one index.
	source[len(source)-1] = 0          // Erase last element (write zero value).
	return source[:len(source)-1]      // Truncate slice.
}

func createTestPage() *page.Page {
	p, err := page.NewPage("test-1")
	if err != nil {
		log.Fatal(err)
	}

	p.AddControl(page.NewControl("Row", "0", "1"))
	p.AddControl(page.NewControl("Column", "1", "2"))
	p.AddControl(page.NewControl("Column", "1", "3"))

	ctl3 := page.NewControl("Text", "2", "4")
	p.AddControl(ctl3)

	ctl4 := page.NewControl("Button", "3", "5")
	ctl4["text"] = "Click me!"
	p.AddControl(ctl4)

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

	p.AddControl(ctl5)

	return p
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
	c.JSON(http.StatusOK, page.Pages().Get("test-1"))
}

func init() {
	flag.StringVar(&pageName, "page", "", "Page name to create and connect to.")
	flag.StringVar(&serverAddr, "server", "", "Pglet server address.")
	flag.StringVar(&sessionID, "session-id", "", "Client session ID.")
	flag.IntVar(&serverPort, "port", 5000, "The port number to run pglet server on.")
	flag.Parse()

	if pageName == "" {
		isServer = true
	}

	if serverPort < 0 || serverPort > 65535 {
		flag.PrintDefaults()
		os.Exit(1)
	}
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

	_, err1 := page.NewPage("test page 2")
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println(page.Pages())

	p2 := &page.Page{}

	err = json.Unmarshal([]byte(jsonPage), p2)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", p2)

	arr := []int{1, 2, 3, 4, 5, 6}

	arr = removeElementAt(arr, 1)
	fmt.Println(arr)
}

func runProxy() {
	fmt.Printf("Running in proxy mode: %s...\n", sessionID)
	time.Sleep(1 * time.Minute)
}

func runClient() {
	fmt.Printf("Running in client mode: %s...\n", pageName)
	u := url.URL{Scheme: "ws", Host: *&serverAddr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer func() {
			fmt.Println("Closing...")
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from Go: %s", pageName)))

	time.Sleep(5 * time.Second)

	// run proxy
	execPath, _ := os.Executable()
	fmt.Println(execPath)

	cmd := exec.Command(execPath, "--session-id=12345")
	err = cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cmd.Process.Pid)
}

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
	api.GET("/pages/:pageID", pageHandler)

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
