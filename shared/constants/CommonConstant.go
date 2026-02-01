package constants

import (
	"fmt"
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

	HEALTH_CHECK_KEY         string = "HEALTH_CHECK_"
	USER_EXPIRE_TIME_KEY     string = "USER_EXPIRE_TIME_"
	GATEWAY_CHECK_KEY               = "gateway_check_"
	USER_BASE_LIST           string = "USER_BASE_LIST"
	USER_CHILD_BASES         string = "USER_CHILD_BASES"
	USER_TENANTS             string = "USER_TENANTS:"
	USER_PARENT_TENANTS      string = "USER_PARENT_TENANTS:"
	SIGNAGE_CERT             string = "SIGNAGE:CERT_"
	USER_READ_ALARMS         string = "USER_READ_ALARMS"
	CAD_CONF                 string = "CAD_CONF"
	IOT_HUB_METHOD_NAME      string = "cloud"
	IOT_HUB_CTRL_METHOD_NAME string = "control"
	WEATHER_CACHE_KEY        string = "WEATHER_CACHE_"
	TEMPERATURE_CACHE_KEY    string = "TEMPERATURE_CACHE_"
	QUEUE_SEPARATOR          string = "/$$$QUEUE$$$/"
)

func GetHealthCheckKey(id int32) string {
	return fmt.Sprintf("%s%d", HEALTH_CHECK_KEY, id)
}

func GetUserExpireTimeKey(id int32) string {
	return fmt.Sprintf("%s%d", USER_EXPIRE_TIME_KEY, id)
}
