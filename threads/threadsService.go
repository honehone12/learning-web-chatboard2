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
	threadsTable     = "threads"
	postsTable       = "posts"
	descendingUpdate = "last_update"
)

func handleErrorInternal(
	loggerErrorMsg string,
	ctx *gin.Context,
) {
	common.LogError(logger).Println(loggerErrorMsg)
	ctx.JSON(http.StatusBadRequest, gin.H{"status": "error"})
}

func createThread(ctx *gin.Context) {
	var newThre common.Thread
	err := createThreadInternal(ctx, &newThre)
	if err != nil {
		handleErrorInternal("failed to create thread", ctx)
		return
	}
	ctx.JSON(http.StatusOK, &newThre)
}

func createThreadInternal(ctx *gin.Context, newThre *common.Thread) (err error) {
	err = ctx.Bind(newThre)
	if err != nil {
		return
	}
	if common.IsEmpty(newThre.Topic, newThre.Owner) {
		err = errors.New("contains empty string")
		return
	}
	now := time.Now()
	newThre.UuId = common.NewUuIdString()
	newThre.LastUpdate = now
	newThre.CreatedAt = now
	err = createThreadSQLInternal(newThre)
	return
}

func createPost(ctx *gin.Context) {
	var post common.Post
	err := createPostInternal(ctx, &post)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &post)
}

func createPostInternal(ctx *gin.Context, post *common.Post) (err error) {
	err = ctx.Bind(post)
	if err != nil {
		return
	}
	if common.IsEmpty(
		post.Body,
		post.Contributor,
	) {
		err = errors.New("contains empty string")
		return
	}
	post.UuId = common.NewUuIdString()
	post.CreatedAt = time.Now()
	err = createPostSQLInternal(post)
	return
}

func readAThread(ctx *gin.Context) {
	var thre common.Thread
	err := readAThreadInternal(ctx, &thre)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &thre)
}

func readAThreadInternal(ctx *gin.Context, thre *common.Thread) (err error) {
	err = ctx.Bind(thre)
	if err != nil {
		return
	}
	if common.IsEmpty(thre.UuId) {
		err = errors.New("need uuid for finding thread")
		return
	}
	err = readAThreadSQLInternal(thre)
	return
}

func updateThread(ctx *gin.Context) {
	var thre common.Thread
	err := updateThreadInternal(ctx, &thre)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &thre)
}

func updateThreadInternal(ctx *gin.Context, thre *common.Thread) (err error) {
	err = ctx.Bind(thre)
	if err != nil {
		return
	}
	if common.IsEmpty(
		thre.UuId,
		thre.Topic,
		thre.Owner,
	) {
		err = errors.New("contains empty string")
		return
	}
	thre.LastUpdate = time.Now()
	err = updateThreadSQLInternal(thre)
	return
}

func readPostsInThread(ctx *gin.Context) {
	var thre common.Thread
	err := ctx.Bind(&thre)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	// is there a way to check valid id before?
	posts, err := readPostsInThreadSQLInternal(&thre)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &posts)
}

func readThreads(ctx *gin.Context) {
	thres, err := readThreadsSQLInternal()
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, &thres)
	}
}

func createThreadSQLInternal(newThre *common.Thread) (err error) {
	affected, err := dbEngine.
		Table(threadsTable).
		InsertOne(&newThre)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func createPostSQLInternal(newPost *common.Post) (err error) {
	affected, err := dbEngine.
		Table(postsTable).
		InsertOne(newPost)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func readThreadsSQLInternal() (threads []common.Thread, err error) {
	err = dbEngine.
		Table(threadsTable).
		Desc(descendingUpdate).
		Find(&threads)
	return
}

func readAThreadSQLInternal(thread *common.Thread) (err error) {
	ok, err := dbEngine.
		Table(threadsTable).
		Get(thread)
	if err == nil && !ok {
		err = errors.New("no such thread")
	}
	return
}

func updateThreadSQLInternal(thread *common.Thread) (err error) {
	affected, err := dbEngine.
		Table(threadsTable).
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

func readPostsInThreadSQLInternal(thread *common.Thread) (posts []common.Post, err error) {
	err = dbEngine.
		Table(postsTable).
		Where("thread_id = ?", thread.Id).
		Find(&posts)
	return
}
