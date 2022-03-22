package common

import (
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
	State     string    `xorm:"TEXT 'state'" json:"state"`
	CreatedAt time.Time `xorm:"not null 'created_at'" json:"created_at"`
}

type Thread struct {
	Id         uint      `xorm:"pk autoincr 'id'" json:"id"`
	UuId       string    `xorm:"not null unique 'uu_id'" json:"uuid"`
	Topic      string    `xorm:"TEXT 'topic'" json:"topic"`
	NumReplies uint      `xorm:"num_replies" json:"num_replies"`
	Owner      string    `xorm:"owner" json:"owner"`
	UserId     uint      `xorm:"user_id" json:"user_id"`
	LastUpdate time.Time `xorm:"not null 'last_update'" json:"last_update"`
	CreatedAt  time.Time `xorm:"not null 'created_at'" json:"created_at"`
}

type Post struct {
	Id          uint      `xorm:"ok autoincr 'id'" json:"id"`
	UuId        string    `xorm:"not null unique 'uu_id'" json:"uuid"`
	Body        string    `xorm:"TEXT 'body'" json:"body"`
	Contributor string    `xorm:"contributor" json:"contributor"`
	UserId      uint      `xorm:"user_id" json:"user_id"`
	ThreadId    uint      `xorm:"thread_id" json:"thread_id"`
	CreatedAt   time.Time `xorm:"not null 'created_at'" json:"created_at"`
}

func (thread *Thread) When() string {
	return thread.CreatedAt.Format("2006/Jan/2 at 3:04pm")
}

func (post *Post) When() string {
	return post.CreatedAt.Format("2006/Jan/2 at 3:04pm")
}
