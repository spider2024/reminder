package etc

import (
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

var AppConfig = &Config{}

type Config struct {
	Redis   RedisCnf   `yaml:"redis"`
	ZenTao  ZenTaoCnf  `yaml:"zenTao"`
	Project Project    `yaml:"project"`
	Server  ServerConf `yaml:"server"`
}

type RedisCnf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type ZenTaoCnf struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Url      string `yaml:"url"`
}

type Project struct {
	ProjectName string `yaml:"projectName"`
	Id          string `yaml:"id"`
	Weight      int    `yaml:"weight" default:"0"`
}

type ServerConf struct {
	LogPath  string `yaml:"logPath"`
	LogName  string `yaml:"logName"`
	LogExt   string `yaml:"logExt"`
	LogLevel string `yaml:"logLevel"`
}
