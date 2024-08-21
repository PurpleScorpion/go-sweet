package main

import (
	sweetyml "sweet-common/yaml"
	"sweet-src/main/golang/models"
	"sweet-src/main/golang/service"
	"testing"
)

/*
 * 单元测试使用:
 * 1. 文件名称必须以_test.go结尾
 * 2. 函数名称必须以Test开头
 * 3. 函数参数必须为*testing.T
 * 命令介绍: (如果你有类似于goland这种高级编译器会自动识别)
 * 0. 测试时需导航到测试包下执行以下命令,如果你有高级编译器,则无需理会这条
 * 1. go test 测试当前目录下的所有文件
 * 2. go test -run 函数名 测试指定的函数 , 例如: go test -run TestDemo1
 * 3. go test -bench=.   这会运行所有的性能测试，并输出每个测试的基准结果。
 * 4. go test -cover 测试代码覆盖率
 *
 *
 * 你还可以生成一个详细的覆盖率报告并输出到文件：
 * go test -coverprofile=coverage.out
 * go tool cover -html=coverage.out -o coverage.html
 */
var userService service.UserService

func init() {
	sweetyml.ReadYml()
}

func TestDemo1(t *testing.T) {
	var user models.User
	user.Username = "root"
	user.Password = "root"
	r := userService.Login(user)
	if r.Code == 200 {
		t.Log("测试通过: ", r.Data)
	} else {
		t.Error("测试失败")
	}
}
