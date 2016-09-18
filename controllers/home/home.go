package home

import (
	"blog/models"
	"blog/modules/context"
)

// Index 首页
func Index(ctx *context.Context) {

	page := ctx.ParamsInt("page")
	if page < 1 {
		page = 1
	}
	// 这个参数之后读配置
	limit := 10

	posts, err := models.FindPosts(page, limit)

	if err != nil {
		ctx.Data["Title"] = "查询出错"
		ctx.Data["Content"] = "查询出错: " + err.Error()
		ctx.HTML(500, "error")
		return
	}
	ctx.Data["Posts"] = posts

	p, err := models.PostsPagination(page, limit, false)
	if err != nil {
		ctx.Data["Title"] = "查询出错"
		ctx.Data["Content"] = "查询出错: " + err.Error()
		ctx.HTML(500, "error")
		return
	}
	ctx.Data["Page"] = p

	ctx.Data["Title"] = "Jiayx 的博客"

	ctx.HTML(200, "index")
	return

}
