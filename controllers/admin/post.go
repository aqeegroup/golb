package admin

import (
	"strings"
	"time"

	"blog/models"
	"blog/modules/context"
)

// PostSubmit 提交文章
func PostSubmit(ctx *context.Context) {
	resp := models.RespJSON{}

	post := &models.Post{}
	post.Title = strings.Trim(ctx.Req.PostFormValue("title"), " ")
	post.Content = ctx.Req.PostFormValue("content")
	post.Type = ctx.Req.PostFormValue("type")
	post.Status = ctx.Req.PostFormValue("status")
	// post.Slug = strings.Trim(ctx.Req.PostFormValue("slug"), "")

	if len(post.Title) == 0 {
		post.Title = "未命名"
	}
	post.CreateTime = int(time.Now().Unix())
	post.UpdateTime = post.CreateTime

	id, err := post.Create()
	if err != nil {
		resp.Code = 401
		resp.Msg = err.Error()
		ctx.JSON(200, resp)
		return
	}

	resp.Code = 200
	resp.Msg = "发布成功"
	resp.Redirect = ctx.URLFor("home")
	resp.Data = make(map[string]interface{})
	resp.Data["id"] = id
	ctx.JSON(200, resp)
	return

}

// WritePage 写文章的页面
func WritePage(ctx *context.Context) {
	ctx.Data["HideSidebar"] = true
	ctx.HTMLSet(200, "admin", "post")
}
