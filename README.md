# danmu

轻量弹幕（IM）系统，基于GO

# Require

`GO`

`KAFKA`

# Demo

本项目为你提供了一个在线demo：![demo](danmu.jwlchina.cn)

# Getting Start

到`config`目录下配置好你的kafka地址以及服务监听地址

本项目为提供了一个使用实例，在完成环境部署后，运行以下命令：
```bash
$ cd /path/to/your/application

$ go run example/serv/main.go
```

服务开启之后，直接访问打开 `example/client/room_list.html` 即可看到服务效果