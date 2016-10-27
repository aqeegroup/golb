package admin

import (
	"blog/models"
	"blog/modules/context"
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
