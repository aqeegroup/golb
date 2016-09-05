package home

import (
	"gopkg.in/macaron.v1"

	"blog/models"
)

// Index 首页
func Index(ctx *macaron.Context) {

	posts, err := models.FindPosts(1, 10)

	if err != nil {
		ctx.Data["Title"] = "查询出错"
		ctx.Data["Content"] = "查询出错: " + err.Error()

		ctx.HTML(500, "error")
		return
	}

	ctx.Data["Title"] = "Jiayx 的博客"
	ctx.Data["Posts"] = posts

	ctx.HTML(200, "index")
	return

}
