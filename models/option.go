package models

// Option 配置项
type Option struct {
	Name   string `xorm:"pk"`
	UserID int    `xorm:"pk"`
	Value  string
}

// UserConfig 用户配置
var UserConfig *[]Option

// LoadUserConfig 从数据库加载个人配置
func LoadUserConfig() (*[]Option, error) {
	UserConfig := &[]Option{}

	err := x.Find(UserConfig)
	return UserConfig, err
}
