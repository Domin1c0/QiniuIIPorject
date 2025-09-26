package storage

type User struct {
	Id       int    `xorm:"pk autoincr 'user_id'"`
	Name     string `xorm:"'user_name' not null"`
	Password string `xorm:"'user_password' not null"`
	CreateAt int64  `xorm:"'create_at' not null"`
	// TODO more info
}

// type
