package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sweet-common/constants"
	"time"
)

/*
数据处理工具类
*/

func GetMD5(str string, salt string) string {
	encryptedPassword := md5.Sum([]byte(str + salt))
	hashedPassword := hex.EncodeToString(encryptedPassword[:])
	return hashedPassword
}

/*
获取随机字符串
*/
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetFileBasePath() string {
	path := os.Getenv(constants.IMG_BASE_PATH)

	if IsEmpty(path) {
		path = "http://localhost:26666"
	}

	return path
}

/*
获取随机字符串
*/
func RandomImageName(fileName string) string {
	re := regexp.MustCompile(`\.(.*?)$`)
	match := re.FindStringSubmatch(fileName)
	var suffix string
	if len(match) == 2 {
		// 获取匹配到的内容
		suffix = match[1]
	} else {
		suffix = ""
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	timestamp := time.Now().Unix()

	return fmt.Sprintf("%d_%s.%s", timestamp, string(b), suffix)
}

/*
判断字符串是否为空

	特别注意 空字符串和NULL字符串都是空字符串
*/
func IsEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}
	upperCaseString := strings.ToUpper(str)
	if upperCaseString == "NULL" {
		return true
	}
	return false
}

/*
判断字符串是否不为空
*/
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

/*
判断字符串数组中是否包含某个字符串值
*/
func ArrContains4Str(list []string, value string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}
	return false
}

/*
int64转string
*/
func Int64ToString(num int64) string {
	str := strconv.FormatInt(num, 10)
	return str
}

/*
string转int64
*/
func StringToInt64(str string) int64 {
	num, _ := strconv.ParseInt(str, 10, 64)
	return num
}
