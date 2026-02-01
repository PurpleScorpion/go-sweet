package utils

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"shared/constants"
	"time"
)

/*
时间工具类
*/

func GetNowUTCDateFormat() string {
	currentTime := time.Now().UTC().Truncate(time.Second)
	nowDate := currentTime.Format("2006-01-02T15:04:05Z")
	return nowDate
}

/*
UTC时间转本地时间
utcTime: 2006-01-02T15:04:05.999999Z
localTime: 2006-01-02 15:04:05
*/
func UTCtoLocal(utcTime string) (string, error) {
	inputTime, err := time.Parse(constants.UTC_LAYOUT, utcTime)
	if err != nil {
		return "", err
	}
	localTime := inputTime.Local()
	outputTimeStr := localTime.Format(constants.LOCAL_LAYOUT)
	return outputTimeStr, nil
}

/*
本地时间转UTC时间
localTime: 2006-01-02 15:04:05
utcTime: 2006-01-02T15:04:05.999999Z
*/
func LocaltoUTC(localTime string) (string, error) {
	inputTime, err := time.Parse(constants.LOCAL_LAYOUT, localTime)
	if err != nil {
		return "", err
	}
	outputTimeStr := inputTime.Format(constants.UTC_LAYOUT)
	return outputTimeStr, nil
}

/*
解析UTC时间
utcTime: 2006-01-02T15:04:05.999999Z
*/
func ParseUTC(utcTime string) (time.Time, error) {
	inputTime, err := time.Parse(constants.UTC_LAYOUT, utcTime)
	if err != nil {
		return time.Time{}, err
	}
	return inputTime, nil
}

/*
解析本地时间
localTime: 2006-01-02 15:04:05
*/
func ParseLocal(localTime string) (time.Time, error) {
	inputTime, err := time.Parse(constants.LOCAL_LAYOUT, localTime)
	if err != nil {
		return time.Time{}, err
	}
	return inputTime, nil
}

func FormatLocalTime(date time.Time) string {
	return date.Format(constants.LOCAL_LAYOUT)
}

func FormatUTC(date time.Time) string {
	return date.Format(constants.UTC_LAYOUT)
}

// 字符串时间转为Time
func Str2Time(timeStr string, layout string) time.Time {
	t, _ := time.Parse(layout, timeStr)
	return t
}

func Time2Str(date time.Time, layout string) string {
	formattedTime := date.Format(layout)
	return formattedTime
}

func CompareTimeDate(time1 string, time2 string) (bool, error) {
	layout := "2006-01-02T15:04:05Z"
	t1, err := time.Parse(layout, time1)
	if err != nil {
		return false, fmt.Errorf("Error parsing the date: %w", err)
	}
	t2, err := time.Parse(layout, time2)
	if err != nil {
		return false, fmt.Errorf("Error parsing the date: %w", err)
	}
	subTime := t1.Sub(t2)
	if subTime.Seconds() <= 30 && subTime.Seconds() >= -5 {
		return true, nil
	} else {
		return false, nil
	}

}

/*
true: 过期
false: 未过期
*/
func CompareDate(utcStr string, data string) (bool, error) {
	t, err := time.Parse(constants.UTC_LAYOUT, utcStr)
	if err != nil {
		return false, fmt.Errorf("Error parsing the date: %w", err)
	}
	now := keqing.NowUTCDate()
	if now.After(t) {
		return true, nil
	}
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	id := js.GetFloat64("id")

	expire := GetCache(constants.GetUserExpireTimeKey(int32(id)))
	if keqing.IsEmpty(expire) {
		return true, nil
	}

	utcDate := keqing.ParseUTC(expire)

	if now.After(utcDate) {
		return true, nil
	}

	return false, nil
}
