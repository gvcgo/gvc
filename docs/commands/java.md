## gvc jdk -h
```shell
NAME:
   g java - Java jdk version management.

USAGE:
   g java command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use jdk.
   remote, r               Show available versions.
   local, l                Show installed versions.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

--------

### gvc jdk remote -h
```shell
NAME:
   g java remote - Show available versions.

USAGE:
   g java remote [command options] [arguments...]

OPTIONS:
   --cn, --zh, -z  Use injdk.cn as resource url. (default: false)
   --help, -h      show help
```
- -zh选项可以使用国内源(injdk.cn)，最低支持到jdk8；

### gvc jdk local
- 显示本地已安装的go版本有哪些

### gvc jdk use jdkxx
- 下载并使用版本为xx的jdk，并配置好环境变量

### gvc jdk rm jdkxx
- 卸载版本为xx的jdk

### gvc jdk ru
- 一键卸载目前不用的jdk
