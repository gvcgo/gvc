## gvc julia -h
```shell
NAME:
   g julia - Julia version management.

USAGE:
   g julia command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use julia.
   remote, r               Show available versions.
   local, l                Show installed versions.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

---------

### 版本比较多，所以remote子命令仅支持稳定版本
### 其他子命令与go和jdk类似
### 安装使用国内源加速
