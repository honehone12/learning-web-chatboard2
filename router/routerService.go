package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"learning-web-chatboard2/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	publicNavbar template.HTML = `<div class="navbar navbar-default navbar-static-top" role="navigation">
  <div class="container">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">KEIJIBAN</a>
    </div>
    <div class="nav navbar-nav navbar-right">
      <a href="/user/login">Login</a>
    </div>
  </div>
</div>`

	privateNavbar template.HTML = `<div class="navbar navbar-default navbar-static-top" role="navigation">
  <div class="container">
    <div class="navbar-header">
	  <a class="navbar-brand" href="/">KEIJIBAN</a>
    </div>
    <div class="nav navbar-nav navbar-right">
	  <a href="/user/logout">Logout</a>
    </div>
  </div>
</div>`
)

const Http = "http://"

const (
	readThreadFailMsg   = "failed to read threads"
	signupFailMsg       = "failed to sign-up"
	authenticateFailMsg = "failed to authenticate"
)

func indexGet(ctx *gin.Context) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/read-index",
		),
		nil,
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, readThreadFailMsg)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}
	thres := make([]common.Thread, 0)
	err = json.Unmarshal(body, &thres)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}

	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"navbar":  publicNavbar,
			"threads": thres,
		},
	)
}

func errorRedirect(ctx *gin.Context, msg string) {
	ctx.Redirect(
		http.StatusFound,
		fmt.Sprintf(
			"%s%s",
			"/error?msg=",
			msg,
		),
	)
}

func errorGet(ctx *gin.Context) {
	errMsg := ctx.Query("msg")
	ctx.HTML(
		http.StatusOK,
		"error.html",
		gin.H{
			"navbar": publicNavbar,
			"msg":    errMsg,
		},
	)
}

func loginGet(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"login.html",
		nil,
	)
}

func signupGet(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"signup.html",
		nil,
	)
}

func logoutGet(ctx *gin.Context) {

}

func signupPost(ctx *gin.Context) {
	newUser := common.User{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	req, err := common.MakeRequestFromUser(
		&newUser,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressUsers,
			"/signup-account",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, signupFailMsg)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, signupFailMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, signupFailMsg)
		return
	}
	ctx.Redirect(http.StatusFound, "/user/login")
}

func authenticatePost(ctx *gin.Context) {
	authUser := common.User{
		Email: ctx.PostForm("email"),
	}
	req, err := common.MakeRequestFromUser(
		&authUser,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressUsers,
			"/authenticate",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, authenticateFailMsg)
		return
	}
	authedUser, err := common.MakeUserFromResponse(res)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
	}
	pass := common.Encrypt(ctx.PostForm("password"))
	if strings.Compare(authedUser.Password, pass) != 0 {
		common.LogError(logger).Println("password mismatch")
		errorRedirect(ctx, authenticateFailMsg)
		return
	}

	req, err = common.MakeRequestFromUser(
		authedUser,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressUsers,
			"/create-session",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
		return
	}
	session, err := common.MakeSessionFromResponse(res)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, authenticateFailMsg)
		return
	}

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		"short-time",
		session.UuId,
		0,
		"/",
		"localhost",
		true,
		true,
	)
	ctx.Redirect(http.StatusFound, "/")
}
