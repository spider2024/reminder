package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"reminder/etc"
	"reminder/internal/origin"
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

	// Set the type of the configuration file
	viper.SetConfigType("yml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s\n", err)
	}

	if err := viper.Unmarshal(&etc.AppConfig); err != nil {
		fmt.Printf("Error unmarshal config file, %s\n", err)
		return
	}

	opt, err := redis.ParseURL("rediss://default:" + etc.AppConfig.Redis.Password + "@" + etc.AppConfig.Redis.Addr)
	if err != nil {
		fmt.Printf("redis url parse failed: %v\n", err)
		return
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
