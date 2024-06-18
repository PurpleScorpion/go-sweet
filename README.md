# go-sweet

# 基于beego的web框架
## github地址 : https://github.com/PurpleScorpion/go-sweet
# 代码生成
```text
本项目可以使用go-sweet-generator进行代码生成 (开发中...)
可生成 (增/删/改/根据ID查询/分页列表查询)
    - controller
    - models
    - service
    - routers
同时可生成前端Vue3代码
```

# 使用前注意
```text
go version 1.20
使用前请先执行 go mod tidy 进行模块依赖下载
请仔细阅读配置文件说明
```

## 目录结构
```text
    - common # 公共模块
        - constants     # 常量包
        - logger        # 日志包
        - utils         # 工具包
        - vo            # vo包
        - yaml          # yaml配置文件处理包 - 一般无需理会
    - conf              # 配置文件包 - 内容参考java的SpringBoot
        - application.yml       # 以下配置文件将会在下面详细介绍
        - application-dev.yml
        - application-prod.yml
        - application-test.yml
    - controller                # controller 包
        - BaseController.go     # 基础controller
        - FileController.go     # 图片上传相关模块 可自行添加其他文件上传相关内容
        - SystemController.go   # 系统相关模块(权限控制)
        - UserController.go     # 登录相关模块
    - models            # 数据库模型包
    - routers           # 路由包
    - service           # 业务逻辑包
    - go.mod            # go.mod
    - main.go           # 程序主入口
```

## 关于项目的书写规范以及配置
- constants 常量包
```text
所有的常量应该遵循<阿里巴巴开发手册>中常量的命名规范
需要按照功能进行创建go文件进行管理
```
- logger 日志包
```text
使用beego的日志模块进行记录 , 在打印日志的同时会将日志记录到文件中
使用方式 logger.Info("日志内容")
```
- yaml 配置文件处理包
```text
该包用于处理conf文件夹下的yaml配置文件
通常无需复写
如有特殊需求 , 可在该包下进行相应功能的增加
```
- conf 配置文件包
```text
该包必须有application.yml文件
以及必须拥有至少1个 <application-环境名.yml> 命名的文件
该包的使用方式与java的SpringBoot类似 , 都是利用了习惯优于配置的思想
```
- application.yml
```text
该文件为默认配置文件 , 也可使用本文件进行配置所有内容 , 但是不建议这么做
因为这么做不符合开发规范
该文件下必须拥有以下配置
server:
  name: go-sweet
  active: dev
其中name为项目名称 - beego需要的配置
active为当前环境 - 用来激活其他配置文件使用 , 该配置的值通常为dev/test/prod
但是也可以自由定义 , 只需要符合命名公约即可
例如:
server:
  name: go-sweet
  active: haha
那么 , 程序在启动时会激活 <application-haha.yml> 的配置文件

application.yml文件的配置与优先级
application-环境名.yml配置值的优先级永远大于application.yml配置值
但是在application-环境名.yml的配置值中 , 不可以配置active的值 , 配置了也不会生效
除此之外 , 其他的配置皆可以覆盖application.yml的配置值
以下是全部配置值的释义
server:
  name: go-sweet # 项目名称                           选填: 默认值go-sweet
  active: dev    # 当前环境
sweet:
  mysql:
    # 必填
    active: true  # 是否激活使用mysql配置           
    # 必填       
    host: 192.168.253.130  # mysql地址     
    # 选填: 默认值3306                  
    port: 3306  # mysql端口号
    # 选填: 默认值root                            
    user: root  # mysql用户名 
    # 必填                    
    password: 123456  # mysql密码
    # 必填                              
    dbName: go_sweet_db  # mysql数据库名称
    # 选填: 默认值50
    maxIdleConns: 50 # 最大空闲连接数        
    # 选填: 默认值100        
    maxOpenConns: 100 # 最大打开连接数               

  redis:
    # 必填 
    active: true #是否激活使用redis配置  
    # 必填               
    host: 192.168.253.130 # redis地址
    # 选填: 默认值6379
    port: 6379    # redis端口号     
    # 选填: 若没有密码，则不用填该项
    password: 123456 # redis密码
    # 选填: 默认值0
    database: 2  # redis使用的数据库
  
  adx:
    # 必填 
    active: true #是否激活使用adx配置 
    # 必填               
    host: https://your.adx.com # adx地址 https开头
    # 选填: 默认值AAK 可选值: AAK/SMI
    # AAK : 全称 AadAppKey , 通过AppId/AppKey/AuthorityID认证
    # SMI : 全称 SystemManagedIdentity , 通过系统托管的标识认证
    # 当且仅当 authMethod为AAK时 AppId/AppKey/AuthorityID 必填
    authMethod: AAK
    appId: yourAppId
    appKey: yourAppKey
    authorityID: yourAuthorityId
    # 选填: 默认值false 是否开启日志
    logActive: true
  
  mqtt:
    # 必填 
    active: true #是否激活使用mqtt配置 
     # 必填               
    host: 192.168.253.130 # mqtt地址
    # 选填: 默认值1883
    port: 1883
    # 必填 用户名
    user: yourUser
    # 必填 密码
    password: yourPassword
    # 注意MQTT不光需要发送还需要监听 , 监听需要在common/yaml/MQTTServer/MqttOnline() 处配置
    # 该处已经有一个配置案例, 请根据需要自行修改
    
  # 日志配置
  logging:
    # 选填: 默认值info
    level: info # 日志级别 info/warn/error
    # 选填: 默认值log/go_sweet.log
    file: log/go_sweet.log # 日志文件路径
    # 选填: 默认值10
    maxSize: 10 # 单个日志文件最大大小，单位：MB
    # 选填: 默认值30
    maxDays: 30 # 单个日志文件保存天数
    # 选填: 默认值10
    maxBackups: 10 # 最大备份日志文件数量

  # 图片配置
  img:
    # 选填 : 默认值/static
    mappingUrl: /static # 映射的url 参考Tomcat映射图片地址配置
    # 必填
    path: D:/imgs
    # 选填: 默认值http://localhost:配置的端口号
    baseUrl: http://localhost:20001
  
  # 权限排除配置
  excludeUrl:
    # 前缀匹配排除
    prefix:
      - /static/ # 静态资源 一般必填
    # 完全匹配排除
    full:
      - /login # 登录接口 一般必填
```

