package main

import (
	"fmt"
	"learning-web-chatboard2/common"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ThreadsTable     = "threads"
	PostsTable       = "posts"
	DescendingUpdate = "last_update"
)

const (
	createTreadFailMsg = "failed to create thread"
)

func CreateThread(ctx *gin.Context) {
	var newThre common.Thread
	err := ctx.Bind(&newThre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createTreadFailMsg)
		return
	}
	if common.IsEmpty(newThre.Topic, newThre.Owner) {
		common.LogError(logger).Println("contains empty string")
		ctx.String(http.StatusBadRequest, createTreadFailMsg)
		return
	}
	err = createThreadInternal(
		newThre.Topic,
		newThre.UserId,
		newThre.Owner,
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createTreadFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, newThre)
}

func ReadThreads(ctx *gin.Context) {
	thres, err := readThreadsInternal()
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusInternalServerError, "failed to read threads")
	} else {
		ctx.JSON(http.StatusOK, thres)
	}
}

func createThreadInternal(
	topic string,
	userId uint,
	userName string,
) (err error) {
	now := time.Now()
	newThre := common.Thread{
		UuId:       common.NewUuIdString(),
		Topic:      topic,
		Owner:      userName,
		UserId:     userId,
		LastUpdate: now,
		CreatedAt:  now,
	}
	affected, err := dbEngine.
		Table(ThreadsTable).
		InsertOne(&newThre)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func readThreadsInternal() (threads []common.Thread, err error) {
	err = dbEngine.
		Table(ThreadsTable).
		Desc(DescendingUpdate).
		Find(&threads)
	return
}
