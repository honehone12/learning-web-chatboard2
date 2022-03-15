package common

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	Id        uint      `xorm:"pk autoincr 'id'" json:"id"`
	UuId      string    `xorm:"not null unique 'uu_id'" json:"uuid"`
	Name      string    `xorm:"not null unique 'name'" json:"name"`
	Email     string    `xorm:"not null unique 'email'" json:"email"`
	Password  string    `xorm:"not null 'password'" json:"password"`
	CreatedAt time.Time `xorm:"not null 'created_at'" json:"created_at"`
}

type Session struct {
	Id        uint      `xorm:"pk autoincr 'id'" json:"id"`
	UuId      string    `xorm:"not null unique 'uu_id'" json:"uuid"`
	UserName  string    `xorm:"user_name" json:"user_name"`
	UserId    uint      `xorm:"user_id" json:"user_id"`
	CreatedAt time.Time `xorm:"not null 'created_at'" json:"created_at"`
}

func MakeRequestFromUser(
	user *User,
	method string,
	addr string,
) (req *http.Request, err error) {
	bin, err := json.Marshal(user)
	if err != nil {
		return
	}
	req, err = http.NewRequest(
		method,
		addr,
		bytes.NewBuffer(bin),
	)
	if err != nil {
		return
	}
	req.Header.Add(
		"Content-Type",
		"application/json",
	)
	return
}

func MakeUserFromResponse(res *http.Response) (user *User, err error) {
	user = &User{}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err = json.Unmarshal(body, user); err != nil {
		log.Fatalln(err.Error())
	}
	return
}
