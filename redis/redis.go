package redis

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"zyz.com/m/config"
)

var Client *redis.Client

func InitRedis() {

	Client = redis.NewClient(&redis.Options{
		Addr:     config.DefaultConfig.RedisAddr,
		Password: config.DefaultConfig.RedisPass, // 密码，如果有的话
		DB:       config.DefaultConfig.RedisDB,   // 使用默认的数据库
	})

	// 检查Redis连接是否正常
	_, err := Client.Ping(Client.Context()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		os.Exit(1)
	}
}

func AddCount(key string) {
	// 每次访问计数器加一
	err := Client.Incr(Client.Context(), key).Err()
	if err != nil {
		log.Fatal("Error incrementing visit count:", err)
		return
	}
	AddToday(key)
}
func GetCount(key string) int {
	// 输出当前的访问量
	visits, err := Client.Get(Client.Context(), key).Result()
	if err != nil {
		fmt.Println("Error getting visit count:", err)
		return 0
	}
	v, _ := strconv.Atoi(visits)
	return v
}

func getToday() string {
	currentTime := time.Now()

	// 格式化当前时间为日期字符串
	today := currentTime.Format("2006-01-02")
	return today
}

func AddToday(key string) {
	key = getToday() + key
	err := Client.Incr(Client.Context(), key).Err()
	if err != nil {
		log.Fatal("Error incrementing visit count:", err)
		return
	}
}

func GetTodayCount(key string) int {
	key = getToday() + key
	// 输出当前的访问量
	visits, err := Client.Get(Client.Context(), key).Result()
	if err != nil {
		fmt.Println("Error getting visit count:", err)
		return 0
	}
	v, _ := strconv.Atoi(visits)
	return v
}
