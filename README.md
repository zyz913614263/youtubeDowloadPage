# YouTube视频下载网站

[FreeVideos.cn](https://www.freevideos.cn/) 是一个免费的YouTube视频下载网站。

该项目使用Go语言和原生HTML进行开发。

## 网站启动方式
```bash
chmod +x goweb
sudo ./goweb > log 2>&1 &
```

## 依赖项
- MySQL数据库
- Docker（可选，用于快速安装MySQL）

### 使用Docker安装MySQL

如果你安装了Docker，你可以使用以下命令来快速安装MySQL：

```bash
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql
```

## Linux版本编译

```bash
GOOS=linux go build -o goweb
```
## 致谢
