package server

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"zyz.com/m/mysql"
	"zyz.com/m/redis"
)

func RegisterRouter(r *gin.Engine) {
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
	r.GET("/messages", func(c *gin.Context) {
		data := &IndexInfo{
			UserName: getUserName(c),
		}
		c.HTML(http.StatusOK, "message.html", data)
	})
	r.GET("/messages-list", func(c *gin.Context) {
		var messages []Messages
		rows, err := mysql.DefaultDB.Query("SELECT name, message,created_at  FROM messages order by id desc limit 50")
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Failed to retrieve messages")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m Messages
			if err := rows.Scan(&m.Name, &m.Message, &m.Time); err != nil {
				log.Println(err)
				continue
			}
			messages = append(messages, m)
		}

		c.JSON(http.StatusOK, messages)
	})

	r.POST("/messages", func(c *gin.Context) {
		var m Messages
		if err := c.Bind(&m); err != nil {
			c.String(http.StatusBadRequest, "Invalid request body")
			return
		}

		_, err := mysql.DefaultDB.Exec("INSERT INTO messages (name, message) VALUES (?, ?)", m.Name, m.Message)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Failed to insert message")
			return
		}

		c.Redirect(http.StatusFound, "/messages")
	})

}
