package utils

import (
	"go-sweet/common/constants"
	"go-sweet/common/logger"
	"time"
)

/*
时间工具类
*/

/*
获取当前时间
格式 2006-01-02 15:04:05
*/
func GetNowDate() string {
	currentTime := time.Now().Truncate(time.Second)
	nowDate := currentTime.Format("2006-01-02 15:04:05")
	return nowDate
}

/*
获取当前时间
格式 2006-01-02T15:04:05.999999Z
*/
func GetNowUTCDate() string {
	now := time.Now().UTC()
	formattedTime := now.Format(constants.UTC_LAYOUT)
	return formattedTime
}

/*
UTC时间转本地时间
utcTime: 2006-01-02T15:04:05.999999Z
localTime: 2006-01-02 15:04:05
*/
func UTCtoLocal(utcTime string) (string, error) {
	inputTime, err := time.Parse(constants.UTC_LAYOUT, utcTime)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	outputTimeStr := inputTime.Format(constants.LOCAL_LAYOUT)
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
		logger.Error(err.Error())
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
		logger.Error(err.Error())
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
		logger.Error(err.Error())
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
