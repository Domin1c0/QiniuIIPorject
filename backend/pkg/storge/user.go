package storage

import (
	"time"

	"github.com/LTSlw/QiniuIIPorject/backend/pkg/crypto"
)

type User struct {
	Id       int       `xorm:"pk autoincr 'user_id'" json:"id"`
	Name     string    `xorm:"'user_name' unique not null" json:"username"`
	Password string    `xorm:"'user_password' not null" json:"-"`
	CreateAt time.Time `xorm:"'create_at' not null" json:"create_at"`
	// TODO more info
}

type UserSession struct {
	UserId int    `xorm:"'user_id'" json:"user_id"`
	Id     string `xorm:"'session_id' pk unique not null" json:"token"` // as token
}

// TODO SELECT methods for user

func (db *Storage) GetUserByName(name string) (*User, error) {
	var user User
	has, err := db.engine.Where("user_name = ?", name).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (db *Storage) GetUserSessionByID(id string) (*UserSession, error) {
	var session UserSession
	has, err := db.engine.Where("session_id = ?", id).Get(&session)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &session, nil
}

// INSERT

func (db *Storage) AddUser(user *User) (User, error) {
	if user.CreateAt.IsZero() {
		user.CreateAt = time.Now()
	}
	_, err := db.engine.Insert(user)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

func (db *Storage) UserLogin(id int, password string) (UserSession, error) {
	token, err := crypto.GenerateAccessToken()
	if err != nil {
		return UserSession{}, err
	}
	session := UserSession{
		UserId: id,
		Id:     token,
	}
	if _, err = db.engine.Insert(&session); err != nil {
		return UserSession{}, err
	}
	return session, nil
}

// DELETE

func (db *Storage) UserLogout(session *UserSession) error {
	_, err := db.engine.Delete(session)
	return err
}
