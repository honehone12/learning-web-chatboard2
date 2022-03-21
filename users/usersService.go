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

func handleErrorInternal(
	loggerErrorMsg string,
	ctx *gin.Context,
) {
	common.LogError(logger).Println(loggerErrorMsg)
	ctx.String(http.StatusBadRequest, "error")
}

// better way to send user data??
func createUser(ctx *gin.Context) {
	var newUser common.User
	err := createUserInternal(ctx, &newUser)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &newUser)
}

func createUserInternal(ctx *gin.Context, newUser *common.User) (err error) {
	err = ctx.Bind(newUser)
	if err != nil {
		return
	}
	if common.IsEmpty(
		newUser.Name,
		newUser.Email,
		newUser.Password,
	) {
		err = errors.New("contains empty string")
		return
	}
	err = createUserSQLInternal(newUser)
	return
}

func createSession(ctx *gin.Context) {
	var sessUser common.User
	sess, err := createSessionInternal(ctx, &sessUser)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, sess)
}

func createSessionInternal(ctx *gin.Context, sessUser *common.User) (sess *common.Session, err error) {
	err = ctx.Bind(sessUser)
	if err != nil {
		return
	}
	if common.IsEmpty(sessUser.Name, sessUser.Email) {
		err = errors.New("contains empty string")
		return
	}
	sess, err = createSessionSQLInternal(sessUser)
	return
}

func readUser(ctx *gin.Context) {
	var searchUser common.User
	err := readUserInternal(ctx, &searchUser)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &searchUser)
}

func readUserInternal(ctx *gin.Context, searchUser *common.User) (err error) {
	err = ctx.Bind(searchUser)
	if err != nil {
		return
	}
	if common.IsEmpty(searchUser.Email) && common.IsEmpty(searchUser.UuId) {
		err = errors.New("need email or uuid for finding user")
		return
	}
	err = readUserSQLInternal(searchUser)
	return
}

func readSession(ctx *gin.Context) {
	var searchSess common.Session
	err := readSessionInternal(ctx, &searchSess)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
	}
	ctx.JSON(http.StatusOK, &searchSess)
}

func readSessionInternal(ctx *gin.Context, searchSess *common.Session) (err error) {
	err = ctx.Bind(searchSess)
	if err != nil {
		return
	}
	if common.IsEmpty(searchSess.UuId) {
		err = errors.New("need uuid for finding session")
		return
	}
	err = readSessionSQLInternal(searchSess)
	return
}

func deleteSession(ctx *gin.Context) {
	var delSess common.Session
	err := deleteSessionInternal(ctx, &delSess)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
	}
	ctx.String(http.StatusOK, "deleted")
}

func deleteSessionInternal(ctx *gin.Context, delSess *common.Session) (err error) {
	err = ctx.Bind(delSess)
	if err != nil {
		return
	}
	if common.IsEmpty(delSess.UuId) {
		err = errors.New("need uuid for finding session")
		return
	}
	err = deleteSessionSQLInternal(delSess)
	return
}

func createUserSQLInternal(newUser *common.User) (err error) {
	newUser.UuId = common.NewUuIdString()
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

func createSessionSQLInternal(sessUser *common.User) (session *common.Session, err error) {
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

func readUserSQLInternal(searchUser *common.User) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(UserTable).
		Get(searchUser)
	if err == nil && !ok {
		err = errors.New("no such users")
	}
	return
}

func readSessionSQLInternal(searchSess *common.Session) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(SessionTable).
		Get(searchSess)
	if err == nil && !ok {
		err = errors.New("no such session")
	}
	return
}

func deleteSessionSQLInternal(delSess *common.Session) (err error) {
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
