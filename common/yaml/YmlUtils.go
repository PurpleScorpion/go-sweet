package sweetyml

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"sweet-common/constants"
	"sweet-common/logger"
	"time"
)

func Init() {
	start := time.Now().UnixMilli()
	printBanner()
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""

	if keqing.IsEmpty(profilesActive) {
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
	var conf1 = make(map[string]interface{})
	yaml.Unmarshal(data, &conf1)
	conf2 := readChildConf(conf1)
	keqing.LoadYml(conf1, conf2)
	serverActive := keqing.ValueString("${server.active}")
	if keqing.IsEmpty(serverActive) {
		panic("No active profiles found. Please specify a profile using the 'server.active' property.")
	}

	port := keqing.ValueInt("${server.port}")
	if port <= 0 || port > 65535 {
		panic("Invalid port number: " + fmt.Sprintf("%d", port))
	}

	logs.Info("The following profiles are active: %s", serverActive)
	logs.Info("Golang server initialized with port(s): %d (http)", port)
	initServer()
	logs.Info("Starting service [Golang server]")
	end := time.Now().UnixMilli()
	logs.Info("Server started on port(s): %d (http) with context path '/'", port)
	logs.Info("Started Application in %d millisecond", (end - start))
}

func initServer() {
	logger.LoggerInit()
	initMySQL()
	initRedis()
	initAdx()
	initMqtt()
	beegoInit()
}

func printBanner() {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""
	if keqing.IsEmpty(profilesActive) {
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

func readChildConf(parentConf map[string]interface{}) map[string]interface{} {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
	confPath := ""
	serverActive := getEnvActive(parentConf)

	if keqing.IsEmpty(profilesActive) {
		confPath = "src/main/resources/application-" + serverActive + ".yml"
		absPath, err := filepath.Abs(confPath)
		if err != nil {
			log.Fatalf("failed to get absolute path: %v", err)
		}
		confPath = absPath
	} else {
		confPath = "conf/application-" + serverActive + ".yml"
	}
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic("Error reading configuration file: " + err.Error())
	}
	var conf1 = make(map[string]interface{})
	yaml.Unmarshal(data, &conf1)
	return conf1
}

func getEnvActive(parentConf map[string]interface{}) string {
	profilesActive := os.Getenv(constants.PROFILES_ACTIVE)

	server := parentConf["server"].(map[string]interface{})
	active := server["active"].(string)

	confPath := ""
	if keqing.IsEmpty(profilesActive) {
		// 如果docker环境变量为空, 则认为是本地环境
		confPath = "src/main/resources/application-" + active + ".yml"
	} else {
		// 否则是docker环境
		confPath = "conf/application-" + profilesActive + ".yml"
	}

	if keqing.IsNotEmpty(profilesActive) {
		absPath, err := filepath.Abs(confPath)
		if err == nil {
			_, err1 := os.ReadFile(absPath)
			if err1 == nil {
				return profilesActive
			}
		}
	}
	return active
}
