package sweetyml

import (
	"fmt"
	cache "github.com/PurpleScorpion/go-sweet-cache"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sweet-common/constants"
	"time"
)

var yamlConf YmlConfig

func ReadYml() {
	start := time.Now().UnixMilli()
	printBanner()
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""

	if IsEmpty(profilesActive) {
		confPath = "src/main/resources/application.yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
		}
		confPath = absPath
	} else {
		confPath = "conf/application.yml"
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
	var conf1 = make(map[string]interface{})
	yaml.Unmarshal(data, &conf1)
	cache.New(cache.NoExpiration, cache.NoExpiration)
	cache.SweetCache.Set("ymlConf", conf1, cache.NoExpiration)

	readChildConf()
	logs.Info("The following profiles are active: %s", yamlConf.Server.Active)
	logs.Info("Golang server initialized with port(s): %d (http)", yamlConf.Server.Port)
	initServer()
	logs.Info("Starting service [Golang server]")
	end := time.Now().UnixMilli()
	logs.Info("Server started on port(s): %d (http) with context path '/'", yamlConf.Server.Port)
	logs.Info("Started Application in %d millisecond", (end - start))
}

func initServer() {
	initMySQL()
	initRedis()
	initAdx()
	initMqtt()
}

func printBanner() {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""
	if IsEmpty(profilesActive) {
		confPath = "src/main/resources/banner.txt"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
		}
		confPath = absPath
	} else {
		confPath = "conf/banner.txt"
	}

	file, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Println("===========================================================================")
		fmt.Println(" ")
		var str = "   _____                            _   \n  / ____|                          | |  \n | (___   __      __   ___    ___  | |_ \n  \\___ \\  \\ \\ /\\ / /  / _ \\  / _ \\ | __|\n  ____) |  \\ V  V /  |  __/ |  __/ | |_ \n |_____/    \\_/\\_/    \\___|  \\___|  \\__|\n                                        \n                                        "
		fmt.Println(str)
		fmt.Println("===========================================================================")
		return
	}
	fmt.Println(string(file))
}

