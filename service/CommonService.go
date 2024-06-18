package service

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

func MqttOnlineTest(client MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	fmt.Println(payload)
}
