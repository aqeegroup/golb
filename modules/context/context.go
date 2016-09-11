package context

import (
	"fmt"
	"strings"
	"time"

	"github.com/Unknwon/log"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"blog/models"
	"blog/modules/setting"
	"blog/modules/utility"
)

// Context 自定义的 Context
type Context struct {
	*macaron.Context
	Session session.Store
	Flash   *session.Flash
}

// Contexter 初始化一个自定义Context
// 这里设置了一些公共信息 扩展一些ctx的方法
func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Session: sess,
			Flash:   f,
		}

		ctx.Data["PageStartTime"] = time.Now()

		log.Debug("Session ID: %s", sess.ID())
		ctx.Data["Version"] = "0.1"
		ctx.Data["VersionDate"] = "2016-09-05"

		ctx.Data["Uid"] = sess.Get("uid")
		ctx.Data["Username"] = sess.Get("username")

		// 导航 active 控制 - 可以不要
		ctx.Data["MainActive"] = ""
		ctx.Data["ThemeActive"] = ""
		ctx.Data["PostActive"] = ""
		ctx.Data["PageActive"] = ""
		ctx.Data["CommentActive"] = ""
		ctx.Data["CateActive"] = ""
		ctx.Data["TagActive"] = ""
		ctx.Data["FileActive"] = ""
		ctx.Data["UserActive"] = ""
		ctx.Data["LinkActive"] = ""

		c.Map(ctx)
	}
}

// Handle 处理错误
func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		log.Error("%s: %v", title, err)
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	switch status {
	case 404:
		ctx.Data["Title"] = "Page Not Found"
		ctx.Data["ErrMsg"] = "Page Not Found."
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
		ctx.Data["ErrMsg"] = "Internal Server Error."
	}
	// fmt.Println(ctx.HasTemplate("error"))
	if ctx.HasTemplate("error") {
		ctx.HTML(status, "error")
		return
	}
	ctx.HTMLSet(status, "admin", "error")
}

// HasTemplate 模板文件是否存在
func (ctx *Context) HasTemplate(t string) bool {

	theme := models.Options.Get("theme")
	tPath := fmt.Sprintf("%s/%s/%s.html", setting.WorkDir(), theme, t)
	// fmt.Println(tPath)
	if utility.FileExist(tPath) {
		return true
	}

	return false
}

// RemoteIP 获取客户端 ip 地址
func (ctx *Context) RemoteIP() string {
	addr := strings.Split(ctx.Req.RemoteAddr, ":")
	return addr[0]
}

// RespJSON 成功返回
func (ctx *Context) RespJSON(r ...string) {
	resp := models.RespJSON{}

	len := len(r)
	if len == 0 {
		resp.Msg = r[0]
		resp.Redirect = "OK"
	} else if len == 1 {
		resp.Msg = r[0]
		resp.Redirect = ""
	} else {
		resp.Msg = r[0]
		resp.Redirect = r[1]
	}

	ctx.JSON(200, resp)
	return
}
