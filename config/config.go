package config

import (
	"log"

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
}

var DefaultConfig = &Config{}

func InitConfig() {
	// 初始化 Viper
	viper.SetConfigName("config") // 设置配置文件名（不含扩展名）
	viper.SetConfigType("yaml")   // 设置配置文件类型
	viper.AddConfigPath(".")      // 设置配置文件路径（当前目录）

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 从配置文件中读取 MySQL 相关字段
	DefaultConfig.MysqlUsername = viper.GetString("mysql.username")
	DefaultConfig.MysqlPassword = viper.GetString("mysql.password")
	DefaultConfig.MysqlHost = viper.GetString("mysql.host")
	DefaultConfig.MysqlPort = viper.GetInt("mysql.port")
	DefaultConfig.MysqlDBName = viper.GetString("mysql.dbname")
	DefaultConfig.IsProxy = viper.GetBool("isproxy")
	DefaultConfig.Proxy = viper.GetString("proxy")
	DefaultConfig.CookiesFile = viper.GetString("cookie")
	DefaultConfig.Online = viper.GetBool("online")
	DefaultConfig.OnlineCSR = viper.GetString("OnlineCSR")
	DefaultConfig.OnlineKEY = viper.GetString("OnlineKEY")
	DefaultConfig.LocalCSR = viper.GetString("LocalCSR")
	DefaultConfig.LocalKey = viper.GetString("LocalKey")

	log.Printf("config load %v", DefaultConfig)
}
