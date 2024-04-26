// @Author Bing
// @Date 2024-04-24 1:00:00
// @Desc
package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config app inited")
}

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("failed to connect to MySQL:", err)
		return
	}
	fmt.Println("MySQL inited")
}
func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	var ctx = context.Background()      // 创建一个背景上下文
	pong, err := Red.Ping(ctx).Result() // 将ctx作为参数传递给Ping方法
	if err != nil {
		fmt.Println("failed to connect to Redis:", err)
	} else {
		fmt.Println("Redis inited", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish...", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Printf("Failed to publish: %s\n", err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	defer sub.Close()

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Printf("Failed to receive message: %s\n", err)
		return "", err
	}
	fmt.Printf("Subscribe: Received message from channel '%s': %s\n", channel, msg.Payload)
	return msg.Payload, err
}
