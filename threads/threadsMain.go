package main

import (
	"learning-web-chatboard2/common"
	"log"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var dbEngine *xorm.Engine
var config *common.Configuration
var logger *log.Logger

func main() {
	var err error
	// config
	config, err = common.LoadConfig()
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
	//database
	dbEngine, err = common.OpenDb(
		config.DbName,
		config.ShowSQL,
		0,
	)
	if err != nil {
		common.LogError(logger).Fatalln(err.Error())
	}
	//router
	routeEngine := gin.Default()
	routeEngine.POST("/create", createThread)
	routeEngine.POST("/create-post", createPost)
	routeEngine.POST("/read", readAThread)
	routeEngine.POST("/read-posts", readPostsInThread)
	routeEngine.GET("/read-index", readThreads)
	routeEngine.POST("/update", updateThread)

	routeEngine.Run(config.AddressThreads)
}
