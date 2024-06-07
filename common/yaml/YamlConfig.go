package sweetyml

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
	MappingUrl string `yaml:"mappingUrl"`
	Path       string `yaml:"path"`
	BaseUrl    string `yaml:"baseUrl"`
}

type ExcludeUrl struct {
	Prefix []string `yaml:"prefix"` // 前缀正则匹配
	Full   []string `yaml:"full"`   // 全路径匹配
}
