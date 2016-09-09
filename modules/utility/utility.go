package utility

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"
)

// Md5Encrypt MD5加密
func Md5Encrypt(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

// Date 时间格式化函数
func Date(format string, timestamp int) string {
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
}

// FileExist 检查文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil || os.IsExist(err)
}

// Str2Int64 字符串转 int64
func Str2Int64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
