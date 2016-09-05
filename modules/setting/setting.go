package setting

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

var (
	// AppPath 可执行文件路径
	AppPath string
	// LogRootPath 日志路径
	LogRootPath string

	// CookieName Cookie 名称
	CookieName string

	// GcLifetime cookie 存活时间
	GcLifetime int64
	// Cfg 存储全局配置文件
	Cfg *ini.File
)

func init() {

	var err error
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Println("获取文件可执行文件路径失败")
	}

	if AppPath, err = filepath.Abs(file); err != nil {
		log.Println("获取文件可执行文件路径失败")
	}
}

// NewContext 初始化配置文件上下文
func NewContext() {
	workDir := WorkDir()
	confFilePath := path.Join(workDir, "conf/app.ini")
	var err error
	if Cfg, err = ini.Load(confFilePath); err != nil {
		log.Println("加载配置文件失败, 文件路径: " + confFilePath)
	}

	LogRootPath = Cfg.Section("log").Key("root_path").MustString(path.Join(workDir, "log"))
	sessionSec := Cfg.Section("session")
	CookieName = sessionSec.Key("cookie_name").MustString("golbSession")
	GcLifetime = sessionSec.Key("gc_lifetime").MustInt64(86400)

}

// WorkDir 当前工作木目录
func WorkDir() string {
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath
	}
	return AppPath[:i]
}
