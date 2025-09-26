package storage

import (
	"time"
)

type Session struct {
	UserId   int       `xorm:"'user_id'" json:"-"`
	Id       int       `xorm:"pk autoincr 'session_id'" json:"-"`
	CreateAt time.Time `xorm:"'create_at' not null" json:"create_at"`
	UpdateAt time.Time `xorm:"'update_at'" json:"update_at"`
}

type Message struct {
	Session_id int       `xorm:"session_id" json:"-"`
	Id         int       `xorm:"pk autoincr 'message_id'" json:"-"`
	Role       string    `xorm:"'role' not null" json:"role"`
	Content    string    `xorm:"'content' not null" json:"content"`
	CreateAt   time.Time `xorm:"'create_at' not null" json:"create_at"`
	// TODO store message audio
}

// SELECT

func (db *Storage) GetSessionsIDByUserID(userID int, num int, from ...time.Time) ([]int, error) {
	var res []int

	q := db.engine.Table(new(Session)).
		Where("user_id = ?", userID).
		Cols("session_id").
		Desc("update_at")

	if len(from) > 0 {
		q = q.Where("update_at >= ?", from[0])
	}

	if num > 0 {
		q = q.Limit(num)
	}

	err := q.Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *Storage) GetSessionsByUserID(userID int, num int, from ...time.Time) ([]Session, error) {
	var res []Session

	q := db.engine.Where("user_id = ?", userID).Desc("update_at")

	if len(from) > 0 {
		q = q.Where("update_at >= ?", from[0])
	}

	if num > 0 {
		q = q.Limit(num)
	}

	err := q.Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *Storage) GetSessionByID(sessionID int) (Session, error) {
	var res Session
	has, err := db.engine.Where("session_id = ?", sessionID).Get(&res)
	if err != nil {
		return res, err
	}
	if !has {
		return res, ErrSessionNotFound
	}
	return res, nil
}

func (db *Storage) GetMessagesBySessionID(sessionID int) ([]Message, error) {
	var res []Message
	err := db.engine.Where("session_id = ?", sessionID).Asc("create_at").Find(&res)
	if err != nil {
		return nil, err
	}
	// Check session existence skipped
	return res, nil
}

// INSERT

// func (db *Storage) AddSession(userID int) (Session, error) {
func (db *Storage) AddSession(session Session) (Session, error) {
	if session.UserId == 0 {
		return session, ErrInvalidUserID
	}
	if session.CreateAt.IsZero() {
		session.CreateAt = time.Now()
	}
	session.UpdateAt = session.CreateAt
	_, err := db.engine.Insert(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

// func (db *Storage) AddMessage(sessionID int, role, content string) (Message, error) {
func (db *Storage) AddMessage(message Message) (Message, error) {
	if message.CreateAt.IsZero() {
		message.CreateAt = time.Now()
	}
	_, err := db.engine.Insert(&message)
	if err != nil {
		return message, err
	}
	// update session update_at
	_, err = db.engine.Where("session_id = ?", message.Session_id).Cols("update_at").Update(&Session{UpdateAt: time.Now()})
	if err != nil {
		return message, err
	}
	return message, nil
}
