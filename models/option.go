package models

import "strings"

// Option 配置项
type Option struct {
	Name   string `xorm:"pk"`
	UserID int64  `xorm:"pk"`
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
	Options.OptionMap["site_url"] = "http://127.0.0.1:4000"

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

// SiteURL 配置的网站url
func SiteURL() string {
	return strings.Trim(Options.Get("site_url"), "/")
}

// AssetURL 生成链接地址
func AssetURL(s string) string {

	return SiteURL() + "/" + s
}
