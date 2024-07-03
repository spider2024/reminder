package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"os"
	"reminder/etc"
	"reminder/internal/origin"
	"strings"
)

func main() {
	ctx := context.Background()
	userId, token, err := origin.Login(ctx, "", "")
	if err != nil {
		fmt.Printf("login failed: %v", err)
		return
	}
	err = origin.Bugs(token, "16", userId)
	if err != nil {
		fmt.Printf("bugs:%+v\n", err)
		return
	}
	fmt.Println(token)
}

func init() {
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
		log.Fatalf("Error reading config file, %s\n", err)
	}

	if err := viper.Unmarshal(&etc.AppConfig); err != nil {
		fmt.Printf("Error unmarshal config file, %s\n", err)
		return
	}
	replaceEnvVariables(&etc.AppConfig)

	fmt.Printf("ZenTao Config: %+v\n", etc.AppConfig.ZenTao)
	fmt.Printf("Redis Config: %+v\n", etc.AppConfig.Redis)

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

func replaceEnvVariables(config *etc.Config) {
	// Replace Redis config environment variables
	config.Redis.Addr = os.ExpandEnv(config.Redis.Addr)
	config.Redis.Password = os.ExpandEnv(config.Redis.Password)

	// Replace ZenTao config environment variables
	config.ZenTao.UserName = os.ExpandEnv(config.ZenTao.UserName)
	config.ZenTao.Password = os.ExpandEnv(config.ZenTao.Password)
	config.ZenTao.Url = os.ExpandEnv(config.ZenTao.Url)

	// Replace Project config environment variables if any
	// config.Project.ProjectName = os.ExpandEnv(config.Project.ProjectName)
	// config.Project.Id = os.ExpandEnv(config.Project.Id)
}
