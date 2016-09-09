package admin

import (
	"blog/modules/context"
)

// Cate 分类列表页面
func Cate(ctx *context.Context) {

	ctx.Data["Title"] = "分类管理"
	ctx.HTML(200, "admin", "cate")
	return
}

// DoCreateCate 提交修改分类
func DoCreateCate(ctx *context.Context) {
	ctx.RespJSON("添加分类成功")
}

// DoUpdateCate 提交更新分类
func DoUpdateCate(ctx *context.Context) {
	ctx.RespJSON("修改分类成功")
}

// CreateOrUpdateCate 创建或者更新分类
func CreateOrUpdateCate(ctx *context.Context) {
	ctx.RespJSON("添加分类成功")
}

// DoDeleteCate 删除分类
func DoDeleteCate(ctx *context.Context) {
	ctx.RespJSON("分类已删除")
}

// FindPostByCate 根据分类查找文章
func FindPostByCate(ctx *context.Context) {

}
