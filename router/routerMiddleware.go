package main

import (
	"errors"
	"fmt"
	"learning-web-chatboard2/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	LoggedIn   = "logged-in"
	SessionPtr = "session-ptr"
)

func LoggedInCheckerMiddleware(ctx *gin.Context) {
	err := checkLoggedIn(ctx)
	ctx.Set(LoggedIn, err == nil)
	ctx.Next()
}

func checkLoggedIn(ctx *gin.Context) (err error) {
	uuid, err := ctx.Cookie(ShortTimeSession)
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
			Http,
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
	ctx.Set(SessionPtr, sess)
	return
}

func ConfirmLoggedIn(ctx *gin.Context) (loggedIn bool) {
	loggedInVal, ok := ctx.Get(LoggedIn)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
	}
	loggedIn, ok = loggedInVal.(bool)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
	}
	return
}

func GetSessionPtr(ctx *gin.Context) (ptr *common.Session, err error) {
	val, ok := ctx.Get(SessionPtr)
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
