package constants

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
)

var (
	NO_DELETE_CODE int32 = 0 // 未删除
	DELETE_CODE    int32 = 1 // 已删除

	NORMAL_STATUS int32 = 1 // 正常
	FAIL_STATUS   int32 = 0 // 失效

	// beego需要的环境变量
	IMG_BASE_PATH   string = "IMG_BASE_PATH" // 图片的基础路径 https开头
	BEEGO_RUNMODE   string = "BEEGO_RUNMODE" // 当前环境 固定值 prod
	CONF_PATH       string = "CONF_PATH"
	PROFILES_ACTIVE string = "PROFILES_ACTIVE"

	HEALTH_CHECK_KEY     string = "HEALTH_CHECK_"
	USER_EXPIRE_TIME_KEY string = "USER_EXPIRE_TIME_"
)

var (
	IMG_BASE_URL string = ""
	IMG_PATH     string = ""
	MAPPING_URL  string = ""
)

func Init() {
	IMG_BASE_URL = keqing.ValueString("${sweet.img.baseUrl}")
	if keqing.IsEmpty(IMG_BASE_URL) {
		IMG_BASE_URL = fmt.Sprintf("http://localhost:%d", keqing.ValueInt("${server.port}"))
	}

	IMG_PATH = keqing.ValueString("${sweet.img.path}")
	MAPPING_URL = keqing.ValueString("${sweet.img.mappingUrl}")

}

func GetHealthCheckKey(id int32) string {
	return fmt.Sprintf("%s%d", HEALTH_CHECK_KEY, id)
}

func GetUserExpireTimeKey(id int32) string {
	return fmt.Sprintf("%s%d", USER_EXPIRE_TIME_KEY, id)
}
