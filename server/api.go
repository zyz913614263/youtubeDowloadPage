package server

import (
	"bufio"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"zyz.com/m/redis"
)

func LoadCookies(jar *cookiejar.Jar, filename string) error {
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

func RequestCounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		redis.AddCount(RequestKey)

		// 处理请求
		c.Next()
	}
}
