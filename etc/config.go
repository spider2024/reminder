package etc

import "github.com/redis/go-redis/v9"

var AppConfig *Config
var Rdb *redis.Client

type Config struct {
	Redis   *RedisCnf
	ZenTao  *ZenTaoCnf
	Project *Project
}

type RedisCnf struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type ZenTaoCnf struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

type Project struct {
	ProjectName string `json:"project_name"`
	Id          string `json:"id"`
	Weight      int    `json:"weight" default:"0"`
}
