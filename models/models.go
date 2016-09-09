package models

import (
	"fmt"
	"log"
	"os"
	"path"

	// 导入mysql包
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"blog/modules/setting"
)

// RespJSON 返回 json
type RespJSON struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Redirect string `json:"redirect"`
}

// RespWithData 返回 json 并且包含data
type RespWithData struct {
	RespJSON
	Data map[string]interface{} `json:"data"`
}

var (
	x      *xorm.Engine
	tables []interface{}

	// DbConfig 数据库配置文件
	DbConfig struct {
		Type, Host, Port, DbName, Username, Password string
	}
)

func init() {

	tables = append(tables, new(Post), new(Meta), new(Relationship), new(Option),
		new(User))
}

// Init 数据库初始化
func Init() error {

	loadConfigs()
	if err := setEngine(); err != nil {
		return err
	}

	// 必须先ping一下才能连接成功
	if err := Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 同步数据库
	log.Println("开始同步数据库")
	if err := x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}

	return nil
}

// loadConfigs 加载时数据库配置
func loadConfigs() {
	section := setting.Cfg.Section("database")
	DbConfig.Type = section.Key("type").String()
	DbConfig.Host = section.Key("host").String()
	DbConfig.Port = section.Key("port").String()
	DbConfig.DbName = section.Key("dbname").String()
	DbConfig.Username = section.Key("username").String()
	DbConfig.Password = section.Key("password").String()
}

func getEngine() (*xorm.Engine, error) {
	dsn := ""
	switch DbConfig.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", DbConfig.Username, DbConfig.Password, DbConfig.Host, DbConfig.Port, DbConfig.DbName)
	}

	return xorm.NewEngine(DbConfig.Type, dsn)
}

func setEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置数据库映射规则
	x.SetMapper(core.GonicMapper{})

	logPath := path.Join(setting.LogRootPath, "xorm.log")
	os.MkdirAll(path.Dir(logPath), os.ModePerm)

	f, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("创建数据库日志文件 xorm.log 失败: %s", err)
	}
	x.SetLogger(xorm.NewSimpleLogger(f))
	x.ShowSQL(true)

	return nil
}

// Ping ping 数据库
// 有必要的话可以隔一段时间 ping 一下, 防止数据库断开连接
func Ping() error {
	return x.Ping()
}
