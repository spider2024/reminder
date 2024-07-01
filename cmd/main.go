package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"reminder/etc"
	"reminder/internal/origin"
)

func main() {
	ctx := context.Background()
	userId, token, err := origin.Login(ctx, "", "")
	if err != nil {
		return
	}
	origin.Bugs("", "1", userId)
	fmt.Println(token)
}

func init() {
	// Set the file name of the configuration file
	viper.SetConfigName("reminder")

	// Set the path to look for the configuration file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Set the type of the configuration file
	viper.SetConfigType("yml")

	// Read the configuration file
	if err := viper.Unmarshal(etc.AppConfig); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return
	}

	etc.Rdb = redis.NewClient(&redis.Options{
		Addr:     etc.AppConfig.Redis.Addr,
		Password: etc.AppConfig.Redis.Password,
		DB:       etc.AppConfig.Redis.DB,
	})

}
