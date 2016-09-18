package utility

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
func Date(format string, timestamp int64) string {
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

// Str2Int 字符串转 int
func Str2Int(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}

// SlugNameFormat 处理 SlugName
func SlugNameFormat(s string) string {

	old := []string{
		"'", ":", "\\", "/", "\"",
	}
	s = Replace(s, old, "")

	old = []string{
		"+", ",", " ", "，", "　", ".", "?", "=", "&", "!", "<", ">", "(", ")", "[", "]", "{", "}",
	}
	s = Replace(s, old, "-")
	fmt.Println(s)

	return strings.Trim(s, "-_")
}

// Replace 支持用新字符串替换多个旧字符串
func Replace(s string, o interface{}, n string) string {

	switch v := o.(type) {
	case []string:
		var replace []string
		for _, o := range v {
			replace = append(replace, o, n)
		}
		return strings.NewReplacer(replace...).Replace(s)
	case string:
		return strings.Replace(s, v, n, -1)
	}
	return s
}

// StringSplitInt64 给一个字符串 指定分隔符 返回分割之后的 int64 数组
func StringSplitInt64(s, sep string) []int64 {
	i := []int64{}
	str := strings.Split(s, sep)
	for _, temp := range str {
		iTemp := Str2Int64(temp)
		if iTemp > 0 {
			i = append(i, iTemp)
		}
	}

	return i
}

// InArray 判断数组里是否包含某个值
func InArray(k interface{}, arr []interface{}) bool {
	for _, v := range arr {
		if k == v {
			return true
		}
	}
	return false
}
