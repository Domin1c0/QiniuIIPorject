package storage

type User struct {
	id       int    `xorm:"pk autoincr 'user_id'"`
	name     string `xorm:"'user_name' not null"`
	password string `xorm:"'user_password' not null"`
	createAt int64  `xorm:"'create_at' not null"`
	// TODO more info
}

// type
