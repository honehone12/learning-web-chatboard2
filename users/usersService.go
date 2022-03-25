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
	visitTable   = "visits"
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

func createVisit(ctx *gin.Context) {
	var newVis common.Visit
	err := createVisitInternal(ctx, &newVis)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &newVis)
}

func createVisitInternal(ctx *gin.Context, newVis *common.Visit) (err error) {
	now := time.Now()
	newVis.UuId = common.NewUuIdString()
	newVis.CreatedAt = now

	err = createVisitSQLInternal(newVis)
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

func readVisit(ctx *gin.Context) {
	var searchVis common.Visit
	err := readVisitInternal(ctx, &searchVis)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &searchVis)
}

func readVisitInternal(ctx *gin.Context, searchVis *common.Visit) (err error) {
	err = ctx.Bind(searchVis)
	if err != nil {
		return
	}
	if common.IsEmpty(searchVis.UuId) {
		err = errors.New("need uuid for finding visit")
		return
	}
	err = readVisitSQLInternal(searchVis)
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
	) {
		err = fmt.Errorf("contains empty string %s %s", sess.UuId, sess.UserName)
		return
	}
	sess.LastUpdate = time.Now()
	err = updateSessionSQLInternal(sess)
	return
}

func updateVisit(ctx *gin.Context) {
	var vis common.Visit
	err := updateVisitInternal(ctx, &vis)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, &vis)
}

func updateVisitInternal(ctx *gin.Context, vis *common.Visit) (err error) {
	err = ctx.Bind(vis)
	if err != nil {
		return
	}
	if common.IsEmpty(
		vis.UuId,
	) {
		err = fmt.Errorf("contains empty string %s %s", vis.UuId, vis.State)
		return
	}
	err = updateVisitSQLInternal(vis)
	return
}

func deleteSession(ctx *gin.Context) {
	var delSess common.Session
	err := deleteSessionInternal(ctx, &delSess)
	if err != nil {
		handleErrorInternal(err.Error(), ctx)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"deleted": "ok",
	})
}

func deleteSessionInternal(ctx *gin.Context, delSess *common.Session) (err error) {
	err = ctx.Bind(delSess)
	if err != nil {
		return
	}

	err = deleteSessionSQLInternal(delSess)
	return
}

func createVisitSQLInternal(newVis *common.Visit) (err error) {
	affected, err := dbEngine.
		Table(visitTable).
		InsertOne(newVis)
	if err == nil && affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
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

func readVisitSQLInternal(searchVis *common.Visit) (err error) {
	var ok bool
	ok, err = dbEngine.
		Table(visitTable).
		Get(searchVis)
	if err == nil && !ok {
		err = errors.New("no such viz")
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

func updateVisitSQLInternal(vis *common.Visit) (err error) {
	affected, err := dbEngine.
		Table(visitTable).
		ID(vis.Id).
		Update(vis)
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

	common.LogInfo(logger).Printf("deleted %d linked %v", affected, *delSess)
	return
}
