package service

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"go-sweet/common/constants"
)

// yml配置文件内容
var ymlConf = constants.YmlConf

func ServiceInit() {
	ymlConf = constants.YmlConf
}

func MqttOnlineTest(client MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	fmt.Println(payload)
}
