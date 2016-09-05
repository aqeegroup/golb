package admin

import (
	"log"

	"gopkg.in/macaron.v1"
)

// Login 后台管理登录页面
func Login(ctx *macaron.Context) {
	ctx.HTMLSet(200, "admin", "login")
}

// DoLogin 登录请求处理
func DoLogin(ctx *macaron.Context) {

	log.Println(ctx.Req.FormValue("username"))
}
