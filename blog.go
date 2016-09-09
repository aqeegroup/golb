package main

import (
	"html/template"

	"github.com/go-macaron/session"
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
	m.Run()
}

func setup() {
	// 初始化文件配置
	setting.NewContext()

	err := models.Init()
	if err != nil {
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
			"URLFor": m.URLFor,     // url 生成函数
			"date":   utility.Date, // 时间格式化函数
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

	// 初始化路由
	routesInit()
}

// 初始化路由配置
func routesInit() {
	m.NotFound(admin.NotFound)

	// 前端页面路由
	m.Get("/", home.Index).Name("home")
	m.Get("/post/:slug([\\w-]+)", home.Detail).Name("postDetail")

	// 后台路由组
	m.Group("/admin", func() {
		m.Get("/", admin.CheckLogin, admin.Index)
		m.Get("/login", admin.Login).Name("login")
		m.Post("/login", admin.DoLogin).Name("doLogin")
		m.Get("/logout", admin.Logout).Name("logout")

		m.Group("/post", func() {
			m.Get("/", admin.WritePage).Name("writePost")
			m.Post("/", admin.PostSubmit)
		})

	})
}
