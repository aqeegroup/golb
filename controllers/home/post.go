package home

import (
	"blog/models"
	"fmt"

	"gopkg.in/macaron.v1"
)

// Detail 查询文章详情
func Detail(ctx *macaron.Context) {

	slug := ctx.Params("slug")
	fmt.Println(slug)
	post, err := models.FindPostBySlug(slug)
	if err != nil {
		ctx.Data["Title"] = "error"
		ctx.Data["Content"] = err.Error()
		ctx.HTML(200, "404")
		return
	}

	ctx.Data["Title"] = post.Title
	ctx.Data["Post"] = post
	ctx.HTML(200, "detail")
	return
}
