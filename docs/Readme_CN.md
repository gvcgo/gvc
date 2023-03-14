## [En](https://github.com/moqsien/gvc)
---------

## 关于[gvc](https://github.com/moqsien/gvc)的一些美好的事情
---------
- 愿每一天都是不用加班的一天
- 愿每一份代码都不会成为“屎山”
- 愿写代码成为一件快乐的事情
- 愿你不在为傻X无脑的需求点头哈腰
- 愿你能在每一台机器上轻松搞定你需要的开发环境
- 愿你即使是一个小白也能快速上手，轻松学习
- 愿你致力于提高效率，而不是无意义地卷
### 基于上述这些个美好的愿景，gvc诞生了！！！
# GVC 是一个跨平台多机器开发环境配置管理工具，让你轻松使用vscode进行多语言开发。
目前，gvc拥有以下功能或特点：
- go编译器自动安装和添加环境变量，多版本轻松切换；
- java自动安装和添加环境变量，版本切换(jdk17 or jdk19)；
- vscode自动安装，一键安装插件(需要配置，也可以使用默认配置)，一键备份和同步插件信息、用户设置、快捷键设置到webdav网盘；
- neovim自动安装和配置，默认与vscode-neovim插件配合，有默认配置可以使用；
- hosts文件更新，加速github访问，对国内用户友好；
- 所有上述需要下载的地方，如果在国内较慢的，一般都有加速；
- 下载源可配置，如果你有更快的下载源，可以在gvc-config.yml中配置并注意保存；
- WebDAV网盘同步配置信息，可以一键将本地的包括gvc-config.yml在内的必要配置同步到网盘，在新机器上只需要使用这些配置就能重新搭建一样的开发环境；
- MacOS、Windows、Linux(暂未测试)全平台支持

gvc将要提供的功能或特点：
- Windows下的git.exe下载；
- Rust自动安装和加速；
- Java版本管理和加速；
- HomeBrew安装和加速；
- Python安装包加速；
- NodeJS自动安装和加速；
- Flutter自动安装；

## gvc具体功能展示
---------
### gvc -h
```shell
moqsien@iMac gvc % gvc -h  
NAME:
   gvc - gvc <Command> <SubCommand>...

USAGE:
   gvc [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   host, h, hosts        Manage system hosts file.
   go, g                 Go version control.
   vscode, vsc, vs, v    VSCode management.
   config, conf, cnf, c  GVC config file management.
   nvim, neovim, nv, n   GVC neovim management.
   java, jdk, j          GVC jdk management.
   help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
### gvc go -h
```shell
moqsien@iMac-Pro gvc % gvc go help
NAME:
   gvc go - Go version control.

USAGE:
   gvc go command [command options] [arguments...]

COMMANDS:
   remote, r                   Show remote versions.
   use, u                      Download and use version.
   local, l                    Show installed versions.
   remove-unused, ru           Remove unused versions.
   remove-version, rm          Remove a version.
   add-envs, env, e, ae        Add envs for go.
   search-package, sp, search  Search for third-party packages.
   help, h                     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc vs -h
```shell
moqsien@iMac-Pro gvc % gvc vs -h
NAME:
   gvc vscode - VSCode management.

USAGE:
   gvc vscode command [command options] [arguments...]

COMMANDS:
   install, i, ins               Automatically install vscode.
   install-extensions, ie, iext  Automatically install extensions for vscode.
   sync-extensions, se, sext     Push local installed vscode extensions info to remote webdav.
   get-settings, gs, gset        Get vscode settings(keybindings include) info from remote webdav.
   push-settings, ps, pset       Push vscode settings(keybindings include) info to remote webdav.
   help, h                       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc nv -h
```shell
moqsien@iMac-Pro gvc % gvc nv -h
NAME:
   gvc nvim - GVC neovim management.

USAGE:
   gvc nvim command [command options] [arguments...]

COMMANDS:
   install, ins, i  Install neovim.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc cnf -h
```shell
moqsien@iMac-Pro gvc % gvc cnf -h
NAME:
   gvc config - GVC config file management.

USAGE:
   gvc config command [command options] [arguments...]

COMMANDS:
   webdav, dav, w   Setup webdav account info to backup local settings for gvc, vscode, neovim etc.
   pull, pl         Pull settings to local backup dir from your remote webdav.
   push, ph         Push settings from local backup dir to your remote webdav.
   show, sh, s      Show path to conf files.
   reset, rs, r     Reset config file to default values.
   download, dl, d  Download example config files from gitee when backup dir is empty.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc host -h
```shell
moqsien@iMac-Pro gvc % gvc host -h
NAME:
   gvc host - Manage system hosts file.

USAGE:
   gvc host command [command options] [arguments...]

COMMANDS:
   fetch, f      Fetch github hosts info.
   fetchall, fa  Get all github hosts info with no ping filters.
   show, s       Show hosts file path.
   help, h       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc java -h
```shell
moqsien@iMac gvc % gvc java -h
NAME:
   gvc java - GVC jdk management.

USAGE:
   gvc java command [command options] [arguments...]

COMMANDS:
   use, u   Download and use jdk.
   show, s  Show available versions.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### 下载和安装
下载文件，解压，双击或者在命令行运行(不带子任何命令和参数)，即可安装到默认文件夹。
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc/releases/tag/v2)
