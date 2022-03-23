package main

import (
	"learning-web-chatboard2/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var httpClient *http.Client
var config *common.Configuration
var logger *log.Logger

func main() {
	var err error
	// config
	config, err = common.LoadConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//processor data
	err = newProcessor()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//log
	logger, err = common.OpenLogger(
		config.LogToFile,
		config.LogFileNameThreads,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	//gin
	webEngine := gin.Default()
	// setup templates
	webEngine.Static("/static", "./public")
	webEngine.Delims("{{", "}}")
	webEngine.LoadHTMLGlob("./templates/*")
	//setup routes
	webEngine.GET("/", LoggedInCheckerMiddleware, indexGet)
	webEngine.GET("/error", LoggedInCheckerMiddleware, errorGet)

	usersRoute := webEngine.Group("/user")
	usersRoute.Use(LoggedInCheckerMiddleware)
	usersRoute.GET("/login", loginGet)
	usersRoute.GET("/signup", signupGet)
	usersRoute.GET("logout", logoutGet)
	usersRoute.POST("/signup-account", signupPost)
	usersRoute.POST("/authenticate", authenticatePost)

	threadsRoute := webEngine.Group("/thread")
	threadsRoute.Use(LoggedInCheckerMiddleware)
	threadsRoute.GET("/read", threadGet)
	threadsRoute.GET("/new", newThreadGet)
	threadsRoute.POST("/create", newThreadPost)
	threadsRoute.POST("/post", newReplyPost)

	httpClient = http.DefaultClient
	webEngine.Run(config.AdressRouter)
}
