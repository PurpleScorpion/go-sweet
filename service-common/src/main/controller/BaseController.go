package controller

import (
	"shared/logger"
	"shared/utils"
	"shared/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func GetUserInfo(c *gin.Context) vo.UserVO {
	info := vo.UserVO{}
	info.Username = c.GetHeader("X-Current-User-Name")
	return info
}

func OK(c *gin.Context, msg string, data interface{}) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}
func GetInt64(c *gin.Context, key string) int64 {
	i, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		logger.Error("GetInt64 error: {}", err)
		return 0
	}
	return i
}

func GetInt(c *gin.Context, key string) int {
	i, err := strconv.Atoi(c.Param(key))
	if err != nil {
		logger.Error("GetInt error: {}", err)
		return 0
	}
	return i
}

func GetFloat(c *gin.Context, key string) float64 {
	float, err := strconv.ParseFloat(c.Param(key), 64)
	if err != nil {
		logger.Error("GetFloat error: {}", err)
		return 0
	}
	return float
}

func GetString(c *gin.Context, key string) string {
	return c.Param(key)
}

func BindForm[T any](c *gin.Context) T {
	var form T
	if err := c.ShouldBind(&form); err != nil {
		return form
	}
	return form
}

func GetFile(c *gin.Context) vo.UploadForm {
	var form vo.UploadForm
	if err := c.ShouldBind(&form); err != nil {
		return form
	}
	return form
}

func BindJSON[T any](c *gin.Context) T {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req
	}
	return req
}

func BindQuery[T any](c *gin.Context) T {
	var q T
	if err := c.ShouldBindQuery(&q); err != nil {
		return q
	}
	return q
}

func Fail(c *gin.Context, code int32, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func HttpFail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"message": msg,
	})
}

func Result(c *gin.Context, rs utils.R) {
	c.JSON(200, gin.H{
		"code": rs.Code,
		"msg":  rs.Msg,
		"data": rs.Data,
	})
}
