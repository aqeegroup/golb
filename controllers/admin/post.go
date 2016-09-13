package admin

import (
	"fmt"
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
		r := regexp.MustCompile("^[\\w-]+$")
		if !r.MatchString(post.Slug) {
			ctx.RespJSON("400", "缩略名格式有误 只能包含大小写字母、_和-")
			return
		}
	}

	publishTime := strings.Trim(ctx.Req.PostFormValue("publish_time"), " ")
	if len(publishTime) > 0 {
		r := regexp.MustCompile("^14\\d{8}$")
		if !r.MatchString(publishTime) {
			ctx.RespJSON("400", "发布日期格式有误")
			return
		}
	}

	post.Title = strings.Trim(ctx.Req.PostFormValue("title"), " ")
	post.Content = ctx.Req.PostFormValue("content")
	post.Type = ctx.Req.PostFormValue("type")
	post.Status = ctx.Req.PostFormValue("status")

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
		resp.Code = "401"
		resp.Msg = err.Error()
		ctx.JSON(200, resp)
		return
	}

	ctx.RespJSON("200", "发布成功", ctx.URLFor("home"))
	return
}

// WritePage 写文章的页面
func WritePage(ctx *context.Context) {
	ctx.Data["HideSidebar"] = true
	ctx.Data["Title"] = "文章管理"

	cates, err := models.FindAllCates()
	if err != nil {
		ctx.Handle(500, "", nil)
	}

	ctx.Data["Cates"] = cates

	ctx.Data["Scripts"] = []string{"admin/js/index.js"}

	ctx.HTMLSet(200, "admin", "post")
	return
}

// PostManage 管理文章页面
func PostManage(ctx *context.Context) {
	ctx.Data["Title"] = "文章管理"
	ctx.Data["PostActive"] = "active"
	ctx.Data["ManageActive"] = "active toggle"

	posts, err := models.FindPostsDetail(1, 10)
	fmt.Println(posts)
	if err != nil {
		ctx.Handle(500, "", err)
		return
	}
	ctx.Data["Posts"] = posts

	ctx.Data["Styles"] = []string{"admin/css/post_list.css"}
	ctx.Data["Scripts"] = []string{"admin/js/index.js", "admin/js/util.js"}
	ctx.HTMLSet(200, "admin", "post_list")
	return
}
