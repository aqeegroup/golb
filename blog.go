package main

import (
	"html/template"
	"log"
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	"blog/controllers/admin"
	"blog/controllers/home"
	"blog/models"
	"blog/modules/setting"
)

func main() {

	// macaron 框架初始化
	m := macaron.Classic()
	m.Use(macaron.Renderers(macaron.RenderOptions{
		Directory: "templates/themes/default",
		Funcs: []template.FuncMap{map[string]interface{}{
			"URLFor": m.URLFor,
			// 考虑下这个函数该放哪里
			"date": func(format string, timestamp int) string {
				dateReplace := []string{
					"Y", "2006",
					"m", "01",
					"d", "02",
					"H", "15",
					"i", "04",
					"s", "05",
				}
				r := strings.NewReplacer(dateReplace...)
				format = r.Replace(format)
				return time.Unix(int64(timestamp), 0).Format(format)
			},
		}},
	}, "admin:templates/admin"))

	// 一些配置初始化
	setup()

	// 初始化路由
	routesInit(m)

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

// 初始化路由配置
func routesInit(m *macaron.Macaron) {
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
	})

}
