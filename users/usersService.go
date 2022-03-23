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
	userTable    = "users"
	sessionTable = "sessions"
)

func handleErrorInternal(
	loggerErrorMsg string,
	ctx *gin.Context,
) {
	common.LogError(logger).Println(loggerErrorMsg)
	ctx.JSON(http.StatusBadRequest, gin.H{"status": "error"})
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
	newUser.UuId = common.NewUuIdString()
	newUser.CreatedAt = time.Now()
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
	now := time.Now()
	sess = &common.Session{
		UuId:       common.NewUuIdString(),
		UserName:   sessUser.Name,
		UserId:     sessUser.Id,
		LastUpdate: now,
		CreatedAt:  now,
	}
	err = createSessionSQLInternal(sess)
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
		return
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

func updateSession(ctx *gin.Context) {
	var sess common.Session
	err := updateSessionInternal(ctx, &sess)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &sess)
}

func updateSessionInternal(ctx *gin.Context, sess *common.Session) (err error) {
	err = ctx.Bind(sess)
	if err != nil {
		return
	}
	if common.IsEmpty(
		sess.UuId,
		sess.UserName,
		sess.State,
	) {
		err = errors.New("contains empty string")
		return
	}
	sess.LastUpdate = time.Now()
	err = updateSessionSQLInternal(sess)
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
	affected, err := dbEngine.
		Table(userTable).
		InsertOne(newUser)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func createSessionSQLInternal(session *common.Session) (err error) {
	affected, err := dbEngine.
		Table(sessionTable).
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
		Table(userTable).
		Get(searchUser)
	if err == nil && !ok {
		err = errors.New("no such users")
	}
	return
}

func readSessionSQLInternal(searchSess *common.Session) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(sessionTable).
		Get(searchSess)
	if err == nil && !ok {
		err = errors.New("no such session")
	}
	return
}

func updateSessionSQLInternal(session *common.Session) (err error) {
	affected, err := dbEngine.
		Table(sessionTable).
		ID(session.Id).
		Update(session)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}

func deleteSessionSQLInternal(delSess *common.Session) (err error) {
	affected, err := dbEngine.
		Table(sessionTable).
		Delete(delSess)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something's wrong. returned value was %d",
			affected,
		)
	}
	return
}
