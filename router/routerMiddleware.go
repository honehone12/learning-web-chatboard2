package main

import (
	"errors"
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
	if err != nil {
		common.LogWarning(logger).Println(err.Error())
	}
	ctx.Set(loggedInLabel, err == nil)
	ctx.Next()
}

func checkLoggedIn(ctx *gin.Context) (err error) {
	uuid, err := pickupSessionCookie(ctx)
	if err != nil {
		return
	}
	sess := common.Session{
		UuId: uuid,
	}
	req, err := common.MakeRequestFromSession(
		&sess,
		http.MethodPost,
		buildHTTP_URL(config.AddressUsers, "/check-session"),
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
	sessPtr, err := common.MakeSessionFromResponse(res)
	if err != nil {
		return
	}
	// update session

	// if time.Since(sessPtr.LastUpdate) > time.Second*1 {
	// 	storeSessionCookie(ctx, sess.UuId)
	// 	requestSessionUpdate(sessPtr)
	// }
	ctx.Set(sessionPtrLabel, sessPtr)
	return
}

func confirmLoggedIn(ctx *gin.Context) (isLoggedIn bool) {
	loggedInVal, ok := ctx.Get(loggedInLabel)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
		return
	}
	isLoggedIn, ok = loggedInVal.(bool)
	if !ok {
		common.LogError(logger).Fatalln("middleware not working")
	}
	return
}

func getSessionPtr(ctx *gin.Context) (ptr *common.Session, err error) {
	val, ok := ctx.Get(sessionPtrLabel)
	if !ok {
		err = errors.New("not logged in")
		return
	}
	if ptr, ok = val.(*common.Session); !ok {
		common.LogError(logger).Fatalln("middleware not working")
		err = errors.New("middleware not working")
	}
	return
}
