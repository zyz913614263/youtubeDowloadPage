# YouTube视频下载网站
  该项目使用Go语言和原生HTML进行开发。
  [go.dev](https://go.dev/) 1.21.7 版本
  使用yt-dlp开源工具下载youtube视频
  [yt-dlp ](https://github.com/yt-dlp/yt-dlp/releases/tag/2025.02.19)
  自己实现了下载资源界面
  ![Uploading image.png…]()


## go代码，使用下面的命令编译

  ```bash
  GOOS=linux go build -o goweb
  ```
## 网站启动方式
  配置文件 config.yaml 包含了与数据库相关的配置信息，可自行修改。
  
  ```bash
  # 将整个文件夹拷贝到Linux服务器上
  cd youtubeDowloadPage
  chmod +x goweb 
  sudo ./goweb > log 2>&1 &
  ```



## 下载自己专用的cookie 文件
要将 YouTube 的 Cookie 导出到文件以便与 `yt-dlp` 一起使用，您可以按照以下步骤操作。请注意，导出 Cookie 的目的是为了访问需要登录的内容（如私人播放列表、年龄限制视频或会员专属内容），但请谨慎使用，以避免账户被封禁。

---

### **步骤 1：使用隐私/无痕模式登录 YouTube**
1. 打开浏览器的隐私/无痕窗口（例如 Chrome 的“无痕模式”或 Firefox 的“隐私窗口”）。
2. 访问 [YouTube](https://www.youtube.com) 并登录您的账户。
3. 登录后，打开一个新的标签页（不要关闭 YouTube 标签页）。
4. 关闭 YouTube 标签页，但保持隐私窗口打开。

---

### **步骤 2：安装 Cookie 导出扩展**
大多数浏览器没有内置的 Cookie 导出功能，因此您需要安装一个扩展程序来导出 Cookie。以下是常用的扩展：
- **Chrome/Firefox**: [Get cookies.txt] (https://www.chajianxw.com/product-tool/29038.html)

安装扩展后，确保它在隐私/无痕窗口中启用（某些浏览器需要手动允许扩展在无痕模式下运行）。

---

### **步骤 3：导出 YouTube 的 Cookie**
1. 在隐私/无痕窗口中，访问 [YouTube](https://www.youtube.com)。
2. 点击浏览器工具栏中的扩展图标（如 Get cookies.txt）。
3. 选择导出 `youtube.com` 的 Cookie。
4. 将导出的 Cookie 文件保存为 `cookies.txt`。

---

### **步骤 4：使用 Cookie 文件与 yt-dlp**
将导出的 `cookies.txt` 文件与 `yt-dlp` 一起使用。例如：
```bash
yt-dlp --cookies cookies.txt <YouTube视频链接>
```

---

### **注意事项**
1. **Cookie 有效期**：导出的 Cookie 可能会在一段时间后失效，尤其是当 YouTube 检测到异常活动时。
2. **账户安全**：避免频繁使用账户下载内容，以免触发 YouTube 的安全机制。
3. **隐私模式**：确保在隐私/无痕窗口中操作，以防止 Cookie 被轮换或失效。
4. **备用账户**：建议使用备用账户进行操作，避免主账户被封禁。

---

## 其他依赖项（可选）
  - Docker（可选，用于快速安装MySQL）
  ```bash
          sudo yum install -y yum-utils
          sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo    
          sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
          sudo systemctl start docker
  ```
  - MySQL数据库
  ```bash
          如果你安装了Docker，你可以使用以下命令来快速安装MySQL：
          docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
  ```
      
  - Redis localhost:6379
    ```bash
        可以通过docker 安装
        sudo docker run -itd --name myredis -v /home/ec2-user/redis:/data -p 6379:6379 redis
    ```
  - yt-dlp 
    ```bash
        sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
        sudo chmod a+rx /usr/local/bin/yt-dlp
        sudo yt-dlp -U
    ```


通过以上步骤，您可以成功导出 YouTube 的 Cookie 并与 `yt-dlp` 一起使用。
## 前端界面代码参考 https://github.com/develon2015/Youtube-dl-REST
