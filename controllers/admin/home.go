package admin

import (
	"github.com/go-macaron/session"

	"blog/models"
	"blog/modules/context"
)

// Index 后台管理首页
func Index(ctx *context.Context, sess session.Store) {
	uid := sess.Get("uid")

	if uid == nil {
		ctx.Redirect("/admin/login")
	}

	postCount, err := models.CountPost()
	if err != nil {
		ctx.Handle(500, "Internal Server Error", err)
		return
	}

	ctx.Data["Uid"] = uid
	ctx.Data["Username"] = sess.Get("username")
	ctx.Data["PostCount"] = postCount
	ctx.HTMLSet(200, "admin", "home")
}

// NotFound 404
func NotFound(ctx *context.Context) {
	ctx.Handle(404, "Page Not Found", nil)
	return
}
