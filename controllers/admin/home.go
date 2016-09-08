package admin

import (
	"github.com/go-macaron/session"

	"blog/models"
	"blog/modules/context"
)

// Index 后台管理首页
func Index(ctx *context.Context, sess session.Store) {
	uid := sess.Get("uid")
	postCount, err := models.CountPost()
	// postCount, err := models.CountCate()
	// postCount, err := models.CountComment()
	if err != nil {
		ctx.Handle(500, "Internal Server Error", err)
		return
	}

	ctx.Data["Uid"] = uid
	ctx.Data["Username"] = sess.Get("username")
	ctx.Data["Title"] = "网站概要"
	ctx.Data["PostCount"] = postCount
	ctx.HTMLSet(200, "admin", "index")
}

// NotFound 404
func NotFound(ctx *context.Context) {
	ctx.Handle(404, "Page Not Found", nil)
	return
}
