package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"os"
	"reminder/etc"
	"reminder/internal/logger"
	"reminder/internal/origin"
)

func main() {
	ctx := context.Background()
	userId, token, err := origin.Login(ctx, etc.AppConfig.ZenTao.UserName, etc.AppConfig.ZenTao.Password)
	if err != nil {
		logger.ErrorF("login failed: %v", err)
		return
	}

	err = origin.Bugs(token, "53", userId)
	if err != nil {
		logger.ErrorF("bugs:%+v\n", err)
		return
	}
}

func init() {
	viper.SetConfigName("reminder")
	viper.AddConfigPath("/Users/lgc/code/reminder/etc/")
	viper.AutomaticEnv()
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s\n", err)
	}

	if err := viper.Unmarshal(etc.AppConfig); err != nil {
		log.Fatalf("Error unmarshal config file, %s\n", err)
	}
	replaceEnvVariables(etc.AppConfig)

	logger.InitLogger(os.Stdout, logger.OpenLogFile())
	logger.InfoF("server config: %+v", *etc.AppConfig)

	opt, err := redis.ParseURL("rediss://default:" + etc.AppConfig.Redis.Password + "@" + etc.AppConfig.Redis.Addr)
	if err != nil {
		logger.FatalF("redis url parse failed: %v\n", err)
	}
	etc.Rdb = redis.NewClient(opt)
	//etc.Rdb = redis.NewClient(&redis.Options{
	//	Addr:     etc.AppConfig.Redis.Addr,
	//	Password: etc.AppConfig.Redis.Password,
	//	DB:       etc.AppConfig.Redis.DB,
	//})
	ping := etc.Rdb.Ping(context.Background())
	logger.InfoF("redis.test:%s\n", ping.String())
}

func replaceEnvVariables(config *etc.Config) {
	config.Redis.Addr = os.ExpandEnv(config.Redis.Addr)
	config.Redis.Password = os.ExpandEnv(config.Redis.Password)

	config.ZenTao.UserName = os.ExpandEnv(config.ZenTao.UserName)
	config.ZenTao.Password = os.ExpandEnv(config.ZenTao.Password)
	config.ZenTao.Url = os.ExpandEnv(config.ZenTao.Url)
}
