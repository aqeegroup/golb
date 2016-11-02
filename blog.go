package main

import (
	"html/template"
	"strings"

	"github.com/go-macaron/session"
	"github.com/jiayx/go/random"
	"gopkg.in/macaron.v1"

	"blog/controllers/admin"
	"blog/controllers/home"
	"blog/models"
	"blog/modules/context"
	"blog/modules/setting"
	"blog/modules/utility"
)

var m *macaron.Macaron

func main() {
	// 一些配置初始化
	setup()
	// macaron 框架初始化
	newMacaron()

	// 初始化路由
	routesInit()

	m.Run()
}

func setup() {
	// 初始化文件配置
	setting.NewContext()

	if err := models.Init(); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}

	if err := models.LoadOptions(); err != nil {
		panic("用户配置加载失败: " + err.Error())
	}
}

func newMacaron() {
	m = macaron.Classic()

	// 模板引擎
	m.Use(macaron.Renderers(macaron.RenderOptions{
		// 这个配置将来要从数据库取
		Directory: models.Options.Get("theme"),
		Funcs: []template.FuncMap{map[string]interface{}{
			"URLFor":     m.URLFor,        // url 生成函数
			"date":       utility.Date,    // 时间格式化函数
			"asset":      models.AssetURL, // 生成静态文件链接
			"join":       strings.Join,    // 生成静态文件链接
			"trim":       strings.Trim,
			"InArray":    utility.InArray,
			"JSONEncode": utility.JSONEncode,
			"RandInt":    random.Int,
		}},
	}, "admin:templates/admin"))

	// session 中间件
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "runtime/sessions",
		CookieName:     setting.CookieName,
		Gclifetime:     setting.GcLifetime,
	}))

	m.Use(context.Contexter())
}

// 初始化路由配置
func routesInit() {
	m.NotFound(admin.NotFound)

	// 前端页面路由
	m.Get("/", home.Index).Name("home")
	m.Get("/page/:page:int", home.Index).Name("page")
	m.Get("/post/:slug([\\w-]+)", home.Detail).Name("postDetail")

	// 后台路由组
	m.Group("/admin", func() {
		m.Get("/", admin.CheckLogin, admin.Index).Name("admin")
		m.Get("/login", admin.Login).Name("login")
		m.Post("/login", admin.DoLogin).Name("doLogin")
		m.Get("/logout", admin.Logout).Name("logout")

		m.Group("/post", func() {
			m.Get("/", admin.WritePage).Name("writePost")
			m.Post("/", admin.PostSubmit)
			m.Get("/manage/", admin.PostManage).Name("postManage")
			m.Post("/del/", admin.PostDelete).Name("PostDelete")
			m.Get("/:id:int", admin.PostUpdate).Name("postUpdate")

		}, admin.CheckLogin)

		m.Group("/cate", func() {
			m.Get("/", admin.Cate).Name("cateManage")
			m.Post("/", admin.CreateOrUpdateCate).Name("cateCreateOrUpdate")
			m.Post("/del/", admin.DeleteCate).Name("cateDel")
		}, admin.CheckLogin)

		m.Group("/tag", func() {
			m.Get("/", admin.Tag).Name("tagManage")
			m.Post("/", admin.CreateOrUpdateTag).Name("tagCreateOrUpdate")
			m.Post("/del", admin.DeleteTag).Name("tagDel")
		}, admin.CheckLogin)

	})
}
