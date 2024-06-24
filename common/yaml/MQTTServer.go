package sweetyml

import (
	"crypto/tls"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go-sweet/common/utils"
	"go-sweet/service"
	"math/rand"
	"net/url"
	"time"
)

var (
	connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
		logs.Info("MQTT Connected")
		MqttOnline()
	}
	connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
		logs.Error("MQTT Connect lost: %s", err.Error())
	}
	reconnectHandler MQTT.ReconnectHandler = func(client MQTT.Client, co *MQTT.ClientOptions) {
		logs.Info("ReconnectHandler: %v", co)
	}

	connectionAttemptHandler MQTT.ConnectionAttemptHandler = func(aa *url.URL, tlsCfg *tls.Config) *tls.Config {
		logs.Info("ConnectionAttemptHandler: %v", tlsCfg)
		return tlsCfg
	}
)

func initMqtt() {
	conf := GetYmlConf()
	if !conf.Sweet.Mqtt.Active {
		return
	}
	logs.Info("Init MQTT....")
	mqttServer := conf.Sweet.Mqtt.Host
	mqttPort := conf.Sweet.Mqtt.Port
	mqttUsername := conf.Sweet.Mqtt.User     // cooldesign
	mqttPassword := conf.Sweet.Mqtt.Password // Q0^@RG4DqzvK3#vd

	str := fmt.Sprintf("mqtt://%s:%d", mqttServer, mqttPort)
	options := MQTT.NewClientOptions().AddBroker(str) // Replace with your MQTT broker details
	options.SetClientID(randomString(10))
	options.SetUsername(mqttUsername) // Replace with your MQTT username (if required)
	options.SetPassword(mqttPassword) // Replace with your MQTT password (if required)
	options.SetAutoReconnect(true)
	options.SetWriteTimeout(3 * time.Second)
	options.SetMaxReconnectInterval(10 * time.Second)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler
	options.OnConnectAttempt = connectionAttemptHandler
	options.OnReconnecting = reconnectHandler
	// Create and start the MQTT client
	utils.MQTTClient = MQTT.NewClient(options)
	if token := utils.MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		logs.Error("Error connecting to broker: %s", token.Error().Error())
	}
	for {
		if !utils.MQTTClient.IsConnected() {
			logs.Info("Disconnected from broker. Reconnecting...")
			token := utils.MQTTClient.Connect()
			token.Wait()
			time.Sleep(5 * time.Second)
		} else {
			logs.Info("break connect mqtt")
			break
		}

	}
}

func MqttOnline() {
	conf := GetYmlConf()
	if !conf.Sweet.Mqtt.Active {
		return
	}
	// 保证日志模块加载完毕
	time.Sleep(1 * time.Second)
	if token := utils.MQTTClient.Subscribe("info/#", 0, service.MqttOnlineTest); token.Wait() && token.Error() != nil {
		return
	}

}

/*
获取随机字符串
*/
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
