package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	//"github.com/smartwalle/alipay/v3"

	"flag"

	"zyz.com/m/config"
	"zyz.com/m/mysql"
	"zyz.com/m/redis"

	_ "zyz.com/m/redis"
)

func main() {

	configFile := flag.String("config", "config.yaml", "path to configuration file")
	// 解析命令行参数
	flag.Parse()
	log.Printf("%v", *configFile)

	config.InitConfig(*configFile)
	mysql.InitMysql()

	r := gin.Default()

	// 初始化一个 cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println("Error creating cookie jar:", err)
		return
	}

	// 从文件中读取 cookie 信息并添加到 cookie jar 中
	if err := loadCookies(jar, "cookies.txt"); err != nil {
		log.Println("Error loading cookies:", err)
		return
	}

	// 添加中间件来记录请求
	r.Use(requestCounterMiddleware())

	// 将 cookie jar 设置为全局中间件
	r.Use(func(c *gin.Context) {
		c.Set("cookiejar", jar)
		c.Next()
	})
	// 使用 cookie 作为会话存储引擎，密钥为随机字符串
	store := cookie.NewStore([]byte("myweb"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")

	registerRouter(r)

	go func() {
		r.Run(":80")
	}()
	if config.DefaultConfig.Online {
		r.RunTLS(":443", config.DefaultConfig.OnlineCSR, config.DefaultConfig.OnlineKEY)
	} else {
		r.RunTLS(":443", config.DefaultConfig.LocalCSR, config.DefaultConfig.LocalKey)
	}

}

func registerRouter(r *gin.Engine) {
	// 定义路由
	r.GET("/", IndexGet)
	//r.POST("/", handleIndexNew)
	r.GET("/y2b/parse", parseHandler)
	r.Static("/static", "./static")

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/pxy", ProxyHandler)

	// 访问指定文件
	r.GET("/.well-known/pki-validation/B5DCEF6CDAF508E79398C3354A6602F4.txt", func(c *gin.Context) {
		// 返回指定的 txt 文件
		c.File("templates/B5DCEF6CDAF508E79398C3354A6602F4.txt")
	})
	r.POST("/register", registerHandler)

	r.POST("/login", loginHandler)

	// 登出路由处理程序
	r.GET("/logout", func(c *gin.Context) {
		// 从会话中删除用户信息或将会话重置为初始状态
		session := sessions.Default(c)
		session.Clear() // 清空会话数据
		session.Save()  // 保存会话

		// 重定向用户到目标页面（例如主页或登录页面）
		c.Redirect(http.StatusFound, "/") // 重定向到主页
	})

	// Profile 路由处理程序
	r.GET("/profile", func(c *gin.Context) {
		// 从会话中获取用户名
		session := sessions.Default(c)
		username := session.Get("username")

		// 从数据库或其他存储中获取用户个人信息
		// 这里只是一个示例，您需要根据实际情况从数据库中获取用户信息
		userProfile := UserProfile{
			UserName:   username.(string),
			Email:      "user@example.com", // 示例：假设用户邮箱地址
			Total:      int(redis.GetCount(RequestKey)),
			TotalParse: int(redis.GetCount(ParseKey)),
			Day:        int(redis.GetTodayCount(RequestKey)),
			DayParse:   int(redis.GetTodayCount(ParseKey)),

			// 其他用户信息...
		}

		// 渲染 Profile 页面并将用户个人信息传递给模板
		c.HTML(http.StatusOK, "profile.html", gin.H{"userProfile": userProfile})
	})
	r.GET("/contact", func(c *gin.Context) {
		data := &IndexInfo{
			UserName: getUserName(c),
		}
		c.HTML(http.StatusOK, "contact.html", data)
	})

	r.GET("/about", func(c *gin.Context) {
		data := &IndexInfo{
			UserName: getUserName(c),
		}
		c.HTML(http.StatusOK, "about.html", data)
	})

	// 设置支付宝客户端
	/*aliPayClient, err := alipay.New(alipay.Config{
		AppId:        "your-app-id",
		NotifyUrl:    "your-notify-url",
		PrivateKey:   "your-private-key",
		AliPublicKey: "alipay-public-key",
		SignType:     "RSA2",
	})
	if err != nil {
		panic(err)
	}*/

	// 处理充值请求
	r.POST("/recharge", func(c *gin.Context) {
		// 获取用户提交的充值方式
		/*subscription := c.PostForm("subscription")
		if subscription == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择充值方式"})
			return
		}

		// 设置订单金额和标题
		var amount float64
		var subject string
		switch subscription {
		case "monthly":
			amount = 19.9
			subject = "月卡充值"
		case "yearly":
			amount = 199.9
			subject = "年卡充值"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择有效的充值方式"})
			return
		}

		// 发起支付请求
		param := &alipay.TradePagePay{}
		param.Subject = subject
		param.TotalAmount = fmt.Sprintf("%.2f", amount)
		param.OutTradeNo = "your-out-trade-no"
		param.ProductCode = "FAST_INSTANT_TRADE_PAY"
		url, err := aliPayClient.TradePagePay(param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "支付请求失败"})
			return
		}

		// 将用户重定向至支付宝支付页面
		c.Redirect(http.StatusFound, url.String())*/
	})

	// 处理支付回调
	r.POST("/callback", func(c *gin.Context) {
		// 处理支付宝支付成功后的回调
		// 更新用户账户余额等信息
	})
}
