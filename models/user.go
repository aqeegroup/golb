package models

// User 用户
type User struct {
	ID         int    `xorm:"pk"`
	Username   string `xorm:"varchar(50)"`
	Password   string `xorm:"varchar(32)"`
	CreateTime int
	UpdateTime int `xorm:"updated"`
	Count      int
	LastIP     string `xorm:"varchar(50)"`
}

// FindUserByUsername 根据用户名查询用户
func FindUserByUsername(username string) (*User, error) {
	user := &User{}
	_, err := x.Where("username=?", username).Get(user)

	return user, err
}
