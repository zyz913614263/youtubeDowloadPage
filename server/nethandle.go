package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	"golang.org/x/crypto/bcrypt"
	"zyz.com/m/config"
	"zyz.com/m/mysql"
	"zyz.com/m/redis"
)

const RequestKey = "request"
const ParseKey = "parse"

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

func getCookieArg() string {
	if config.DefaultConfig.CookiesFile != "" {
		return fmt.Sprintf("--cookies %s", config.DefaultConfig.CookiesFile)
	}
	return ""
}

func parseHandler(c *gin.Context) {

	//var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request
	username := getUserName(c)
	if username == "" {
		// 如果用户名为空，则将其设置为空字符串
		//c.Redirect(http.StatusSeeOther, "/login")
		//return
	}

	// 获取第一个参数名
	videoURL := r.URL.RawQuery
	videoURL, _ = url.QueryUnescape(videoURL)
	videoURL = strings.Replace(videoURL, "y2b", "youtube", -1)
	videoURL = strings.Replace(videoURL, "y2", "youtu", -1)
	log.Printf("parseHandler %v", videoURL)

	var audios []Audio
	var videos []Video
	var bestAudio Audio
	var bestVideo Video

	var rs Result

	msg := &Message{URL: videoURL, Website: r.RemoteAddr}
	if runtime.GOOS == "windows" {
		cmd := fmt.Sprintf(" --print-json --skip-download %s '%s'  ", getCookieArg(), videoURL)
		fmt.Println("解析视频, 命令:", cmd)
		args := []string{
			"--print-json",
			"--skip-download",
			"--cookies",
			config.DefaultConfig.CookiesFile,
			videoURL,
		}
		output, err := exec.Command("yt-dlp", args...).Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = json.Unmarshal(output, &rs)
		if err != nil {
			cmd := fmt.Sprintf(" --print-json --skip-download %s '%s?p=1'  ", getCookieArg(), videoURL)
			fmt.Println("尝试分P, 命令:", cmd)
			args = []string{
				"--print-json",
				"--skip-download",
				"--cookies",
				config.DefaultConfig.CookiesFile,
				videoURL + "?p=1",
			}
			output, err = exec.Command("yt-dlp", args...).Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = json.Unmarshal(output, &rs)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			msg.P = "1"
			msg.URL = msg.URL + "?p=1"
		}
	} else {
		cmd := fmt.Sprintf("yt-dlp --print-json --skip-download %s '%s' 2> /dev/null", getCookieArg(), videoURL)
		fmt.Println("解析视频, 命令:", cmd)
		output, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = json.Unmarshal(output, &rs)
		if err != nil {
			cmd := fmt.Sprintf("yt-dlp --print-json --skip-download %s '%s?p=1' 2> /dev/null", getCookieArg(), videoURL)
			fmt.Println("尝试分P, 命令:", cmd)
			output, err = exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = json.Unmarshal(output, &rs)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			msg.P = "1"
			msg.URL = msg.URL + "?p=1"
		}
		fmt.Println("解析完成:", rs.Title, msg.URL)
	}

	for _, it := range rs.Formats {
		length := ""
		length += strconv.FormatFloat(float64(it.FileSize+it.FileSizeApprox)/1024/1024, 'f', 2, 64)
		if it.AudioExt != "none" {
			audios = append(audios, Audio{
				Id:     it.FormatID,
				Format: it.URL,
				Rate:   it.ABR,
				Info:   it.FormatNote,
				Size:   length,
			})
		} else if it.VideoExt != "none" {
			videos = append(videos, Video{
				Id:     it.FormatID,
				Format: it.URL,
				//Resolution: it.Resolution,
				Scale: it.Height,
				Rate:  it.VBR,
				Info:  it.FormatNote,
				Size:  length,
			})
		}
	}

	sort.Slice(audios, func(i, j int) bool { return audios[i].Rate < audios[j].Rate })
	sort.Slice(videos, func(i, j int) bool { return videos[i].Rate < videos[j].Rate })

	bestAudio = audios[len(audios)-1]
	bestVideo = videos[len(videos)-1]

	// parseSubtitle(msg) 方法尚未实现

	result := &HResult{
		Website:   msg.Website,
		V:         msg.VideoID,
		P:         msg.P,
		Title:     rs.Title,
		Thumbnail: rs.Thumbnail,
		Best: Best{
			Audio: bestAudio,
			Video: bestVideo,
		},
		Available: Available{
			Audios: audios,
			Videos: videos,
		},
	}
	redis.AddCount(ParseKey)

	log.Printf("data=%v", result)
	c.JSON(http.StatusOK, gin.H{
		"result": result, // 返回数据 result
	})
}

func handleIndex(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request
	username := getUserName(c)
	if username == "" {
		// 如果用户名为空，则将其设置为空字符串
		//c.Redirect(http.StatusSeeOther, "/login")
		//return
	}

	videoURL := r.FormValue("url")
	/*if videoURL == "" || !strings.Contains(videoURL, "youtube") {
		http.Error(w, "Invalid YouTube video URL", http.StatusBadRequest)
		return
	}*/
	log.Printf("handleIndex handleIndex:%v username:%v", videoURL, username)
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
	if err := LoadCookies(jar, config.DefaultConfig.CookiesFile); err != nil {
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
	redis.AddCount(ParseKey)
	data := &IndexInfo{
		VideoURL:   videoURL,
		VideoLinks: videoLinks,
		AudioLinks: audioLinks,
		UserName:   username,
	}
	//log.Printf("data=%v", data)
	c.HTML(http.StatusOK, "index.html", data)
}

func ProxyHandler(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request
	urlP := r.URL.Query().Get("url")
	if !isValidURL(urlP) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
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
	tran := &http.Client{
		Transport: transport,
	}
	resp, err := tran.Get(urlP)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func isValidURL(inputURL string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return false
	}

	if !(parsedURL.Host == "i.ytimg.com" || parsedURL.Host == "i0.hdslb.com" || parsedURL.Host == "i1.hdslb.com" || parsedURL.Host == "i2.hdslb.com" || parsedURL.Host == "i3.hdslb.com" || parsedURL.Host == "i4.hdslb.com" || parsedURL.Host == "i5.hdslb.com" || parsedURL.Host == "i6.hdslb.com" || parsedURL.Host == "i7.hdslb.com" || parsedURL.Host == "i8.hdslb.com" || parsedURL.Host == "i9.hdslb.com") {
		return false
	}

	return true
}
