package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	MysqlUsername string
	MysqlPassword string
	MysqlHost     string
	MysqlPort     int
	MysqlDBName   string

	IsProxy     bool   //是否使用代理
	Proxy       string //代理地址
	CookiesFile string //cookies文件路径
	Online      bool
	OnlineCSR   string
	OnlineKEY   string
	LocalCSR    string
	LocalKey    string
	NouseMysql  bool //测试使用，不加载mysql
	RedisAddr   string
	RedisPass   string
	RedisDB     int
}

var DefaultConfig = &Config{}

func InitConfig(cconfigFile string) {
	str := strings.Split(cconfigFile, ".")
	// 初始化 Viper
	viper.SetConfigName(str[0]) // 设置配置文件名（不含扩展名）
	viper.SetConfigType(str[1]) // 设置配置文件类型
	viper.AddConfigPath(".")    // 设置配置文件路径（当前目录）

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 将配置文件解析到结构体中
	if err := viper.Unmarshal(&DefaultConfig); err != nil {
		log.Fatalf("Error unmarshaling config file: %v", err)
	}

	log.Printf("config load %v", DefaultConfig)
}
