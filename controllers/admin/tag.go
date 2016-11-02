package admin

import (
	"blog/models"
	"blog/modules/context"
	"log"
)

// Tag 列表页面
func Tag(ctx *context.Context) {

	ctx.Data["Title"] = "标签管理"
	ctx.Data["TagActive"] = "active"
	ctx.Data["ManageActive"] = "active toggle"

	tags, err := models.FindAllTags()
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	ctx.Data["Tags"] = tags

	ctx.Data["Styles"] = []string{"admin/css/tag.css"}
	ctx.Data["Scripts"] = []string{"admin/js/index.js", "admin/js/util.js"}

	ctx.HTMLSet(200, "admin", "tag")
	return
}

// CreateOrUpdateTag 创建或者更新标签
func CreateOrUpdateTag(ctx *context.Context) {
	tag := &models.Meta{}
	id := ctx.PostInt64("id")
	tag.ID = id
	tag.Name = ctx.PostString("name")
	tag.Slug = ctx.PostString("slug", tag.Name)
	tag.Type = "tag"

	if len(tag.Name) == 0 {
		ctx.RespJSON("400", "标签名称不能为空")
		return
	}

	exist, err := tag.NameExist("tag")
	log.Println(exist, err)
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	if exist {
		ctx.RespJSON("401", "标签名已存在")
		return
	}
	exist, err = tag.SlugExist("tag")
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	if exist {
		ctx.RespJSON("401", "缩略名已存在")
		return
	}

	log.Println(tag)
	// 执行创建操作
	err = tag.CreateOrUpdate()
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}

	msg := "添加标签成功"
	if id > 0 {
		msg = "修改标签成功"
	}

	ctx.RespJSON("200", msg, ctx.URLFor("tagManage"))
	return
}

// DeleteTag 删除标签操作
func DeleteTag(ctx *context.Context) {
	ids := ctx.PostString("ids")
	if len(ids) == 0 {
		ctx.RespJSON("没有删除任何标签")
		return
	}
	_, err := models.DeleteMetas(ids)
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	// 删除标签以后要删除关联表里的记录
	models.DeleteRelationship(ids, "meta")

	ctx.RespJSON("200", "删除成功", ctx.URLFor("tagManage"))
	return
}
