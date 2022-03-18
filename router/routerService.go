package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"learning-web-chatboard2/common"
	"net/http"
	"strconv"
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

	replyForm template.HTML = `<div class="panel panel-info">
  <div class="panel-body">
    <form id="post" role="form" action="/thread/post" method="post">
	  <div class="form-group">
	    <textarea class="form-control" name="body" id="body" placeholder="Write your reply here" rows="3"></textarea>
	     <!-- get url with javascript? <input type="hidden" name="uuid" value=""> -->
	     <br/>
	     <button class="btn btn-primary pull-right" type="submit">Reply</button>
	  </div>
    </form>
  </div>
</div>`
)

const (
	Http             = "http://"
	ShortTimeSession = "short-time"
)

const (
	readThreadFailMsg   = "failed to read thread"
	postThreadFailMsg   = "failed to post thread"
	signupFailMsg       = "failed to sign-up"
	authenticateFailMsg = "failed to authenticate"
	logoutFailMsg       = "failed to logout"
	replyPostFaileMsg   = "failed to reply"
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

	var navbar template.HTML
	if ConfirmLoggedIn(ctx) {
		navbar = privateNavbar
	} else {
		navbar = publicNavbar
	}
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"navbar":  navbar,
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
	var navbar template.HTML
	if ConfirmLoggedIn(ctx) {
		navbar = privateNavbar
	} else {
		navbar = publicNavbar
	}
	ctx.HTML(
		http.StatusOK,
		"error.html",
		gin.H{
			"navbar": navbar,
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
	if ConfirmLoggedIn(ctx) {
		uuid, _ := ctx.Cookie(ShortTimeSession)
		sess := &common.Session{UuId: uuid}
		req, err := common.MakeRequestFromSession(
			sess,
			http.MethodPost,
			fmt.Sprintf(
				"%s%s%s",
				Http,
				config.AddressUsers,
				"/delete-session",
			),
		)
		if err != nil {
			common.LogError(logger).Println(err.Error())
			errorRedirect(ctx, logoutFailMsg)
			return
		}
		res, err := httpClient.Do(req)
		if err != nil {
			common.LogError(logger).Println(err.Error())
			errorRedirect(ctx, logoutFailMsg)
			return
		} else if res.StatusCode != http.StatusOK {
			common.LogError(logger).Println(res.Status)
			errorRedirect(ctx, logoutFailMsg)
			return
		}
	}
	ctx.Redirect(http.StatusFound, "/")
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
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
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
		ShortTimeSession,
		session.UuId,
		0,
		"/",
		"localhost",
		true,
		true,
	)
	ctx.Redirect(http.StatusFound, "/")
}

func threadGet(ctx *gin.Context) {
	uuid := ctx.Query("id")
	thre := common.Thread{UuId: uuid}
	req, err := common.MakeRequestFromThread(
		&thre,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/read",
		),
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
	threPtr, err := common.MakeThreadFromResponse(res)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}

	req, err = common.MakeRequestFromThread(
		threPtr,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/read-posts",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}
	res, err = httpClient.Do(req)
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
	posts := make([]common.Post, 0)
	err = json.Unmarshal(body, &posts)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, readThreadFailMsg)
		return
	}

	var navbar template.HTML
	var reply template.HTML
	if ConfirmLoggedIn(ctx) {
		navbar = privateNavbar
		reply = replyForm
	} else {
		navbar = publicNavbar
		reply = ""
	}
	ctx.HTML(
		http.StatusOK,
		"thread.html",
		gin.H{
			"navbar": navbar,
			"thread": threPtr,
			"reply":  reply,
			"posts":  posts,
			"token":  "easy-token",
		},
	)
}

func newThreadGet(ctx *gin.Context) {
	if ConfirmLoggedIn(ctx) {
		ctx.HTML(
			http.StatusOK,
			"newthread.html",
			gin.H{
				"navbar": privateNavbar,
			},
		)
	} else {
		ctx.Redirect(http.StatusFound, "/user/login")
	}
}

func newThreadPost(ctx *gin.Context) {
	if !ConfirmLoggedIn(ctx) {
		ctx.Redirect(http.StatusFound, "/user/login")
		return
	}

	sess, err := GetSessionPtr(ctx)
	if err != nil {
		errorRedirect(ctx, postThreadFailMsg)
		return
	}

	thre := common.Thread{
		Topic:  ctx.PostForm("topic"),
		Owner:  sess.UserName,
		UserId: sess.UserId,
	}
	req, err := common.MakeRequestFromThread(
		&thre,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/create",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, postThreadFailMsg)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, postThreadFailMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, postThreadFailMsg)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

func newReplyPost(ctx *gin.Context) {
	if !ConfirmLoggedIn(ctx) {
		ctx.Redirect(http.StatusFound, "/user/login")
		return
	}

	// token must be changed everytime and
	// we have to remember token.
	token := ctx.PostForm("token")
	if strings.Compare(token, "easy-token") != 0 {
		ctx.Abort()
		return
	}

	sess, err := GetSessionPtr(ctx)
	if err != nil {
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}

	body := ctx.PostForm("body")
	/////////////////////////////////////////
	// here means uuid and id are now public info
	// should be encrypted
	// or use cookie
	threUuId := ctx.PostForm("uuid")
	threId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}

	post := common.Post{
		Body:        body,
		Contributor: sess.UserName,
		UserId:      sess.UserId,
		ThreadId:    uint(threId),
	}
	req, err := common.MakeRequestFromPost(
		&post,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/create-post",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}

	thre := common.Thread{UuId: threUuId}
	req, err = common.MakeRequestFromThread(
		&thre,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/read",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}
	threPtr, err := common.MakeThreadFromResponse(res)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}
	threPtr.NumReplies++

	req, err = common.MakeRequestFromThread(
		threPtr,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/update",
		),
	)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		common.LogError(logger).Println(err.Error())
		errorRedirect(ctx, replyPostFaileMsg)
		return
	} else if res.StatusCode != http.StatusOK {
		common.LogError(logger).Println(res.Status)
		errorRedirect(ctx, replyPostFaileMsg)
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprint("/thread/read?id=", threUuId))
}
