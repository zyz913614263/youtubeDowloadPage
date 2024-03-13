package main

import (
	"log"
	"net/http/cookiejar"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	//"github.com/smartwalle/alipay/v3"

	"flag"

	"zyz.com/m/config"
	"zyz.com/m/mysql"
	"zyz.com/m/redis"
	"zyz.com/m/server"
)

func main() {

	configFile := flag.String("config", "config.yaml", "path to configuration file")
	// 解析命令行参数
	flag.Parse()
	log.Printf("%v", *configFile)

	config.InitConfig(*configFile)
	redis.InitRedis()
	mysql.InitMysql()

	r := gin.Default()

	// 初始化一个 cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println("Error creating cookie jar:", err)
		return
	}

	// 从文件中读取 cookie 信息并添加到 cookie jar 中
	if err := server.LoadCookies(jar, "cookies.txt"); err != nil {
		log.Println("Error loading cookies:", err)
		return
	}

	// 添加中间件来记录请求
	r.Use(server.RequestCounterMiddleware())

	// 将 cookie jar 设置为全局中间件
	r.Use(func(c *gin.Context) {
		c.Set("cookiejar", jar)
		c.Next()
	})
	// 使用 cookie 作为会话存储引擎，密钥为随机字符串
	store := cookie.NewStore([]byte("myweb"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")

	server.RegisterRouter(r)

	go func() {
		r.Run(":80")
	}()
	if config.DefaultConfig.Online {
		r.RunTLS(":443", config.DefaultConfig.OnlineCSR, config.DefaultConfig.OnlineKEY)
	} else {
		r.RunTLS(":443", config.DefaultConfig.LocalCSR, config.DefaultConfig.LocalKey)
	}

}
