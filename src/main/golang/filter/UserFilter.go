package filter

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"net"
	"net/http"
	"strings"
	"sweet-common/constants"
	"sweet-common/utils"
	"time"
)

var fullExcludeUrls []string
var prefixExcludeUrls []string
var whitelist []string

func Init() {
	fullExcludeUrls = keqing.ValueStringArr("${sweet.excludeUrl.full}")
	prefixExcludeUrls = keqing.ValueStringArr("${sweet.excludeUrl.prefix}")
	whitelist = keqing.ValueStringArr("${sweet.whitelist}")
}

func RecoverPanic(ctx *context.Context, cfg *web.Config) {
	if err := recover(); err != nil {
		// 在这里处理 panic
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(getErrorToken(fmt.Sprintf("%v", err)), false, false)
	}
}

var UserFilter = func(ctx *context.Context) {
	if excludeUrl(ctx.Input.URL()) {
		return
	}

	checkIpFlag := checkWhiteIp(ctx.Input.IP())
	if checkIpFlag {
		return
	}

	headers := ctx.Request.Header
	values := headers["Token"]

	if len(values) == 0 {
		ctx.Output.JSON(getErrorToken("Token does not exist"), false, false)
		return
	}

	token := values[0]

	if keqing.IsEmpty(token) {
		ctx.Output.JSON(getErrorToken("Token does not exist"), false, false)
		return
	}

	data := keqing.RsaDecrypt(token)
	if keqing.IsEmpty(data) {
		ctx.Output.JSON(getErrorToken("Token error"), false, false)
		return
	}

	flag := hasExpire(data)

	if flag {
		ctx.Output.JSON(getErrorToken("Token has expired"), false, false)
		return
	}
}

/*
true: 过期
false: 未过期
*/
func hasExpire(data string) bool {
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	if !js.HasKey("expirationTime") {
		return true
	}
	expirationTime := js.GetString("expirationTime")

	flag, err := compareDate(expirationTime, data)
	if err != nil {
		return true
	}
	return flag
}

func checkWhiteIp(ip string) bool {
	for _, str := range whitelist {
		if str == "-1.-1.-1.-1" {
			return false
		}
		if str == "0.0.0.0" {
			return true
		}
		if str == ip {
			return true
		}
		if checkIPAround(str, ip) {
			return true
		}
	}
	return false
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
	js := jsonutil.NewJSONObject()
	js.FluentPut("code", constants.TOKEN_ERROR)
	js.FluentPut("msg", msg)
	js.FluentPut("data", "")
	return js.GetData()
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

/*
true: 过期
false: 未过期
*/
func compareDate(utcStr string, data string) (bool, error) {
	t, err := time.Parse(constants.UTC_LAYOUT, utcStr)
	if err != nil {
		return false, fmt.Errorf("Error parsing the date: %w", err)
	}
	now := time.Now().UTC()
	if now.After(t) {
		return true, nil
	}
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	id := js.GetFloat64("id")

	expire := utils.GetCache(constants.GetUserExpireTimeKey(int32(id)))
	if keqing.IsEmpty(expire) {
		return true, nil
	}

	utcDate := keqing.ParseUTC(expire)

	if now.After(utcDate) {
		return true, nil
	}

	return false, nil
}
