# YouTube视频下载网站
  该项目使用Go语言和原生HTML进行开发。

## 依赖项
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

## 网站启动方式
  配置文件 config.yaml 包含了与数据库相关的配置信息，可自行修改。
  
  ```bash
  # 将整个文件夹拷贝到Linux服务器上
  cd youtubeDowloadPage
  chmod +x goweb 
  sudo ./goweb > log 2>&1 &
  ```

## go代码修改后，使用下面的命令编译

  ```bash
  GOOS=linux go build -o goweb
  ```
## 参考 https://github.com/develon2015/Youtube-dl-REST
