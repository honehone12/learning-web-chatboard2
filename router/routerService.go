package main

import (
	"encoding/json"
	"errors"
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

func handleErrorInternal(
	loggerErrorMsg string,
	ctx *gin.Context,
	publicErrorMsg string,
) {
	common.LogError(logger).Println(loggerErrorMsg)
	errorRedirect(ctx, publicErrorMsg)
}

func getHTMLElemntInternal(isLoggedin bool) (template.HTML, template.HTML) {
	if isLoggedin {
		return privateNavbar, replyForm
	} else {
		return publicNavbar, ""
	}
}

func indexGet(ctx *gin.Context) {
	thres, err := indexGetInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to read thread")
		return
	}
	navbar, _ := getHTMLElemntInternal(ConfirmLoggedIn(ctx))
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"navbar":  navbar,
			"threads": thres,
		},
	)
}

func indexGetInternal(ctx *gin.Context) (threads []common.Thread, err error) {
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
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &threads)
	return
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
	navbar, _ := getHTMLElemntInternal(ConfirmLoggedIn(ctx))
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
		err := logoutGetInternal(ctx)
		if err != nil {
			handleErrorInternal(err.Error(), ctx, "failed to logout")
			return
		}
	}
	ctx.Redirect(http.StatusFound, "/")
}

func logoutGetInternal(ctx *gin.Context) (err error) {
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
		return
	}
	res, err := httpClient.Do(req)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}
	return
}

func signupPost(ctx *gin.Context) {
	err := signupPostInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to sign-up")
		return
	}
	ctx.Redirect(http.StatusFound, "/user/login")
}

func signupPostInternal(ctx *gin.Context) (err error) {
	pw := processPassword(ctx.PostForm("password"))
	newUser := common.User{
		Name:     ctx.PostForm("name"),
		Email:    ctx.PostForm("email"),
		Password: pw,
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
		return
	}
	res, err := httpClient.Do(req)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}
	return
}

func authenticatePost(ctx *gin.Context) {
	err := authenticatePostInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to authenticate")
		return
	}
	ctx.Redirect(http.StatusFound, "/")
}

func authenticatePostInternal(ctx *gin.Context) (err error) {
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
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	authedUser, err := common.MakeUserFromResponse(res)
	if err != nil {
		return
	}
	pw := processPassword(ctx.PostForm("password"))
	if strings.Compare(authedUser.Password, pw) != 0 {
		err = errors.New("password mismatch")
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
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	session, err := common.MakeSessionFromResponse(res)
	if err != nil {
		return
	}

	/////////////////////////////////////////////
	//session cookie
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
	return
}

func threadGet(ctx *gin.Context) {
	thre, posts, err := threadGetInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to read thread")
		return
	}
	navbar, reply := getHTMLElemntInternal(ConfirmLoggedIn(ctx))
	ctx.HTML(
		http.StatusOK,
		"thread.html",
		gin.H{
			"navbar": navbar,
			"thread": thre,
			"reply":  reply,
			"posts":  posts,
			"token":  "easy-token",
		},
	)
}

func threadGetInternal(ctx *gin.Context) (thread *common.Thread, posts []common.Post, err error) {
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
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	thread, err = common.MakeThreadFromResponse(res)
	if err != nil {
		return
	}

	req, err = common.MakeRequestFromThread(
		thread,
		http.MethodPost,
		fmt.Sprintf(
			"%s%s%s",
			Http,
			config.AddressThreads,
			"/read-posts",
		),
	)
	if err != nil {
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &posts)
	return
}

func newThreadGet(ctx *gin.Context) {
	loggedin := ConfirmLoggedIn(ctx)
	navbar, _ := getHTMLElemntInternal(loggedin)
	if loggedin {
		ctx.HTML(
			http.StatusOK,
			"newthread.html",
			gin.H{
				"navbar": navbar,
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

	err := newThreadPostInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to post thread")
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

func newThreadPostInternal(ctx *gin.Context) (err error) {
	sess, err := GetSessionPtr(ctx)
	if err != nil {
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
		return
	}
	res, err := httpClient.Do(req)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}
	return
}

func newReplyPost(ctx *gin.Context) {
	if !ConfirmLoggedIn(ctx) {
		ctx.Redirect(http.StatusFound, "/user/login")
		return
	}

	threUuId, err := newReplyPostInternal(ctx)
	if err != nil {
		handleErrorInternal(err.Error(), ctx, "failed to reply")
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprint("/thread/read?id=", threUuId))
}

func newReplyPostInternal(ctx *gin.Context) (threUuId string, err error) {
	////////////////////////////////////////////////
	// token must be changed everytime and
	// we have to remember token.
	token := ctx.PostForm("token")
	if strings.Compare(token, "easy-token") != 0 {
		ctx.Abort()
		return
	}

	sess, err := GetSessionPtr(ctx)
	if err != nil {
		return
	}

	body := ctx.PostForm("body")
	/////////////////////////////////////////
	// here means uuid and id are now public info
	// should be encrypted
	// or use (session) cookie
	threUuId = ctx.PostForm("uuid")
	threId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
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
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
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
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	threPtr, err := common.MakeThreadFromResponse(res)
	if err != nil {
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
		return
	}
	res, err = httpClient.Do(req)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}
	return
}
