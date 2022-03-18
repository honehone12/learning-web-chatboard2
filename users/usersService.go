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
	UserTable    = "users"
	SessionTable = "sessions"
)

const (
	createUserFailMsg    = "failed to create user"
	createSessionFailMsg = "failed to create session"
	readUserFailMsg      = "failed to find user"
	readSessionFailMsg   = "failed to find session"
	deleteSessionFailMsg = "failed to delete session"
)

// don't forget TLS
func createUser(ctx *gin.Context) {
	var newUser common.User
	err := ctx.Bind(&newUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createUserFailMsg)
		return
	}
	if common.IsEmpty(
		newUser.Name,
		newUser.Email,
		newUser.Password,
	) {
		common.LogError(logger).Println("contains empty string")
		ctx.String(http.StatusBadRequest, createUserFailMsg)
		return
	}
	err = createUserInternal(&newUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createUserFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, newUser)
}

func createSession(ctx *gin.Context) {
	sessUser := common.User{}
	err := ctx.Bind(&sessUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createSessionFailMsg)
		return
	}
	if common.IsEmpty(sessUser.Name, sessUser.Email) {
		common.LogError(logger).Println("contains empty string")
		ctx.String(http.StatusBadRequest, createSessionFailMsg)
		return
	}
	sess, err := createSessionInternal(&sessUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, createSessionFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, *sess)
}

func readUser(ctx *gin.Context) {
	var searchUser common.User
	err := ctx.Bind(&searchUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readUserFailMsg)
		return
	}
	if common.IsEmpty(searchUser.Email) && common.IsEmpty(searchUser.UuId) {
		common.LogError(logger).Println("need email or uuid for finding user")
		ctx.String(http.StatusBadRequest, readUserFailMsg)
		return
	}
	err = readUserInternal(&searchUser)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readUserFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, searchUser)
}

func readSession(ctx *gin.Context) {
	var searchSess common.Session
	err := ctx.Bind(&searchSess)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	if common.IsEmpty(searchSess.UuId) {
		common.LogError(logger).Println("need uuid for finding session")
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	err = readSessionInternal(&searchSess)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	ctx.JSON(http.StatusOK, searchSess)
}

func deleteSession(ctx *gin.Context) {
	var delSess common.Session
	err := ctx.Bind(&delSess)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	if common.IsEmpty(delSess.UuId) {
		common.LogError(logger).Println("need uuid for finding session")
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	err = deleteSessionInternal(&delSess)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		ctx.String(http.StatusBadRequest, readSessionFailMsg)
		return
	}
	ctx.String(http.StatusOK, "deleted")
}

// don't forget TLS
func createUserInternal(newUser *common.User) (err error) {
	newUser.UuId = common.NewUuIdString()
	newUser.Password = common.Encrypt(newUser.Password)
	newUser.CreatedAt = time.Now()
	affected, err := dbEngine.
		Table(UserTable).
		InsertOne(newUser)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func createSessionInternal(sessUser *common.User) (session *common.Session, err error) {
	session = &common.Session{
		UuId:      common.NewUuIdString(),
		UserName:  sessUser.Name,
		UserId:    sessUser.Id,
		CreatedAt: time.Now(),
	}
	affected, err := dbEngine.
		Table(SessionTable).
		InsertOne(session)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something's wrong. returned value was %d",
			affected,
		)
	}
	return
}

func readUserInternal(searchUser *common.User) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(UserTable).
		Get(searchUser)
	if err == nil && !ok {
		err = errors.New("no such users")
	}
	return
}

func readSessionInternal(searchSess *common.Session) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(SessionTable).
		Get(searchSess)
	if err == nil && !ok {
		err = errors.New("no such session")
	}
	return
}

func deleteSessionInternal(delSess *common.Session) (err error) {
	affected, err := dbEngine.
		Table(SessionTable).
		Delete(delSess)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something's wrong. returned value was %d",
			affected,
		)
	}
	return
}
