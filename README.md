# danmu

轻量弹幕（IM）系统，基于GO

# Require

`GO` >= `1.8`

`KAFKA` >= `1.0`

# Demo

本项目为你提供了一个在线demo：[demo](http://danmu.jwlchina.cn)

# Getting Start

安装好 `GOLANG` 环境以及 `KAFKA`

获取项目源码
```bash
go get github.com/kong36088/danmu
```

到项目`config`目录下配置好你的kafka地址以及服务监听地址

为 `kafka` 创建一个topic `danmu`
```bash
kafka-topics.sh --create --topic danmu --zookeeper localhost:2181 --replication-factor 1 --partitions 1
```

本项目为提供了一个使用实例，在完成环境部署后，运行以下命令：
```bash
$ cd /path/to/your/application

$ go run example/serv/main.go
```

服务开启之后，用浏览器访问打开html文件 `example/client/room_list.html` 即可看到服务效果