package utils

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient MQTT.Client

/*
topic: 主题 例如: info/user
jsonData: json数据 , 建议使用 jsonutils.JSONObject 来构建

JSONObject使用示例:

	 	js := jsonutil.NewJSONObject()
		js.FluentPut("aa", 12)
		js.FluentPut("bb", "cc")
		utils.MQTTSend("info/user", js.GetData())
*/
func MQTTSend(topic string, jsonData interface{}) {
	jsonBytes, _ := json.Marshal(jsonData)
	jsonString := string(jsonBytes)
	token := MQTTClient.Publish(topic, 0, false, jsonString)
	token.Wait()
}
