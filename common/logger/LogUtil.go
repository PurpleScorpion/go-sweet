package logger

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/beego/beego/v2/core/logs"
)

var log *logs.BeeLogger

func Init() {
	// 创建一个日志器，可以给它指定一个名称，便于区分多个日志器
	log = logs.NewLogger()
	// 设置日志级别，例如：debug、info、warn、error、critical，默认为debug
	js := jsonutil.NewJSONObject()
	level := keqing.ValueString("${sweet.logging.level}")

	if keqing.IsEmpty(level) {
		level = "info"
	}

	switch level {
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

	file := keqing.ValueString("${sweet.logging.file}")
	if keqing.IsEmpty(file) {
		file = "log/go_sweet.log"
	}
	maxSize := keqing.ValueInt("${sweet.logging.maxSize}")
	if maxSize <= 0 || maxSize > 1024 {
		maxSize = 10
	}
	maxDays := keqing.ValueInt("${sweet.logging.maxDays}")
	if maxDays <= 0 {
		maxDays = 7
	}
	maxBackups := keqing.ValueInt("${sweet.logging.maxBackups}")
	if maxBackups <= 0 {
		maxBackups = 10
	}

	// 添加一个文件日志引擎，指定日志文件路径和模式（如按天分割、按大小分割等）
	js.FluentPut("filename", file)
	js.FluentPut("maxSize", maxSize*1024*1024)
	js.FluentPut("maxDays", maxDays)
	js.FluentPut("daily", true)
	js.FluentPut("maxBackups", maxBackups)

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
