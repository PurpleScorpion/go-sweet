package logger

import (
	"go-sweet/common/constants"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type YmlConfig struct {
	Server Server `yaml:"server"`
	Sweet  Sweet  `yaml:"sweet"`
}

type Server struct {
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
	Active string `yaml:"active"`
}

type Sweet struct {
	Log Logging `yaml:"logging"`
}

type Logging struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxDays    int    `yaml:"maxDays"`
}

var yamlConf YmlConfig

func readLog() {
	confPath := os.Getenv(constants.CONF_PATH)

	if IsEmpty(confPath) {
		confPath = "conf/application.yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			panic("Error reading configuration file: " + err.Error())
		}
		confPath = absPath
	} else {
		confPath = confPath + "/conf/application.yml"
	}

	data, err := os.ReadFile(confPath)
	if err != nil {
		panic("Error reading configuration file: " + err.Error())
	}
	err = yaml.Unmarshal(data, &yamlConf)
	if err != nil {
		panic("Error parsing YAML: " + err.Error())
	}
	if IsEmpty(yamlConf.Server.Active) {
		panic("server.active is empty")
	}
	readChildConf()
}

func getEnvActive() string {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := os.Getenv(constants.CONF_PATH)
	if IsEmpty(confPath) {
		confPath = "conf/application-" + profilesActive + ".yml"
	} else {
		confPath = confPath + "/conf/application-" + profilesActive + ".yml"
	}

	if IsNotEmpty(profilesActive) {
		absPath, err := filepath.Abs(confPath)
		if err == nil {
			_, err1 := os.ReadFile(absPath)
			if err1 == nil {
				return profilesActive
			}
		}
	}
	return yamlConf.Server.Active
}

func readChildConf() {
	var yamlConf2 YmlConfig
	confPath := os.Getenv(constants.CONF_PATH)
	yamlConf.Server.Active = getEnvActive()

	if IsEmpty(confPath) {
		confPath = "conf/application-" + yamlConf.Server.Active + ".yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			panic("Error reading configuration file: " + err.Error())
		}
		confPath = absPath
	} else {
		confPath = confPath + "/conf/application-" + yamlConf.Server.Active + ".yml"
	}
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic("Error reading configuration file: " + err.Error())
	}
	err = yaml.Unmarshal(data, &yamlConf2)
	if err != nil {
		panic("Error parsing YAML: " + err.Error())
	}
	//yamlConf2 = defaultData(yamlConf2)
	saveConf(yamlConf2)
}

func saveConf(yamlConf2 YmlConfig) {
	if IsNotEmpty(yamlConf2.Sweet.Log.Level) {
		yamlConf.Sweet.Log.Level = yamlConf2.Sweet.Log.Level
	} else {
		if IsEmpty(yamlConf.Sweet.Log.Level) {
			yamlConf.Sweet.Log.Level = "info"
		}
	}
	switch yamlConf.Sweet.Log.Level {
	case "info":
		break
	case "warn":
		break
	case "error":
		break
	default:
		panic("log.level is error, must be info/warn/error")
	}

	if IsNotEmpty(yamlConf2.Sweet.Log.File) {
		yamlConf.Sweet.Log.File = yamlConf2.Sweet.Log.File
	} else {
		if IsEmpty(yamlConf.Sweet.Log.File) {
			yamlConf.Sweet.Log.File = "logs/go-sweet.log"
		}
	}

	if yamlConf2.Sweet.Log.MaxSize > 0 {
		yamlConf.Sweet.Log.MaxSize = yamlConf2.Sweet.Log.MaxSize
	} else {
		if yamlConf.Sweet.Log.MaxSize == 0 {
			yamlConf.Sweet.Log.MaxSize = 10
		}
	}

	if yamlConf2.Sweet.Log.MaxBackups > 0 {
		yamlConf.Sweet.Log.MaxBackups = yamlConf2.Sweet.Log.MaxBackups
	} else {
		if yamlConf.Sweet.Log.MaxBackups == 0 {
			yamlConf.Sweet.Log.MaxBackups = 10
		}
	}

	if yamlConf2.Sweet.Log.MaxDays > 0 {
		yamlConf.Sweet.Log.MaxDays = yamlConf2.Sweet.Log.MaxDays
	} else {
		if yamlConf.Sweet.Log.MaxDays == 0 {
			yamlConf.Sweet.Log.MaxDays = 7
		}
	}
}

func IsEmpty(str string) bool {
	if len(str) == 0 {
		return true
	}
	upperCaseString := strings.ToUpper(str)
	if upperCaseString == "NULL" {
		return true
	}
	return false
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}
