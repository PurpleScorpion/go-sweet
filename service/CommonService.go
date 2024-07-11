package service

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"go-sweet/common/utils"
)

// yml配置文件内容
var (
	imgBaseUrl string = ""
)

func ServiceInit() {
	imgBaseUrl = utils.ValueString("${sweet.img.baseUrl}")
}

func MqttOnlineTest(client MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	fmt.Println(payload)
}
