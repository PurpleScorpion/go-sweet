package constants

var YmlConf YmlConfig

type YmlConfig struct {
	Server Server `yaml:"server"`
	Sweet  Sweet  `yaml:"sweet"`
}

type Server struct {
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
	Active string `yaml:"active"`
}

type Sweet struct {
	MySqlConfig MySqlConf  `yaml:"mysql"`
	RedisConfig RedisConf  `yaml:"redis"`
	Adx         Adx        `yaml:"adx"`
	Mqtt        Mqtt       `yaml:"mqtt"`
	Log         Logging    `yaml:"logging"`
	Img         Images     `yaml:"img"`
	ExcUrl      ExcludeUrl `yaml:"excludeUrl"`
}

type MySqlConf struct {
	Active       bool   `yaml:"active"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DbName       string `yaml:"dbName"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

type RedisConf struct {
	Active   bool   `yaml:"active"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type Logging struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxDays    int    `yaml:"maxDays"`
}

type Images struct {
	Active     bool   `yaml:"active"`
	MappingUrl string `yaml:"mappingUrl"`
	Path       string `yaml:"path"`
	BaseUrl    string `yaml:"baseUrl"`
}

type Mqtt struct {
	Active   bool   `yaml:"active"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Adx struct {
	Active      bool   `yaml:"active"`
	AuthMethod  string `yaml:"authMethod"` // 认证方式  AAK(AadAppKey): 通过AppId/AppKey/AuthorityID认证  SMI(SystemManagedIdentity): 通过系统托管的标识认证
	Host        string `yaml:"host"`
	AppId       string `yaml:"appId"`
	AppKey      string `yaml:"appKey"`
	AuthorityID string `yaml:"authorityID"`
	LogActive   bool   `yaml:"logActive"` // 是否激活日志
}

type ExcludeUrl struct {
	Prefix []string `yaml:"prefix"` // 前缀正则匹配
	Full   []string `yaml:"full"`   // 全路径匹配
}
