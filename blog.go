package main

import (
	"html/template"
	"log"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"blog/controllers/admin"
	"blog/controllers/home"
	"blog/models"
	"blog/modules/setting"
	"blog/modules/utility"
)

var m *macaron.Macaron

func main() {
	// 一些配置初始化
	setup()
	// macaron 框架初始化
	macaronInit()
	m.Run()
}

func setup() {
	// 初始化文件配置
	setting.NewContext()

	err := models.Init()
	if err != nil {
		log.Println("数据库初始化失败: ", err)
	}

	models.LoadUserConfig()
}

func macaronInit() {
	m = macaron.Classic()

	// session 中间件
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "runtime/sessions",
		CookieName:     setting.CookieName,
	}))

	// 模板引擎
	m.Use(macaron.Renderers(macaron.RenderOptions{
		Directory: "templates/themes/default",
		Funcs: []template.FuncMap{map[string]interface{}{
			"URLFor": m.URLFor,     // url 生成函数
			"date":   utility.Date, // 时间格式化函数
		}},
	}, "admin:templates/admin"))

	// 初始化路由
	routesInit()
}

// 初始化路由配置
func routesInit() {
	m.NotFound(func() string {
		return "404"
	})

	// 前端页面路由
	m.Get("/", home.Index)
	m.Get("/post/:slug:string", home.Detail).Name("postDetail")

	// 后台路由组
	m.Group("/admin", func() {
		m.Get("/", admin.Index)
		m.Get("/login", admin.Login).Name("login")
		m.Post("/login", admin.DoLogin).Name("doLogin")
		m.Get("/logout", admin.Logout).Name("logout")
	})

}
