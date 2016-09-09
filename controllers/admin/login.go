package admin

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"blog/models"
	"blog/modules/context"
	"blog/modules/setting"
	"blog/modules/utility"
)

// Login 后台管理登录页面
func Login(ctx *context.Context, sess session.Store) {
	if uid := sess.Get("uid"); uid != nil {
		ctx.Redirect("/admin")
		return
	}
	ctx.HTMLSet(200, "admin", "login")
	return
}

// DoLogin 登录请求处理
func DoLogin(ctx *context.Context, sess session.Store) {
	username := ctx.Req.FormValue("username")
	password := ctx.Req.FormValue("password")
	auto := ctx.Req.FormValue("auto")

	user, err := models.FindUserByUsername(username)

	if err != nil {
		ctx.Error(500)
		return
	}

	salt := setting.Cfg.Section("security").Key("salt").String()
	password = utility.Md5Encrypt(password + salt)

	if password != user.Password {
		ctx.Data["Title"] = "登录出错"
		ctx.Data["Content"] = "用户名或密码错误."
		ctx.HTMLSet(200, "admin", "error")
		return
	}

	// 登录信息记录
	user.LastIP = ctx.RemoteIP()

	user.Count++
	if err := models.UpdateUser(user); err != nil {
		ctx.Data["Title"] = "登录出错"
		ctx.Data["Content"] = "登录发生错误: " + err.Error()
		ctx.Error(500)
		return
	}

	// session 写入
	sess.Set("uid", user.ID)
	sess.Set("username", user.Username)

	// cookie 写入
	if auto == "on" {
		// 这个时间之后要可配置
		ctx.SetCookie(setting.CookieName, sess.ID(), 7*24*3600)
	}

	ctx.Redirect("/admin")
}

// Logout 退出登录
func Logout(ctx *macaron.Context, sess session.Store) {
	sess.Destory(ctx)
	ctx.Redirect("/admin/login")
}

// CheckLogin 检查登录
func CheckLogin(ctx *context.Context, sess session.Store) {
	if uid := sess.Get("uid"); uid == nil {
		ctx.Redirect("/admin/login")
		return
	}
}