- controller
```text
controller包下所有的controller均继承自BaseController
BaseController中封装了常用的方法 , 如: 
getUserId(ctx *context.Context) 获取当前登录用户的id
Result(r utils.R) 返回json格式的数据

根据编码规约
controller层尽量少编写业务逻辑
为了区分service层调用,所有与controller层相关的service方法都应该先定义一个service变量
例如:
SystemController.go中 
var systemService service.SystemService
```

- models
```text
models包下的所有结构体都应该以文件名进行分割 , 一个表一个结构体 , 一个表一个文件
每个文件中都应该包含以下结构
结构体的名称应与文件名一致 , 且名称为数据库表名的<大驼峰>写法
每一个结构体中都应该包含一个Id字段 , 并且在后方的tag中添加 tableId:"数据库表的主键字段名(列原名)" 
每一个结构体文件中都应该包含一个TableName方法 , 该方法返回表名
例如: 
type 名称 struct {
    Id int32 `json:"id" tableId:"id"`
    名称 string `json:"名称"`
}

func (名称) TableName() string {
	return "表名称"
}
```

- service
```text
service包下的所有结构体都应该以文件名进行分割 , 一种业务一个结构体 , 一个业务一个文件
且每个文件中都应该有与该文件名相同的struct
在service层中所有与controller层直接交互的方法名必须是大写字母开头 , 返回值都应当是 utils.R
且函数名前应当使用service层结构体指针
例如:
type UserService struct {
}
func (that *UserService) HealthCheck(id int32) utils.R {}

所有service层中由service包内部调用的方法均应该以小写字母开头
所有的公共调用的函数都应当写在CommonService.go中 , 且函数名前<不得使用>service层结构体指针
例如:
func GetCache(key string) string {}
```

- main.go
```text
该函数为整个项目的入口,运行项目时应该在cmd中进入该文件所在目录 , 并使用 ./main.go运行
该函数包含跨域配置/权限拦截/恐慌函数配置/yaml配置文件初始化
```