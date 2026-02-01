package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"shared/logger"
	"time"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

/*
获取随机字符串
*/
func RandomFileName(fileName string) string {
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

func SaveFile(centPath, fileName string, src io.Reader) string {

	keqing.MkdirAll(centPath)

	localPath := centPath + "/" + fileName
	// 创建目标文件
	dstFile, err := os.Create(localPath)
	if err != nil {
		logger.Error("创建目标文件失败:", err)
		return ""
	}
	defer dstFile.Close()
	keqing.CopyFile2IO(dstFile, src)
	return fileName
}

func GetCertPath() string {
	return "/cert"
}
