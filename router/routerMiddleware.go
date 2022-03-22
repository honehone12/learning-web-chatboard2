package main

import (
	"errors"
	"fmt"
	"learning-web-chatboard2/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	loggedInLabel   = "logged-in"
	sessionPtrLabel = "session-ptr"
)

func LoggedInCheckerMiddleware(ctx *gin.Context) {
	err := checkLoggedIn(ctx)
	ctx.Set(loggedInLabel, err == nil)
	ctx.Next()
}

func checkLoggedIn(ctx *gin.Context) (err error) {
	uuid, err := pickupSessionCookie(ctx)
	if err != nil {
		return
	}
	sess := &common.Session{
		UuId: uuid,
	}
	req, err := common.MakeRequestFromSession(
		sess,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			httpPrefix,
			config.AddressUsers,
			"/check-session",
		),
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	sess, err = common.MakeSessionFromResponse(res)
	if err != nil {
		return
	}
	ctx.Set(sessionPtrLabel, sess)
	return
}

func ConfirmLoggedIn(ctx *gin.Context) (isLoggedIn bool) {
	loggedInVal, ok := ctx.Get(loggedInLabel)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
	}
	isLoggedIn, ok = loggedInVal.(bool)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
	}
	return
}

func GetSessionPtr(ctx *gin.Context) (ptr *common.Session, err error) {
	val, ok := ctx.Get(sessionPtrLabel)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
		err = errors.New("middleware not working")
	}
	if ptr, ok = val.(*common.Session); !ok {
		common.LogError(logger).Fatalln("middleware not working")
		err = errors.New("middleware not working")
	}
	return
}
