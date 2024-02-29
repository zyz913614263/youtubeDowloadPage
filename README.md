# youtubeDowloadPage

免费 youtube，油管视频下载网站

https://www.freevideos.cn/



开发使用go语言和原生的html

依赖数据库mysql
如果你安装了docker可以使用docker安装
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql


linux版本编译方式

GOOS=linux go build -o goweb
