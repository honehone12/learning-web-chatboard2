package main

import (
	"errors"
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
	createPostFailMsg  = "failed to create post"
	readAThreadFailMsg = "failed to read thread"
	updateTreadFailMsg = "failed to update thread"
)

func createThread(ctx *gin.Context) {
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
	err = createThreadInternal(&newThre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createTreadFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, &newThre)
}

func createPost(ctx *gin.Context) {
	var post common.Post
	err := ctx.Bind(&post)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createPostFailMsg)
		return
	}
	if common.IsEmpty(
		post.Body,
		post.Contributor,
	) {
		common.LogError(logger).Println("contains empty string")
		ctx.String(http.StatusBadRequest, createPostFailMsg)
		return
	}
	err = createPostInternal(&post)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createPostFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, &post)
}

func readAThread(ctx *gin.Context) {
	var thre common.Thread
	err := ctx.Bind(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readAThreadFailMsg)
		return
	}
	if common.IsEmpty(thre.UuId) {
		common.LogError(logger).Println("need uuid for finding thread")
		ctx.String(http.StatusBadRequest, readAThreadFailMsg)
		return
	}
	err = readAThreadInternal(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readAThreadFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, &thre)
}

func updateThread(ctx *gin.Context) {
	var thre common.Thread
	err := ctx.Bind(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, updateTreadFailMsg)
		return
	}
	if common.IsEmpty(thre.UuId, thre.Topic, thre.Owner) {
		common.LogError(logger).Println("contains empty string")
		ctx.String(http.StatusBadRequest, updateTreadFailMsg)
		return
	}
	err = updateThreadInternal(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, updateTreadFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, &thre)
}

func readPostsInThread(ctx *gin.Context) {
	var thre common.Thread
	err := ctx.Bind(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readAThreadFailMsg)
		return
	}
	// is there a way to check valid id before?
	posts, err := readPostsInThreadInternal(&thre)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readAThreadFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, &posts)
}

func readThreads(ctx *gin.Context) {
	thres, err := readThreadsInternal()
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusInternalServerError, readAThreadFailMsg)
	} else {
		ctx.JSON(http.StatusOK, &thres)
	}
}

func createThreadInternal(newThre *common.Thread) (err error) {
	now := time.Now()
	newThre.UuId = common.NewUuIdString()
	newThre.LastUpdate = now
	newThre.CreatedAt = now

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

func createPostInternal(newPost *common.Post) (err error) {
	newPost.UuId = common.NewUuIdString()
	newPost.CreatedAt = time.Now()

	affected, err := dbEngine.
		Table(PostsTable).
		InsertOne(newPost)
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

func readAThreadInternal(thread *common.Thread) (err error) {
	ok, err := dbEngine.
		Table(ThreadsTable).
		Get(thread)
	if err == nil && !ok {
		err = errors.New("no such thread")
	}
	return
}

func updateThreadInternal(thread *common.Thread) (err error) {
	thread.LastUpdate = time.Now()

	affected, err := dbEngine.
		Table(ThreadsTable).
		ID(thread.Id).
		Update(thread)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func readPostsInThreadInternal(thread *common.Thread) (posts []common.Post, err error) {
	err = dbEngine.
		Table(PostsTable).
		Where("thread_id = ?", thread.Id).
		Find(&posts)
	return
}
