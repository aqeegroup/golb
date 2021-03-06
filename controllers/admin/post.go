package admin

import (
	"fmt"
	"regexp"
	"time"

	"blog/models"
	"blog/modules/context"
	"blog/modules/utility"
)

// PostSubmit 提交文章
func PostSubmit(ctx *context.Context) {
	post := &models.Post{}

	post.ID = ctx.PostInt64("id")
	post.Slug = ctx.PostString("slug")
	if len(post.Slug) > 0 {
		r := regexp.MustCompile("^[\\w-]+$")
		if !r.MatchString(post.Slug) {
			ctx.RespJSON("400", "缩略名格式有误 只能包含大小写字母、_和-")
			return
		}
	}

	publishTime := ctx.PostString("publish_time")
	if len(publishTime) > 0 {
		r := regexp.MustCompile("^14\\d{8}$")
		if !r.MatchString(publishTime) {
			ctx.RespJSON("400", "发布日期格式有误")
			return
		}
	}

	post.Title = ctx.PostString("title", "未命名文档")
	post.Content = ctx.PostString("content")
	post.Type = ctx.PostString("type", "post")
	post.Status = ctx.PostString("status", "publish")

	post.CreateTime = time.Now().Unix()
	post.UpdateTime = post.CreateTime

	if len(publishTime) > 0 {
		post.PublishTime = utility.Str2Int64(publishTime)
	} else {
		post.PublishTime = post.CreateTime
	}
	post.AuthorID = ctx.Session.Get("uid").(int64)

	// 分类
	cates := ctx.PostString("cates")
	// 标签
	tags := ctx.PostString("tags")

	err := post.Create(cates, tags)
	if err != nil {
		ctx.RespJSON("500", "写入数据库出错"+err.Error())
		return
	}

	ctx.RespJSON("200", "发布成功", ctx.URLFor("postManage"))
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

	tags, err := models.FindAllTags()
	if err != nil {
		ctx.Handle(500, "", nil)
	}

	ctx.Data["Cates"] = cates
	ctx.Data["Tags"] = tags

	ctx.Data["Scripts"] = []string{"admin/js/index.js"}

	ctx.HTMLSet(200, "admin", "post")
	return
}

// PostUpdate 更新文章
func PostUpdate(ctx *context.Context) {
	id := ctx.ParamsInt64(":id")
	if id <= 0 {
		ctx.Handle(404, "404 Not Found.", nil)
		return
	}
	post, err := models.FindPostByID(id)
	if err != nil {
		ctx.Handle(500, "", nil)
		return
	}
	ctx.Data["Post"] = post
	fmt.Println(post)

	cates, err := models.FindAllCates()
	if err != nil {
		ctx.Handle(500, "", nil)
	}

	tags, err := models.FindAllTags()
	if err != nil {
		ctx.Handle(500, "", nil)
	}

	ctx.Data["Cates"] = cates
	ctx.Data["Tags"] = tags

	ctx.Data["HideSidebar"] = true
	ctx.Data["Title"] = "文章编辑"
	ctx.Data["Scripts"] = []string{"admin/js/index.js"}

	ctx.HTMLSet(200, "admin", "post_update")
	return
}

// PostManage 管理文章页面
func PostManage(ctx *context.Context) {
	ctx.Data["Title"] = "文章管理"
	ctx.Data["PostActive"] = "active"
	ctx.Data["ManageActive"] = "active toggle"

	page := ctx.GetInt("page", 1)
	limit := ctx.GetInt("limit", 10)
	if limit == 0 {
		limit = 10
	}

	posts, err := models.FindPostsDetail(page, limit, true)
	if err != nil {
		ctx.Handle(500, "", err)
		return
	}
	ctx.Data["Posts"] = posts

	p, err := models.PostsPagination(page, limit, true)
	if err != nil {
		ctx.Handle(500, "", err)
		return
	}
	ctx.Data["Page"] = p

	ctx.Data["Styles"] = []string{"admin/css/post_list.css"}
	ctx.Data["Scripts"] = []string{"admin/js/index.js", "admin/js/util.js"}
	ctx.HTMLSet(200, "admin", "post_list")
	return
}

// PostDelete 删除文章
func PostDelete(ctx *context.Context) {
	ids := ctx.PostString("ids")
	if len(ids) == 0 {
		ctx.RespJSON("没有删除任何文章")
		return
	}

	err := models.DeletePosts(ids)
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}

	// 删除关系表中的相关记录
	models.DeleteRelationship(ids, "post")

	ctx.RespJSON("200", "删除文章成功", ctx.URLFor("postManage"))
	return
}
