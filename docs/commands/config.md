## gvc config -h
```shell
NAME:
   g config - Config file management for gvc.

USAGE:
   g config command [command options] [arguments...]

COMMANDS:
   webdav, dav, w  Setup webdav account info.
   pull, pl        Pull settings from remote webdav and apply them to applications.
   push, ph        Gather settings from applications and sync them to remote webdav.
   reset, rs, r    Reset the gvc config file to default values.
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

---------

### gvc cnf dav
- 开始配置gvc所使用的webdav的账号信息
- 默认使用坚果云的webdav，但是需要自行提供已配置好的Webdav账号，坚果云示例：https://github.com/moqsien/easynotes/blob/main/usage.md
- 坚果云官网：https://www.jianguoyun.com/ (请自行注册并配置webdav)

### gvc cnf pull
- 从webdav上拉取配置，并应用到相应的软件；
- 可以在~/.gvc/backup/gvc-config.yml的"files_to_sync"中，按照例子配置对应的需要同步的文件’
- 占位符"$home$"表示家目录，windows也支持；
- 占位符"$appdata$"表示windows的%APPDATA%变量所代表的目录，常用于存放软件的配置等，可自行百度；

### gvc cnf push
- 根据~/.gvc/backup/gvc-config.yml的"files_to_sync"中的信息，将本地的配置文件收集起来，并上传到webdav；
- 具体见gvc cnf pull如何配置；

### gvc cnf reset
- 将gvc自身的配置文件，也就是gvc-config.yml重置成默认状态
