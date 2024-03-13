package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"zyz.com/m/config"
)

func TestDownload(t *testing.T) {
	config.InitConfig()
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

	// TikTok 视频链接
	videoURL := "https://www.tiktok.com/@cat_and_cat23/video/7332468233045331232"

	// 发起 HTTP GET 请求获取页面内容
	resp, err := tran.Get(videoURL)
	if err != nil {
		fmt.Println("Error fetching video page:", err)
		return
	}
	defer resp.Body.Close()

	// 读取页面内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 提取视频链接
	videoURL = extractVideoURL(string(body))
	fmt.Println("Video Download URL:", videoURL)
}

func extractVideoURL(body string) string {
	// 找到视频链接的开始和结束位置
	startTag := `videoObject`
	endTag := `</script>`

	startIndex := strings.Index(body, startTag)
	if startIndex == -1 {
		fmt.Println("Start tag not found")
		return ""
	}
	startIndex = strings.Index(body[startIndex:], "<script>")
	if startIndex == -1 {
		fmt.Println("Start tag not found")
		return ""
	}
	endIndex := strings.Index(body[startIndex:], endTag)
	if endIndex == -1 {
		fmt.Println("End tag not found")
		return ""
	}
	endIndex += startIndex

	// 提取视频链接
	content := body[startIndex:endIndex]
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "videoUrl") {
			parts := strings.Split(line, `"`)
			for i, part := range parts {
				if part == "videoUrl" && i+1 < len(parts) {
					return parts[i+2]
				}
			}
		}
	}

	return ""
}
