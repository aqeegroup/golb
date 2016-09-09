package admin

import (
	"regexp"
	"strings"
	"time"

	"blog/models"
	"blog/modules/context"
	"blog/modules/utility"
)

// PostSubmit 提交文章
func PostSubmit(ctx *context.Context) {
	resp := models.RespJSON{}

	post := &models.Post{}

	post.Slug = strings.Trim(ctx.Req.PostFormValue("slug"), " ")
	if len(post.Slug) > 0 {
		r := regexp.MustCompile("^[\\w-]$")
		if !r.MatchString(post.Slug) {
			ctx.RespJSON("缩略名格式有误 只能包含大小写字母和_-")
			return
		}

	}
	post.Title = strings.Trim(ctx.Req.PostFormValue("title"), " ")
	post.Content = ctx.Req.PostFormValue("content")
	post.Type = ctx.Req.PostFormValue("type")
	post.Status = ctx.Req.PostFormValue("status")
	publishTime := strings.Trim(ctx.Req.PostFormValue("publish_time"), " ")

	if len(post.Title) == 0 {
		post.Title = "未命名"
	}
	post.CreateTime = time.Now().Unix()
	post.UpdateTime = post.CreateTime
	if len(publishTime) > 0 {
		post.PublishTime = utility.Str2Int64(publishTime)
	} else {
		post.PublishTime = post.CreateTime
	}
	post.AuthorID = ctx.Session.Get("uid").(int64)

	err := post.Create()
	if err != nil {
		resp.Code = 401
		resp.Msg = err.Error()
		ctx.JSON(200, resp)
		return
	}

	resp.Code = 200
	resp.Msg = "发布成功"
	resp.Redirect = ctx.URLFor("home")
	ctx.JSON(200, resp)
	return

}

// WritePage 写文章的页面
func WritePage(ctx *context.Context) {
	ctx.Data["HideSidebar"] = true
	ctx.HTMLSet(200, "admin", "post")
}
