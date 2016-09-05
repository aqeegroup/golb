package admin

import "gopkg.in/macaron.v1"

// Index 后台管理首页
func Index(ctx *macaron.Context) {
	ctx.HTMLSet(200, "admin", "home")
}
