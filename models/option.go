package models

// Option 配置项
type Option struct {
	Name   string `xorm:"pk"`
	UserID int    `xorm:"pk"`
	Value  string
}

// UserOptions 用户配置
type UserOptions struct {
	UserOptions *[]Option
	OptionMap   map[string]string
}

// Options 用户配置
var Options *UserOptions

// LoadOptions 从数据库加载个人配置
func LoadOptions() error {
	Options = &UserOptions{}
	Options.UserOptions = &[]Option{}
	Options.OptionMap = make(map[string]string)

	Options.OptionMap["theme"] = "templates/themes/default"

	// return x.Find(UserConfig)

	return nil
}

// Get 获取配置
func (o *UserOptions) Get(key string) string {

	if v, ok := o.OptionMap[key]; ok {
		return v
	}

	return ""
}