func GetYmlConf() YmlConfig {
	return yamlConf
}
func getEnvActive() string {
	//profilesActive := "dev"
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""
	if IsEmpty(profilesActive) {
		confPath = "src/main/resources/application-" + profilesActive + ".yml"
	} else {
		confPath = "conf/application-" + profilesActive + ".yml"
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
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""
	yamlConf.Server.Active = getEnvActive()

	if IsEmpty(profilesActive) {
		confPath = "src/main/resources/application-" + yamlConf.Server.Active + ".yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
		}
		confPath = absPath
	} else {
		confPath = "conf/application-" + yamlConf.Server.Active + ".yml"
	}
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic("Error reading configuration file: " + err.Error())
	}
	err = yaml.Unmarshal(data, &yamlConf2)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}
	var conf1 = make(map[string]interface{})
	yaml.Unmarshal(data, &conf1)
	cache.SweetCache.Set("ymlConf2", conf1, cache.NoExpiration)
	//yamlConf2 = defaultData(yamlConf2)
	saveConf(yamlConf2)
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

	if yamlConf.Sweet.MySqlConfig.Active || yamlConf2.Sweet.MySqlConfig.Active {
		yamlConf.Sweet.MySqlConfig.Active = true
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
	if yamlConf.Sweet.RedisConfig.Active || yamlConf2.Sweet.RedisConfig.Active {
		yamlConf.Sweet.RedisConfig.Active = true
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

	if yamlConf.Sweet.Img.Active || yamlConf2.Sweet.Img.Active {
		yamlConf.Sweet.Img.Active = true
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
	}

	if len(yamlConf2.Sweet.ExcUrl.Prefix) > 0 {
		yamlConf.Sweet.ExcUrl.Prefix = yamlConf2.Sweet.ExcUrl.Prefix
	}

	if len(yamlConf2.Sweet.ExcUrl.Full) > 0 {
		yamlConf.Sweet.ExcUrl.Full = yamlConf2.Sweet.ExcUrl.Full
	}

	if yamlConf.Sweet.Adx.Active || yamlConf2.Sweet.Adx.Active {
		yamlConf.Sweet.Adx.Active = true
		if IsNotEmpty(yamlConf2.Sweet.Adx.Host) {
			yamlConf.Sweet.Adx.Host = yamlConf2.Sweet.Adx.Host
		} else {
			if IsEmpty(yamlConf.Sweet.Adx.Host) {
				panic("adx.host is empty")
			}
		}
		if yamlConf.Sweet.Adx.LogActive || yamlConf2.Sweet.Adx.LogActive {
			yamlConf.Sweet.Adx.LogActive = true
		}

		if IsNotEmpty(yamlConf2.Sweet.Adx.AuthMethod) {
			yamlConf.Sweet.Adx.AuthMethod = yamlConf2.Sweet.Adx.AuthMethod
		} else {
			if IsEmpty(yamlConf.Sweet.Adx.AuthMethod) {
				yamlConf.Sweet.Adx.AuthMethod = "AAK"
			}
		}
		switch yamlConf.Sweet.Adx.AuthMethod {
		case "AAK":
			break
		case "SMI":
			break
		default:
			panic("Adx AuthMethod is error, must be AAK/SMI")
		}

		if yamlConf.Sweet.Adx.AuthMethod == "AAK" {
			if IsNotEmpty(yamlConf2.Sweet.Adx.AppId) {
				yamlConf.Sweet.Adx.AppId = yamlConf2.Sweet.Adx.AppId
			} else {
				if IsEmpty(yamlConf.Sweet.Adx.AppId) {
					panic("adx.appId is empty")
				}
			}

			if IsNotEmpty(yamlConf2.Sweet.Adx.AppKey) {
				yamlConf.Sweet.Adx.AppKey = yamlConf2.Sweet.Adx.AppKey
			} else {
				if IsEmpty(yamlConf.Sweet.Adx.AppKey) {
					panic("adx.appKey is empty")
				}
			}

			if IsNotEmpty(yamlConf2.Sweet.Adx.AuthorityID) {
				yamlConf.Sweet.Adx.AuthorityID = yamlConf2.Sweet.Adx.AuthorityID
			} else {
				if IsEmpty(yamlConf.Sweet.Adx.AuthorityID) {
					panic("adx.authorityID is empty")
				}
			}
		}
	}

	if yamlConf.Sweet.Mqtt.Active || yamlConf2.Sweet.Mqtt.Active {
		yamlConf.Sweet.Mqtt.Active = true
		if IsNotEmpty(yamlConf2.Sweet.Mqtt.Host) {
			yamlConf.Sweet.Mqtt.Host = yamlConf2.Sweet.Mqtt.Host
		} else {
			if IsEmpty(yamlConf.Sweet.Mqtt.Host) {
				panic("mqtt.host is empty")
			}
		}
		if yamlConf2.Sweet.Mqtt.Port > 0 {
			yamlConf.Sweet.Mqtt.Port = yamlConf2.Sweet.Mqtt.Port
		} else {
			if yamlConf.Sweet.Mqtt.Port == 0 {
				yamlConf.Sweet.Mqtt.Port = 1883
			}
		}

		if IsNotEmpty(yamlConf2.Sweet.Mqtt.User) {
			yamlConf.Sweet.Mqtt.User = yamlConf2.Sweet.Mqtt.User
		} else {
			if IsEmpty(yamlConf.Sweet.Mqtt.User) {
				panic("mqtt.user is empty")
			}
		}
		if IsNotEmpty(yamlConf2.Sweet.Mqtt.Password) {
			yamlConf.Sweet.Mqtt.Password = yamlConf2.Sweet.Mqtt.Password
		} else {
			if IsEmpty(yamlConf.Sweet.Mqtt.Password) {
				panic("mqtt.password is empty")
			}
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
