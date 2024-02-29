package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"

	//"github.com/smartwalle/alipay/v3"
	"golang.org/x/crypto/bcrypt"
	"zyz.com/m/config"
	"zyz.com/m/mysql"
)

func main() {

	config.InitConfig()
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

	// 将 cookie jar 设置为全局中间件
	r.Use(func(c *gin.Context) {
		c.Set("cookiejar", jar)
		c.Next()
	})
	// 使用 cookie 作为会话存储引擎，密钥为随机字符串
	store := cookie.NewStore([]byte("myweb"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")

	// 定义路由
	r.GET("/", IndexGet)
	r.POST("/", handleIndex)
	r.Static("/static", "./static")

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

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
			UserName: username.(string),
			Email:    "user@example.com", // 示例：假设用户邮箱地址
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

	go func() {
		r.Run(":80")
	}()
	if config.DefaultConfig.Online {
		r.RunTLS(":443", config.DefaultConfig.OnlineCSR, config.DefaultConfig.OnlineKEY)
	} else {
		r.RunTLS(":443", config.DefaultConfig.LocalCSR, config.DefaultConfig.LocalKey)
	}

}

type UserProfile struct {
	UserName string
	Email    string
	// 其他用户信息字段...
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 查询用户是否存在
	var hashedPassword string
	err := mysql.DefaultDB.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, "登录失败")
		fmt.Println("Error querying user:", err)
		return
	}

	// 比较密码哈希
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		c.String(http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	// 将用户名保存到会话中
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()

	//c.String(http.StatusOK, "登录成功")

	c.Redirect(http.StatusSeeOther, "/")
}

func registerHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	// 检查用户名是否已存在
	var count int
	err := mysql.DefaultDB.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&count)
	if err != nil {
		c.String(http.StatusInternalServerError, "注册失败")
		fmt.Println("Error checking username:", err)
		return
	}
	if count > 0 {
		c.String(http.StatusBadRequest, "用户名已存在")
		return
	}

	err = mysql.DefaultDB.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", email).Scan(&count)
	if err != nil {
		c.String(http.StatusInternalServerError, "注册失败")
		fmt.Println("Error checking username:", err)
		return
	}
	if count > 0 {
		c.String(http.StatusBadRequest, "邮箱已存在")
		return
	}
	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "注册失败")
		fmt.Println("Error hashing password:", err)
		return
	}

	// 插入新用户数据
	_, err = mysql.DefaultDB.Exec("INSERT INTO user (username, password, email) VALUES (?, ?, ?)", username, hashedPassword, email)
	if err != nil {
		c.String(http.StatusInternalServerError, "注册失败")
		fmt.Println("Error inserting user:", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

type Link struct {
	Quality  string
	URL      string
	MimeType string
	Width    int
	Height   int
}
type IndexInfo struct {
	VideoURL   string
	VideoLinks []*Link
	AudioLinks []*Link
	UserName   string
}

func getUserName(c *gin.Context) string {
	// 从请求上下文中获取会话对象
	session := sessions.Default(c)

	// 使用 Get 方法从会话中获取存储的值（这里是用户名）
	usernameInterface := session.Get("username")

	// 将 interface{} 类型的 username 转换为 string
	var username string
	if usernameInterface != nil {
		username = usernameInterface.(string)
	} else {
		// 如果用户名为空，则将其设置为空字符串
		username = ""
	}

	return username
}

func IndexGet(c *gin.Context) {

	data := &IndexInfo{
		UserName: getUserName(c),
	}
	c.HTML(http.StatusOK, "index.html", data)
	return
}

func handleIndex(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request
	username := getUserName(c)
	if username == "" {
		// 如果用户名为空，则将其设置为空字符串
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	videoURL := r.FormValue("url")
	if videoURL == "" || !strings.Contains(videoURL, "youtube") {
		http.Error(w, "Invalid YouTube video URL", http.StatusBadRequest)
		return
	}
	//是否选择使用代理
	transport := &http.Transport{}
	if config.DefaultConfig.IsProxy {
		transport.Proxy = func(request *http.Request) (*url.URL, error) {
			// 设置代理服务器的地址
			proxyURL, err := url.Parse(config.DefaultConfig.Proxy)
			if err != nil {
				return nil, err
			}
			return proxyURL, nil
		}
	}

	// 创建一个 cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println("Error creating cookie jar:", err)
		return
	}
	// 从文件中读取cookie信息
	if err := loadCookies(jar, config.DefaultConfig.CookiesFile); err != nil {
		log.Println("Error loading cookies:", err)
		return
	}

	tran := &http.Client{
		Jar:       jar,
		Transport: transport,
	}
	client := youtube.Client{
		HTTPClient: tran,
	}
	var video *youtube.Video

	for retries := 3; retries > 0; retries-- {
		video, err = func() (*youtube.Video, error) {
			video, err := client.GetVideo(videoURL)
			if err != nil {
				return nil, err
			}
			return video, nil

		}()
		if err == nil {
			break
		}
		log.Printf("Error getting video info: %v", err)
		time.Sleep(1 * time.Second) // 延迟1秒后重试
	}
	if err != nil {
		log.Printf("Error getting video info: %v", err)
		http.Error(w, "Error getting video info", http.StatusInternalServerError)
		return
	}

	var videoLinks, audioLinks []*Link
	for _, format := range video.Formats {
		stype := strings.Split(format.MimeType, ";")
		if len(stype) < 2 {
			continue
		}
		link := &Link{format.Quality, format.URL, stype[0], format.Width, format.Height}
		if strings.Contains(stype[0], "audio") {
			audioLinks = append(audioLinks, link)
		} else {
			videoLinks = append(videoLinks, link)
		}
	}
	/*sort.Slice(videoLinks, func(i, j int) bool {
		return videoLinks[i].MimeType > videoLinks[j].MimeType // 降序排列
	})

	sort.Slice(audioLinks, func(i, j int) bool {
		return audioLinks[i].MimeType > audioLinks[j].MimeType // 降序排列
	})*/

	data := &IndexInfo{
		VideoURL:   videoURL,
		VideoLinks: videoLinks,
		AudioLinks: audioLinks,
		UserName:   username,
	}
	//log.Printf("data=%v", data)
	c.HTML(http.StatusOK, "index.html", data)
}

func loadCookies(jar *cookiejar.Jar, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 7 {
			continue
		}

		cookie := &http.Cookie{
			Name:   parts[5],
			Value:  parts[6],
			Path:   parts[2],
			Domain: parts[0],
		}
		url := &url.URL{
			Scheme: "http", // 您可以根据需要调整协议
			Host:   parts[0],
		}
		jar.SetCookies(url, []*http.Cookie{cookie})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
