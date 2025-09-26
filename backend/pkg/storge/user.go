package storage

import "time"

type User struct {
	Id       int       `xorm:"pk autoincr 'user_id'"`
	Name     string    `xorm:"'user_name' not null"`
	Password string    `xorm:"'user_password' not null"`
	CreateAt time.Time `xorm:"'create_at' not null"`
	// TODO more info
}

// TODO SELECT methods for user

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
