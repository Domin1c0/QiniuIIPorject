package storage

import "time"

type Session struct {
	user_id  int       `xorm:"pk 'user_id'" json:"-"`
	id       int       `xorm:"pk autoincr 'session_id'" json:"-"`
	createAt time.Time `xorm:"'create_at' not null" json:"create_at"`
	Messages []Message `xorm:"-" json:"messages"`
}

type Message struct {
	session_id int    `xorm:"session_id" json:"-"`
	id         int    `xorm:"pk autoincr 'message_id'" json:"-"`
	role       string `xorm:"'role' not null" json:"role"`
	content    string `xorm:"'content' not null" json:"content"`
	createAt   int64  `xorm:"'create_at' not null" json:"create_at"`
	// TODO store message audio
}
