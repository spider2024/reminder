package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"os"
	"reminder/etc"
	"reminder/internal/logger"
	"reminder/internal/origin"
	"strings"
)

func main() {
	ctx := context.Background()
	userId, token, err := origin.Login(ctx, etc.AppConfig.ZenTao.UserName, etc.AppConfig.ZenTao.Password)
	if err != nil {
		logger.Log.ErrorF("login failed: %v", err)
		return
	}
	err = origin.Bugs(token, "16", userId)
	if err != nil {
		logger.Log.ErrorF("bugs:%+v\n", err)
		return
	}
	logger.Log.InfoF(token)
}

func init() {
	logSet()

	// Set the file name of the configuration file
	viper.SetConfigName("reminder")

	// Set the path to look for the configuration file
	viper.AddConfigPath("/Users/lgc/code/reminder/etc/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set the type of the configuration file
	viper.SetConfigType("yaml")

	viper.GetString("")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		logger.Log.ErrorF("Error reading config file, %s\n", err)
	}

	fmt.Printf("%+v\n", etc.AppConfig)
	if err := viper.Unmarshal(etc.AppConfig); err != nil {
		fmt.Printf("Error unmarshal config file, %s\n", err)
		return
	}
	replaceEnvVariables(etc.AppConfig)

	opt, err := redis.ParseURL("rediss://default:" + etc.AppConfig.Redis.Password + "@" + etc.AppConfig.Redis.Addr)
	if err != nil {
		fmt.Printf("redis url parse failed: %v\n", err)
		panic("redis parse url failed")
	}
	etc.Rdb = redis.NewClient(opt)
	//etc.Rdb = redis.NewClient(&redis.Options{
	//	Addr:     etc.AppConfig.Redis.Addr,
	//	Password: etc.AppConfig.Redis.Password,
	//	DB:       etc.AppConfig.Redis.DB,
	//})
	ping := etc.Rdb.Ping(context.Background())
	fmt.Printf("redis.test:%s\n", ping.String())
}

func logSet() {
	logFile, err := os.OpenFile("application.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.FatalF("无法打开日志文件: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			logger.Log.ErrorF("close file failed: %v", err)
		}
	}(logFile)

	// 创建自定义Logger，日志输出到控制台和文件
	logger.Log = logger.NewLogger(logger.INFO, os.Stdout, logFile)
}

func replaceEnvVariables(config *etc.Config) {
	// Replace Redis config environment variables
	config.Redis.Addr = os.ExpandEnv(config.Redis.Addr)
	config.Redis.Password = os.ExpandEnv(config.Redis.Password)

	// Replace ZenTao config environment variables
	config.ZenTao.UserName = os.ExpandEnv(config.ZenTao.UserName)
	config.ZenTao.Password = os.ExpandEnv(config.ZenTao.Password)
	config.ZenTao.Url = os.ExpandEnv(config.ZenTao.Url)

}
