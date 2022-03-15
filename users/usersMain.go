package main

import (
	"learning-web-chatboard2/common"
	"log"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var dbEngine *xorm.Engine

func main() {
	// config
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	//log
	logger, err := common.OpenLogger(
		config.LogToFile,
		config.LogFileNameUsers,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	//database
	dbEngine, err = common.OpenDb(
		config.DbNameUsers,
		config.ShowSQL,
		0,
	)
	if err != nil {
		common.LogError(logger).Fatalln(err.Error())
	}
	//router
	routeEngine := gin.Default()
	usersRoute := routeEngine.Group("/users")
	usersRoute.POST("/post", CreateUser)
	//should address be in envs or args ??
	routeEngine.Run(config.AddressUsers)
}
