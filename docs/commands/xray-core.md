## gvc xray -h
```shell
NAME:
   g xray - Start Xray Shell for free VPN.

USAGE:
   g xray [command options] [arguments...]

OPTIONS:
   --start, --st, -s  Start Xray Client. (default: false)
   --help, -h         show help
```

----------

### gvc xray
```shell
*** Xray Shell Start ***
>>> help

Commands:
  clear        clear the screen
  exit         exit the program
  help         display help
  omega        Download Switchy-Omega for GoogleChrome. # 自动下载Google浏览器插件Switchy-Omega到~/.gvc/proxy_files目录，添加到浏览器插件即可使用，具体如何添加和配置请百度之
  restart      restart xray client. # 重启xray-core进程，怀疑 xray-core进程由于某些未知原因假死，可以以此重启
  start        Start an Xray Client. # 当当前系统中还没有启动任何xray-core进程时，启动xray-core进程
  stop         Stop an Xray Client. # 结束当前系统中正在后台运行xray-core进程
  vmess        Fetch proxies from vmess sources.  # 刷新免费vpn列表
```

- 开启gvc实现的xray-core的交互式shell工具，在shell工具内可以控制xray-core的启、停以及强制刷新免费vpn列表等；
- -s选项用户无需关心，是用于与shell之间交互的，作用是启动一个后台运行的xray-core实例，默认为false就行；
- 可以到~/.gvc/backup/gvc-config.yml中的proxy下suburls中添加你自己的免费或者收费的vmess订阅地址；
- gvc会每隔一段时间之后自动帮你筛选和切换可用的响应比较快的免费vpn，你也可以在shell中手动触发筛选和切换；
