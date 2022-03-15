package main

import (
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

func CreateUser(ctx *gin.Context) {
	var newUser common.User
	err := ctx.Bind(&newUser)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if common.IsEmpty(
		newUser.Name,
		newUser.Email,
		newUser.Password,
	) {
		ctx.String(http.StatusBadRequest, "contains empty string")
		return
	}
	err = createUserInternal(&newUser)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, newUser)
}

//////////////////////////////////////
// never forget TLS
func createUserInternal(newUser *common.User) (err error) {
	newUser.UuId = common.NewUuIdString()
	newUser.Password = common.Encrypt(newUser.Password)
	newUser.CreatedAt = time.Now()
	affected, err := dbEngine.
		Table(UserTable).
		InsertOne(newUser)
	if err != nil {
		return
	} else if affected != 1 {
		err = fmt.Errorf(
			"something wrong. returned value was %d",
			affected,
		)
	}
	return
}
