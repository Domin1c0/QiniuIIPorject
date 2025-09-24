package storage

import "time"

type Session struct {
	UserId   int       `xorm:"pk 'user_id'" json:"-"`
	Id       int       `xorm:"pk autoincr 'session_id'" json:"-"`
	CreateAt time.Time `xorm:"'create_at' not null" json:"create_at"`
	Messages []Message `xorm:"-" json:"messages"`
}

type Message struct {
	Session_id int    `xorm:"session_id" json:"-"`
	Id         int    `xorm:"pk autoincr 'message_id'" json:"-"`
	Role       string `xorm:"'role' not null" json:"role"`
	Content    string `xorm:"'content' not null" json:"content"`
	createAt   int64  `xorm:"'create_at' not null" json:"create_at"`
	// TODO store message audio
}
