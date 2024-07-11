package utils

import (
	cache "github.com/PurpleScorpion/go-sweet-cache"
	"reflect"
	"strings"
)

func ValueObject(key string) interface{} {
	return getYamlValue(key)
}

func ValueInt(key string) int {
	val := getYamlValue(key)
	return val.(int)
}

func ValueInt64(key string) int64 {
	return int64(ValueInt(key))
}

func ValueInt32(key string) int32 {
	return int32(ValueInt(key))
}

func ValueFloat32(key string) float32 {
	return float32(ValueFloat64(key))
}

func ValueFloat64(key string) float64 {
	val := getYamlValue(key)
	return val.(float64)
}

func ValueBool(key string) bool {
	val := getYamlValue(key)
	return val.(bool)
}

func ValueString(key string) string {
	val := getYamlValue(key)
	return val.(string)
}

func ValueStringArr(key string) []string {
	val := getYamlValue(key)
	val2 := val.([]interface{})
	var arr []string
	for i := 0; i < len(val2); i++ {
		arr = append(arr, val2[i].(string))
	}
	return arr
}

func getYamlValue(key string) interface{} {
	if !(strings.HasPrefix(key, "${") && strings.HasSuffix(key, "}")) {
		panic("key must start with ${ and end with }")
	}
	key = key[2 : len(key)-1]
	arr := strings.Split(key, ".")
	ymlConf := getYmlConf("ymlConf")
	ymlConf2 := getYmlConf("ymlConf2")

	val := ymlConf[arr[0]]
	val2 := ymlConf2[arr[0]]
	if len(arr) == 1 {
		if val2 == nil {
			return val
		}
		return val2
	}

	for i := 1; i < len(arr); i++ {
		tmp := arr[i]
		if val != nil {
			v := val.(map[string]interface{})
			val = v[tmp]
		}
		if val2 != nil {
			v := val2.(map[string]interface{})
			val2 = v[tmp]
		}
	}
	if val2 == nil {
		return val
	}
	return val2
}

func getYamlValType(val interface{}) string {
	if val == nil {
		return "NULL"
	}
	return reflect.TypeOf(val).String()
}

func getYmlConf(key string) map[string]interface{} {
	val, err := cache.SweetCache.Get(key)
	if !err {
		return nil
	}
	return val.(map[string]interface{})
}
