package filter

import (
	"net"
	"shared/constants"
	"strings"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/gin-gonic/gin"
)

var fullExcludeUrls []string
var prefixExcludeUrls []string
var whitelist []string

func Init() {
	fullExcludeUrls = keqing.ValueStringArr("${sweet.excludeUrl.full}")
	prefixExcludeUrls = keqing.ValueStringArr("${sweet.excludeUrl.prefix}")
	whitelist = keqing.ValueStringArr("${sweet.whitelist}")
}

func AuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {

		if excludeUrl(c.Request.URL.String()) {
			c.Next()
			return
		}

		checkIpFlag := checkWhiteIp(c.ClientIP())
		if checkIpFlag {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(200, getErrorToken("Authorization does not exist"))
			c.Abort() //
			return
		}

		// 校验通过，继续
		c.Next()
	}
}

func checkWhiteIp(ip string) bool {
	flag := false
	for _, str := range whitelist {
		if str == "-1.-1.-1.-1" {
			return false
		}
		if str == "0.0.0.0" {
			return true
		}
		if str == ip {
			flag = true
			break
		}
		if checkIPAround(str, ip) {
			flag = true
			break
		}
	}
	return flag
}

func checkIPAround(ipNetStr string, ipStr string) bool {
	_, ipNet, err := net.ParseCIDR(ipNetStr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {

		return false
	}
	if ipNet.Contains(ip) {
		return true
	}
	return false
}

func getErrorToken(msg string) interface{} {
	return gin.H{
		"code": constants.TOKEN_ERROR,
		"msg":  msg,
		"data": "",
	}
}

func excludeUrl(url string) bool {
	if len(prefixExcludeUrls) > 0 {
		for _, v := range prefixExcludeUrls {
			if strings.HasPrefix(url, v) {
				return true
			}
		}
	}

	if len(fullExcludeUrls) > 0 {
		for _, v := range fullExcludeUrls {
			if url == v {
				return true
			}
			tmp := "/" + v
			if url == tmp {
				return true
			}
		}
	}
	return false
}
