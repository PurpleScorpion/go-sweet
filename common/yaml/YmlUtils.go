package sweetyml

import (
	"fmt"
	"go-sweet/common/constants"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var yamlConf YmlConfig

func ReadYml() {
	confPath := os.Getenv(constants.CONF_PATH)

	if IsEmpty(confPath) {
		confPath = "conf/application.yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
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
		log.Fatalf("error parsing YAML: %v", err)
	}

	if IsEmpty(yamlConf.Server.Active) {
		panic("server.active is empty")
	}
	readChildConf()
}

func GetYmlConf() YmlConfig {
	return yamlConf
}

func readChildConf() {
	var yamlConf2 YmlConfig
	confPath := os.Getenv(constants.CONF_PATH)

	if IsEmpty(confPath) {
		confPath = "conf/application-" + yamlConf.Server.Active + ".yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
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
		log.Fatalf("error parsing YAML: %v", err)
	}
	//yamlConf2 = defaultData(yamlConf2)
	saveConf(yamlConf2)
}

func defaultData(yamlConf2 YmlConfig) YmlConfig {
	if yamlConf2.Server.Port == 0 {
		yamlConf2.Server.Port = 8080
	}
	if IsEmpty(yamlConf2.Server.Name) {
		yamlConf2.Server.Name = "go-sweet"
	}
	if yamlConf2.Sweet.MySqlConfig.Active {
		if IsEmpty(yamlConf2.Sweet.MySqlConfig.Host) {
			panic("mysql.host is empty")
		}
		if yamlConf2.Sweet.MySqlConfig.Port == 0 {
			yamlConf2.Sweet.MySqlConfig.Port = 3306
		}
		if IsEmpty(yamlConf2.Sweet.MySqlConfig.User) {
			yamlConf2.Sweet.MySqlConfig.User = "root"
		}
		if IsEmpty(yamlConf2.Sweet.MySqlConfig.Password) {
			panic("mysql.password is empty")
		}
		if IsEmpty(yamlConf2.Sweet.MySqlConfig.DbName) {
			panic("mysql.dbName is empty")
		}
		if yamlConf2.Sweet.MySqlConfig.MaxOpenConns == 0 {
			yamlConf2.Sweet.MySqlConfig.MaxOpenConns = 100
		}
		if yamlConf2.Sweet.MySqlConfig.MaxIdleConns == 0 {
			yamlConf2.Sweet.MySqlConfig.MaxIdleConns = 50
		}
	}

	if yamlConf2.Sweet.RedisConfig.Active {
		if IsEmpty(yamlConf2.Sweet.RedisConfig.Host) {
			panic("redis.host is empty")
		}
		if yamlConf2.Sweet.RedisConfig.Port == 0 {
			yamlConf2.Sweet.RedisConfig.Port = 6379
		}
		if yamlConf2.Sweet.RedisConfig.Database == 0 {
			yamlConf2.Sweet.RedisConfig.Database = 0
		}
	}

	if IsEmpty(yamlConf2.Sweet.Log.Level) {
		yamlConf2.Sweet.Log.Level = "info"
	}
	switch yamlConf2.Sweet.Log.Level {
	case "info":
		break
	case "warn":
		break
	case "error":
		break
	default:
		panic("log.level is error, must be info/warn/error")
	}

	if IsEmpty(yamlConf2.Sweet.Log.File) {
		yamlConf2.Sweet.Log.File = "logs/go-sweet.log"
	}
	if yamlConf2.Sweet.Log.MaxSize == 0 {
		yamlConf2.Sweet.Log.MaxSize = 10
	}
	if yamlConf2.Sweet.Log.MaxBackups == 0 {
		yamlConf2.Sweet.Log.MaxBackups = 10
	}
	if yamlConf2.Sweet.Log.MaxDays == 0 {
		yamlConf2.Sweet.Log.MaxDays = 7
	}

	if yamlConf2.Sweet.Img.Path == "" {
		panic("img.path is empty")
	}

	if IsEmpty(yamlConf2.Sweet.Img.MappingUrl) {
		yamlConf2.Sweet.Img.MappingUrl = "/static"
	}
	if IsEmpty(yamlConf2.Sweet.Img.BaseUrl) {
		yamlConf2.Sweet.Img.BaseUrl = fmt.Sprintf("http://localhost:%d", yamlConf2.Server.Port)
	}

	if IsEmpty(yamlConf2.Sweet.Img.Path) {
		panic("img.path is empty")
	}
	return yamlConf2
}

func saveConf(yamlConf2 YmlConfig) {
	if yamlConf2.Server.Port > 0 {
		yamlConf.Server.Port = yamlConf2.Server.Port
	} else {
		if yamlConf.Server.Port == 0 {
			yamlConf.Server.Port = 8080
		}
	}
	if IsNotEmpty(yamlConf2.Server.Name) {
		yamlConf.Server.Name = yamlConf2.Server.Name
	} else {
		if IsEmpty(yamlConf.Server.Name) {
			yamlConf.Server.Name = "go-sweet"
		}
	}

	yamlConf.Sweet.MySqlConfig.Active = yamlConf2.Sweet.MySqlConfig.Active

	if yamlConf.Sweet.MySqlConfig.Active {
		if IsNotEmpty(yamlConf2.Sweet.MySqlConfig.Host) {
			yamlConf.Sweet.MySqlConfig.Host = yamlConf2.Sweet.MySqlConfig.Host
		} else {
			if IsEmpty(yamlConf.Sweet.MySqlConfig.Host) {
				panic("mysql.host is empty")
			}
		}

		if yamlConf2.Sweet.MySqlConfig.Port > 0 {
			yamlConf.Sweet.MySqlConfig.Port = yamlConf2.Sweet.MySqlConfig.Port
		} else {
			if yamlConf.Sweet.MySqlConfig.Port == 0 {
				yamlConf.Sweet.MySqlConfig.Port = 3306
			}
		}
		if IsNotEmpty(yamlConf2.Sweet.MySqlConfig.User) {
			yamlConf.Sweet.MySqlConfig.User = yamlConf2.Sweet.MySqlConfig.User
		} else {
			if IsEmpty(yamlConf.Sweet.MySqlConfig.User) {
				yamlConf.Sweet.MySqlConfig.User = "root"
			}
		}

		if IsNotEmpty(yamlConf2.Sweet.MySqlConfig.Password) {
			yamlConf.Sweet.MySqlConfig.Password = yamlConf2.Sweet.MySqlConfig.Password
		} else {
			if IsEmpty(yamlConf.Sweet.MySqlConfig.Password) {
				panic("mysql.password is empty")
			}
		}

		if IsNotEmpty(yamlConf2.Sweet.MySqlConfig.DbName) {
			yamlConf.Sweet.MySqlConfig.DbName = yamlConf2.Sweet.MySqlConfig.DbName
		} else {
			if IsEmpty(yamlConf.Sweet.MySqlConfig.DbName) {
				panic("mysql.dbName is empty")
			}
		}

		if yamlConf2.Sweet.MySqlConfig.MaxOpenConns > 0 {
			yamlConf.Sweet.MySqlConfig.MaxOpenConns = yamlConf2.Sweet.MySqlConfig.MaxOpenConns
		} else {
			if yamlConf.Sweet.MySqlConfig.MaxOpenConns == 0 {
				yamlConf.Sweet.MySqlConfig.MaxOpenConns = 100
			}
		}
		if yamlConf2.Sweet.MySqlConfig.MaxIdleConns > 0 {
			yamlConf.Sweet.MySqlConfig.MaxIdleConns = yamlConf2.Sweet.MySqlConfig.MaxIdleConns
		} else {
			if yamlConf.Sweet.MySqlConfig.MaxIdleConns == 0 {
				yamlConf.Sweet.MySqlConfig.MaxIdleConns = 50
			}
		}

	}
	yamlConf.Sweet.RedisConfig.Active = yamlConf2.Sweet.RedisConfig.Active
	if yamlConf.Sweet.RedisConfig.Active {
		if IsNotEmpty(yamlConf2.Sweet.RedisConfig.Host) {
			yamlConf.Sweet.RedisConfig.Host = yamlConf2.Sweet.RedisConfig.Host
		} else {
			if IsEmpty(yamlConf.Sweet.RedisConfig.Host) {
				panic("redis.host is empty")
			}
		}

		if yamlConf2.Sweet.RedisConfig.Port > 0 {
			yamlConf.Sweet.RedisConfig.Port = yamlConf2.Sweet.RedisConfig.Port
		} else {
			if yamlConf.Sweet.RedisConfig.Port == 0 {
				yamlConf.Sweet.RedisConfig.Port = 6379
			}
		}

		if yamlConf2.Sweet.RedisConfig.Database > 0 {
			yamlConf.Sweet.RedisConfig.Database = yamlConf2.Sweet.RedisConfig.Database
		}

		if IsNotEmpty(yamlConf2.Sweet.RedisConfig.Password) {
			yamlConf.Sweet.RedisConfig.Password = yamlConf2.Sweet.RedisConfig.Password
		}
	}

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

	if IsNotEmpty(yamlConf2.Sweet.Img.MappingUrl) {
		yamlConf.Sweet.Img.MappingUrl = yamlConf2.Sweet.Img.MappingUrl
	} else {
		if IsEmpty(yamlConf.Sweet.Img.MappingUrl) {
			yamlConf.Sweet.Img.MappingUrl = "/static"
		}
	}

	if IsNotEmpty(yamlConf2.Sweet.Img.Path) {
		yamlConf.Sweet.Img.Path = yamlConf2.Sweet.Img.Path
	} else {
		if IsEmpty(yamlConf.Sweet.Img.Path) {
			panic("img.path is empty")
		}
	}

	if IsNotEmpty(yamlConf2.Sweet.Img.BaseUrl) {
		yamlConf.Sweet.Img.BaseUrl = yamlConf2.Sweet.Img.BaseUrl
	} else {
		if IsEmpty(yamlConf.Sweet.Img.BaseUrl) {
			yamlConf.Sweet.Img.BaseUrl = fmt.Sprintf("http://localhost:%d", yamlConf.Server.Port)
		}
	}

	if len(yamlConf2.Sweet.ExcUrl.Prefix) > 0 {
		yamlConf.Sweet.ExcUrl.Prefix = yamlConf2.Sweet.ExcUrl.Prefix
	}

	if len(yamlConf2.Sweet.ExcUrl.Full) > 0 {
		yamlConf.Sweet.ExcUrl.Full = yamlConf2.Sweet.ExcUrl.Full
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
