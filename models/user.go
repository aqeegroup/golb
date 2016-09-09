package models

// User 用户
type User struct {
	ID         int64  `xorm:"pk autoincr"`
	Username   string `xorm:"varchar(50)"`
	Password   string `xorm:"varchar(32)"`
	CreateTime int64
	UpdateTime int64 `xorm:"updated"`
	Count      int
	LastIP     string `xorm:"varchar(50)"`
}

// FindUserByUsername 根据用户名查询用户
func FindUserByUsername(username string) (*User, error) {
	user := &User{}
	_, err := x.Where("username=?", username).Get(user)
	// if !has {
	// 	return nil, err
	// }

	return user, err
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	_, err := x.Update(user)
	return err
}
