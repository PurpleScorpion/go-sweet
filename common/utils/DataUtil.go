package utils

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"math/rand"
	"os"
	"regexp"
	"sweet-common/constants"
	"time"
)

/*
数据处理工具类
*/

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

func LoadRsaKey() {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	privateKey := "conf/privateKey.pem"
	publicKey := "conf/publicKey.pem"
	if keqing.IsEmpty(profilesActive) {
		privateKey = "src/main/resources/privateKey.pem"
		publicKey = "src/main/resources/publicKey.pem"
	}
	keqing.RsaLoadKey(publicKey, privateKey, keqing.RSA_KEY_FILE_TYPE)
}
