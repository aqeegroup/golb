package admin

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

// Index 后台管理首页
func Index(ctx *macaron.Context, sess session.Store) {
	uid := sess.Get("uid")

	if uid == nil {
		ctx.Redirect("/admin/login")
	}

	ctx.HTMLSet(200, "admin", "home")
}
