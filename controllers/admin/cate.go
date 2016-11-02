package admin

import (
	"blog/models"
	"blog/modules/context"
	"log"
)

// Cate 分类列表页面
func Cate(ctx *context.Context) {

	ctx.Data["Title"] = "分类管理"
	ctx.Data["CateActive"] = "active"
	ctx.Data["ManageActive"] = "active toggle"

	cates, err := models.FindAllCates()
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}

	ctx.Data["Cates"] = cates

	ctx.Data["Styles"] = []string{"admin/css/post_list.css"}
	ctx.Data["Scripts"] = []string{"admin/js/index.js", "admin/js/util.js"}

	ctx.HTMLSet(200, "admin", "cate")
	return
}

// DeleteCate 删除分类操作
func DeleteCate(ctx *context.Context) {
	ids := ctx.PostString("ids")
	if len(ids) == 0 {
		ctx.RespJSON("没有删除任何分类")
		return
	}
	_, err := models.DeleteMetas(ids)
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}

	// 删除关系表里的相关记录
	models.DeleteRelationship(ids, "meta")

	ctx.RespJSON("200", "删除成功", ctx.URLFor("cateManage"))
	return
}

// UpdateCate 提交更新分类
func UpdateCate(ctx *context.Context) {
	ctx.RespJSON("修改分类成功")
}

// CreateOrUpdateCate 创建或者更新分类
func CreateOrUpdateCate(ctx *context.Context) {
	cate := &models.Meta{}
	id := ctx.PostInt64("id")
	cate.ID = id
	cate.Name = ctx.PostString("name")
	cate.Slug = ctx.PostString("slug", cate.Name)
	cate.ParentID = ctx.PostInt64("parent_id")
	cate.Type = "category"

	if len(cate.Name) == 0 {
		ctx.RespJSON("400", "分类名称不能为空")
		return
	}

	exist, err := cate.NameExist("category")
	log.Println(exist, err)
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	if exist {
		ctx.RespJSON("401", "分类名已存在")
		return
	}
	exist, err = cate.SlugExist("category")
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}
	if exist {
		ctx.RespJSON("401", "缩略名已存在")
		return
	}

	log.Println(cate)
	// 执行创建操作
	err = cate.CreateOrUpdate()
	if err != nil {
		ctx.RespJSON("500", "内部服务错误")
		return
	}

	msg := "添加分类成功"
	if id > 0 {
		msg = "修改分类成功"
	}
	ctx.RespJSON("200", msg, ctx.URLFor("cateManage"))
	return
}

// FindPostByCate 根据分类查找文章
func FindPostByCate(ctx *context.Context) {

}
