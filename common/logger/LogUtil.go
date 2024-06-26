package logger

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/beego/beego/v2/core/logs"
)

type LogUtil struct {
}

var log *logs.BeeLogger

func init() {
	readLog()
	conf := yamlConf
	// 创建一个日志器，可以给它指定一个名称，便于区分多个日志器
	log = logs.NewLogger()
	// 设置日志级别，例如：debug、info、warn、error、critical，默认为debug
	js := jsonutil.NewJSONObject()

	switch conf.Sweet.Log.Level {
	case "info":
		log.SetLevel(logs.LevelInfo)
		js.FluentPut("level", logs.LevelInfo)
		break
	case "warn":
		log.SetLevel(logs.LevelWarn)
		js.FluentPut("level", logs.LevelWarn)
		break
	case "error":
		log.SetLevel(logs.LevelError)
		js.FluentPut("level", logs.LevelError)
		break
	default:
		log.SetLevel(logs.LevelInfo)
		js.FluentPut("level", logs.LevelInfo)
	}

	// 添加一个文件日志引擎，指定日志文件路径和模式（如按天分割、按大小分割等）
	js.FluentPut("filename", conf.Sweet.Log.File)
	js.FluentPut("maxSize", conf.Sweet.Log.MaxSize*1024*1024)
	js.FluentPut("maxDays", conf.Sweet.Log.MaxDays)
	js.FluentPut("daily", true)
	js.FluentPut("maxBackups", conf.Sweet.Log.MaxBackups)

	// 如果需要按小时分割文件，可以设置HourlyRolling为true
	// fileLogConfig.HourlyRolling = true
	log.Async()
	// 将文件日志引擎添加到日志器中
	if err := log.SetLogger(logs.AdapterFile, js.ToJsonString()); err != nil {
		panic("Failed to set file logger: " + err.Error())
	}
}

func Info(format string, data ...interface{}) {
	// 控制台打印
	logs.Info(fmt.Sprintf(format, data...))
	// 文件记录
	log.Info(fmt.Sprintf(format, data...))
}

func Warn(format string, data ...interface{}) {
	logs.Warn(fmt.Sprintf(format, data...))
	log.Warn(fmt.Sprintf(format, data...))
}

func Error(format string, data ...interface{}) {
	logs.Error(fmt.Sprintf(format, data...))
	log.Error(fmt.Sprintf(format, data...))
}
