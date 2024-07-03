package etc

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

var AppConfig Config

type Config struct {
	Redis   RedisCnf  `yaml:"redis"`
	ZenTao  ZenTaoCnf `yaml:"zen_tao"`
	Project Project   `yaml:"project"`
}

type RedisCnf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type ZenTaoCnf struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Url      string `yaml:"url"`
}

type Project struct {
	ProjectName string `yaml:"project_name"`
	Id          string `yaml:"id"`
	Weight      int    `yaml:"weight" default:"0"`
}
